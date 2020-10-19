package routers

import (
	"fmt"
	"context"
	"time"
	"path"
	"os"
	"log"
	// "errors"

	"net/http"
	"encoding/json"
	"html/template"

    // "github.com/twinj/uuid"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"andrewViscu/golang1/pkg/db"
	"andrewViscu/golang1/pkg/middleware"

)

type User struct{
	Id primitive.ObjectID `bson:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	City string `json:"city,omitempty"`
	Age int `json:"age,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
  }



var userCollection = db.DBConnect().Database("foo").Collection("bar")

var tmpl = template.Must(template.ParseFiles(path.Join("public", "Index.html")))

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func CheckToken(tokenString string) {
	hmacSampleSecret := []byte(os.Getenv("ACCESS_SECRET"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		return hmacSampleSecret, nil
	})
	if err  != nil {
		fmt.Println(err)
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Claims: ", claims["user_id"],  claims["authorized"])
	}
}

// (POST /login)
func Login(w http.ResponseWriter, r *http.Request) {

	var body, user User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Print(err)
	}

	opts := options.FindOne().SetCollation(&options.Collation{})	
	res := userCollection.FindOne(ctx, bson.M{"username": body.Username}, opts)
	// _id, err := primitive.ObjectIDFromHex(result.Id)
	err = res.Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"` + err.Error() + `", "response": 500}`))
		return
	}
	res.Decode(&user)
	
	if !CheckPasswordHash(body.Password, user.Password) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "Wrong password", "response": 500 }`))
		return
	}

	stringID := user.Id.Hex()
	fmt.Println(stringID)
	token, err := CreateToken(stringID)
	if err != nil {
	   w.WriteHeader(http.StatusUnprocessableEntity)
	   w.Write([]byte(`{"message": "Something's wrong, I can feel it.", "response":422}`))
	   return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token": ` + token + `, "response":200}`))
	CheckToken(token)
}

func CreateToken(userId string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
	   return "", err
	}
	return token, nil
}



func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results []primitive.M                                   //slice for multiple documents
	cur, err := userCollection.Find(ctx, bson.M{}) //returns a *mongo.Cursor
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `",
		 "response": 500}`))
		return
	}
	defer cur.Close(ctx) // close the cursor once stream of documents has exhausted
	for cur.Next(ctx) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	json.NewEncoder(w).Encode(results)
	// w.Write([]byte())
	// tmpl.Execute(w, results)
	// fmt.Print(results)
}


// (POST /register)
func CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for adding Content-type

	var body User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := json.NewDecoder(r.Body).Decode(&body) // storing in user variable of type user
	if err != nil {
		fmt.Print(err)
	}


	if !middleware.Password(body.Password){ //
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Impossible password", "description":"Password should have: 8+ characters, an uppercase letter, a lowercase letter and a number.", "response": 500`))
		return
	}

	body.Password, err = HashPassword(body.Password)
	body.CreatedAt = time.Now()
	body.Id = primitive.NewObjectID()
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := userCollection.InsertOne(ctx, body)
	if err != nil {
		log.Print("An error occured while INSERTING!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `", "response": 500}`))
		return
	}

	fmt.Printf("Created user '%v': %v\n", body.Username, insertResult)
	json.NewEncoder(w).Encode(insertResult) // return the mongodb ID of generated document

}

// (GET /users/{id})
func GetUser(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(req)["id"] //get Parameter value as string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	opts := options.FindOne().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res := userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: _id}}, opts)
	err = res.Err()
	if err != nil {
		log.Print("An error occured while GETTING USER!\n")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `",
		 "response": 500}`))
		 return
	}
	var result primitive.M //  an unordered representation of a BSON //document which is a Map
	res.Decode(&result)
    json.NewEncoder(w).Encode(result)

}
// (PUT /users/{id})
func UpdateUser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		Name string `json:"name,omitempty"` 
		City string `json:"city,omitempty"`
		Age int 	`json:"age,omitempty"`
	}	
	var body updateBody
	e := json.NewDecoder(req.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	params := mux.Vars(req)["id"] //get Parameter value as string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	filter := bson.D{{Key:"_id",Value: _id}} // converting value to BSON type
	after := options.After          // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.M{
		"$set": body,
	}
	log.Print(update)
	updateResult := userCollection.FindOneAndUpdate(ctx, filter, update, &returnOpt)
	err = updateResult.Err()
	if err != nil {
		log.Print("An error occured while UPDATING!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `",
		 "response": 500}`))
		 return
	}
	var result primitive.M
	_ = updateResult.Decode(&result)
	json.NewEncoder(w).Encode(result)
}

// (DELETE /users/{id})
func DeleteUser(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)["id"] //get Parameter value as string


	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	opts := options.FindOneAndDelete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res := userCollection.FindOneAndDelete(context.TODO(), bson.D{{Key:"_id",Value: _id}}, opts)
	err = res.Err()
	if err != nil {
		log.Print("An error occured while DELETING!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `",
		 "response": 500}`))
		 return
	}
	var deletedDocument bson.M
	res.Decode(&deletedDocument)

	fmt.Printf("Document: %v\n", deletedDocument)
	w.Write([]byte(`{"message":"User ` + deletedDocument["username"].(string) + ` deleted.", "response":200}`))

}


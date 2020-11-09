package routers

import (
	// "os"
	"context"
	"log"
	"time"

	"encoding/json"
	"net/http"

	"andrewViscu/golang1/pkg/db"
	mw "andrewViscu/golang1/pkg/middleware"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//User struct
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	City      string             `json:"city,omitempty"`
	Age       int                `json:"age,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

var userCollection = db.Connect().Database("foo").Collection("bar")

//Login - (POST /login)
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body, user User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Print(err)
	}

	opts := options.FindOne().SetCollation(&options.Collation{})
	res := userCollection.FindOne(ctx, bson.M{"username": body.Username}, opts)
	err = res.Err()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	res.Decode(&user)

	if !mw.CheckPasswordHash(body.Password, user.Password) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "error": "Wrong password", "code": 500 }`))
		return
	}

	stringID := user.ID.Hex()
	token, err := mw.CreateToken(stringID)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"error": "Something's wrong, I can feel it.", "code":422}`))
		return
	}
	log.Println("Login successful")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token": "` + token + `", "code":200}`))
	// http.Redirect(w, r, `/users/` + userId, http.StatusSeeOther)
}

//GetAllUsers - (GET /users)
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var results []primitive.M                      //slice for multiple documents
	cur, err := userCollection.Find(ctx, bson.M{}) //returns a *mongo.Cursor
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	defer cur.Close(ctx) // close the cursor once stream of documents has exhausted
	for cur.Next(ctx) {  //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

//CreateUser - (POST /register)
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := json.NewDecoder(r.Body).Decode(&body) // storing in user variable of type user
	if err != nil {
		log.Print(err)
	}

	if !mw.Password(body.Password) { //
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Impossible password", "description":"Password should have: 8+ characters, an uppercase letter, a lowercase letter and a number.", "code": 500`))
		return
	}

	body.Password, err = mw.HashPassword(body.Password) //Hash password and store it
	body.CreatedAt = time.Now()                         // Get current time
	body.ID = primitive.NewObjectID()                   // and new ID and store them
	if err != nil {
		log.Print(err)
	}
	insertResult, err := userCollection.InsertOne(ctx, body)
	if err != nil {
		log.Print("An error occured while INSERTING!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "error": "` + err.Error() + `", "code": 500}`))
		return
	}
	id := insertResult.InsertedID
	stringID := id.(primitive.ObjectID).Hex()
	token, err := mw.CreateToken(stringID)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"error": "Something's wrong, I can feel it.", "code":422}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token": "` + token + `", "code":200}`))
}

//GetUser - (GET /users/{id})
func GetUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)["id"] //get Parameter value as string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		log.Printf(err.Error())
	}
	opts := options.FindOne().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res := userCollection.FindOne(ctx, bson.D{{Key: "_id", Value: _id}}, opts)
	err = res.Err()
	if err != nil {

		log.Print("An error occured while GETTING USER!\n")
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "document not found: ` + err.Error() + `", "code":404}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "error": "` + err.Error() + `",
		 "code": 500}`))
		return
	}
	var result primitive.M
	res.Decode(&result)
	json.NewEncoder(w).Encode(result)

}

//UpdateUser - (PUT /users/{id})
func UpdateUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type updateBody struct {
		Name string `json:"name,omitempty"`
		City string `json:"city,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
	var body updateBody

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {

		log.Print(err)
	}
	params := mux.Vars(req)["id"] //get Parameter value as string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		log.Printf(err.Error())
	}
	filter := bson.D{{Key: "_id", Value: _id}} // converting value to BSON type
	after := options.After                     // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := bson.M{
		"$set": body,
	}
	updateResult := userCollection.FindOneAndUpdate(ctx, filter, update, &returnOpt)
	err = updateResult.Err()
	if err != nil {
		log.Print("An error occured while UPDATING!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "error": "` + err.Error() + `",
		 "code": 500}`))
		return
	}
	var result primitive.M
	w.WriteHeader(http.StatusOK)
	updateResult.Decode(&result)
	json.NewEncoder(w).Encode(result)

}

//DeleteUser (DELETE /users/{id})
func DeleteUser(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)["id"] //get Parameter value as string

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		log.Printf(err.Error())
	}
	opts := options.FindOneAndDelete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res := userCollection.FindOneAndDelete(context.TODO(), bson.D{{Key: "_id", Value: _id}}, opts)
	err = res.Err()
	if err != nil {
		log.Print("An error occured while DELETING!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "error": "` + err.Error() + `",
		 "code": 500}`))
		return
	}
	var deletedDocument bson.M
	res.Decode(&deletedDocument)
	w.Write([]byte(`{"message":"User ` + deletedDocument["username"].(string) + ` deleted.", "code":200}`))

}

//Index - (GET /)
func Index(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`Index page was added! Hooray!`))
}

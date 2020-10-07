package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"context"
	"time"
	"path"
	// "io/ioutil"
	"html/template"
	"log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "github.com/gorilla/mux"
)

type User struct{
	Name string `json:"name"`
	City string `json:"city"`
	Age int `json:"age"`
}

var userCollection = db().Database("foo").Collection("bar")


var tmpl = template.Must(template.ParseFiles(path.Join("public", "Index.html")))

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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
	// encodedResult := json.NewEncoder(w).Encode(results)
	tmpl.Execute(w, results)
	// fmt.Print(results)

}

func createProfile(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // for adding Content-type

	var person User
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := json.NewDecoder(r.Body).Decode(&person) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := userCollection.InsertOne(ctx, person)
	if err != nil {
		log.Print("An error occured while INSERTING!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `",
		 "response": 500}`))
		return
	}

	fmt.Println("Inserted a single document: ", insertResult)
	json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document

}


func getUserProfile(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(req)["id"] //get Parameter value as string
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	opts := options.FindOne().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res := userCollection.FindOne(ctx, bson.D{{"_id", _id}}, opts)
	err = res.Err()
	if err != nil {
		log.Print("An error occured while GETTING USER!")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "message": "` + err.Error() + `",
		 "response": 500}`))
		 return
	}
	var result primitive.M //  an unordered representation of a BSON //document which is a Map
	res.Decode(&result)
    json.NewEncoder(w).Encode(result)

}

func updateProfile(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		Name string `json:"name"` 
		City string `json:"city"`
		Age int 	`json:"age"`
	}	
	var body updateBody
	e := json.NewDecoder(req.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	params := mux.Vars(req)["id"] //get Parameter value as string
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	filter := bson.D{{"_id", _id}} // converting value to BSON type
	after := options.After          // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	update := body
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

// //Delete Profile of User

func deleteProfile(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)["id"] //get Parameter value as string

	_id, err := primitive.ObjectIDFromHex(params) // convert params to mongodb Hex ID
	if err != nil {
		fmt.Printf(err.Error())
	}
	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", _id}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	json.NewEncoder(w).Encode(res.DeletedCount) // return number of documents deleted

}


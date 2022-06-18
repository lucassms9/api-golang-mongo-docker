package controllers

import (
	mongodb "api-mongo/mongodb"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// struct for storing data
type User struct {
	Name string `json:name`
	Age  int    `json:age`
	City string `json:city`
}

var userCollection = mongodb.Db().Database("goTest").Collection("users") // get collection "users" from db() which returns *mongo.Client

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // for adding       //Content-type

	var person User

	err := json.NewDecoder(r.Body).Decode(&person) // storing in person   //variable of type user

	if err != nil {
		fmt.Println(err)
	}

	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
	json.NewEncoder(w).Encode(insertResult.InsertedID) // return the mongodb ID of generated document

}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body User

	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Println(e)
	}

	var result primitive.M = make(primitive.M, 0)

	err := userCollection.FindOne(context.TODO(), bson.D{{"name", body.Name}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(result)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var results []primitive.M = make([]primitive.M, 0)

	//slice for multiple documents
	cur, err := userCollection.Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {

		fmt.Println(err)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	fmt.Println(results)
	json.NewEncoder(w).Encode(results)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type updateBody struct {
		Name string `json:"name"` //value that has to be matched
		City string `json:"city"` // value that has to be modified
	}

	var body updateBody
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {

		fmt.Print(e)
	}
	filter := bson.D{{"name", body.Name}}
	after := options.After
	returnOpt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	update := bson.D{{"$set", bson.D{{"city", body.City}}}}

	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)["id"]
	fmt.Println(params)
	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		fmt.Print(err)
	}

	res, e := userCollection.DeleteOne(context.TODO(), bson.D{{"_id", _id}})
	if err != nil {
		log.Fatal(e)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	json.NewEncoder(w).Encode(res.DeletedCount)
}

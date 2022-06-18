package server

import (
	controllers "api-mongo/controllers"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	group string
}

func Init() server {
	var newServer server
	newServer.group = "/api"
	return newServer
}

func (e *server) StartServer() {
	fmt.Println("[API] Started on")

	route := mux.NewRouter()
	s := route.PathPrefix(e.group).Subrouter() //Base Path

	//Routes
	s.HandleFunc("/createProfile", controllers.CreateProfile).Methods("POST")
	s.HandleFunc("/getAllUsers", controllers.GetAllUsers).Methods("GET")
	s.HandleFunc("/getUserProfile", controllers.GetUserProfile).Methods("POST")
	s.HandleFunc("/updateProfile", controllers.UpdateProfile).Methods("PUT")
	s.HandleFunc("/deleteProfile/{id}", controllers.DeleteProfile).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", s)) // Run Server
}

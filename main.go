package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//User ...
type User struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}

var users []User

func main() {

	router := mux.NewRouter()

	users = append(users, User{ID: "1", Email: "majidmokht@gmail.com"})
	users = append(users, User{ID: "2", Email: "saeedmokht@gmail.com"})

	router.HandleFunc("/users", GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", CreateUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

//GetAllUsers ...
func GetAllUsers(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(users)
}

//GetUser ...
func GetUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range users {
		if item.ID == params["id"] {
			fmt.Println(item)
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

//CreateUser ...
func CreateUser(w http.ResponseWriter, req *http.Request) {
	//get params from request
	params := mux.Vars(req)
	//create new user type and add body to it
	user := User{}
	user.ID = params["id"]
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}
	//add new user to the rest of users
	users = append(users, user)
	//send updated users back to browser
	json.NewEncoder(w).Encode(users)
}

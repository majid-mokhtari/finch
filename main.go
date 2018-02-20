package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//User ...
type User struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

//MyError ...
type MyError struct {
	Err string `json:"err"`
}

var users []User

func main() {

	router := mux.NewRouter()

	users = append(users, User{ID: "1", Email: "A", Password: "A"})
	users = append(users, User{ID: "2", Email: "B", Password: "B"})

	router.HandleFunc("/users", LoginUser).Methods("POST")
	router.HandleFunc("/users/{id}", GetUserByID).Methods("GET")
	router.HandleFunc("/users/{id}", CreateUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

//LoginUser ...
func LoginUser(w http.ResponseWriter, req *http.Request) {
	body, error := ioutil.ReadAll(req.Body)
	if error != nil {
		fmt.Println(error)
	}
	user := make(map[string]string)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users {
		if u.Email == user["email"] && u.Password == user["password"] {
			err := json.NewEncoder(w).Encode(user)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	myError := MyError{}
	myError.Err = "User not found"
	e := json.NewEncoder(w).Encode(&myError)
	if e != nil {
		log.Fatal(err)
	}
}

//GetUserByID ...
func GetUserByID(w http.ResponseWriter, req *http.Request) {
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

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
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
	Birthday string `json:"birthday,omitempty"`
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

	router.HandleFunc("/user/login", LoginUser).Methods("POST")
	router.HandleFunc("/user/signup", SignupUser).Methods("POST")
	router.HandleFunc("/user/{id}", GetUserByID).Methods("GET")

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
	myError.Err = "Email or password is incorrect!"
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

//SignupUser ...
func SignupUser(w http.ResponseWriter, req *http.Request) {
	body, error := ioutil.ReadAll(req.Body)
	if error != nil {
		log.Fatal(error)
	}
	//create new user type and add body to it
	newUser := User{}
	json.Unmarshal(body, &newUser)
	users = append(users, newUser)
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Fatal(err)
	}
}

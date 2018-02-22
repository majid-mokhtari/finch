package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//User ...
type User struct {
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	City     string    `json:"city,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
	Active   bool      `json:"active,omitempty"`
}

//MyError ...
type MyError struct {
	Err string `json:"err"`
}

var users []User

var dataSourceName = "root:root@tcp(127.0.01:3306)/finch"

func main() {

	//Routes
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
	myError.Err = "Incorrect email or password!"
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

	//DB connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println(err)
	}

	//Query Statements
	stmt, err := db.Prepare("insert into finch.users(name, email, password, birthday, city, active) values(?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, errr := stmt.Exec(newUser.Name, newUser.Email, newUser.Password, newUser.Birthday, newUser.City, 1)
	if errr != nil {
		log.Fatal(errr)
	}
	//send back new user
	json.NewEncoder(w).Encode(newUser)
	defer db.Close()
}

package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/finch-app/finch/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//MyError ...
type MyError struct {
	Err string `json:"err"`
}

var dataSourceName = "root:root@tcp(127.0.01:3306)/finch?parseTime=true"

var users []models.User

func main() {

	GetAllUsers()

	//Routes
	router := mux.NewRouter()

	router.HandleFunc("/user/login", LoginUser).Methods("POST")
	router.HandleFunc("/user/signup", SignupUser).Methods("POST")
	router.HandleFunc("/user/{id}", GetUserByID).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))

}

//GetAllUsers ...
func GetAllUsers() {
	//DB connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println(err)
	}

	rows, errr := db.Query("select id,name,birthdate,city,email,password,status from finch.users where active = 1")
	if errr != nil {
		log.Fatal(errr)
	}

	defer db.Close()

	for rows.Next() {
		user := models.User{}

		var Birthdate string

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&Birthdate,
			&user.City,
			&user.Email,
			&user.Password,
			&user.Status,
		)
		if err != nil {
			log.Fatalln(err)
		}
		//user.Birthday = time.Parse("2006-01-02", Birthday)
		users = append(users, user)
	}
}

//LoginUser ...
func LoginUser(w http.ResponseWriter, req *http.Request) {
	body, error := ioutil.ReadAll(req.Body)
	if error != nil {
		log.Fatal(error)
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
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&models.User{})
}

//SignupUser ...
func SignupUser(w http.ResponseWriter, req *http.Request) {
	body, error := ioutil.ReadAll(req.Body)
	if error != nil {
		log.Fatal(error)
	}
	//create new user type and add body to it
	newUser := models.User{}
	json.Unmarshal(body, &newUser)

	//DB connection
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Println(err)
	}

	//Query Statements
	stmt, err := db.Prepare("insert into finch.users(name, email, password, gender, birthdate, city, status, active) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, errr := stmt.Exec(newUser.Name, newUser.Email, newUser.Password, newUser.Gender, newUser.Birthdate, newUser.City, newUser.Status, 1)
	if errr != nil {
		log.Fatal(errr)
	}
	//send back new user
	json.NewEncoder(w).Encode(newUser)
	defer db.Close()
}

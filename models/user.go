package models

import "time"

//User ...
type User struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Birthdate time.Time `json:"birthdate,omitempty"`
	Gender    string    `json:"gender,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	City      string    `json:"city,omitempty"`
	Status    string    `json:"status,omitempty"`
}

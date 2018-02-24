package models

import "time"

//User ...
type User struct {
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
	City     string    `json:"city,omitempty"`
	Status   string    `json:"status,omitempty"`
}

package models

type User struct {
	Username string `json:"username"`
	Gmail    string `json:"gmail"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Password string `json:"password"`
}
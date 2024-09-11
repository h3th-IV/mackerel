package models

type User struct {
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	Location  string `json:"location"`
	IpAddress string `json:"ip_address"`
}

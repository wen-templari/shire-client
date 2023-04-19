package model

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type LoginResponse struct {
	User
	Token string `json:"token"`
}

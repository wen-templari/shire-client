package model

type User struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Address   string `json:"address,omitempty"`
	Port      int    `json:"port,omitempty"`
	Password  string `json:"password,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type UpdateUserRequest struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type LoginResponse struct {
	User
	Token string `json:"token"`
}

package model

type User struct {
	Id        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Address   string `json:"address,omitempty"`
	Port      int    `json:"port,omitempty"`
	RPCPort   int    `json:"rpcPort,omitempty" db:"rpcPort"`
	Password  string `json:"password,omitempty"`
	CreatedAt string `json:"createdAt,omitempty" db:"createdAt"`
	UpdatedAt string `json:"updatedAt,omitempty" db:"createdAt"`
}

type Group struct {
	Id        int    `json:"id,omitempty"`
	Users     []User `json:"users,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type CreateGroupRequest struct {
	Ids []int `json:"ids"`
}

type UpdateUserRequest struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	RPCPort int    `json:"rpcPort"`
}

type LoginResponse struct {
	User
	Token string `json:"token"`
}

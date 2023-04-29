package model

type Message struct {
	From      int    `json:"from"`
	To        int    `json:"to"`
	Content   string `json:"content"`
	GroupId   int    `json:"groupId" db:"groupId"`
	Time      string `json:"time"`
	CreatedAt string `json:"createdAt,omitempty" db:"createdAt"`
	UpdatedAt string `json:"updatedAt,omitempty" db:"updatedAt"`
}

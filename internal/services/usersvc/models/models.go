package models

type User struct {
	ID       int32  `json:"id,omitempty"`
	UserName string `json:"username,omitempty"`
}

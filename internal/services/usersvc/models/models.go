package models

type User struct {
	ID       int32  `json:"id,omitempty"`
	UserName string `json:"username,omitempty"`
}

type Updates struct {
	Message string `json: "message,omitempty"`
}

// type StreamRequest struct {

// }

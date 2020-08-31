package models

type Todo struct {
	UserId int64  `json:"userId"`
	Title  string `json:"title"`
}

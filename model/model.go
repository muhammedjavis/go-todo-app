package model

type ToDo struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

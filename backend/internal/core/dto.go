package core

type AbstractModelResponce struct {
	ID uint `json:"id"`
}

type AbstractNameModelResponce struct {
	AbstractModelResponce
	Name string `json:"name"`
}

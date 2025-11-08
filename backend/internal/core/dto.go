package core

type AbstractModelResponce struct {
	ID uint `json:"id"`
}

type AbstractNameModelResponce struct {
	AbstractModelResponce
	Name string `json:"name"`
}

type AbstractTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

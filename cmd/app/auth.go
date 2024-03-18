package main

var secretKey = []byte("supersecrettoken1488")

type RegisterRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Password string `json:"password" example:"Pass123"`
	Role     string `json:"role" example:"user"`
}

type LoginRequest struct {
	Name     string `json:"name" example:"Posos"`
	Password string `json:"password" example:"123"`
}

type LoginToken struct {
	AccessToken string `json:"access_token"`
}

type contextKey string

var roleKey = contextKey("role")

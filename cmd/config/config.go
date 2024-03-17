package config

import (
	"os"
	"strconv"
)

type Config struct{
	Host string
	Port int
	User string
	Password string 
	DB_name string
}

func GetConfig() Config {
	c := Config{}
	port := os.Getenv("PORT")
	if port != ""{
		c.Port, _ = strconv.Atoi(port)
	}
	c.Host = os.Getenv("HOST") 
	c.User = os.Getenv("DB_USER")
	c.Password = os.Getenv("DB_PASSWORD")
	c.DB_name = os.Getenv("DB_NAME")
	return c
}
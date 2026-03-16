package main

import (
	"log"
	"login/database"
	"login/handlers"
	"net/http"

	_ "login/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Login App API
// @version         1.0
// @description     A simple authentication API with register and login endpoints
// @host            localhost:8080
// @BasePath        /
func main() {

	client, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}

	http.HandleFunc("/register", handlers.RegisterHandler(client))
	http.HandleFunc("/login", handlers.LoginHandler(client))
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server is running on port :8080")
	log.Println("http://localhost:8080/login")
	log.Println("http://localhost:8080/register")
	log.Println("http://localhost:8080/swagger/index.html")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
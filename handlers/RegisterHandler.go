package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"login/database"
	"login/model"
	"login/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterHandler godoc
// @Summary      Register a new user
// @Description  Creates a new user with a hashed password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      model.User        true  "User credentials"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /register [post]
func RegisterHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user model.User

		// 1. Decode request body
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		// 2. Validate fields
		if user.Email == "" || user.Password == "" {
			http.Error(w, "Email and password are required", http.StatusBadRequest)
			return
		}

		collection := database.GetCollection(client, "Forms")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 3. Check if email already exists
		var existing model.User
		err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existing)
		if err == nil {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}

		// 4. Hash the password
		hashedPassword, err := services.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword

		// 5. Insert into MongoDB
		_, err = collection.InsertOne(ctx, user)
		if err != nil {
			http.Error(w, "Error registering user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
	}
}
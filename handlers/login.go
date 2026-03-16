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

// LoginHandler godoc
// @Summary      Login a user
// @Description  Validates credentials and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      model.User        true  "User credentials"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /login [post]
func LoginHandler(client *mongo.Client) http.HandlerFunc {
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

		// 3. Find user by email in MongoDB
		var storedUser model.User
		err = collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&storedUser)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// 4. Compare password with stored hash
		if !services.CheckPassword(storedUser.Password, user.Password) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// 5. Generate JWT using MongoDB ObjectID as user identifier
		tokenString, err := services.GenerateToken(storedUser.ID.Hex())
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	}
}
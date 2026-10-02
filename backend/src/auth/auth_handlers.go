package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"ad/jwtToken"
	"ad/postgre"
	"ad/utils"
)

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LogIn(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error in log in request : %v\n", err)
		return
	}

	hashedPassword := utils.Hash(req.Password)
	log.Printf("Log in request username(%s)\n", req.Username)

	// Check users in postgres by username and hashed Password
	// get id from that
	userId, err := postgre.GetUserIdByLogIn(req.Username, hashedPassword)
	if err != nil {
		http.Error(w, "Invalid login or password", http.StatusBadRequest)
		log.Printf("Error in log in request username(%s) : %v\n", req.Username, err)
		return
	}

	tkn, err := jwtToken.GenToken(userId)
	if err != nil {
		http.Error(w, "Cannot generate a token", http.StatusInternalServerError)
		log.Printf("Error in log in request username(%s) : %v\n", req.Username, err)
		return
	}
	SetCookie(w, tkn)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tkn)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error in sign up request : %v\n", err)
		return
	}

	hashedPassword := utils.Hash(req.Password)
	log.Printf("Sign up request username(%s)\n", req.Username)

	// generate user id, and insert to postgres if username is unique
	userId, err := utils.GenerateId()
	if err != nil {
		http.Error(w, "Error during user creation", http.StatusInternalServerError)
		log.Printf("Error in sign up request username(%s) : %v\n", req.Username, err)
		return
	}
	adToken, err := utils.GenerateBase64(32)
	if err != nil {
		http.Error(w, "Error during user creation", http.StatusInternalServerError)
		log.Printf("Error in sign up request username(%s) : %v\n", req.Username, err)
		return
	}

	err = postgre.InsertUser(userId, req.Username, hashedPassword, adToken)
	if err != nil {
		http.Error(w, "Error during user insertion", http.StatusInternalServerError)
		log.Printf("Error in sign up request username(%s) : %v\n", req.Username, err)
		return
	}

	jwt, err := jwtToken.GenToken(userId)
	if err != nil {
		http.Error(w, "Cannot generate a token", http.StatusInternalServerError)
		log.Printf("Error in sign up request username(%s) : %v\n", req.Username, err)
		return
	}

	SetCookie(w, jwt)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwt)
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error in remove user request : user_id field is not found\n")
		return
	}
	log.Printf("Remove user request userId(%s)\n", userId)

	// delete from postgres by usedId
	err := postgre.RemoveByUserId(userId)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error in remove user request userId(%s) : %v\n", userId, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "The user was successfully deleted",
		"user_id": userId,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error in get user request : user_id field is not found\n")
		return
	}
	log.Printf("Get user request userId(%s)\n", userId)

	// get user from postgres by usedId
	username, token, err := postgre.GetUserByUserId(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		log.Printf("Error in get user request userId(%s) : %v\n", userId, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"username": username,
		"token":    token,
	})
}

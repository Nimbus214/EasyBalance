package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type User struct {
	ID      string
	Balance float64
}

type UserService struct {
	users map[string]*User
	mutex sync.Mutex
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*User),
	}
}

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.users[user.ID]; ok {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Error: User with ID %s already exists", user.ID)
		return
	}

	s.users[user.ID] = &user

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User created successfully")
}

func (s *UserService) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, ok := s.users[userID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error: User with ID %s not found", userID)
		return
	}

	balance := user.Balance
	response := map[string]float64{"balance": balance}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *UserService) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	type UpdateBalanceRequest struct {
		UserID  string  `json:"user_id"`
		Amount  float64 `json:"amount"`
		IsDebit bool    `json:"is_debit"`
	}

	var req UpdateBalanceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, ok := s.users[req.UserID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error: User with ID %s not found", req.UserID)
		return
	}

	if req.IsDebit {
		if user.Balance < req.Amount {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error: User with ID %s has insufficient balance", req.UserID)
			return
		}

		user.Balance -= req.Amount
	} else {
		user.Balance += req.Amount
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Balance updated successfully")
}

func main() {
	userService := NewUserService()

	http.HandleFunc("/users", userService.CreateUser)
	http.HandleFunc("/balance", userService.GetBalance)
	http.HandleFunc("/balance/update", userService.UpdateBalance)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rinha-backend/api/dto"
	"rinha-backend/api/repository"
	"strconv"
)

type Server struct {
	router   *http.ServeMux
	database *repository.Database
}

func NewHttpServer(database *repository.Database) *Server {
	server := http.NewServeMux()
	return &Server{server, database}
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		request dto.CreateUserRequest
	)

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.response(w, dto.TransactionErrorResponse{Message: "Could not create user", Error: err.Error()}, http.StatusUnprocessableEntity)
		return
	}

	user, err := s.database.CreateUser(ctx, request)
	if err != nil {
		s.response(w, dto.TransactionErrorResponse{Message: "Could not create user", Error: err.Error()}, http.StatusUnprocessableEntity)
		return
	}

	s.response(w, user, http.StatusOK)
}
func (s *Server) MakeTransact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	var transaction dto.Transaction

	_, err := s.database.GetUser(ctx, userID)
	if err != nil {
		s.response(w, dto.TransactionErrorResponse{Message: "Could not find user", Error: fmt.Sprintf("User id %d not found", userID)}, http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		s.response(w, dto.TransactionErrorResponse{Message: "Could not update funds", Error: err.Error()}, http.StatusUnprocessableEntity)
		return
	}

	done := s.database.MakeTransaction(ctx, userID, transaction)
	if !done {
		s.response(w, dto.TransactionErrorResponse{Message: "Could not update funds", Error: "Insufficient funds"}, http.StatusUnprocessableEntity)
		return
	}

	balance, err := s.database.GetBalance(ctx, userID)
	if err != nil {
		s.response(w, dto.TransactionErrorResponse{Message: "Could not get user balance", Error: "Could not get user balance"}, http.StatusUnprocessableEntity)
		return
	}

	s.response(w, balance, http.StatusOK)
}

func (s *Server) ListTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	statement, err := s.database.GetStatement(ctx, userID)
	if err != nil {
		s.response(w, dto.TransactionErrorResponse{Message: "Could not update funds", Error: err.Error()}, http.StatusUnprocessableEntity)
		return
	}
	s.response(w, statement, http.StatusOK)
}

func (s *Server) response(w http.ResponseWriter, data interface{}, status int) {
	response, _ := json.Marshal(data)
	w.WriteHeader(status)
	w.Write(response)
}

func (s *Server) Start() {
	s.router.HandleFunc("POST /clientes", s.CreateUser)
	s.router.HandleFunc("POST /clientes/{id}/transacoes", s.MakeTransact)
	s.router.HandleFunc("GET /clientes/{id}/extrato", s.ListTransactions)
	http.ListenAndServe(":8000", s.router)
}

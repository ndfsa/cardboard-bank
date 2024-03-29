package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ndfsa/cardboard-bank/cmd/api/repository"
	"github.com/ndfsa/cardboard-bank/internal/middleware"
)

const (
	// research a better way to share this key
	tokenKey = "2bbb515c1311dd69a609a0d553dc7ac1ac8eadc2b22daa9aaa99483d2f381374"
)

func main() {
	db, err := sql.Open("pgx", "postgres://back:root@db:5432/cardboard_bank")
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer db.Close()

	userRepo := repository.NewUsersRepository(db)
	serviceRepo := repository.NewServicesRepository(db)
	transactionRepo := repository.NewTransactionsRepository(db)

	basicAuth := middleware.BasicAuth(tokenKey)

	http.Handle("GET /user", basicAuth(getUser(userRepo)))
	http.Handle("PUT /user", basicAuth(updateUser(userRepo)))
	http.Handle("DELETE /user", basicAuth(deleteUser(userRepo)))

	http.Handle("GET /service", basicAuth(getAllServices(serviceRepo)))
	http.Handle("GET /service/{id}", basicAuth(getService(serviceRepo)))
	http.Handle("GET /service/{id}/transaction", basicAuth(getServiceTransactions(serviceRepo)))
	http.Handle("POST /service", basicAuth(createService(serviceRepo)))
	http.Handle("DELETE /service/{id}", basicAuth(cancelService(serviceRepo)))

	http.Handle("GET /transaction", basicAuth(getAllTransactions(transactionRepo)))
	http.Handle("GET /transaction/{id}", basicAuth(getTransaction(transactionRepo)))
	http.Handle("POST /transaction", basicAuth(executeTransaction(transactionRepo)))
	http.Handle("DELETE /transaction/{id}", basicAuth(rollbackTransaction(transactionRepo)))

	log.Println("starting API server")
	if err = http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

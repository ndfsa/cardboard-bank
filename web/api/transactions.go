package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/ndfsa/cardboard-bank/common/model"
	"github.com/ndfsa/cardboard-bank/common/repository"
	"github.com/ndfsa/cardboard-bank/web/dto"
	"github.com/ndfsa/cardboard-bank/web/middleware"
)

type TransactionsHandlerFactory struct {
	repo    repository.TransactionsRepository
	mdf     middleware.MiddlewareFactory
	srvRepo repository.ServicesRepository
}

func NewTransactionsHandlerFactory(
	repo repository.TransactionsRepository,
	mdf middleware.MiddlewareFactory,
	srvRepo repository.ServicesRepository,
) TransactionsHandlerFactory {
	return TransactionsHandlerFactory{repo, mdf, srvRepo}
}

func (factory *TransactionsHandlerFactory) CreateTransaction() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth)
	f := func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateTransactionRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		transaction, err := req.Parse()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}

		if err := factory.repo.CreateTransaction(r.Context(), transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := json.NewEncoder(w).Encode(dto.CreateTransactionResponseDTO{
			Id: transaction.Id.String(),
		}); err != nil {
			w.WriteHeader(http.StatusCreated)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *TransactionsHandlerFactory) ReadSingleTransaction() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.ClearanceOrOwnership(model.UserClearanceTeller, middleware.OwnershipTrs))
	f := func(w http.ResponseWriter, r *http.Request) {
		transactionId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
			return
		}

		transaction, err := factory.repo.FindTransaction(r.Context(), transactionId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if err := json.NewEncoder(w).Encode(dto.ReadTransactionResponseDTO{
			Id:          transaction.Id.String(),
			State:       transaction.State,
			Time:        transaction.Time,
			Currency:    transaction.Currency,
			Amount:      transaction.Amount.String(),
			Source:      transaction.Source.String(),
			Destination: transaction.Destination.String(),
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *TransactionsHandlerFactory) ReadMultipleTransactions() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.Clearance(model.UserClearanceTeller))
	f := func(w http.ResponseWriter, r *http.Request) {
		cursorString := r.URL.Query().Get("cursor")
		var cursor uuid.UUID
		if cursorString != "" {
			var err error
			cursor, err = uuid.Parse(cursorString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				return
			}
		} else {
			cursor = uuid.UUID{}
		}

		transactionsIt, err := factory.repo.FindAllTransactions(r.Context(), cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "")
		for transaction, err := range transactionsIt {
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}

			if err := encoder.Encode(dto.ReadTransactionResponseDTO{
				Id:          transaction.Id.String(),
				State:       transaction.State,
				Time:        transaction.Time,
				Currency:    transaction.Currency,
				Amount:      transaction.Amount.String(),
				Source:      transaction.Source.String(),
				Destination: transaction.Destination.String(),
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}
	}
	return mid(http.HandlerFunc(f))
}

func (factory *TransactionsHandlerFactory) ReadServiceTransactions() http.Handler {
	mid := middleware.Chain(
		factory.mdf.Logger,
		factory.mdf.UploadLimit(1000),
		factory.mdf.Auth,
		factory.mdf.ClearanceOrOwnership(model.UserClearanceTeller, middleware.OwnershipSrv))
	f := func(w http.ResponseWriter, r *http.Request) {
		serviceId, _ := uuid.Parse(r.PathValue("id"))
		cursorString := r.URL.Query().Get("cursor")
		var cursor uuid.UUID
		if cursorString != "" {
			var err error
			cursor, err = uuid.Parse(cursorString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Println(err)
				return
			}
		} else {
			cursor = uuid.UUID{}
		}

		transactionsIt, err := factory.repo.FindServiceTransactions(r.Context(), serviceId, cursor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

        encoder := json.NewEncoder(w)
        encoder.SetIndent("", "")
		for transaction, err := range transactionsIt {
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}

			if err := encoder.Encode(dto.ReadTransactionResponseDTO{
				Id:          transaction.Id.String(),
				State:       transaction.State,
				Time:        transaction.Time,
				Currency:    transaction.Currency,
				Amount:      transaction.Amount.String(),
				Source:      transaction.Source.String(),
				Destination: transaction.Destination.String(),
			}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}
	}
	return mid(http.HandlerFunc(f))
}

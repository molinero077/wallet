package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"wallet/internal/model"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type App struct {
	storage model.WalletProvider
	ctx     context.Context
}

func New() App {
	return App{
		ctx: context.Background(),
	}
}

func (app *App) GetContext() *context.Context {
	return &app.ctx
}

func (app *App) AssignStorage(storage *model.WalletProvider) {
	app.storage = *storage
}

func (app *App) Run(addr string) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/wallet", app.CurryOperation)
	mux.HandleFunc("GET /api/v1/wallets/{wallet_uuid}", app.getWallet)

	log.Info("Server is started. Listen on ", addr)
	log.Fatal(http.ListenAndServe(addr, middleWareLogger(mux)))
}

func middleWareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Info(fmt.Sprintf("%s %s %s", req.RemoteAddr, req.Proto, req.URL))
		next.ServeHTTP(w, req)
	})
}

func (app *App) getWallet(w http.ResponseWriter, req *http.Request) {

	wallet_uuid := req.PathValue("wallet_uuid")

	err := uuid.Validate(wallet_uuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	balance, err := app.storage.GetBalance(wallet_uuid)
	if err != nil {
		// не выводить пользователю ошибко о несуществующем кошельке, только статус
		if errors.Is(err, model.ErrNonExistentWallet) {
			log.Error("attempt to get balance of non-existent wallet")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug("trying to marshal struct to json")
	payload, err := json.Marshal(balance)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("the structure is successfully marshaled to json")

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(payload))
}

func (app *App) CurryOperation(w http.ResponseWriter, req *http.Request) {
	var operation model.Operation

	log.Debug("trying to marshal struct to json")

	operationDecoder := json.NewDecoder(req.Body)
	err := operationDecoder.Decode(&operation)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug("the structure is successfully marshaled to json")
	err = app.storage.CarryOperation(operation)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet/internal/model"
)

type FakeDatabase struct {
	data map[string]model.WalletBalance
}

func (db *FakeDatabase) Data() model.WalletProvider {
	return &FakeDatabase{
		data: map[string]model.WalletBalance{
			"99785d88-d0d0-44c1-9d85-3f187147d11a": {WalletId: "99785d88-d0d0-44c1-9d85-3f187147d11a", Balance: 1000},
			"dde8af2f-4a23-4025-acef-48f21977b2ec": {WalletId: "dde8af2f-4a23-4025-acef-48f21977b2ec", Balance: 1000},
			"ad8c367a-dd01-4642-b83b-640b5e44dc7f": {WalletId: "ad8c367a-dd01-4642-b83b-640b5e44dc7f", Balance: 1000},
			"af4e55f2-756b-4e72-be18-06fb45378e5a": {WalletId: "af4e55f2-756b-4e72-be18-06fb45378e5a", Balance: 1000},
		}}
}

func (fdb *FakeDatabase) GetWalletBalance(id string) (*model.WalletBalance, error) {
	if walletBalance, ok := fdb.data[id]; ok {
		return &walletBalance, nil
	}

	return &model.WalletBalance{}, model.ErrNonExistentWallet
}

func (fdb *FakeDatabase) WalletOperation(ops model.Operation) error {
	if _, ok := fdb.data[ops.WalletId]; !ok {
		return model.ErrNonExistentWallet
	}

	return nil
}

func TestGetBalance(t *testing.T) {
	var db FakeDatabase

	app := App{
		storage: db.Data(),
	}

	type test struct {
		uuid       string
		needStatus int
	}

	tests := []test{
		{
			uuid:       "",
			needStatus: http.StatusBadRequest,
		},
		{
			uuid:       "some string",
			needStatus: http.StatusBadRequest,
		},
		{
			uuid:       "99785d88-d0d0-44c1-9d85-3f187147d11f",
			needStatus: http.StatusNotFound,
		},
		{
			uuid:       "99785d88-d0d0-44c1-9d85-3f187147d11a",
			needStatus: http.StatusOK,
		},
	}

	for _, testItem := range tests {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/", nil)
		req.SetPathValue("wallet_uuid", testItem.uuid)
		w := httptest.NewRecorder()

		app.getWalletBalance(w, req)

		if w.Result().StatusCode != testItem.needStatus {
			t.Errorf("expected status %d, got %d", testItem.needStatus, w.Result().StatusCode)
		}
	}
}

func TestWalletOperation(t *testing.T) {
	var db FakeDatabase

	app := App{
		storage: db.Data(),
	}

	type test struct {
		payload    string
		needStatus int
	}

	tests := []test{
		{payload: "", needStatus: http.StatusBadRequest},
		// сумма строкой
		{payload: "{\"valletId\":\"99785d88-d0d0-44c1-9d85-3f187147d11a\", \"operationType\":\"DEPOSIT\", \"amount\":\"100.5\" }", needStatus: http.StatusBadRequest},
		// ошибки в именах параметров
		{payload: "{\"walletId\":\"99785d88-d0d0-44c1-9d85-3f187147d11a\", \"operationType\":\"DEPOSIT\", \"ammount\":100.5 }", needStatus: http.StatusBadRequest},
		{payload: "{\"valletId\":\"99785d88-d0d0-44c1-9d85-3f187147d11f\", \"operationType\":\"DEPOSIT\", \"amount\":100.5 }", needStatus: http.StatusNotFound},
		{payload: "{\"valletId\":\"99785d88-d0d0-44c1-9d85-3f187147d11a\", \"operationType\":\"DEPOSIT\", \"amount\":100.5 }", needStatus: http.StatusOK},
	}

	for _, testItem := range tests {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet/", strings.NewReader(testItem.payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		app.walletOperation(w, req)
		if w.Result().StatusCode != testItem.needStatus {
			t.Errorf("expected status %d, got %d", testItem.needStatus, w.Result().StatusCode)
		}
	}
}

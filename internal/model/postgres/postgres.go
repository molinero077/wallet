package postgres

import (
	"context"
	"fmt"
	"wallet/internal/model"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool struct {
	ctx *context.Context
	*pgxpool.Pool
}

type ConnectionParameters struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func New(ctx *context.Context, cp *ConnectionParameters) (*PgxPool, error) {
	var err error

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cp.User, cp.Password, cp.Host, cp.Port, cp.Database)
	log.Info("connect to ", connectionString)

	pg := new(PgxPool)
	pg.Pool, err = pgxpool.New(*ctx, connectionString)
	if err != nil {
		log.Error("не удалось подключить ся к БД")
		return nil, err
	}

	pg.ctx = ctx

	return pg, pg.Ping(*ctx)
}

func (pool *PgxPool) GetWalletBalance(walletId string) (*model.WalletBalance, error) {
	log.Debug("get balance of the wallet ", walletId)
	log.Debug("SELECT SUM(amount) as balance FROM public.operations WHERE wallet_id=", walletId)

	rows, err := pool.Query(*pool.ctx, "SELECT SUM(amount) as balance FROM public.operations WHERE wallet_id=$1", walletId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var balance interface{}

	for rows.Next() {
		err := rows.Scan(&balance)
		if err != nil {
			return nil, err
		}

	}

	if balance != nil {
		return &model.WalletBalance{
			WalletId: walletId,
			Balance:  balance.(float64),
		}, nil
	}

	return nil, model.ErrNonExistentWallet
}

func (pool *PgxPool) WalletOperation(op model.Operation) error {
	var amount float64

	if amount = op.GetAmount(); amount == 0 {
		return model.ErrZeroAmount
	}

	rows, err := pool.Query(*pool.ctx, "SELECT wallet_id FROM public.operations WHERE wallet_id=$1", op.WalletId)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return model.ErrNonExistentWallet
	}

	log.Debug("operations with the wallet ", op.WalletId)
	log.Debug(fmt.Sprintf("INSERT INTO public.operations(wallet_id, amount) VALUES(%s, %f)", op.WalletId, amount))

	isnertResult, err := pool.Exec(*pool.ctx, "INSERT INTO public.operations(wallet_id, amount) VALUES($1, $2)", op.WalletId, fmt.Sprintf("%f", amount))
	if err != nil {
		return err
	}

	log.Debug(isnertResult)

	return nil
}

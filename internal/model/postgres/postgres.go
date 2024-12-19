package postgres

import (
	"context"
	"fmt"
	"wallet/internal/model"

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
	// log.Info(fmt.Sprintf("подключение к БД postgres://%s:%s/%s", cp.Host, cp.Port, cp.Database))

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cp.User, cp.Password, cp.Host, cp.Port, cp.Database)

	pg := new(PgxPool)
	pg.Pool, err = pgxpool.New(*ctx, connectionString)
	if err != nil {
		// log.Error("не удалось подключить ся к БД")
		// log.Debug(err)

		return nil, fmt.Errorf("не удалось подключить ся к БД")
	}

	pg.ctx = ctx

	return pg, nil
}

func (pool *PgxPool) GetBalance(walletId string) (float32, error) {
	rows, err := pool.Query(*pool.ctx, "SELECT wallet_id, amount FROM public.wallets WHERE wallet_id=$1", walletId)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var balance float32

	for rows.Next() {
		err := rows.Scan(balance)
		if err != nil {
			return 0, err
		}

	}

	return balance, nil
}

func (pool *PgxPool) CarryOperation(op model.Operation) error {
	return nil
}

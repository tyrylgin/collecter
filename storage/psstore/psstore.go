package psstore

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tyrylgin/collecter/model"
	"github.com/tyrylgin/collecter/storage"
)

var _ storage.MetricStorer = (*PsStore)(nil)

type PsStore struct {
	db *pgxpool.Pool
}

func Init(ctx context.Context, dsn string) (*PsStore, error) {
	dbPoll, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &PsStore{db: dbPoll}, nil
}

func (s PsStore) GetAll() model.MetricMap {
	return nil
}

func (s PsStore) Get(name string) model.Metric {
	return nil
}

func (s PsStore) Save(name string, metric model.Metric) error {
	return nil
}

func (s PsStore) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

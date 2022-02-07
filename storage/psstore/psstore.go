package psstore

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/tyrylgin/collecter/model"
	"github.com/tyrylgin/collecter/storage"
)

var _ storage.MetricStorer = (*PsStore)(nil)

type PsStore struct {
	db *sqlx.DB
}

func Init(dsn string) (*PsStore, error) {
	store := &PsStore{}

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	store.db = db
	store.SetupSchema()

	return store, nil
}

func (s PsStore) SetupSchema() {
	s.db.MustExec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT FROM pg_type WHERE typname = 'metric_type') THEN
				CREATE TYPE metric_type AS ENUM ('counter', 'gauge');
			END IF;

			CREATE TABLE IF NOT EXISTS metrics
			(
				id text CONSTRAINT firstkey PRIMARY KEY,
				type metric_type NOT NULL,
				delta bigint,
				value double precision
			);
		END$$;
 	`)
}

type metricDB struct {
	ID         string
	MetricType string `db:"type"`
	Delta      sql.NullInt64
	Value      sql.NullFloat64
}

func convertToDB(name string, metric model.Metric) metricDB {
	mDB := metricDB{
		ID:         name,
		MetricType: string(metric.Type()),
	}

	switch metric.Type() {
	case model.MetricTypeCounter:
		mDB.Delta = sql.NullInt64{
			Int64: metric.(model.Counter).Delta,
			Valid: true,
		}
	case model.MetricTypeGauge:
		mDB.Value = sql.NullFloat64{
			Float64: metric.(model.Gauge).Value,
			Valid:   true,
		}
	}

	return mDB
}

func convertToModel(mDB metricDB) model.Metric {
	var metric model.Metric

	switch model.MetricType(mDB.MetricType) {
	case model.MetricTypeCounter:
		metric = model.Counter{Delta: mDB.Delta.Int64}
	case model.MetricTypeGauge:
		metric = model.Gauge{Value: mDB.Value.Float64}
	}

	return metric
}

func (s PsStore) GetAll() model.MetricMap {
	metrics := map[string]model.Metric{}
	var mDBs []metricDB

	err := s.db.Select(&mDBs, `SELECT * FROM metrics ORDER BY id ASC`)
	if err != nil {
		log.Printf("can't get metrics: %v", err)
		return nil
	}

	for _, item := range mDBs {
		metrics[item.ID] = convertToModel(item)
	}

	return metrics
}

func (s PsStore) Get(name string) model.Metric {
	var mDB metricDB

	err := s.db.Get(&mDB, `SELECT * FROM metrics WHERE id=$1`, name)
	if err != nil {
		log.Printf("can't get metric by name: %v", err)
		return nil
	}

	return convertToModel(mDB)
}

func (s PsStore) Save(name string, metric model.Metric) error {
	mDB := convertToDB(name, metric)

	rawQuery := `
		INSERT INTO metrics(id, type, delta, value) VALUES(:id, :type, :delta, :value)
		ON CONFLICT (id) DO UPDATE SET delta=:delta, value=:value
	`

	_, err := s.db.NamedExec(rawQuery, mDB)
	if err != nil {
		return fmt.Errorf("can't upsert metric: %v", err)
	}

	return nil
}

func (s PsStore) SaveAll(metrics model.MetricMap) error {
	tx := s.db.MustBegin()

	rawQuery := `
		INSERT INTO metrics (id, type, delta, value) VALUES (:id, :type, :delta, :value)
		ON CONFLICT (id) DO UPDATE SET delta = metrics.delta + :delta, value=:value
	`

	stmt, err := tx.PrepareNamed(rawQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}

	for name, metric := range metrics {
		_, err := stmt.Exec(convertToDB(name, metric))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute statement: %v", err)
		}
	}

	tx.Commit()

	return nil
}

func (s PsStore) Ping() error {
	return s.db.Ping()
}

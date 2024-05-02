package db

import (
	"context"
	"log"

	"github.com/Eevangelion/ewallet/config"
	"github.com/Eevangelion/ewallet/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var pool *pgxpool.Pool = nil

func GetPool() (db *pgxpool.Pool, err error) {
	if pool == nil {
		conf := config.GetConfig()

		dbConfig, err := pgxpool.ParseConfig("")

		if err != nil {
			log.Fatal("Error parse db config:", err)
			return nil, err
		}

		dbConfig.ConnConfig.Host = conf.DB.Host
		dbConfig.ConnConfig.Port = uint16(conf.DB.Port)
		dbConfig.ConnConfig.Database = conf.DB.Name
		dbConfig.ConnConfig.User = conf.DB.User
		dbConfig.ConnConfig.Password = conf.DB.Password
		dbConfig.MinConns = int32(conf.DB.MinOpenConns)
		dbConfig.MaxConns = int32(conf.DB.MaxOpenConns)
		dbConfig.ConnConfig.ConnectTimeout = conf.DB.ConnTimeout

		pool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
		if err != nil {
			log.Fatal("Connection Error:", err)
			return nil, err
		}
		err = CreateTables(pool)
		if err != nil {
			return nil, err
		}
	}
	return pool, err
}

const createTables = `
BEGIN;
	
CREATE TABLE IF NOT EXISTS public."wallet"
(
	id varchar(50) NOT NULL,
	balance float,
	CONSTRAINT "wallet_pkey" PRIMARY KEY (id),
	CONSTRAINT "balance_positive" CHECK (balance >= 0) NOT VALID
);

CREATE TABLE IF NOT EXISTS public."transaction" 
(
	id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
	sender_id varchar(50) NOT NULL,
	receiver_id varchar(50) NOT NULL,
	amount float NOT NULL,
	time_stamp timestamp NOT NULL,
	CONSTRAINT "transaction_pkey" PRIMARY KEY (id),
	FOREIGN KEY (sender_id) REFERENCES public."wallet" (id) MATCH SIMPLE ON DELETE CASCADE,
	FOREIGN KEY (receiver_id) REFERENCES public."wallet" (id) MATCH SIMPLE ON DELETE CASCADE
);

COMMIT;`

func CreateTables(db *pgxpool.Pool) (err error) {

	logger := logger.GetLogger()
	_, err = db.Exec(context.TODO(), createTables)

	if err != nil {
		logger.Error(
			"Error while creating DB tables:",
			zap.String("event", "create_database_tables"),
			zap.String("error", err.Error()),
		)
		return
	}

	return
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq" //nolint
)

type Postgres struct {
	db *sql.DB
}

type Config struct {
	Host         string
	Port         int
	User         string
	Pass         string
	DBName       string
	MaxOpenConns int
}

func New(ctx context.Context, c Config) (*Postgres, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Pass, c.Host, c.Port, c.DBName)

	db, err := sql.Open("postgres", connStr)

	//db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		err := db.Close()
		if err != nil {
			log.Error("close sqldb connection error:", err.Error())
			return
		}
		log.Info("close sqldb connection")
	}()

	db.SetMaxOpenConns(c.MaxOpenConns)

	return &Postgres{
		db: db,
	}, nil
}

func (s *Postgres) Ping() error {
	return s.db.Ping()
}

package db

import (
	"github.com/go-pg/pg/extra/pgotel/v10"
	"github.com/go-pg/pg/v10"
	"source.golabs.io/engineering-platforms/lens/trace-app-golang/config"
)

type Proxy struct {
	db *pg.DB
}

func New(cfg config.Config) *Proxy {
	pgDB := pg.Connect(&pg.Options{
		Addr:     cfg.DBAddr,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBName,
	})

	Run(cfg)

	pgDB.AddQueryHook(pgotel.NewTracingHook())

	return &Proxy{db: pgDB}
}

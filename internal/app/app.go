package app

import (
	"github.com/jmoiron/sqlx"
)

type App struct {
	db *sqlx.DB
}

func NewApp(db *sqlx.DB) *App {
	return &App{
		db: db,
	}
}

func (a *App) Start() error {
	return nil
}

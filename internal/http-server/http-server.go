package httpserver

import (
	"database/sql"
	"gopher/internal/config"
	"gopher/internal/storage/sqlstorage"
	"net/http"

	"github.com/gorilla/sessions"
)

func Start(config *config.Config) error {
	db, err := newDatabase(config.Storage)

	if err != nil {
		return err
	}

	defer db.Close()

	storage := sqlstorage.New(db)

	session := sessions.NewCookieStore([]byte(config.SessionKey))

	server := newServer(storage, session)

	return http.ListenAndServe(config.Address, server.router)
}

func newDatabase(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databaseUrl)

	if err != nil {
		return nil, err
	}

	return db, nil
}

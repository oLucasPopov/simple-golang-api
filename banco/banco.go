package banco

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Conectar() (*sql.DB, error) {
	stringConexao := "user=golang password=golang dbname=devbook sslmode=disable"

	db, erro := sql.Open("postgres", stringConexao)

	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}

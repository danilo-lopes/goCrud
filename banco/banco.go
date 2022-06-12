package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Required
)

// Conectar abre a conexao com o banco de dados
func Conectar() (*sql.DB, error) {
	stringConexao := "admin:password@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, erro := sql.Open("mysql", stringConexao)

	if erro != nil {
		return nil, erro
	}

	if erro := db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}

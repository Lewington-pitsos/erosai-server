package database

import (
	"bitbucket.org/lewington/autoroller/assist"
	"bitbucket.org/lewington/erosai/shared"
	"github.com/jmoiron/sqlx"
)

type Archivist struct {
	connection
}

func (a *Archivist) execute(stmt *sqlx.Stmt, args ...interface{}) {
	stmt.MustExec(args...)
	stmt.Close()
}

func (a *Archivist) getId(row *sqlx.Row) int {
	var id int
	err := row.Scan(&id)
	assist.Check(err)

	return id
}

func (a *Archivist) DoesUserExist(reg shared.Registration) bool {
	var ID int
	stmt := a.MakeStmt(`
		SELECT id from users WHERE name = $1;
	`)
	defer stmt.Close()

	stmt.Get(&ID, reg.Username)

	return ID != 0
}

func (a *Archivist) RegisterUser(reg shared.Registration) int {
	stmt := a.SprintfStmt(`
		INSERT INTO users (name, password)
		VALUES ($1, $2)
			ON CONFLICT ON CONSTRAINT unique_name
			DO UPDATE SET password = $2
			RETURNING id;
	`)

	return a.ExecuteAndReturnId(stmt, reg.Username, reg.Password)
}

func (a *Archivist) ExecuteAndReturnId(stmt *sqlx.Stmt, args ...interface{}) int {
	row := stmt.QueryRowx(args...)
	stmt.Close()

	return a.getId(row)
}

func (a *Archivist) MustExec(stmt string) {
	a.db.MustExec(stmt)
}

func NewArchivist() Archivist {
	return Archivist{
		newConnection(),
	}
}

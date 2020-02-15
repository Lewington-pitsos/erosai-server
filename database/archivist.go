package database

import (
	"bitbucket.org/lewington/autoroller/assist"
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

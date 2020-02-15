package database

import (
	"fmt"
	"sync"

	"bitbucket.org/lewington/autoroller/assist"
	"bitbucket.org/lewington/autoroller/lg"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var driverName string = "postgres"
var allConnections = []*connection{}
var mutex = &sync.Mutex{}

type connection struct {
	db *sqlx.DB
}

func (c *connection) MakeStmt(query string) *sqlx.Stmt {
	stmt, err := c.db.Preparex(query)
	assist.Check(err)
	return stmt
}

func (c *connection) Get(dest interface{}, query string, args ...interface{}) {
	err := c.db.Get(dest, query, args...)
	assist.Check(err)
}

func (c *connection) SprintfStmt(query string, args ...interface{}) *sqlx.Stmt {
	queryString := fmt.Sprintf(query, args...)
	return c.MakeStmt(queryString)
}

func (c *connection) GetRows(query string, args ...interface{}) *sqlx.Rows {
	stmt := c.MakeStmt(query)
	defer stmt.Close()

	rows, err := stmt.Queryx(args...)
	assist.Check(err)
	return rows
}

func (c *connection) LogStats() {
	lg.L.Info("Number of open connections: %v", c.db.Stats().OpenConnections)
}

func (c *connection) Close() {
	c.db.Close()
}

// newConnection opens a connection to the adlerbot database (or throws an error) and returns a new connection instance wrapped around that connection.
func newConnection() connection {
	database := sqlx.MustOpen(driverName, "user=erosai dbname=links password=Erosai11!! sslmode=disable")
	c := connection{
		db: database,
	}

	registerConnection(&c)

	return c
}

func registerConnection(c *connection) {
	mutex.Lock()
	allConnections = append(allConnections, c)
	mutex.Unlock()
}

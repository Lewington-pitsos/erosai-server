package database

import (
	"bitbucket.org/lewington/autoroller/assist"
	"bitbucket.org/lewington/erosai/globals"
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

func (a *Archivist) DoesUserExist(reg shared.Details) bool {
	var ID int
	stmt := a.MakeStmt(`
		SELECT id from users WHERE name = $1;
	`)
	defer stmt.Close()

	stmt.Get(&ID, reg.Username)

	return ID != 0
}

func (a *Archivist) GetIdForToken(token string) int {
	var ID int
	stmt := a.MakeStmt(`
		SELECT id from users WHERE token = $1;
	`)
	defer stmt.Close()

	stmt.Get(&ID, token)

	return ID
}

func (a *Archivist) SetUserToken(reg shared.Details, token string) {
	stmt := a.MakeStmt(`
		UPDATE users set token = $1 WHERE name = $2;
	`)

	a.execute(stmt, token, reg.Username)
}

func (a *Archivist) RegisterUser(reg shared.Details) int {
	stmt := a.SprintfStmt(`
		INSERT INTO users (name, password)
		VALUES ($1, $2)
			ON CONFLICT ON CONSTRAINT unique_name
			DO UPDATE SET password = $2
			RETURNING id;
	`)

	return a.ExecuteAndReturnId(stmt, reg.Username, reg.Password)
}

func (a *Archivist) URLIsNew(URL string) bool {
	var ID int
	stmt := a.MakeStmt(`
		SELECT id from links WHERE url = $1;
	`)
	defer stmt.Close()

	stmt.Get(&ID, URL)

	return ID == 0
}

func (a *Archivist) AddURL(URL string) int {
	stmt := a.SprintfStmt(`
		INSERT INTO links (url)
		VALUES ($1)
			ON CONFLICT ON CONSTRAINT unique_url
			DO UPDATE SET url = $1
			RETURNING id;
	`)

	return a.ExecuteAndReturnId(stmt, URL)
}

func (a *Archivist) GetReccomendations(userID int) []shared.Link {
	var links []shared.Link

	stmt := a.MakeStmt(`
		SELECT l.id, l.url, l.scanned, l.porn from links as l JOIN visits as v on l.id = v.link_id 
		WHERE v.user_id != $1 AND l.porn > $2 ORDER BY l.id LIMIT 8;
	`)
	defer stmt.Close()

	err := stmt.Select(&links, userID, globals.PornCutoff)

	assist.Check(err)

	return links
}

func (a *Archivist) UpdateLink(link shared.Link) {
	stmt := a.MakeStmt(`
		UPDATE links 
		SET scanned = $1, 
			porn = $2
			WHERE id = $3;
	`)

	a.execute(stmt, link.Scanned, link.Porn, link.ID)
}

func (a *Archivist) AddVisit(userID int, URLID int) int {
	stmt := a.SprintfStmt(`
		INSERT INTO visits (user_id, link_id)
		VALUES ($1, $2)
			ON CONFLICT ON CONSTRAINT unique_user_link
			DO UPDATE SET user_id = $1
			RETURNING id;
	`)

	return a.ExecuteAndReturnId(stmt, userID, URLID)
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

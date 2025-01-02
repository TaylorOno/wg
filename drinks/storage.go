package drinks

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"golang.org/x/net/context"
)

type Turso struct {
	db *sql.DB
}

func NewTursoDatabase() *Turso {
	database := os.Getenv("TURSO_DATABASE")
	token := os.Getenv("TURSO_AUTH_TOKEN")
	if len(database) < 1 || len(token) < 1 {
		log.Fatalf("TURSO_DATABASE and TURSO_AUTH_TOKEN are required")
	}

	dbName := fmt.Sprintf("%s?authToken=%s", strings.TrimSpace(database), strings.TrimSpace(token))
	db, err := sql.Open("libsql", dbName)
	if err != nil {
		log.Fatalf("failed to open db %s: %s", dbName, err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS history (id TEXT PRIMARY KEY, created TEXT , history BLOB)")
	if err != nil {
		log.Fatal(err)
	}

	return &Turso{
		db: db,
	}
}

func (t *Turso) Save(ctx context.Context, id string, history []byte) error {
	stmt, err := t.db.Prepare("INSERT INTO history (id, created, history) VALUES (?, ? ,?);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id, time.Now(), history)
	if err != nil {
		return err
	}

	return nil
}

func (t *Turso) GetById(ctx context.Context, id string) ([]byte, error) {
	stmt, err := t.db.Prepare("SELECT history from history where id=?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []byte
	for rows.Next() {
		err = rows.Scan(&history)
		if err != nil {
			return nil, err
		}

		return history, nil
	}

	return nil, errors.New("not found")
}

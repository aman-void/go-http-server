package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/aman-void/go-http-server/internal/config"
	"github.com/aman-void/go-http-server/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE if not exists users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    age INTEGER NOT NULL
);`,
	)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil

}

func (s *Sqlite) CreateUser(name string, email string, age int) (int64, error) {

	statement, err := s.Db.Prepare("INSERT into users (name, email, age) VALUES (?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer statement.Close()

	result, err := statement.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return lastId, nil
}

func (s *Sqlite) GetUserById(id int64) (types.User, error) {

	statement, err := s.Db.Prepare(`SELECT id, name, email, age FROM users WHERE id = ? LIMIT 1`)
	if err != nil {
		return types.User{}, nil
	}

	defer statement.Close()

	var user types.User

	err = statement.QueryRow(id).Scan(&user.Id, &user.Name, &user.Email, &user.Age)
	if err != nil {

		if err == sql.ErrNoRows {
			return types.User{}, fmt.Errorf("no user found with id %d", id)
		}
		return types.User{}, fmt.Errorf("query error %w", err)
	}

	return user, nil
}

func (s *Sqlite) GetUsers() ([]types.User, error) {

	statement, err := s.Db.Prepare("SELECT id, name, email, age FROM users")
	if err != nil {
		return nil, err
	}

	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.User

	for rows.Next() {
		var user types.User

		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Age)
		if err != nil {
			return nil, err
		}

		users = append(users, user)

	}
	return users, nil
}

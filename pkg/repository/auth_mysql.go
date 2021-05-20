package repository

import (
	"fmt"
	"zhashkRestApi"

	"github.com/jmoiron/sqlx"
)

type AuthMysql struct {
	db *sqlx.DB
}

func newAuthMysql(db *sqlx.DB) *AuthMysql {
	return &AuthMysql{db: db}
}

func (r *AuthMysql) CreateUser(user zhashkRestApi.User) (int, error) {
	var idInt64 int64 = 2 << 32

	query := fmt.Sprintf(`
		insert into %s (name, username, password_hash) values (?, ?, ?)
	`, usersTable)

	result, err := r.db.Exec(query, user.Name, user.Username, user.Password)
	if err != nil {
		return 0, err
	}

	idInt64, err = result.LastInsertId()
	id := int(idInt64)

	return id, nil
}

func (r *AuthMysql) GetUser(username, password string) (zhashkRestApi.User, error) {
	var user zhashkRestApi.User
	query := fmt.Sprintf(`select id from %s where username=? and password_hash=?`, usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}

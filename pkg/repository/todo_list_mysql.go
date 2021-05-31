package repository

import (
	"fmt"
	"zhashkRestApi"

	"github.com/jmoiron/sqlx"
)

type TodoListMysql struct {
	db *sqlx.DB
}

func NewTodoListMysql(db *sqlx.DB) *TodoListMysql {
	return &TodoListMysql{db: db}
}

func (r *TodoListMysql) Create(userId int, list zhashkRestApi.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// inserting new list
	var idInt64 int64 = 2 << 32
	createListQuery := fmt.Sprintf(`
		insert into %s (title, description) values (?, ?)
	`, todoListsTable)
	row, err := tx.Exec(createListQuery, list.Title, list.Description)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	idInt64, err = row.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	id := int(idInt64)

	// inserting new row to users_lists table
	createUsersListsQuery := fmt.Sprintf(`
		insert into %s (user_id, list_id) values (?, ?)
	`, usersListsTable)
	_, err = tx.Exec(createUsersListsQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListMysql) GetAll(userId int) ([]zhashkRestApi.TodoList, error) {
	var lists []zhashkRestApi.TodoList

	query := fmt.Sprintf(`
	SELECT tl.id, tl.title, tl.description FROM %s AS tl INNER JOIN %s as ul ON tl.id = ul.list_id WHERE ul.user_id = ?
	`, todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListMysql) GetById(userId, listId int) (zhashkRestApi.TodoList, error) {
	var list zhashkRestApi.TodoList

	query := fmt.Sprintf(`
		SELECT tl.id, tl.title, tl.description FROM %s AS tl 
		INNER JOIN %s as ul ON tl.id = ul.list_id WHERE ul.user_id = ? and ul.list_id = ?
	`, todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}
func (r *TodoListMysql) Delete(userId, listId int) error {
	query := fmt.Sprintf(`
		DELETE tl, ul FROM %s tl, %s ul
		WHERE tl.id = ul.list_id 
		AND ul.user_id = ?
		AND ul.list_id = ? 
	`, todoListsTable, usersListsTable)

	_, err := r.db.Exec(query, userId, listId)

	return err
}

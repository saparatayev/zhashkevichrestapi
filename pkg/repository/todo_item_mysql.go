package repository

import (
	"fmt"
	"zhashkRestApi"

	"github.com/jmoiron/sqlx"
)

type TodoItemMysql struct {
	db *sqlx.DB
}

func NewTodoItemMysql(db *sqlx.DB) *TodoItemMysql {
	return &TodoItemMysql{db: db}
}

func (r *TodoItemMysql) Create(listId int, item zhashkRestApi.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// inserting new item
	var idInt64 int64 = 2 << 32
	createItemQuery := fmt.Sprintf(`
		insert into %s (title, description) values (?, ?)
	`, todoItemsTable)
	row, err := tx.Exec(createItemQuery, item.Title, item.Description)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	idInt64, err = row.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	itemId := int(idInt64)

	// inserting new row to lists_items table
	createListsItemsQuery := fmt.Sprintf(`
		insert into %s (list_id, item_id) values (?, ?)
	`, listsItemsTable)
	_, err = tx.Exec(createListsItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemMysql) GetAll(userId, listId int) ([]zhashkRestApi.TodoItem, error) {
	var items []zhashkRestApi.TodoItem

	query := fmt.Sprintf(`
		select ti.id, ti.title, ti.description, ti.done from %s as ti inner join %s as li on li.item_id = ti.id
		inner join %s as ul on ul.list_id = li.list_id where li.list_id = ? and ul.user_id = ?
	`, todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

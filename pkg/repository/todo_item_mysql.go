package repository

import (
	"fmt"
	"strings"
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

func (r *TodoItemMysql) GetById(userId, itemId int) (zhashkRestApi.TodoItem, error) {
	var item zhashkRestApi.TodoItem

	query := fmt.Sprintf(`
		select ti.id, ti.title, ti.description, ti.done from %s as ti inner join %s as li on li.item_id = ti.id
		inner join %s as ul on ul.list_id = li.list_id where ti.id = ? and ul.user_id = ?
	`, todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItemMysql) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`
		DELETE ti, li FROM %s ti 
		LEFT JOIN %s li ON ti.id = li.item_id 
		LEFT JOIN %s ul on li.list_id = ul.list_id 
		WHERE ul.user_id = ? 
		AND ti.id = ?
	`, todoItemsTable, listsItemsTable, usersListsTable)

	_, err := r.db.Exec(query, userId, itemId)

	return err
}

func (r *TodoItemMysql) Update(userId, itemId int, input zhashkRestApi.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	if input.Title != nil {
		setValues = append(setValues, "title = ?")
		args = append(args, *input.Title)
	}

	if input.Description != nil {
		setValues = append(setValues, "description = ?")
		args = append(args, *input.Description)
	}

	if input.Done != nil {
		setValues = append(setValues, "done = ?")
		args = append(args, *input.Done)
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`
		update %s ti, %s li, %s ul
		set %s
		where
		ti.id = li.item_id
		and li.list_id = ul.list_id
		and ul.user_id = ? 
		and ti.id = ?
	`, todoItemsTable, listsItemsTable, usersListsTable, setQuery)

	// list_id: 3 список programms
	// item_id: 8 title: learn golang decr: rest api golang
	// user_id: 2 Maksim

	args = append(args, userId, itemId)

	_, err := r.db.Exec(query, args...)

	return err
}

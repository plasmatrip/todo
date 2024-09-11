package repository

import (
	"database/sql"

	"todo/model"
)

func (d *Repository) GetTask(id int) (model.Task, error) {
	t := model.Task{}

	row := d.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return model.Task{}, err
	}

	if err := row.Err(); err != nil {
		return model.Task{}, err
	}

	return t, nil
}

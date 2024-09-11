package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"todo/configs"
	"todo/model"
)

func (d *Repository) GetTasks(search string) ([]model.Task, error) {
	var res []model.Task
	var rows *sql.Rows
	var err error
	var date time.Time

	if len(search) > 0 {
		date, err = time.Parse(configs.SearchLayout, search)
		if err != nil {
			rows, err = d.db.Query("SELECT * FROM scheduler WHERE UPPER(title) LIKE :search OR UPPER(comment) LIKE :search ORDER BY date LIMIT 25",
				sql.Named("search", fmt.Sprintf("%%%s%%", strings.ToUpper(search))))
		} else {
			rows, err = d.db.Query("SELECT * FROM scheduler WHERE date = :date ORDER BY date LIMIT 25",
				sql.Named("date", date.Format(configs.DateLayout)))
		}
	} else {
		rows, err = d.db.Query("SELECT * FROM scheduler ORDER BY date LIMIT 25")
	}

	if err != nil {
		return []model.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		t := model.Task{}
		err := rows.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return []model.Task{}, err
		}
		res = append(res, t)
	}

	if err := rows.Err(); err != nil {
		return []model.Task{}, err
	}

	if res == nil {
		res = []model.Task{}
	}

	return res, nil
}

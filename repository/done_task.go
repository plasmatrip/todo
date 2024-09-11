package repository

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"todo/model"
	"todo/service"
)

func (d *Repository) DoneTask(id int) error {
	t := model.Task{}

	row := d.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))

	err := row.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return err
	}

	if err := row.Err(); err != nil {
		return err
	}

	if len(t.Repeat) == 0 {
		id, err := strconv.Atoi(t.Id)
		if err != nil {
			return err
		}

		err = d.Delete(id)
		if err != nil {
			return err
		}

		log.Printf("задача выполнена и удалена id=%d", id)
	} else {
		t.Date, err = service.NextDate(time.Now(), t.Date, t.Repeat)
		if err != nil {
			return err
		}

		err = d.Update(t)
		if err != nil {
			return err
		}
		log.Printf("выполнена задача id=%d, время выполнения изменено согласно праввилу повторения", id)
	}

	return nil
}

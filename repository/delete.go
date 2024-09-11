package repository

import (
	"database/sql"
	"fmt"
)

func (d *Repository) Delete(id int) error {
	res, err := d.db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("запись с id=%d в таблице не найдена", id)
	}

	return nil
}

package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
	// и т.д.
}

func AddTask(task *Task) (int64, error) {
	var id int64

	//dbFile := "scheduler.db"

	db, err := Open()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func checkDateFormat(dateStr string) bool {
	_, err := time.Parse("02.01.2006", dateStr)
	return err == nil
}

func Tasks(limit int, param string) ([]*Task, error) {

	var err error
	db, err := Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var query string
	var rows *sql.Rows

	if param == "" {

		query = `SELECT  * FROM scheduler ORDER BY date LIMIT :limit;`
		rows, err = db.Query(query, sql.Named("limit", limit))

	} else if checkDateFormat(param) {

		t, err := time.Parse("02.01.2006", param)
		if err != nil {

			return nil, fmt.Errorf("Ошибка парсинга даты: " + err.Error())
		}
		date_f := t.Format("20060102")

		query = `SELECT * FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit;`
		rows, err = db.Query(query,
			sql.Named("limit", limit),
			sql.Named("date", date_f))

	} else {

		query = `SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit;`

		rows, err = db.Query(query,
			sql.Named("limit", limit),
			sql.Named("search", "%"+param+"%"))
	}

	if err != nil {
		return nil, fmt.Errorf("error query db: %w", err)
	}
	defer rows.Close()

	var ret_tasks []*Task

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		ret_tasks = append(ret_tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in next rows: %w", err)
	}

	if len(ret_tasks) == 0 {
		ret_tasks = []*Task{}
	}

	return ret_tasks, nil
}

func UpdateTask(task *Task) error {

	db, err := Open()
	if err != nil {
		return err
	}
	defer db.Close()

	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET date =:date, title =:title, comment =:comment, repeat =:repeat WHERE id = :id;`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func UpdateDate(next string, id string) error {

	//dbFile := "scheduler.db"

	db, err := Open()
	if err != nil {
		return err
	}
	defer db.Close()

	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET date =:date WHERE id = :id;`
	res, err := db.Exec(query,
		sql.Named("date", next),
		sql.Named("id", id))
	if err != nil {
		return err
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil

}

func GetTask(id string) (*Task, error) {

	//dbFile := "scheduler.db"
	var err error
	db, err := Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id;`
	row := db.QueryRow(query, sql.Named("id", id))

	task := Task{}
	err = row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func DeleteTask(id string) error {

	//dbFile := "scheduler.db"
	var err error
	db, err := Open()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM scheduler WHERE id = :id;`
	_, err = db.Exec(query, sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil

}

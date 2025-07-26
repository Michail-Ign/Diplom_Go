package api

import (
	"Diplom_Go/pkg/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type SucessRet struct {
	ID string `json:"id"`
}

type ErrorRet struct {
	Error string `json:"error"`
}

func addTaskHandler(w http.ResponseWriter, req *http.Request) {

	var err_ret ErrorRet

	//1 Нужно десериализовать полученный в запросе JSON в переменную var task db.Task.
	bodyBytes, err := io.ReadAll(req.Body)

	if err != nil {

		err_ret.Error = fmt.Sprintf("ошибка чтения тела (%v)", err)
		writeJson(w, err_ret)
		return
	}

	var task db.Task

	if err := json.Unmarshal(bodyBytes, &task); err != nil {

		err_ret.Error = fmt.Sprintf("ошибка десериализации (%v)", err)
		writeJson(w, err_ret)
		return
	}

	//2 Проверить, что поле task.Title не пустое.
	if task.Title == "" {

		err_ret.Error = "пустое значение поля title"
		writeJson(w, err_ret)
		return
	}

	//3 Проверить на корректность полученное значение task.Date.
	// Это лучше сделать в отдельной функции checkDate(task *db.Task) error
	if err := checkDate(&task); err != nil {

		err_ret.Error = fmt.Sprintf("неверная дата (%v)", err)
		writeJson(w, err_ret)

		return
	}

	//4 Пришла очередь вызвать функцию db.AddTask(task),
	// чтобы добавить задачу в базу данных.
	id, err := db.AddTask(&task)
	if err != nil {

		err_ret.Error = fmt.Sprintf("ошибка добавления задачи (%v)", err)
		writeJson(w, err_ret)
		return
	}

	//5 Осталось возвратить идентификатор добавленной задачи в виде JSON.
	var id_ret SucessRet
	id_ret.ID = strconv.FormatInt(id, 10)

	writeJson(w, id_ret)
}

func checkDate(task *db.Task) error {

	now := time.Now()

	//3_1 Если поле date не указано или содержит пустую строку, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = now.Format(Format_date)
	}

	t, err := time.Parse(Format_date, task.Date)
	if err != nil {
		return err
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil && len(task.Repeat) != 0 {
		return err
	}

	// если сегодня (now) больше task.Date (t)
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(Format_date)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}

	return nil
}

// Вспомогательная функция для сериализации и отправки JSON-ответа
/*func writeJson(w http.ResponseWriter, data any) error {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//log.Println("data = ", data) //
	err := json.NewEncoder(w).Encode(data)

	return err
}*/

func writeJson(w http.ResponseWriter, data any) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//log.Println("2data = ", data) //
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		//log.Println("2err = ", err.Error()) //
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//log.Println("2err = nil") //

}

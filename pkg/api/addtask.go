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

	fmt.Println("----вызов addTaskHandler")
	//
	var err_ret ErrorRet

	//1 Нужно десериализовать полученный в запросе JSON в переменную var task db.Task.
	bodyBytes, err := io.ReadAll(req.Body)
	fmt.Println("    bodyBytes = " + string(bodyBytes))
	if err != nil {
		//http.StatusInternalServerError
		err_ret.Error = fmt.Sprintf("ошибка чтения тела (%v)", err)
		err = writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	var task db.Task
	fmt.Println("    2 = ")
	if err := json.Unmarshal(bodyBytes, &task); err != nil {

		//http.StatusBadRequest
		err_ret.Error = fmt.Sprintf("ошибка десериализации (%v)", err)
		err = writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	//2 Проверить, что поле task.Title не пустое.
	fmt.Println("    2 = ")
	if task.Title == "" {

		fmt.Println("    2 err ")
		err_ret.Error = "пустое значение поля title"
		//w.WriteHeader(http.StatusInternalServerError)
		err = writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	//3 Проверить на корректность полученное значение task.Date.
	// Это лучше сделать в отдельной функции checkDate(task *db.Task) error
	fmt.Println("    3 = ")
	if err := checkDate(&task); err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		err_ret.Error = fmt.Sprintf("неверная дата (%v)", err)
		err = writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	//4 Пришла очередь вызвать функцию db.AddTask(task),
	// чтобы добавить задачу в базу данных.
	id, err := db.AddTask(&task)
	if err != nil {

		//http.Error(w, err.Error(), http.StatusBadRequest)
		err_ret.Error = fmt.Sprintf("ошибка добавления задачи (%v)", err)
		err = writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	//5 Осталось возвратить идентификатор добавленной задачи в виде JSON.
	var id_ret SucessRet
	id_ret.ID = strconv.FormatInt(id, 10)

	err = writeJson(w, id_ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func checkDate(task *db.Task) error {
	fmt.Println("------вызов checkDate")
	//
	now := time.Now()
	fmt.Println("      task.Date= " + task.Date)
	fmt.Println("      now= ", now.String())
	fmt.Println("      repeat= ", task.Repeat)
	//3_1 Если поле date не указано или содержит пустую строку, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = now.Format(Format_date)
	}

	fmt.Println("      2_1")
	t, err := time.Parse(Format_date, task.Date)
	if err != nil {
		return err
	}

	fmt.Println("      2_2")

	next, err := NextDate(now, task.Date, task.Repeat)
	fmt.Println("      next= ", next)
	fmt.Println("      err= ", err)
	if err != nil && len(task.Repeat) != 0 {
		return err
	}

	fmt.Println("      2_3")
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
	fmt.Println("      task.Date= ", task.Date)
	return nil
}
//
/*func CheckDate(task *db.Task) error {
	fmt.Println("------вызов checkDate")
	//
	now := time.Now()
	fmt.Println("      task.Date= " + task.Date)
	fmt.Println("      now= ", now.String())
	fmt.Println("      repeat= ", task.Repeat)
	//3_1 Если поле date не указано или содержит пустую строку, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = now.Format(Format_date)
	}

	fmt.Println("      2_1")
	t, err := time.Parse(Format_date, task.Date)
	if err != nil {
		return err
	}

	fmt.Println("      2_2")

	next, err := NextDate(now, task.Date, task.Repeat)
	fmt.Println("      next= ", next)
	fmt.Println("      err= ", err)
	if err != nil && len(task.Repeat) != 0 {
		return err
	}

	fmt.Println("      2_3")
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
	fmt.Println("      task.Date= ", task.Date)
	return nil
}*/
//

// Вспомогательная функция для сериализации и отправки JSON-ответа
func writeJson(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	return json.NewEncoder(w).Encode(data)
}

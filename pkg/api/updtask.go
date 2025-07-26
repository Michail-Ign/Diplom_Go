package api

import (
	"Diplom_Go/pkg/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func getTaskHandler(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")

	task, err := db.GetTask(id)

	var err_ret ErrorRet
	if err != nil {

		err_ret.Error = fmt.Sprintf("ошибка получения задачи (%v)", err)
		writeJson(w, err_ret)
		return
	}

	writeJson(w, task)
}
func updTaskHandler(w http.ResponseWriter, req *http.Request) {

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
	err = db.UpdateTask(&task)
	if err != nil {

		err_ret.Error = fmt.Sprintf("ошибка обновления задачи (%v)", err)
		writeJson(w, err_ret)
		return
	}

	//5 Осталось возвратить идентификатор добавленной задачи в виде JSON.
	var id_ret SucessRet
	id_ret.ID = task.ID

	writeJson(w, id_ret)
}

func delTaskHandler(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")
	var err_ret ErrorRet

	if id == "" {
		err_ret.Error = "не задан параметр id"
		writeJson(w, err_ret)
		return
	}
	_, err := strconv.Atoi(id)
	if err != nil {

		err_ret.Error = fmt.Sprintf("параметр  id не валиден  (%v)", err)
		writeJson(w, err_ret)
		return
	}

	err = db.DeleteTask(id)

	if err != nil {

		err_ret.Error = fmt.Sprintf("ошибка при удалении задачи id  (%v)", err)
		writeJson(w, err_ret)
		return
	}

	emptyJSON := map[string]interface{}{}
	writeJson(w, emptyJSON)
}

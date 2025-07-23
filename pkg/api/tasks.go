package api

import (
	"Diplom_Go/pkg/db"
	"fmt"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {

	var err_ret ErrorRet

	param := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, param) // в параметре максимальное количество записей

	if err != nil {
		// здесь вызываете функцию, которая возвращает ошибку в JSON
		// её желательно было реализовать на предыдущем шаге
		err_ret.Error = fmt.Sprintf("ошибка получения задачи (%v)", err)

		err = writeJson(w, err_ret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

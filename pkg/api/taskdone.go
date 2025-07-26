package api

import (
	"Diplom_Go/pkg/db"
	"fmt"
	"net/http"
	"time"
)

func taskDoneHandler(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")

	task, err := db.GetTask(id)

	var err_ret ErrorRet
	if err != nil {

		err_ret.Error = fmt.Sprintf("ошибка получения задачи по id(%v)", err)
		writeJson(w, err_ret)
		return
	}

	if task.Repeat == "" {

		err = db.DeleteTask(id)
		if err != nil {
			
			err_ret.Error = fmt.Sprintf("ошибка удаления задачи по id(%v)", err)
			writeJson(w, err_ret)
			return
		}

	} else {

		//NextDate(now, task.Date, task.Repeat)
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {

			err_ret.Error = fmt.Sprintf("ошибка удаления задачи по id(%v)", err)
			writeJson(w, err_ret)
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {

			err_ret.Error = fmt.Sprintf("ошибка обновления удаления задачи по id(%v)", err)
			writeJson(w, err_ret)
			return
		}

	}

	emptyJSON := map[string]interface{}{}
	writeJson(w, emptyJSON)
}

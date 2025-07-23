package api

import (
	"fmt"
	"net/http"
	"time"
)

func nextDayHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		http.Error(w, "Сервер не поддерживает "+req.Method, http.StatusMethodNotAllowed)
		return
	}

	s_now := req.FormValue("now")

	now, err := time.Parse(Format_date, s_now)
	if err != nil {

		str := fmt.Sprintf("ошибка(%v) при преобразовании даты", err)
		http.Error(w, str, http.StatusInternalServerError)

		return
	}

	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	str_d, err := NextDate(now, date, repeat)
	if err != nil {

		str := fmt.Sprintf("ошибка(%v) при вычеслении даты", err)
		http.Error(w, str, http.StatusInternalServerError)
		//http.Error(w, err, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte(str_d))
}

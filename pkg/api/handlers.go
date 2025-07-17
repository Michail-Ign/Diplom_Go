package api

import (
	"fmt"
	"net/http"
	"time"
)

func nextDayHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Println("вызов nextDayHandler")
	if req.Method != http.MethodGet {
		http.Error(w, "Сервер не поддерживает "+req.Method, http.StatusMethodNotAllowed)
		return
	}

	s_now := req.FormValue("now")
	fmt.Println("now = " + s_now)

	now, err := time.Parse(Format_date, s_now)
	if err != nil {
		//fmt.Println("_err = " + err)

		str := fmt.Sprintf("ошибка(%v) при преобразовании даты", err)
		http.Error(w, str, http.StatusInternalServerError)
		//

		return
	}

	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	//
	fmt.Println("date = " + date)
	fmt.Println("repeat = " + repeat)

	str_d, err := NextDate(now, date, repeat)
	if err != nil {

		str := fmt.Sprintf("ошибка(%v) при вычеслении даты", err)
		http.Error(w, str, http.StatusInternalServerError)
		//http.Error(w, err, http.StatusMethodNotAllowed)
		return
	}
	//
	fmt.Println("    str_d = " + str_d)

	w.Write([]byte(str_d))
}

package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const Format_date = "20060102"

func contains(d time.Weekday, m []string) (bool, error) {

	for _, n := range m {

		interval, err := strconv.Atoi(n)
		if err != nil {
			return false, fmt.Errorf("[Parse int] %w", err)
		}

		if interval >= 8 {
			return false, fmt.Errorf("недопустимое значение " + n)
		}

		if int(d) == 0 && interval == 7 {
			return true, nil
		} else if int(d) == interval {
			return true, nil
		}
	}

	return false, nil
}

func getLastDate(date time.Time, d_interval int) time.Time {

	y := date.Year()
	m := date.Month()

	date_return := time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)

	if d_interval == -2 {
		date_return = date_return.AddDate(0, 0, -1)
	}

	return date_return
}

func searchDay(date time.Time, m_interval []string) (bool, error) {

	for _, d := range m_interval {

		day, err := strconv.Atoi(d)
		if err != nil {
			return false, fmt.Errorf("[parse int] %w", err)
		}

		if day > 31 {
			return false, fmt.Errorf("недопустимый день месяца" + d)
		}

		if day == -1 || day == -2 {
			//Получить последний день месяца
			l_day := getLastDate(date, day)
			if int(l_day.Day()) == int(date.Day()) {
				return true, nil
			}

		} else if day <= -3 {
			return false, fmt.Errorf("недопустимый день месяца" + d)

		} else if day == int(date.Day()) {
			return true, nil
		}
	}

	return false, nil //fmt.Errorf("не входит в период")
}

func containsPeriod(date time.Time, m_interval, month_interval []string) (bool, error) {

	//Проверка и по месяцу
	if len(month_interval) == 0 {
		return searchDay(date, m_interval)
	}

	for _, m := range month_interval {

		month, err := strconv.Atoi(m)
		if err != nil {
			return false, fmt.Errorf("[parse int] %w", err)
		}
		if month > 12 {
			return false, fmt.Errorf("недопустимый месяц " + m)
		}

		if int(date.Month()) == month {

			//Совпал месяц проверка по дате
			return searchDay(date, m_interval)

		}
	}

	return false, nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("пустая строка repeat")
	}

	//step 2
	date, err := time.Parse(Format_date, dstart)
	if err != nil {
		return "", err
	}

	//step3
	interval := 0
	massive := strings.Split(repeat, " ")

	switch massive[0] {

	case "d":
		if len(massive) == 1 {
			return "", fmt.Errorf("не указан интервал в днях")
		}

		interval, err = strconv.Atoi(massive[1])
		if err != nil {
			return "", fmt.Errorf("[Parse int] %w", err)
		}

		if interval > 400 {
			return "", fmt.Errorf("превышен максимально допустимый интервал")
		} else if interval == 0 {
			return "", fmt.Errorf("не указан интервал в днях")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(massive) == 1 {
			return "", fmt.Errorf("не указан день недели")
		}

		m_interval := strings.Split(massive[1], ",")

		for {
			date = date.AddDate(0, 0, 1)
			day_w := date.Weekday()

			ok, err := contains(day_w, m_interval)
			if ok {
				if afterNow(date, now) {
					break
				}
			} else if err != nil {
				return "", err
			}
		}

	case "m":

		if len(massive) == 1 {
			return "", fmt.Errorf("не указан день месяца")
		}
		m_interval := strings.Split(massive[1], ",")

		var month_interval []string

		if len(massive) == 3 {
			month_interval = strings.Split(massive[2], ",")
		}

		for {
			date = date.AddDate(0, 0, 1)

			ok, err := containsPeriod(date, m_interval, month_interval)
			if ok {

				if afterNow(date, now) {
					break
				}

			} else if err != nil {
				return "", err
			}
		}

	case "y":
		//date = date.AddDate(1, 0, 0)
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

	default:
		return "", fmt.Errorf("указан неверный формат repeat" + massive[0])
	}

	//step6
	to := Format_date
	return date.Format(to), nil
}

// Функция возвращает true, если первая дата(date) больше второй(now)
func afterNow(date, now time.Time) bool {

	dy, dm, dd := date.Date()
	ny, nm, nd := now.Date()

	if dy != ny {
		return dy > ny
	}

	if dm != nm {
		return dm > nm
	}

	return dd > nd
}

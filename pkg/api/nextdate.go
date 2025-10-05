package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения пустое")
	}

	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", errors.New("неверный формат даты dstart")
	}

	parts := strings.Split(repeat, " ")

    switch parts[0] {
    case "y":
        // Всегда сдвигаем как минимум на 1 год
        date = addOneYearAdjustingLeap(date)
        for !afterNow(date, now) {
            date = addOneYearAdjustingLeap(date)
        }
    case "d":
        if len(parts) != 2 {
            return "", errors.New("неверный формат правила d <число>")
        }
        interval, err := strconv.Atoi(parts[1])
        if err != nil || interval < 1 || interval > 400 {
            return "", errors.New("недопустимое количество дней")
        }
        // Всегда сдвигаем как минимум на 1 интервал
        date = date.AddDate(0, 0, interval)
        for !afterNow(date, now) {
            date = date.AddDate(0, 0, interval)
        }
    default:
        return "", errors.New("неизвестный формат repeat")
    }

	return date.Format(DateFormat), nil
}

// addOneYearAdjustingLeap добавляет 1 год и корректирует 29 февраля на 1 марта в невисокосные годы
func addOneYearAdjustingLeap(d time.Time) time.Time {
    next := d.AddDate(1, 0, 0)
    if d.Month() == time.February && d.Day() == 29 {
        // time.AddDate(1,0,0) даст 28 фев в невисокосный год — сдвигаем на 1 марта
        if next.Month() == time.February && next.Day() == 28 {
            next = next.AddDate(0, 0, 1)
        }
    }
    return next
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем GET-параметры через FormValue
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")
	nowStr := r.FormValue("now")

	if dateStr == "" || repeat == "" {
		http.Error(w, "параметры 'date' и 'repeat' обязательны", http.StatusBadRequest)
		return
	}

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("неверный формат now: %v", err), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка вычисления следующей даты: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}

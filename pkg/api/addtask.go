package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Vorobey112/go-final/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "invalid JSON: " + err.Error()}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "title is required"}, http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": "invalid date/repeat: " + err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "failed to add task: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]int64{"id": id}, http.StatusOK)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

    if task.Repeat != "" {
        // Если дата в прошлом — вычисляем следующее в будущем.
        // Если дата сегодня — оставляем сегодня (не сдвигаем вперёд при создании),
        // т.к. тест ожидает сегодняшнюю дату для правила d 1.
        if t.Before(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)) {
            next, err := NextDate(now, task.Date, task.Repeat)
            if err != nil {
                return err
            }
            task.Date = next
        } else if !afterNow(t, now) {
            // t == today: оставить как есть
            task.Date = now.Format("20060102")
        }
    } else if !afterNow(t, now) {
		task.Date = now.Format("20060102")
	}

	return nil
}

func writeJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

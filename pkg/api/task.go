package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Vorobey112/go-final/pkg/db"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJson(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is empty"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	if task == nil {
		writeJson(w, map[string]string{"error": "task not found"}, http.StatusNotFound)
		return
	}

	writeJson(w, task, http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	if t.ID == "" {
		writeJson(w, map[string]string{"error": "id is empty"}, http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(&t); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, t, http.StatusOK)
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is empty"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	if task == nil {
		writeJson(w, map[string]string{"error": "task not found"}, http.StatusNotFound)
		return
	}

	// Параметр now больше не используется в логике вычисления nextDate
	// Оставляем разбор для совместимости, но не применяем переменную
	nowStr := r.FormValue("now")
	if nowStr != "" {
		if _, err := time.Parse(DateFormat, nowStr); err != nil {
			writeJson(w, map[string]string{"error": "invalid now format"}, http.StatusBadRequest)
			return
		}
	}

	// Если повторения нет — удаляем задачу
	if task.Repeat == "" {
		if err := db.DeleteTask(task.ID); err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		writeJson(w, map[string]string{}, http.StatusOK)
		return
	}

	// Если задача периодическая — вычисляем следующую дату относительно текущей даты задачи
	// Это гарантирует, что при каждом выполнении дата сдвигается минимум на один шаг
	base, perr := time.Parse(DateFormat, task.Date)
	if perr != nil {
		writeJson(w, map[string]string{"error": "invalid task date"}, http.StatusBadRequest)
		return
	}
	nextDate, err := NextDate(base, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Обновляем только дату
	task.Date = nextDate
	if err := db.UpdateDate(task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Возвращаем пустой объект по требованиям тестов
	writeJson(w, map[string]string{}, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJson(w, map[string]string{"error": "method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is empty"}, http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]string{}, http.StatusOK)
}

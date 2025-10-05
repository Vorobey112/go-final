package api

import (
	"net/http"

	"github.com/Vorobey112/go-final/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		// здесь вызываете функцию, которая возвращает ошибку в JSON
		// её желательно было реализовать на предыдущем шаге
		writeJson(w, map[string]string{
			"error": err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}

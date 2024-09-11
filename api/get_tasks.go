package api

import (
	"encoding/json"
	"log"
	"net/http"

	"todo/model"
	"todo/service"
)

func (h *TodoHandlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	log.Printf("получен запрос [%s]", r.RequestURI)

	search := r.FormValue("search")

	tasks, err := h.Repo.GetTasks(search)
	if err != nil {
		log.Printf("ошибка получения данных: %s", err.Error())
		service.ErrorResponse(w, "ошибка получения данных", err)
		return
	}

	foundTasks, err := json.Marshal(model.Tasks{Tasks: tasks})
	if err != nil {
		log.Println(err.Error())
		service.ErrorResponse(w, "ошибка сериализации данных", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(foundTasks)
	if err != nil {
		log.Println(err.Error())
		service.ErrorResponse(w, "внутренняя ошибка сервера", err)
		return
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"todo/model"
	"todo/service"
)

func (h *TodoHandlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var buf bytes.Buffer

	log.Printf("получен запрос [%s]", r.RequestURI)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		log.Println(err.Error())
		service.ErrorResponse(w, "ошибка десериализации JSON", err)
		return
	}

	log.Printf("данные в запросе: task[%s]", task.String())

	err = service.CheckTask(&task)
	if err != nil {
		log.Println(err.Error())
		service.ErrorResponse(w, "ошибка в данных", err)
		return
	}

	err = h.Repo.Update(task)
	if err != nil {
		log.Println(err.Error())
		service.ErrorResponse(w, "внутренняя ощибка сервера", err)
		return
	}

	log.Printf("обновлена задача id=%s", task.Id)

	res, err := json.Marshal(map[string]any{})
	if err != nil {
		log.Println(err.Error())
		service.ErrorResponse(w, "ошибка сериализации данных", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(res)
	if err != nil {
		log.Println(err.Error())
		service.ErrorResponse(w, "внутренняя ошибка сервера", err)
		return
	}
}

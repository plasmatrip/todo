package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"todo/service"
)

func (h *TodoHandlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	log.Printf("получен запрос [%s]", r.RequestURI)

	v := r.FormValue("id")
	if len(v) == 0 {
		log.Println("не указан идентификатор")
		service.ErrorResponse(w, "не указан идентификатор: ", errors.New("идентификатор не может быть пустым"))
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		log.Println("идентификатор не может быть преобразован в число")
		service.ErrorResponse(w, "ошибка в запросе: ", errors.New("идентификатор не может быть преобразован в число"))
		return
	}

	err = h.Repo.Delete(id)
	if err != nil {
		log.Println(err.Error())
		service.ErrorResponse(w, "внутренняя ощибка сервера", err)
		return
	}

	log.Printf("удалена задача id=%d", id)

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

package api

import (
	"log"
	"net/http"
	"time"

	"todo/configs"
	"todo/service"
)

func (h *TodoHandlers) NextDate(w http.ResponseWriter, r *http.Request) {
	log.Printf("получен запрос [%s]", r.RequestURI)

	now, err := time.Parse(configs.DateLayout, r.FormValue("now"))
	if err != nil {
		log.Printf("%s [now=%s]", err.Error(), r.FormValue("now"))
		http.Error(w, "переданное значение не может быть преобразовано в дату", http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nextDate, err := service.NextDate(now, date, repeat)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("отправлен ответ [%s]", nextDate)
}

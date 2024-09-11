package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"

	"todo/api"
	"todo/api/middleware"
	"todo/configs"
	"todo/repository"
)

func main() {
	log.Println("загружаю переменные окружения")
	configs.LoadEnv()

	log.Println("перенаправляю логирование в файл /log/app.log")
	configs.StartLog()
	defer configs.StopLog()

	log.Println("инициализирую базу данных /db/scheduller.db")
	db := repository.NewToDo()
	defer db.Close()

	log.Println("запускаю роутер запросов")
	r := chi.NewRouter()

	log.Println("подключаю зависимости")
	todoHandlers := api.NewTodoHandlers(db)

	log.Println("подключаю хэндлеры запросов")
	r.Mount("/", http.FileServer(http.Dir(configs.WebDir)))
	r.Get("/api/nextdate", todoHandlers.NextDate)
	r.Post("/api/task", middleware.Auth(todoHandlers.AddTask))
	r.Get("/api/task", middleware.Auth(todoHandlers.GetTask))
	r.Put("/api/task", middleware.Auth(todoHandlers.UpdateTask))
	r.Get("/api/tasks", middleware.Auth(todoHandlers.GetTasks))
	r.Post("/api/task/done", middleware.Auth(todoHandlers.TaskDone))
	r.Delete("/api/task", middleware.Auth(todoHandlers.DeleteTask))
	r.Post("/api/signin", todoHandlers.Login)

	err := http.ListenAndServe(":"+configs.Port, func(next http.Handler) http.Handler {
		log.Printf("сервер запущен, порт: %s\n", configs.Port)
		return next
	}(r))
	if err != nil {
		log.Panic(err)
	}
}

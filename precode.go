package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	//Сериализируем данные запроса
	resp, err := json.Marshal(tasks)
	//Обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Установливается заголовок JSON
	w.Header().Set("Content-type", "application/json")
	//Установка статуса 200, успешный запрос
	w.WriteHeader(http.StatusOK)
	//Отправка данных клиенту
	w.Write(resp)
}

func postTasks(w http.ResponseWriter, r *http.Request) {
	//Создаются переменные: task - для хранения десериализованных данных. buf - это буфер для временного хранения данных из тела запроса.
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getTaskId(w http.ResponseWriter, r *http.Request) {
	//Получаем id нужно задачи из URL
	id := chi.URLParam(r, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNoContent)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	//первая переменная пропускается, так как мы ничего не запрашиваем у сервера и ничего не добавляем. Следовательно и сериализация не нужна.
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	//Используется для удаления
	delete(tasks, id)
	//нужно задать статус 204
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getAllTasks)

	r.Post("/tasks", postTasks)

	r.Get("/tasks/{id}", getTaskId)

	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

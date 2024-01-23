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

// обработчик для вывода всех задач при запросе GET /tasks
func getAllTasks(res http.ResponseWriter, req *http.Request) {
	// переводим мапу в слайс байт
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// определяем заголовок с типом контента
	res.Header().Set("Content-Type", "application/json")

	// так как все успешно, то статус OK
	res.WriteHeader(http.StatusOK)

	res.Write(resp)
}

// обработчик для вывода задачи по индексу id при запросе GET /tasks/{id}
func getTask(res http.ResponseWriter, req *http.Request) {
	// считываем параметр id
	id := chi.URLParam(req, "id")

	// проверяем наличие данного id в мапе
	task, ok := tasks[id]
	if !ok {
		http.Error(res, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// переводим структуру в слайс байт
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// определяем заголовок с типом контента
	res.Header().Set("Content-Type", "application/json")

	// так как все успешно, то статус ОК
	res.WriteHeader(http.StatusOK)

	res.Write(resp)
}

// обработчик для добавления структуры в мапу при запросе POST /tasks
func postTask(res http.ResponseWriter, req *http.Request) {
	// объявляем переменные
	var task Task
	var buf bytes.Buffer

	// считываем тело запроса
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// переводим слайс байт в структуру
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// добавляем структуру в мапу
	tasks[task.ID] = task

	// определяем заголовок с типом контента
	res.Header().Set("Content-Type", "application/json")

	// сохраняем статус 201 Created
	res.WriteHeader(http.StatusCreated)
}

// обработчик для удаления записи при запросе DELETE /tasks
func deleteTask(res http.ResponseWriter, req *http.Request) {
	// считываем параметр id
	id := chi.URLParam(req, "id")

	// проверяем наличие в мапе
	_, ok := tasks[id]
	if !ok {
		http.Error(res, "Задача не найдена", http.StatusBadRequest)
		return
	}

	// удаляем запись
	delete(tasks, id)

	// определяем заголовок с типом контента
	res.Header().Set("Content-Type", "application/json")

	// так как все успешно, то статус ОК
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// регистрируем обработчики
	r.Get("/tasks", getAllTasks)
	r.Get("/tasks/{id}", getTask)
	r.Post("/tasks", postTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

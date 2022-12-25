package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	userId := fmt.Sprintf("%s", r.Context().Value("id"))
	if len(userId) == 0 {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "invalid user id"}`))
		return
	}
	convUserId, err := strconv.Atoi(userId)
	if err != nil {
		panic(err)
	}

	taskID := r.URL.Query().Get("task_id")
	if taskID == "" {
		tasks, err := t.taskService.GetTasks(r.Context(), convUserId)
		if err != nil {
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "error internal server"}`))
			return
		}
		jsonInBytes, err := json.Marshal(tasks)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonInBytes)
	} else {
		convTaskID, err := strconv.Atoi(taskID)
		if err != nil {
			panic(err)
		}
		tasks, err := t.taskService.GetTaskByID(r.Context(), convTaskID)
		if err != nil {
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "error internal server"}`))
			return
		}
		jsonInBytes, err := json.Marshal(tasks)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonInBytes)
	}
	// TODO: answer here
}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	// fmt.Println("ini api task Crate Task : ", task)
	newTask := entity.Task{}
	newTask.Title = task.Title
	newTask.Description = task.Description
	newTask.CategoryID = task.CategoryID
	if len(newTask.Title) == 0 || len(newTask.Description) == 0 || newTask.CategoryID == 0 {
		w.WriteHeader(400) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "invalid task request"}`))
		return
	}

	taskUserId := fmt.Sprintf("%s", r.Context().Value("id"))
	if len(taskUserId) == 0 {
		w.WriteHeader(400) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "invalid user id"}`))
		return
	}
	convTaskUserId, err := strconv.Atoi(taskUserId)
	if err != nil {
		panic(err)
	}
	newTask.UserID = convTaskUserId
	// fmt.Println("ini api task Crate Task (new task): ", newTask)
	resultTask, err := t.taskService.StoreTask(r.Context(), &newTask)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "error internal server"}`))
	}
	// fmt.Println("ini api task Crate Task (result repo task): ", resultTask)
	res := make(map[string]interface{})
	res["message"] = "success create new task"
	res["user_id"] = resultTask.UserID
	res["task_id"] = resultTask.ID
	jsonInBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(201) // WriteHeader ditulis di awal!
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
	// TODO: answer here
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskUserId := fmt.Sprintf("%s", r.Context().Value("id"))
	convUserID, err := strconv.Atoi(taskUserId)
	if err != nil {
		panic(err)
	}
	// fmt.Println("ini api task DeleteTask (id user) : ", convTaskUserId)
	taskID := r.URL.Query().Get("task_id")
	convTaskID, err := strconv.Atoi(taskID)
	if err != nil {
		panic(err)
	}
	// fmt.Println("ini api task DeleteTask (id task) : ", taskID)
	err = t.taskService.DeleteTask(r.Context(), convTaskID)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "error internal server"}`))
		return
	}
	res := make(map[string]interface{})
	res["message"] = "success delete task"
	res["user_id"] = convUserID
	res["task_id"] = convTaskID
	jsonInBytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200) // WriteHeader ditulis di awal!
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
	// TODO: answer here
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	taskUserId := fmt.Sprintf("%s", r.Context().Value("id"))
	if len(taskUserId) == 0 {
		w.WriteHeader(400) // WriteHeader ditulis di awal!
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "invalid user id"}`))
		return
	}
	convTaskUserId, err := strconv.Atoi(taskUserId)
	if err != nil {
		panic(err)
	}
	sentTask := entity.Task{}
	sentTask.ID = task.ID
	sentTask.Title = task.Title
	sentTask.Description = task.Description
	sentTask.CategoryID = task.CategoryID
	sentTask.UserID = convTaskUserId
	// fmt.Println("ini api task UpdateTask (idcategory) : ", sentTask.CategoryID)
	// fmt.Println("ini api task UpdateTask : ", sentTask)
	updateTask, err := t.taskService.UpdateTask(r.Context(), &sentTask)
	if err != nil {
		if err != nil {
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error": "error internal server"}`))
		}
	}
	res := make(map[string]interface{})
	res["message"] = "success update task"
	res["user_id"] = updateTask.UserID
	res["task_id"] = updateTask.ID
	jsonInBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200) // WriteHeader ditulis di awal!
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInBytes)
	// TODO: answer here
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}

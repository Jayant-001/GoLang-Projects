package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"task-manager/model"
	"task-manager/service"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Application struct {
	TaskManager *service.TaskManager
	scheduler   *service.Scheduler
}

func main() {

	router := httprouter.New()

	taskManager := service.NewTaskManager()
	scheduler := service.NewScheduler(100, taskManager)
	defer scheduler.Stop()

	app := &Application{
		TaskManager: taskManager,
		scheduler:   scheduler,
	}

	router.GET("/tasks", app.getAllTasks)
	router.POST("/tasks", app.createTask)
	router.GET("/tasks/:id", app.getTaskByID)
	router.GET("/test", app.TestHandler)

	srv := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		fmt.Println("Server is running on port 8000")
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("Server failed to start: ", err.Error())
		}
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shutdown.Done()

	fmt.Println("Shutting down the server ")
	// ctx. stop := context.Background()
	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Println("Server shutdown with error: ", err.Error())
	}
	fmt.Println("Server stopped gracefully")
}

func (app *Application) getAllTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	response, err := json.Marshal(app.TaskManager.GetAllTasks())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching tasks: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", response)
}

func (app *Application) createTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	var task model.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding task: %v", err), http.StatusBadRequest)
		return
	}

	taskModel := model.Task{
		Title:       task.Title,
		Description: task.Description,
	}

	newTask := app.TaskManager.CreateTask(taskModel)

	app.scheduler.AddTask(newTask)

	response, err := json.Marshal(map[string]any{
		"task": newTask,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating task: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (app *Application) getTaskByID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	w.Header().Set("Content-Type", "application/json")
	taskId := params.ByName("id")

	task := app.TaskManager.GetTaskByID(taskId)
	response, err := json.Marshal(map[string]any{"task_id": task})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (app *Application) TestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	time.Sleep(time.Second * 10)
	w.Write([]byte("Test handler executed after 10 seconds delay"))
}

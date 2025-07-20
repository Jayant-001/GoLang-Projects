package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"task-manager/model"
	"time"

	"github.com/google/uuid"
)

var (
	tasksDBPath = "./data/tasks.json"
)

type TaskManager struct {
	tasks []model.Task
	mu    sync.RWMutex
}

func NewTaskManager() *TaskManager {
	// check if file does not exist, create it
	if _, err := os.Stat(tasksDBPath); os.IsNotExist(err) {
		_, err := os.Create(tasksDBPath)
		if err != nil {
			panic("Failed to create tasks database file: " + err.Error())
		}
	}

	// Load users from the database
	tasks, err := loadTasksFromFile()
	if err != nil {
		panic("Failed to unmarshal tasks data: " + err.Error())
	}

	return &TaskManager{
		tasks: tasks,
	}
}

func (tm *TaskManager) GetAllTasks() []model.Task {
	return tm.tasks
}

func (tm *TaskManager) CreateTask(task model.Task) model.Task {

	task.ID = uuid.New().String()
	min := 10
	max := 20
	delayedTime := rand.Intn(max-min) + min
	task.Delayed = time.Duration(delayedTime) * time.Second
	task.Status = "pending"

	tm.mu.Lock()
	tm.tasks = append(tm.tasks, task)
	tm.mu.Unlock()

	err := writeTasksToFile(tm.tasks)
	if err != nil {
		panic("Failed to write tasks to file: " + err.Error())
	}
	return task
}

func (tm *TaskManager) UpdateTask(task model.Task) error {
	for i, t := range tm.tasks {
		if t.ID == task.ID {

			tm.mu.Lock()
			tm.tasks[i] = task
			tm.mu.Unlock()

			err := writeTasksToFile(tm.tasks)
			if err != nil {
				return fmt.Errorf("Failed to write tasks to file: %v", err)
			}
		}
	}
	return nil
}

func (tm *TaskManager) GetTaskByID(id string) *model.Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	for _, task := range tm.tasks {
		if id == task.ID {
			return &task
		}
	}
	return nil
}

func writeTasksToFile(tasks []model.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to marshal tasks data: " + err.Error())
	}

	// os.WriteFile is simpler and handles file creation and permissions.
	return os.WriteFile(tasksDBPath, data, 0644)
}

func loadTasksFromFile() ([]model.Task, error) {
	data, err := os.ReadFile(tasksDBPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read tasks database file: " + err.Error())
	}

	if len(data) == 0 {
		return []model.Task{}, nil
	}

	var tasks []model.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal tasks data: " + err.Error())
	}
	return tasks, nil
}

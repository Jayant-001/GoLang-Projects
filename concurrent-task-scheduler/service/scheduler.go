package service

import (
	"sync"
	"task-manager/model"
	"time"
)

type Scheduler struct {
	taskManager *TaskManager
	requestChan chan model.Task
	mu          sync.RWMutex
	wg          sync.WaitGroup
	isShutdown  bool
}

func NewScheduler(bufferSize int, taskManager *TaskManager) *Scheduler {

	s := &Scheduler{
		requestChan: make(chan model.Task, bufferSize),
		taskManager: taskManager,
	}

	go s.Start()

	return s
}

func (s *Scheduler) Start() {

	for i := 0; i < 5; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for task := range s.requestChan {
				processedTask := processTask(task)
				s.taskManager.UpdateTask(processedTask)
			}
		}()
	}

	s.wg.Wait()
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isShutdown {
		s.mu.Unlock()
		return
	}
	s.isShutdown = true
	s.mu.Unlock()

	close(s.requestChan)
	s.wg.Wait()
	
}

func (s *Scheduler) AddTask(task model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestChan <- task
}

func processTask(task model.Task) model.Task {
	time.Sleep(time.Duration(task.Delayed) * time.Second)
	task.Status = "completed"
	task.CreatedAt = time.Now()
	return task
}

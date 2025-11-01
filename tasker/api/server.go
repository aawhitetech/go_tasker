package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"tasker/task"
)

type Server struct {
	mu    sync.Mutex
	tasks []task.Task
}

func NewServer() (*Server, error) {
	s := &Server{}
	if err := s.loadFromDisk(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Start(address string) error {
	http.HandleFunc("/tasks", s.handleTasksEntry)
	http.HandleFunc("/tasks/", s.handleTasksSubroute)

	return http.ListenAndServe(address, nil)
}

func (s *Server) loadFromDisk() error {
	t, err := task.Load()
	if err != nil {
		return err
	}

	s.tasks = t

	return nil
}

func (s *Server) saveToDisk() error {
	return task.Save(s.tasks)
}

func (s *Server) handleTasksEntry(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tasks" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.handleList(w,r)
	case http.MethodPost:
		s.handleAdd(w,r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTasksSubroute(w http.ResponseWriter, r *http.Request) {
	s.handleMarkDone(w,r)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.tasks)
}

type AddRequest struct {
	Description string `json:"description"`
}

func (s *Server) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if req.Description == "" {
		http.Error(w, "Description Required", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks = task.Add(s.tasks, req.Description)

	if err := s.saveToDisk(); err != nil {
		http.Error(w, "Failed to Save", http.StatusInternalServerError)
		return
	}

	newTask := s.tasks[len(s.tasks)-1]

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func (s *Server) handleMarkDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	if len(parts) != 2 || parts[1] != "done" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(parts[0])
	if err != nil || id <= 0 {
		http.Error(w, "Task Id Invalid", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ok, updated := task.MarkDone(s.tasks, id)
	if ok == false {
		http.Error(w, "Task Not Found", http.StatusNotFound)
		return
	}
	s.tasks = updated

	if err := s.saveToDisk(); err != nil {
		http.Error(w, "Failed to Save", http.StatusInternalServerError)
		return
	}

	var marked task.Task
	for _, v := range s.tasks {
		if v.ID == id {
			marked = v
			break
		}
	}
	
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(marked)
}

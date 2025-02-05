package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos_pointer, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	var todos []model.TODO

	if len(todos_pointer) == 0 {
		return &model.ReadTODOResponse{TODOs: []model.TODO{}}, nil
	}

	for _, todo := range todos_pointer {
		todos = append(todos, *todo)
	}

	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Decode the request body
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request: Invalid JSON format", http.StatusBadRequest)
			return
		}

		if req.Subject == "" {
			http.Error(w, "Bad Request: Subject cannot be empty", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		todo, err := h.Create(ctx, &req)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Content-Type ヘッダを application/json に設定
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(todo); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		// Decode the request body
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request: Invalid JSON format", http.StatusBadRequest)
			return
		}

		if req.ID == 0 || req.Subject == "" {
			http.Error(w, "Bad Request: Subject cannot be empty or ID cannot be 0", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		todo, err := h.Update(ctx, &req)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Content-Type ヘッダを application/json に設定
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(todo); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case http.MethodGet:
		// URLのクエリパラメーターから取得
		prevID, _ := strconv.ParseInt(r.URL.Query().Get("prev_id"), 10, 64)
		size, _ := strconv.ParseInt(r.URL.Query().Get("size"), 10, 64)

		req := model.ReadTODORequest{
			PrevID: prevID,
			Size:   size,
		}

		ctx := r.Context()
		todos, err := h.Read(ctx, &req)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Content-Type ヘッダを application/json に設定
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(todos); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		// Decode the request body
		var req model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request: Invalid JSON format", http.StatusBadRequest)
			return
		}

		if len(req.IDs) == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		res, _ := h.Delete(ctx, &req)
		if res == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// Content-Type ヘッダを application/json に設定
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

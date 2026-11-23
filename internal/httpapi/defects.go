package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/dxssert/lab4KR/internal/service"
)

type DefectController struct {
	service *service.Defects
}

func NewRouter(defects *service.Defects) http.Handler {
	controller := &DefectController{service: defects}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/defects", controller.list)
	mux.HandleFunc("GET /api/v1/defects/{id}", controller.get)
	mux.HandleFunc("POST /api/v1/defects", controller.create)
	mux.HandleFunc("PUT /api/v1/defects/{id}", controller.update)
	mux.HandleFunc("DELETE /api/v1/defects/{id}", controller.delete)
	return mux
}

func (c *DefectController) list(w http.ResponseWriter, r *http.Request) {
	items, err := c.service.List(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	result := make([]DefectResponse, 0, len(items))
	for _, item := range items {
		result = append(result, defectResponse(item))
	}
	writeJSON(w, http.StatusOK, result)
}

func (c *DefectController) get(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	item, err := c.service.Get(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, defectResponse(*item))
}

func (c *DefectController) create(w http.ResponseWriter, r *http.Request) {
	var input CreateDefectRequest
	if !readJSON(w, r, &input) {
		return
	}
	item, err := c.service.Create(r.Context(), input.data(), input.CreatedBy)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Location", "/api/v1/defects/"+strconv.FormatInt(item.ID, 10))
	writeJSON(w, http.StatusCreated, defectResponse(*item))
}

func (c *DefectController) update(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	var input DefectRequest
	if !readJSON(w, r, &input) {
		return
	}
	item, err := c.service.Update(r.Context(), id, input.data())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, defectResponse(*item))
}

func (c *DefectController) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := readID(w, r)
	if !ok {
		return
	}
	if err := c.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func readID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "ID должен быть положительным целым числом"})
		return 0, false
	}
	return id, true
}

func readJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err == nil {
		if err = decoder.Decode(new(json.RawMessage)); errors.Is(err, io.EOF) {
			return true
		}
	}
	writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "некорректное тело JSON"})
	return false
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, service.ErrNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	default:
		log.Printf("defects: %v", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "ошибка при работе с журналом брака"})
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write JSON: %v", err)
	}
}

package httpapi

import (
	"time"

	"github.com/dxssert/lab4KR/internal/model"
	"github.com/dxssert/lab4KR/internal/service"
)

type DefectRequest struct {
	ProductionOutputID int64  `json:"production_output_id"`
	ReasonID           int64  `json:"reason_id"`
	DefectGrams        int64  `json:"defect_grams"`
	Comment            string `json:"comment"`
}

type CreateDefectRequest struct {
	DefectRequest
	CreatedBy int64 `json:"created_by"`
}

type DefectResponse struct {
	ID int64 `json:"id"`
	DefectRequest
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (r DefectRequest) data() service.DefectData {
	return service.DefectData{
		ProductionOutputID: r.ProductionOutputID,
		ReasonID:           r.ReasonID,
		DefectGrams:        r.DefectGrams,
		Comment:            r.Comment,
	}
}

func defectResponse(d model.Defect) DefectResponse {
	return DefectResponse{
		ID: d.ID,
		DefectRequest: DefectRequest{
			ProductionOutputID: d.ProductionOutputID,
			ReasonID:           d.ReasonID,
			DefectGrams:        d.DefectGrams,
			Comment:            d.Comment,
		},
		CreatedBy: d.CreatedBy,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

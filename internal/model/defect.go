package model

import "time"

type Defect struct {
	ID                 int64     `json:"id"`
	ProductionOutputID int64     `json:"production_output_id"`
	ReasonID           int64     `json:"reason_id"`
	DefectGrams        int64     `json:"defect_grams"`
	Comment            string    `json:"comment"`
	CreatedBy          int64     `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

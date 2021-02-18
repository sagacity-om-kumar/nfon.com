package models

type AddTemplateRequest struct {
	TemplateID int               `db:"id" json:"templateid"`
	TemplateTitle string               `db:"name" json:"templatetitle"`
	Headers       []TemplateHeaderData `json:"headers"`
	Type          uint                 `db:"type" json:"templatetype"`
	Description   string               `db:"description" json:"description"`
	CategoryID    int                  `db:"category_id" json:"categoryid"`
	CreatedBy    int                  `db:"created_by" json:"createdby"`
}

type TemplateHeaderData struct {
	CategoryName     string `db:"name" json:"name"`
	CategoryHeaderId int    `db:"category_header_id" json:"categoryHeaderId"`
	TemplateID       int    `db:"template_id" json:""`
	Datatype         string `db:"datatype" json:"datatype"`
}

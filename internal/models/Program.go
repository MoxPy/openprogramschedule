package models

// Program Pointer allows null value
type Program struct {
	Id           *uint  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Host         string `json:"host"`
	Category     string `json:"category"`
	InProduction *bool  `json:"in_production"`
}

package models

type Schedule struct {
	Id          *uint  `json:"id"`
	ProgramId   uint   `json:"program_id"`
	Description string `json:"description"`
	Day         string `json:"day"`
	Date        string `json:"date"`
}

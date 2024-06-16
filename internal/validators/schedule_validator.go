package validators

import (
	"errors"
	"openprogramschedule/internal/models"
)

func ValidateSchedule(schedule *models.Schedule) error {

	// Schedule description validation
	if len(schedule.Description) == 0 {
		return errors.New("invalid input: schedule description is missing")
	}
	if len(schedule.Description) < 3 {
		return errors.New("invalid input: schedule description must be at least 3 characters")
	}
	if len(schedule.Description) > 100 {
		return errors.New("invalid input: schedule description must be less than 100 characters")
	}

	return nil
}

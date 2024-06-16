package validators

import (
	"errors"
	"openprogramschedule/internal/models"
)

func ValidateProgram(program *models.Program) error {
	// Program name validation
	if len(program.Name) == 0 {
		return errors.New("invalid input: program name is required")
	}
	if len(program.Name) < 3 {
		return errors.New("invalid input: program name must be at least 3 characters")
	}
	if len(program.Name) > 52 {
		return errors.New("invalid input: program name must be less than 52 characters")
	}

	// Program description validation
	if len(program.Description) == 0 {
		return errors.New("invalid input: program description is missing")
	}
	if len(program.Description) < 3 {
		return errors.New("invalid input: program description must be at least 3 characters")
	}
	if len(program.Description) > 100 {
		return errors.New("invalid input: program description must be less than 100 characters")
	}

	// Program host validation
	if len(program.Host) == 0 {
		return errors.New("invalid input: program host is missing")
	}
	if len(program.Host) < 3 {
		return errors.New("invalid input: program host must be at least 3 characters")
	}
	if len(program.Host) > 52 {
		return errors.New("invalid input: program host must be less than 52 characters")
	}

	// Program category validation
	if len(program.Category) == 0 {
		return errors.New("invalid input: program category is missing")
	}
	if len(program.Category) < 3 {
		return errors.New("invalid input: program category must be at least 3 characters")
	}
	if len(program.Category) > 52 {
		return errors.New("invalid input: program category must be less than 52 characters")
	}

	// Program inProduction validation
	if program.InProduction == nil {
		return errors.New("invalid input: program is missing inProduction")
	}

	return nil
}

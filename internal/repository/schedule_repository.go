package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"openprogramschedule/internal/models"
	"time"
)

// IT
var daysOfTheWeek = map[int]string{
	1: "Lunedi",
	2: "Martedi",
	3: "Mercoledi",
	4: "Giovedi",
	5: "Venerdi",
	6: "Sabato",
	7: "Domenica",
}

// AddSchedule Create a schedule
func AddSchedule(schedule *models.Schedule, db *sql.DB) (uint, error) {
	_, err := GetProgramByID(schedule.ProgramId, db)
	if err != nil {
		return 0, errors.New("could not get program")
	}
	query := `INSERT INTO schedules (program_id, description, day, date)
			VALUES (@p1, @p2, @p3, @p4);
			SELECT SCOPE_IDENTITY() AS id`

	row := db.QueryRow(query,
		sql.Named("p1", schedule.ProgramId),
		sql.Named("p2", schedule.Description),
		sql.Named("p3", schedule.Day),
		sql.Named("p4", schedule.Date),
	)

	var id uint
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	log.Printf("Added schedule with id: %d", id)
	return id, nil
}

// GetAllSchedules Get all schedules
func GetAllSchedules(db *sql.DB) ([]models.Schedule, error) {
	query := `SELECT * FROM schedules`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(&schedule.Id, &schedule.ProgramId, &schedule.Description, &schedule.Day, &schedule.Date)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return schedules, nil
}

// GetScheduleByID Get a schedule by its ID
func GetScheduleByID(scheduleID uint, db *sql.DB) (*models.Schedule, error) {
	query := `SELECT id, program_id, description, day, date FROM schedules WHERE id = @p1;`
	row := db.QueryRow(query, sql.Named("p1", scheduleID))
	var schedule models.Schedule
	err := row.Scan(&schedule.Id, &schedule.ProgramId, &schedule.Description, &schedule.Day, &schedule.Date)
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

// GetScheduleByProgramID Get the schedule of a program using its ID
func GetScheduleByProgramID(programId uint, db *sql.DB) (*[]models.Schedule, error) {
	_, err := GetProgramByID(programId, db)
	if err != nil {
		return nil, errors.New("could not get program")
	}
	query := `SELECT id, program_id, description, day, date FROM schedules WHERE program_id = @p1`
	rows, err := db.Query(query, sql.Named("p1", programId))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(&schedule.Id, &schedule.ProgramId, &schedule.Description, &schedule.Day, &schedule.Date)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &schedules, nil
}

// The day parameter should be an integer representing the day of the week (1 for Monday, 7 for Sunday). Days are in italian (daysOfTheWeek)
func GetScheduleByDay(day int, db *sql.DB) (*[]models.Schedule, error) {
	dayName, exists := daysOfTheWeek[day]
	if !exists {
		return nil, errors.New("invalid day number")
	}

	query := `SELECT id, program_id, description, day, date FROM schedules WHERE day = @p1;`
	rows, err := db.Query(query, sql.Named("p1", dayName))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	var schedules []models.Schedule

	for rows.Next() {
		var schedule models.Schedule
		err = rows.Scan(&schedule.Id, &schedule.ProgramId, &schedule.Description, &schedule.Day, &schedule.Date)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &schedules, nil
}

// GetScheduleByDate Get the schedule of a date (es. 2024-06-30)
func GetScheduleByDate(date string, db *sql.DB) (*[]models.Schedule, error) {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}
	start := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	query := `SELECT id, program_id, description, day, date FROM schedules WHERE date >= @p1 AND date <= @p2;`
	rows, err := db.Query(query, sql.Named("p1", start), sql.Named("p2", end))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err = rows.Scan(&schedule.Id, &schedule.ProgramId, &schedule.Description, &schedule.Day, &schedule.Date)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &schedules, nil
}

// UpdateScheduleByID
func UpdateScheduleByID(scheduleID uint, updatedSchedule models.Schedule, db *sql.DB) error {
	query := `UPDATE schedules SET program_id = @program_id, description = @description, day = @day, date = @date WHERE id = @id;`

	_, err := db.Exec(query,
		sql.Named("program_id", updatedSchedule.ProgramId),
		sql.Named("description", updatedSchedule.Description),
		sql.Named("day", updatedSchedule.Day),
		sql.Named("date", updatedSchedule.Date),
		sql.Named("id", scheduleID),
	)

	if err != nil {
		return err
	}

	log.Println("Updated schedule with id:", scheduleID)

	return nil
}

// DeleteScheduleByID
func DeleteScheduleByID(scheduleID uint, db *sql.DB) error {
	query := `DELETE FROM schedules WHERE id = @p1;`
	_, err := db.Exec(query, sql.Named("p1", scheduleID))
	if err != nil {
		return err
	}
	log.Printf("Deleted schedule: %+v\n", scheduleID)
	return nil
}

// DeleteAllSchedules
func DeleteAllSchedules(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	query := `DELETE FROM schedules;`
	result, err := tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute delete query: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	log.Printf("Deleted %d schedules", rowsAffected)
	return nil
}

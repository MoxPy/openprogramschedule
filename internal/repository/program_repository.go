package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"openprogramschedule/internal/models"
)

// AddProgram Create new program
func AddProgram(program *models.Program, db *sql.DB) (uint, error) {
	query := `INSERT INTO programs (name, description, host, category, in_production)
             VALUES (@p1, @p2, @p3, @p4, @p5);
             SELECT SCOPE_IDENTITY() AS id`

	row := db.QueryRow(query, sql.Named("p1", program.Name),
		sql.Named("p2", program.Description),
		sql.Named("p3", program.Host),
		sql.Named("p4", program.Category),
		sql.Named("p5", program.InProduction))

	var id uint
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error while creating the program: %v", err)
	}

	fmt.Printf("Added new program: %+v\n", program.Name)
	return id, nil
}

// GetProgramByID Get Program by id
func GetProgramByID(programID uint, db *sql.DB) (*models.Program, error) {
	query := `SELECT * FROM programs WHERE id = @p1;`
	row := db.QueryRow(query, sql.Named("p1", programID))
	var program models.Program
	err := row.Scan(&program.Id, &program.Name, &program.Description, &program.Host, &program.Category, &program.InProduction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &program, errors.New("program not found")
		}
		return &program, err
	}

	return &program, nil
}

// GetProgramByName Get Program by name
func GetProgramByName(programName string, db *sql.DB) (*models.Program, error) {
	query := `SELECT * FROM programs WHERE name = @p1;`
	row := db.QueryRow(query, sql.Named("p1", programName))
	var program models.Program
	err := row.Scan(&program.Id, &program.Name, &program.Description, &program.Host, &program.Category, &program.InProduction)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &program, errors.New("program not found")
		}
		return &program, err
	}

	return &program, nil
}

// GetProgramsByCategory Get programs by category
func GetProgramsByCategory(category string, db *sql.DB) ([]models.Program, error) {
	query := `SELECT * FROM programs WHERE category = @p1;`
	rows, err := db.Query(query, sql.Named("p1", category))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	var programs []models.Program
	for rows.Next() {
		var program models.Program
		err = rows.Scan(&program.Id, &program.Name, &program.Description, &program.Host, &program.Category, &program.InProduction)
		if err != nil {
			return nil, err
		}
		programs = append(programs, program)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return programs, nil
}

// GetAllPrograms Get all programs
func GetAllPrograms(db *sql.DB) ([]models.Program, error) {
	query := `SELECT * FROM programs;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	var programs []models.Program
	for rows.Next() {
		var program models.Program
		err = rows.Scan(&program.Id, &program.Name, &program.Description, &program.Host, &program.Category, &program.InProduction)
		if err != nil {
			return nil, err
		}
		programs = append(programs, program)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return programs, nil
}

// UpdateProgramByID Update program by id
func UpdateProgramByID(programID uint, updatedProgram models.Program, db *sql.DB) error {
	query := `UPDATE programs SET name = @name, description = @description, host = @host, category = @category, in_production = @in_production WHERE id = @id;`

	_, err := db.Exec(query,
		sql.Named("name", updatedProgram.Name),
		sql.Named("description", updatedProgram.Description),
		sql.Named("host", updatedProgram.Host),
		sql.Named("category", updatedProgram.Category),
		sql.Named("in_production", updatedProgram.InProduction),
		sql.Named("id", programID),
	)

	if err != nil {
		return err
	}

	log.Printf("Program updated: %+v\n", updatedProgram.Name)

	return nil
}

// DeleteProgram Delete program
func DeleteProgram(programID uint, db *sql.DB) error {
	query := `DELETE FROM programs WHERE id = @p1;`
	_, err := db.Exec(query, sql.Named("p1", programID))
	if err != nil {
		return err
	}
	log.Printf("Deleted program: %+v\n", programID)
	return nil
}

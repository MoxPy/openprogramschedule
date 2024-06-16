package db

import (
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"os"
)

var db *sql.DB

type EnvDBConfig struct {
	user     string
	password string
	server   string
	name     string
	port     string
}

func NewEnvDBConfig() *EnvDBConfig {
	return &EnvDBConfig{
		user:     os.Getenv("DB_USER"),     // Es. "myuser"
		password: os.Getenv("DB_PASSWORD"), // Es. "mypassword"
		server:   os.Getenv("DB_HOST"),     // Es. "myserver.mysql.database.azure.com"
		name:     os.Getenv("DB_NAME"),     // Es. "mydatabase"
		port:     os.Getenv("DB_PORT"),     // Es. "3306"
	}
}

// ConnectDB Remember to allow your ip on SQL Server, or it doesn't work
func ConnectDB() *sql.DB {
	config := NewEnvDBConfig()

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		config.server, config.user, config.password, config.port, config.name)
	var err error
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	_, err = db.Exec(`IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='programs' AND xtype='U')
BEGIN
    CREATE TABLE programs (
        id INT IDENTITY(1,1) PRIMARY KEY,
        name NVARCHAR(255) NOT NULL,
        description NTEXT,
        host NVARCHAR(255),
        category NVARCHAR(255),
        in_production BIT
    )
END`)
	if err != nil {
		log.Fatalf("Error while creating table programs: %v", err)
	}

	_, err = db.Exec(`
IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'schedules')
BEGIN
    CREATE TABLE schedules (
        id INT IDENTITY(1,1) PRIMARY KEY,
        program_id INT NOT NULL,
        description NTEXT,
        day NVARCHAR(20) NOT NULL,
        date DATETIME NOT NULL,
        FOREIGN KEY (program_id) REFERENCES programs(id)
    )
END
`)
	if err != nil {
		log.Fatalf("Error while creating table schedules: %v", err)
	}
	log.Println("Successfully connected to DB")
	return db
}

func CloseDB() error {
	return db.Close()
}

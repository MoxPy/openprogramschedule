# OpenProgramSchedule
This is a Go-based application that provides APIs for managing programs and schedules using a SQL Database hosted on Microsoft Azure SQL Server.

This project utilizes only the standard library of Go, with the following exceptions:

    godotenv v1.5.1: A Go package for loading environment variables from a .env file.
    go-mssqldb v1.7.2: Microsoft SQL Server driver for Go.

These dependencies ensure efficient handling of environment variables and database interactions with Microsoft SQL Server.

## Configuration

Configure the application using environment variables. Create a .env file in the root directory with the following variables:

    DB_USER=myuser
    DB_PASSWORD=mypassword
    DB_HOST=myserver.database.windows.net
    DB_NAME=mydatabase
    DB_PORT=1433  # Default port for SQL Server
    PRIVATE_KEY=your_private_key
    PUBLIC_API_KEY=your_public_api_key

## Features

- Add, update, retrieve, and delete programs
- Add, update, retrieve, and delete schedules
- Query programs and schedules based on various filters

## APIs

### Program APIs

API Endpoint: http://127.0.0.1:8080

- `POST /programs/add`: Add a new program
- `GET /programs/get-by-id?id={id}`: Retrieve a program by its ID
- `GET /programs/get-by-name?name={name}`: Retrieve a program by its name
- `GET /programs/get-by-category?category={category}`: Retrieve programs by category
- `GET /programs/all`: Retrieve all programs
- `PUT /programs/update?id={id}`: Update a program by its ID
- `DELETE /programs/delete-by-id?id={id}`: Delete a program by its ID

### Schedule APIs

- `POST /schedules/add`: Add a new schedule
- `GET /schedules/all`: Retrieve all schedules
- `GET /schedules/get-by-id?id={id}`: Retrieve a schedule by its ID
- `GET /schedules/get-by-program-id?programId={programId}`: Retrieve schedules by program ID
- `GET /schedules/get-by-day?day={day}`: Retrieve schedules by day
- `GET /schedules/get-by-date?date={date}`: Retrieve schedules by date
- `PUT /schedules/update?id={id}`: Update a schedule by its ID
- `DELETE /schedules/delete-by-id?id={id}`: Delete a schedule by its ID
- `DELETE /schedules/delete-all`: Delete all schedules

## Middleware

The application includes an authentication middleware to protect endpoints. The middleware checks the Authorization header for a valid token.
Public URLs

The following URLs are accessible with a public API key:

    /programs/all
    /programs/get-by-id
    /programs/get-by-name
    /programs/get-by-category
    /schedules/all
    /schedules/get-by-program-id
    /schedules/get-by-id
    /schedules/get-by-day
    /schedules/get-by-date

All other endpoints require a private API key.

## License
This project is licensed under the Mozilla Public License 2.0. For more details, refer to the LICENSE file in the repository.

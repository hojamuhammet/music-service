# Pseudo Music Service API

Repository contains the source code for the Music Service API, that allows you to manage a music library by adding, updating, retrieving, and deleting songs. (Simple CRUD with some filtering logic)

## Table of Content
- Installation
- Dependencies
- Environmental Configurations
- Project Setup
- Running the Application
- Swagger Documentation
  
## Installation

1. Clone the repository:

```bash
git clone https://github.com/hojamuhammet/music-service.git
cd music-service
```

2. Download and install Go dependencies:

```bash
go mod tidy
```

## Dependencies
This project relies on several external dependencies, including Goose for database migrations and Swaggo for Swagger documentation. Make sure you have these tools installed.

### Installing Goose for Migrations

Goose is used to manage database migrations. To install Goose, run:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Add Goose to your system's $PATH by adding the following line to your .bashrc or .zshrc:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Installing Swaggo for Swagger Documentation

Swaggo is used to generate Swagger API documentation. To install it, run:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Similarly, make sure Swaggo is in your $PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Environmental Configurations

1. Configure your environment variables in the .env file. Copy .env.example to .env and modify the values as necessary:

```bash
cp .env.example .env
```

You will need to provide values for the following:
- DB_HOST: Database host (e.g., localhost).
- DB_PORT: Port where the database is running (default is usually 5432).
- DB_USER: Username to connect to the database.
- DB_PASSWORD: Password for the database user.
- DB_NAME: The database name.
- HTTP_PORT: The port where the API will be served.
- LOG_LEVEL: The logging level (e.g., debug, info).

### Example .env file:
```makefile
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=music
HTTP_PORT=8080
LOG_LEVEL=debug
```

## Project Setup
After configuring the environment file, follow these steps to set up the project:

1. Goose Migrations: Migrations will be applied automatically during the service's startup. You do not need to run goose up manually.
  
3. Swagger Documentation: Generate the Swagger documentation:

```bash
swag init --parseDependency --parseInternal -g cmd/main.go
```

3. Build the Application:
To build the application, run:

```bash
go build -o music-service cmd/main.go
```

## Running the Application

To start the application, after building it, run:

```bash
./music-service
```

### Swagger Documentation
This project uses Swaggo to generate and serve Swagger documentation. Once the application is running, access the API documentation by navigating to:

```bash
http://localhost:8080/swagger/index.html
```


# GoFiber API Base Structure

This project serves as a robust base structure for building RESTful APIs using the **Fiber** web framework in Go, integrated with **GORM** for database interactions. It emphasizes a clean, modular architecture, making it easy to extend and maintain.




  * **GoFiber Framework**: Fast, unopinionated, and low-allocation HTTP framework for Go.
  * **GORM ORM**: Elegant ORM for Go, with support for PostgreSQL (and easily extendable to other databases).
  * **Structured Project Layout**: Organized into logical directories for clear separation of concerns.
  * **Environment Configuration**: Manages application settings using environment variables.
  * **Database Migrations**: Simple setup for managing database schema changes.
  * **Request/Response Helpers**: Standardized methods for handling API responses.
  * **Centralized Error Handling**: Custom middleware for consistent error responses.
  * **Logging**: Basic logging setup for application monitoring.
  * **Docker Support**: `Dockerfile` and `docker-compose.yml` for easy containerization and local development setup.
  * **Swagger/OpenAPI Integration**: (Planned/Placeholder) for API documentation.
  * **Authentication & Authorization (Placeholder)**: Structure for implementing JWT-based authentication.
  * **Generic CRUD Operations**: Reusable helper functions for common database operations (`FindAll`, `FindByID`).



## Project Structure Overview

```
.
├── Dockerfile                  # Defines the Docker image for the application
├── LICENSE                     # Project license (e.g., MIT)
├── Makefile                    # Utility for common development tasks (e.g., build, run, test)
├── README.md                   # This file
├── config                      # Application configuration (environment variables)
│   ├── config.go               # Loads environment variables
│   └── swagger.go              # Swagger configuration (placeholder)
├── controllers                 # API logic handlers (where business logic resides)
├── database                    # Database related files
│   ├── base.go                 # Database connection setup (GORM DB instance)
│   ├── init                    # Database initialization (e.g., auto-migrate)
│   ├── models                  # GORM models (structs representing database tables)
│   │   ├── init.go             # Model initialization
│   │   └── users.go            # Example User model
│   └── views                   # Database views (if any)
│       └── init.go             # View initialization
├── docker-compose.yml          # Defines multi-container Docker applications (app + postgres)
├── docs                        # API documentation (Swagger/OpenAPI spec)
│   ├── docs.go                 # Generated Swagger docs
│   ├── swagger.json            # Generated Swagger JSON
│   └── swagger.yaml            # Generated Swagger YAML
├── go.mod                      # Go modules file (dependencies)
├── go.sum                      # Go modules checksums
├── logs                        # Directory for application logs
├── main.go                     # Main entry point of the application
├── middleware                  # Fiber middleware
│   ├── errorhandler.go         # Custom error handling middleware
│   └── logger.go               # Logging middleware (e.g., FiberAccessLogger)
├── pkg                         # Reusable packages/modules
│   ├── auth                    # Authentication related logic (placeholder)
│   └── jwt                     # JWT specific logic (placeholder)
├── routes                      # API route definitions
│   ├── api                     # API versioning or logical grouping
│   │   └── user.go             # User-specific API routes
│   └── routes.go               # Central route registration
├── scripts                     # Utility scripts
│   └── migration.go            # Database migration script (e.g., for auto-migrating models)
├── tmp                         # Temporary files (e.g., build logs)
│   └── build-errors.log        # Log for build errors
└── utils                       # Utility functions
    ├── helper                  # General helper functions
    │   └── get.go              # Generic `FindAll`, `FindByID` functions
    └── response                # Standardized API response handlers
        ├── response.go         # `OK`, `Fail` response functions
        └── swagger.go          # Swagger response definitions (placeholder)
```



## Getting Started

### Prerequisites

  * Go (version 1.20 or higher recommended)
  * Docker & Docker Compose (for local development with PostgreSQL)

### Setup and Run (Local)

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/phonsing-Hub/GoLang.git # (Replace with actual repo URL)
    cd GoLang
    ```

2.  **Configure Environment Variables**:
    Create a `.env` file in the root directory based on `config/config.go` or `config/config.example.env` (if provided).
    Example `.env` file:

    ```dotenv
    APP_PORT=8080
    DB_URL="host=localhost user=youruser password=yourpassword dbname=yourdb port=5432 sslmode=disable TimeZone=Asia/Bangkok"
    # Other environment variables like JWT_SECRET, etc.
    ```

    *Make sure your `DB_URL` points to a running PostgreSQL instance.*

3.  **Run with Docker Compose (Recommended for Local Dev)**:
    This will start both the PostgreSQL database and your Go application.

    ```bash
    docker-compose up --build
    ```

    This command will:

      * Build the Docker image for your Go application.
      * Start a PostgreSQL container.
      * Start your Go application container, connected to the PostgreSQL container.

4.  **Run Natively (without Docker Compose)**:

      * **Start PostgreSQL Database**: Ensure you have a PostgreSQL server running locally and configured as per your `DB_URL`.
      * **Install Dependencies**:
        ```bash
        go mod tidy
        ```
      * **Run Migrations**:
        If you have a migration script to create tables:
        ```bash
        go run scripts/migration.go
        ```
        (You might need to adjust this command based on your `migration.go` content.)
      * **Run the Application**:
        ```bash
        go run main.go
        ```

## API Endpoints (Examples)

Once the server is running (e.g., on `http://localhost:8080`), you can test the following endpoints:

### User Management (`/users`)

  * **`GET /users`**: Get a list of all users with pagination, sorting, and filtering.
      * **Example**: `GET /users?page=1&limit=10&sort_by=name&sort_order=desc&search[name]=john&status=active`
  * **`GET /users/:id`**: Get a single user by ID.
      * **Example**: `GET /users/123`



## Contributing

Feel free to fork this repository and adapt it to your needs. Contributions are welcome\!



## License

This project is licensed under the [MIT License](https://www.google.com/search?q=LICENSE).


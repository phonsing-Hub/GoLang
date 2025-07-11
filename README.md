# GoFiber API Base Structure

This project serves as a robust and scalable base structure for building RESTful APIs using the **Fiber** web framework in Go. It emphasizes a clean, layered architecture with a clear separation of concerns, making it easy to extend and maintain.

This structure follows modern Go best practices, including the use of an `internal` directory to ensure proper encapsulation of application-specific logic.

## Key Features

  * **Fast & Efficient Framework**: Built on [GoFiber](https://gofiber.io/), a high-performance HTTP framework.
  * **Powerful ORM**: Integrated with [GORM](https://gorm.io/) for elegant database interactions with PostgreSQL.
  * **Secure Project Layout**: Utilizes an `internal` directory to hide core application logic and expose only public, reusable packages.
  * **Layered Architecture**: Clear separation of concerns between routing, middleware, business logic, and data access layers.
  * **Ready-to-Use Authentication**: JWT-based authentication middleware and login/profile endpoints are pre-configured.
  * **Configuration Management**: Centralized configuration using environment variables and a `.env` file.
  * **Database Migrations**: Simple script to automatically migrate GORM models.
  * **Standardized Responses**: JSON response helpers for consistent success and error handling.
  * **Dockerized Environment**: Includes `docker-compose.yml` for easy setup of the application and a PostgreSQL database.
  * **Live Reloading**: Configured with `.air.toml` for automatic recompilation during development.

## Project Structure

The project structure is designed to be scalable and maintainable by separating public packages from internal application logic.

```
.
├── internal/                   # Core application logic, not importable by other projects.
│   ├── config/                 # Application configuration (environment variables)
│   ├── database/               # Database connection, models, and schemas
│   ├── middleware/             # Fiber middleware (e.g., JWT Auth, Error Handler)
│   ├── routes/                 # API route definitions and handlers
│   └── utils/                  # Internal helper and response functions
├── pkg/                        # Public, reusable libraries that can be used by other projects.
│   ├── auth/                   # Generic authentication helpers (e.g., bcrypt)
│   └── jwt/                    # Generic JWT generation and validation
├── scripts/                    # Utility scripts (e.g., database migration)
├── .air.toml                   # Configuration for live-reloading (Air)
├── .gitignore                  # Files and directories to ignore in git
├── docker-compose.yml          # Defines services for local development (app + postgres)
├── go.mod                      # Go modules file (dependencies)
├── main.go                     # Main entry point of the application
└── README.md                   # This file
```

## Getting Started

Follow these steps to get the project up and running on your local machine.

### Prerequisites

  * Go (version 1.20 or higher recommended)
  * Docker & Docker Compose

### 1\. Setup

**Clone the repository:**

```bash
git clone https://github.com/phonsing-Hub/GoLang.git
cd GoLang
```

**Create the environment file:**
Create a `.env` file in the root directory and populate it with the necessary variables. You can use the example below as a starting point.

### 2\. Running the Application

#### Option A: With Docker Compose (Recommended)

This is the easiest way to start the entire stack, including the PostgreSQL database.

```bash
docker-compose up --build
```

The API will be available at `http://localhost:3000` (or the port you specify in `.env`).

#### Option B: Natively (Without Docker)

1.  **Start PostgreSQL Database**: Make sure you have a running PostgreSQL instance that matches the `DATABASE_URL` in your `.env` file.
2.  **Install Dependencies**:
    ```bash
    go mod tidy
    ```
3.  **Run Database Migrations**: This script will create the necessary tables in your database.
    ```bash
    go run scripts/migration.go
    ```
4.  **Run the Application**:
    ```bash
    go run main.go
    ```
    Or, for development with live reloading (requires [Air](https://github.com/cosmtrek/air)):
    ```bash
    air
    ```



## API Endpoints (Examples)

Once the server is running (e.g., on `http://localhost:8080`), you can test the following endpoints:


  * **`GET /users`**: Get a list of all users with pagination, sorting, and filtering.
      * **Example**: `GET /users?page=1&limit=10&sort_by=name&sort_order=desc&search[name]=john&status=active`
  * **`GET /users/:id`**: Get a single user by ID.
      * **Example**: `GET /users/123`



## License

This project is licensed under the MIT License.
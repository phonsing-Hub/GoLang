# Go Fiber MVC Pattern API

🚀 RESTful API built with Go Fiber using MVC (Model-View-Controller) architecture pattern with JWT authentication and PostgreSQL database.

## 📋 Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Environment Variables](#environment-variables)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)

## ✨ Features

- 🏗️ **MVC Architecture** - Clean and maintainable code structure
- 🔐 **JWT Authentication** - Secure user authentication and authorization
- 👤 **User Management** - Complete CRUD operations for users
- 🗃️ **PostgreSQL** - Robust database with GORM ORM
- 🔒 **Password Hashing** - Secure password storage with bcrypt
- 📝 **Request Validation** - Input validation using struct tags
- 🚦 **Middleware** - Custom authentication and logging middleware
- 📊 **Error Handling** - Centralized error handling
- 🐳 **Docker Support** - Containerized application
- 📚 **API Documentation** - Comprehensive API documentation

## 📁 Project Structure

```
project-name/
├── controllers/          # Request handlers and business logic
│   ├── auth_controller.go
│   └── user_controller.go
├── models/              # Data models and database schemas
│   └── user.go
├── routes/              # Route definitions and grouping
│   └── routes.go
├── middleware/          # Custom middleware functions
│   └── auth.go
├── database/           # Database connection and configuration
│   ├── connection.go
│   └── migration.go
├── config/             # Application configuration
│   └── config.go
├── utils/              # Utility functions
│   └── validator.go
├── docker/             # Docker configuration files
│   └── Dockerfile
├── migrations/         # Database migration files
├── .env.example        # Environment variables example
├── .gitignore
├── go.mod
├── go.sum
├── main.go             # Application entry point
└── README.md
```

## 🔧 Prerequisites

- **Go** 1.19 or higher
- **PostgreSQL** 12 or higher
- **Git**
- **Docker** (optional)

## 🚀 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/phonsing-Hub/GoLang.git
   cd GoLang
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

4. **Edit environment variables** (see [Environment Variables](#environment-variables))

## 🌍 Environment Variables

Create a `.env` file in the root directory:

```env
# Server Configuration
PORT=3000
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRE_HOURS=72

# CORS Configuration
CORS_ORIGINS=http://localhost:3000,http://localhost:3001
```

## 🗄️ Database Setup

1. **Create PostgreSQL database**
   ```sql
   CREATE DATABASE your_db_name;
   ```

2. **Run migrations** (automatic on first run)
   ```bash
   go run main.go
   ```

   Or manually:
   ```bash
   go run database/migration.go
   ```

## 🏃‍♂️ Running the Application

### Development Mode
```bash
# Run with auto-reload (using air)
go install github.com/cosmtrek/air@latest
air

# Or run normally
go run main.go
```

### Production Mode
```bash
# Build the application
go build -o app main.go

# Run the binary
./app
```

### Using Docker
```bash
# Build and run with Docker Compose
docker-compose up --build

# Or build Docker image manually
docker build -t fiber-mvc-api .
docker run -p 3000:3000 --env-file .env fiber-mvc-api
```

## 📚 API Documentation

### Base URL
```
http://localhost:3000/api/v1
```

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Login User
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### User Management Endpoints

> 🔒 **Note:** All user endpoints require authentication. Include the JWT token in the Authorization header:
> ```
> Authorization: Bearer YOUR_JWT_TOKEN
> ```

#### Get All Users
```http
GET /api/v1/users
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Get User by ID
```http
GET /api/v1/users/{id}
Authorization: Bearer YOUR_JWT_TOKEN
```

#### Create User
```http
POST /api/v1/users
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "password123"
}
```

#### Update User
```http
PUT /api/v1/users/{id}
Authorization: Bearer YOUR_JWT_TOKEN
Content-Type: application/json

{
  "name": "Jane Smith",
  "email": "jane.smith@example.com"
}
```

#### Delete User
```http
DELETE /api/v1/users/{id}
Authorization: Bearer YOUR_JWT_TOKEN
```

### Error Responses

```json
{
  "error": "Error message description"
}
```

**Common HTTP Status Codes:**
- `200` - OK
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### API Testing with curl

```bash
# Register a new user
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Get all users (replace TOKEN with actual JWT token)
curl -X GET http://localhost:3000/api/v1/users \
  -H "Authorization: Bearer TOKEN"
```

## 🐳 Deployment

### Docker Deployment

1. **Build Docker image**
   ```bash
   docker build -t fiber-mvc-api .
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

### Cloud Deployment

#### Heroku
```bash
# Login to Heroku
heroku login

# Create app
heroku create your-app-name

# Set environment variables
heroku config:set JWT_SECRET=your-jwt-secret
heroku config:set DATABASE_URL=your-postgres-url

# Deploy
git push heroku main
```

#### Railway
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login and deploy
railway login
railway init
railway up
```

## 🔒 Security Best Practices

- 🔐 **JWT Secret**: Use a strong, randomly generated JWT secret
- 🛡️ **Password Hashing**: Passwords are hashed using bcrypt
- 🚫 **SQL Injection**: Protected by GORM ORM
- 🌐 **CORS**: Configure CORS for production domains
- 📝 **Input Validation**: All inputs are validated using struct tags
- 🔍 **Error Handling**: Sensitive information is not exposed in error messages

## 📈 Performance Optimization

- ⚡ **Connection Pooling**: Database connection pooling configured
- 🗜️ **Compression**: Gzip compression enabled
- 📊 **Logging**: Structured logging with configurable levels
- 🚀 **Caching**: Ready for Redis integration
- 📏 **Rate Limiting**: Can be easily added with Fiber middleware

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style Guidelines

- Follow Go conventions and best practices
- Use `gofmt` for code formatting
- Write tests for new features
- Update documentation as needed
- Use meaningful commit messages

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙋‍♂️ Support

If you have any questions or need help, please:

1. Check the [Issues](https://github.com/yourusername/fiber-mvc-api/issues) page
2. Create a new issue if your problem isn't already reported
3. Provide detailed information about your environment and the issue

## 🚀 Roadmap

- [ ] Add Redis caching
- [ ] Implement rate limiting
- [ ] Add file upload functionality
- [ ] Create admin dashboard
- [ ] Add email verification
- [ ] Implement role-based access control (RBAC)
- [ ] Add API versioning
- [ ] Create comprehensive test suite
- [ ] Add monitoring and metrics
- [ ] Implement GraphQL support

---

**Happy Coding!** 🎉

Made with ❤️ using Go Fiber
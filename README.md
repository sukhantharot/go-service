# Go Service with Gin

A RESTful API service built with Go and Gin framework, featuring user authentication, JWT middleware, and role-based permissions.

## Features

- User Authentication (Register/Login)
- JWT-based Authorization
- Role-based Access Control
- PostgreSQL Database Integration
- RESTful API Endpoints

## Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Git

## Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go-service.git
cd go-service
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
   - Copy `.env.example` to `.env`
   - Update the database and JWT configuration in `.env`

4. Create PostgreSQL database:
```sql
CREATE DATABASE go_service_db;
```

5. Run the application:
```bash
go run main.go
```

## API Endpoints

### Public Endpoints

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### Protected Endpoints

- `GET /api/users/me` - Get current user info
- `GET /api/admin/users` - Get all users (admin only)
- `POST /api/admin/roles` - Create new role (admin only)
- `POST /api/admin/permissions` - Create new permission (admin only)

## Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_token>
```

## Database Schema

The application uses the following main tables:
- Users
- Roles
- Permissions
- Role_Permissions (junction table)

## License

MIT 
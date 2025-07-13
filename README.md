# Placement Log Backend

A Go-based backend API for the Placement Log application with JWT authentication, user management, and post review system.

## Features

- **JWT Authentication**: Secure token-based authentication for users and admins
- **Post Management**: CRUD operations for placement posts
- **Admin Review System**: Admin approval workflow for posts
- **Role-based Access**: Separate user and admin permissions
- **RESTful API**: Clean REST API design

## Prerequisites

- Go 1.21+
- Git

## Setup

1. **Clone the repository** (if not already done):
```bash
git clone <repository-url>
cd Backend
```

2. **Set up environment variables**:
Create a `.env` file in the backend directory:
```bash
# Required for JWT token generation
SECRET=your-secure-secret-key-here

# Optional: Database configuration
# DB_HOST=localhost
# DB_PORT=5432
# DB_NAME=placementlog
# DB_USER=postgres
# DB_PASSWORD=password
```

**Important**: The `SECRET` environment variable is **required** for JWT token generation. Without it, the server will generate malformed tokens that cause authentication errors.

3. **Install dependencies**:
```bash
go mod download
```

4. **Build and run the server**:
```bash
# Using Makefile
make runServer

# Or manually
go build -o server ./cmd/main.go
./server
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /auth/login` - User login
- `POST /auth/register` - User registration
- `POST /auth/logout` - User logout
- `POST /admin/login` - Admin login
- `POST /admin/logout` - Admin logout

### Posts (User)
- `GET /posts` - Get all approved posts
- `GET /posts/user` - Get user's posts
- `POST /posts` - Create new post
- `PUT /posts` - Update post
- `DELETE /posts` - Delete post

### Posts (Admin)
- `GET /admin/posts` - Get all posts for review
- `PUT /admin/posts/review` - Review post (approve/reject)
- `DELETE /admin/posts` - Delete post as admin

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SECRET` | Yes | JWT signing secret (must be set) |
| `DB_HOST` | No | Database host |
| `DB_PORT` | No | Database port |
| `DB_NAME` | No | Database name |
| `DB_USER` | No | Database user |
| `DB_PASSWORD` | No | Database password |

## Troubleshooting

### "unauthorized: invalid token: token is malformed"

This error occurs when the `SECRET` environment variable is not set or is empty. The JWT library requires a secret key to sign tokens.

**Solution**: Set the `SECRET` environment variable:
```bash
export SECRET=your-secure-secret-key-here
```

Or create a `.env` file:
```bash
echo "SECRET=your-secure-secret-key-here" > .env
```

### Token Validation Issues

If you're getting token validation errors, ensure:
1. The `SECRET` environment variable is set
2. The secret is the same across server restarts
3. The frontend is sending the token in the correct format: `Bearer <token>`

## Development

### Running with Docker

```bash
# Build Docker image
make buildDocker

# Run with environment file
make runDocker
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Security Notes

- **Never commit the actual SECRET value** to version control
- Use a strong, random secret in production
- Consider using environment-specific secrets
- Rotate secrets periodically in production

## License

This project is licensed under the MIT License. 
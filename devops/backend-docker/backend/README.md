# Backend Project

This is a simple backend project with the following features:

- **Server**: A basic HTTP server running on port 8082.
- **API**: Endpoints to fetch and insert user data.
- **SQLite Database**: A lightweight database to store user information.

## Getting Started

### Prerequisites

- Install [Go](https://golang.org/dl/) (version 1.20 or later).

### Setup

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd backend-docker/backend
   ```

2. Initialize Go modules:
   ```bash
   go mod tidy
   ```

3. Run the server:
   ```bash
   go run main.go
   ```

### API Endpoints

#### `GET /api/users`
Fetch all users from the database.

#### `POST /api/users`
Insert a new user into the database.

- **Parameters**:
  - `name` (string): Name of the user.
  - `email` (string): Email of the user.

- **Example**:
  ```bash
  curl -X POST -d "name=John Doe&email=john.doe@example.com" http://localhost:8082/api/users
  ```

### Database

The SQLite database file `data.db` will be created automatically when the server starts.
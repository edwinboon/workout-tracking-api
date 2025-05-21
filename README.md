# workout-tracking-api

A simple REST API for tracking workouts, built with Go.

## Prerequisites

- Go 1.18 or later
- PostgreSQL 12 or later
- Docker (optional, for running PostgreSQL in a container)
- Docker Compose (optional, for running the application in a container)
- Postman or curl (for testing the API)

## curl requests

### Create a new user

```bash
curl -X POST "http://localhost:8080/users" \
     -H "Content-Type: application/json" \
     -d '{
          "username": "johndoe",
          "email": "johndoe@example.com",
          "password": "SecureP@ssword123",
          "bio": "Fitness enthusiast and software developer"
        }'
```

### Authenticate a user

```bash
curl -X POST "http://localhost:8080/auth/token" \
     -H "Content-Type: application/json" \
     -d '{
          "username": "johndoe",
          "password": "SecureP@ssword123"
        }'
```

### Create a new workout

copy and past the token from the previous request and replace it in the Authorization header

```bash
url -X POST "http://localhost:8080/workouts" \
     -H "Authorization: Bearer {token}" \
     -H "Content-Type: application/json" \
     -d '{
          "title": "Morning Cardio",
          "description": "A light 30-minute jog to start the day.",
          "duration_minutes": 30,
          "calories_burned": 300,
          "entries": [
              {
                  "exercise_name": "Jogging",
                  "sets": 1,
                  "duration_seconds": 1800,
                  "weight": 0,
                  "notes": "Maintain a steady pace",
                  "order_index": 1
              }
          ]
        }'
```

### Get a specific workout

replace {id} with the workout ID you want to retrieve

```bash
curl -X GET "http://localhost:8080/workouts/{id}"
```

### Update a workout

copy and past the token from the previous request and replace it in the Authorization header
replace {id} with the workout ID you want to update

```bash
curl -X PUT "http://localhost:8080/workouts/{id}" \
     -H "Authorization: Bearer {token}" \
     -H "Content-Type: application/json" \
     -d '{
          "title": "Updated Cardio",
          "description": "A relaxed 45-minute walk after dinner.",
          "duration_minutes": 45,
          "calories_burned": 250,
          "entries": [
              {
                  "exercise_name": "Walking",
                  "sets": 1,
                  "duration_seconds": 2700,
                  "weight": 0,
                  "notes": "Keep a steady pace",
                  "order_index": 1
              }
          ]
        }'
```

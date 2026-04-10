# Scoreboard API Documentation

This document outlines the available endpoints for retrieving the scoreboards of the CTF and LMS portions of the platform.

## Base URL
`/scoreboard`

---

## 1. Get CTF Scoreboard (Competing)

Retrieves the aggregated scores for users based on the challenges they have solved in the CTF module. It sums up the points from the `challanges` table for each challenge a user has successfully solved (tracked in the `solves` table). Admin users are excluded from the scoreboard.

**Endpoint:**  
`GET /ctf`

**Response Example:**
```json
{
  "status": "success",
  "scoreboard": [
    {
      "user_id": 1,
      "username": "johndoe",
      "name": "John Doe",
      "score": 1500
    },
    {
      "user_id": 4,
      "username": "janedoe",
      "name": "Jane Doe",
      "score": 1200
    }
  ]
}
```

---

## 2. Get LMS Scoreboard (Learning)

Retrieves the aggregated scores for users based on the lesson questions they have answered correctly in the LMS module. It sums up the points from the `questions` table for each question a user has successfully answered (tracked in the `question_solves` table). Admin users are excluded from the scoreboard.

**Endpoint:**  
`GET /lms`

**Response Example:**
```json
{
  "status": "success",
  "scoreboard": [
    {
      "user_id": 2,
      "username": "student_alpha",
      "name": "Alpha Student",
      "score": 300
    },
    {
      "user_id": 1,
      "username": "johndoe",
      "name": "John Doe",
      "score": 250
    }
  ]
}
```

---

### Scoreboard Entry Model

Both endpoints return an array of `ScoreboardEntry` objects in the `scoreboard` field.

| Field      | Type     | Description                                      |
|------------|----------|--------------------------------------------------|
| `user_id`  | Integer  | The unique identifier of the user                |
| `username` | String   | The login username of the user                   |
| `name`     | String   | The full name of the user                        |
| `score`    | Integer  | The total aggregated points for the given module |

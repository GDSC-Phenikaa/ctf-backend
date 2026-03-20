# GDSC CTF Backend API Documentation

**Base URL**: `http://localhost:3333` (by default)  
**Swagger API Docs**: `/swagger/index.html` (Available internally when running in `IsDebug` mode)

*Note: All `[Protected]` routes assume you have authenticated via `/user/login` and that you are sending the session cookie with your requests.*

---

## 1. Authentication & Profile 🧑‍💻

### Login a User
- **URL**: `/user/login`
- **Method**: `POST`
- **Access**: Public
- **Body**:
  ```json
  {
      "email": "user@example.com", 
      "password": "mysecurepassword"
  }
  ```
- **Response**: `200 OK`
  ```json
  {
      "message": "Logged in",
      "token": "jwt_token..."
  }
  ```

### Register a User
- **URL**: `/user/register`
- **Method**: `POST`
- **Access**: Public
- **Body**:
  ```json
  {
      "name": "Jane Doe",
      "username": "janedoe",
      "email": "jane@example.com",
      "password": "SecurePassword1" // Min 8 chars, 1 letter, 1 number
  }
  ```
- **Response**: `201 Created`

### Get Current User Profile
- **URL**: `/profile`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Response**: `200 OK`
  ```json
  {
      "id": 1,
      "name": "Jane Doe",
      "username": "janedoe",
      "email": "jane@example.com"
  }
  ```

---

## 2. CTF Challenges (User Facing) 🚩

### List Challenges
- **URL**: `/user/challenges/list`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Response**: `200 OK`
  ```json
  {
      "challenges": [
          {
              "id": 1,
              "title": "Easy Web Challenge",
              "description": "Find the flag.",
              "difficulty": "Easy",
              "type": "Web",
              "points": 100,
              "solved": false
          }
      ]
  }
  ```

### Submit a Flag
- **URL**: `/user/challenges/submit`
- **Method**: `POST`
- **Access**: `[Protected]`
- **Body**:
  ```json
  {
      "challenge_id": 1,
      "flag": "GDSC{fake_flag}"
  }
  ```
- **Response**: `200 OK` on Correct, `400 Bad Request` on Incorrect.

---

## 3. LMS System (User Facing) 📚

### List Learning Modules
- **URL**: `/user/lms/modules`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Response**: `200 OK` (Array of nested Modules and simplified list of enclosed Lessons)

### Get Lesson Content & Quiz
- **URL**: `/user/lms/lessons/{id}`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Response**: `200 OK` (Lesson title, Markdown/HTML content, and Questions. The backend strips `CorrectAnswer`s dynamically so players cannot cheat!)

### Submit Question Answer
- **URL**: `/user/lms/questions/{id}/submit`
- **Method**: `POST`
- **Access**: `[Protected]`
- **URL Params**: `id` corresponds to the `QuestionID`
- **Body**:
  ```json
  {
      "answer": "Option 3" // or literal text if it's a fill-in-the-blank question
  }
  ```
- **Response**: `200 OK`
  ```json
  {
      "status": "success",
      "correct": true // or false
  }
  ```

### Get Learning Progress
- **URL**: `/user/lms/progress`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Response**: `200 OK` (List of user's submitted question answers and their correctness tracker).

---

## 4. Admin Management Routes 🔒

These routes are protected by `AdminMiddleware` meaning the request session must be attached to a User with `IsAdmin = true`.

### Challenges
- **`[GET] /admin/challenges/list`** - Fetch all challenges
- **`[POST] /admin/challenges/create`** - Create a new puzzle
  ```json
  {
      "Title": "SQL Injection Basics",
      "Description": "Find the flaw in the login form",
      "Difficulty": "Easy",
      "Type": "Web",
      "Points": 100,
      "Flag": "GDSC{sql_inj3ct}",
      "Hidden": false,
      "Docker": false,
      "DockerImage": ""
  }
  ```
  *(Note that fields must be Capitalized depending on the model struct bindings, though case-insensitive forms like "title" usually work)*

- **`[PUT] /admin/challenges/{id}`** - Modify an existing puzzle (send similarly mapped payload as above).

### LMS Content Management
LMS endpoints support full RESTful CRUD operations (GET, POST, PUT, DELETE). 
Base paths:
- `/admin/lms/modules`
- `/admin/lms/lessons`
- `/admin/lms/questions`

- `[GET] /admin/lms/<resource>` - List all instances of `<resource>`
- `[GET] /admin/lms/<resource>/{id}` - Fetch single instance
- `[POST] /admin/lms/<resource>` - Create a new instance (JSON payload)
  - For `questions`, you can specify `"type": "mcq"` or `"type": "text"`. If it's `"text"`, the frontend should render a text box and pass the literal text as the answer!
- `[PUT] /admin/lms/<resource>/{id}` - Update an instance (JSON payload)
- `[DELETE] /admin/lms/<resource>/{id}` - Soft-delete an instance


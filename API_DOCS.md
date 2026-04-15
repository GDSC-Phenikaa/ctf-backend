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

## 3. LMS System v2 (Implementation Spec) 📚

This section defines the **target LMS structure** to implement.

### 3.1 Goals
- Store YouTube embeds as iframe HTML (no extra media service required).
- Support questions in two placements:
  - Lesson-level (applies to whole lesson)
  - Video-segment-level (applies to a specific segment)
- Support all practical question types via one extensible schema.
- Grade answers case-insensitively.
- Migrate all existing LMS data without losing progress.

### 3.2 Data Model (Recommended)

#### `modules`
- `id` (pk)
- `title`
- `description`
- `order`

#### `lessons`
- `id` (pk)
- `module_id` (fk -> modules.id)
- `title`
- `body` (rename from old `content`, still Markdown/HTML)
- `video_iframe` (text, nullable)
- `order`

#### `video_segments`
- `id` (pk)
- `lesson_id` (fk -> lessons.id)
- `title`
- `description` (nullable)
- `start_seconds` (int, >= 0)
- `end_seconds` (int, > start_seconds)
- `order`

#### `questions`
- `id` (pk)
- `lesson_id` (fk -> lessons.id)
- `video_segment_id` (fk -> video_segments.id, nullable)
- `placement` (enum: `lesson`, `segment`)
- `prompt` (text)
- `type` (enum: `single_choice`, `multi_choice`, `true_false`, `short_text`, `long_text`, `numeric`, `code`)
- `options` (jsonb/text JSON, nullable)
- `answer_key` (jsonb/text JSON)
- `points` (int)
- `order`

#### `question_solves`
- keep existing table
- keep relation to `question_id`, `user_id`
- recommended additions:
  - `normalized_answer` (text/json, nullable)
  - `attempt_no` (int, default 1)

### 3.3 Question Schema

#### Single Choice
```json
{
  "type": "single_choice",
  "options": ["A", "B", "C"],
  "answer_key": {"correct": "b"}
}
```

#### Multi Choice
```json
{
  "type": "multi_choice",
  "options": ["A", "B", "C", "D"],
  "answer_key": {"correct": ["a", "c"]}
}
```

#### True / False
```json
{
  "type": "true_false",
  "answer_key": {"correct": true}
}
```

#### Short / Long Text
```json
{
  "type": "short_text",
  "answer_key": {
    "accepted": ["sql injection", "sqli"]
  }
}
```

#### Numeric
```json
{
  "type": "numeric",
  "answer_key": {
    "value": 3.14159,
    "tolerance": 0.01
  }
}
```

#### Code
```json
{
  "type": "code",
  "answer_key": {
    "accepted": ["nmap -sV target", "nmap -sv target"]
  }
}
```

### 3.4 Grading Rules (Case-Insensitive)

Apply normalization before comparison:
- convert to lowercase
- trim leading/trailing spaces
- collapse repeated internal whitespace to a single space

For each type:
- `single_choice`: compare normalized submitted value with normalized correct value
- `multi_choice`: normalize each choice, deduplicate, sort, compare exact set equality
- `true_false`: parse to boolean (`true/false`, `1/0`, `yes/no`) then compare
- `short_text` / `long_text` / `code`: normalize and compare against each accepted answer
- `numeric`: parse float and validate `abs(input - value) <= tolerance`

### 3.5 User-Facing LMS Endpoints (v2)

### List Learning Modules
- **URL**: `/user/lms/modules`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Response**: `200 OK` (modules with lesson summaries)

### Get Lesson Content, Video, Segments, and Questions
- **URL**: `/user/lms/lessons/{id}`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Response**: `200 OK`
  ```json
  {
    "status": "success",
    "lesson": {
      "id": 10,
      "title": "Intro to Web Security",
      "body": "# Markdown body",
      "video_iframe": "<iframe src=\"https://www.youtube.com/embed/VIDEO_ID\" allowfullscreen></iframe>",
      "segments": [
        {
          "id": 1,
          "title": "SQLi Demo",
          "start_seconds": 35,
          "end_seconds": 140
        }
      ],
      "questions": [
        {
          "id": 100,
          "placement": "lesson",
          "video_segment_id": null,
          "type": "single_choice",
          "prompt": "What does SQL stand for?",
          "options": ["Structured Query Language", "Secure Query Layer"],
          "points": 75
        },
        {
          "id": 101,
          "placement": "segment",
          "video_segment_id": 1,
          "type": "short_text",
          "prompt": "Name this attack type",
          "options": null,
          "points": 125
        }
      ]
    }
  }
  ```

### Submit Question Answer
- **URL**: `/user/lms/questions/{id}/submit`
- **Method**: `POST`
- **Access**: `[Protected]`
- **Body**:
  ```json
  {
    "answer": "Structured query language"
  }
  ```
- **Response**: `200 OK`
  ```json
  {
    "status": "success",
    "correct": true,
    "awarded_points": 75,
    "normalized_answer": "structured query language"
  }
  ```

### Get Learning Progress
- **URL**: `/user/lms/progress`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Response**: `200 OK` (all attempts and correctness status, grouped by lesson/question in service layer)

### 3.6 Admin LMS Endpoints (v2)

Keep existing module/lesson/question CRUD and add segment CRUD:
- `/admin/lms/modules`
- `/admin/lms/lessons`
- `/admin/lms/video-segments`
- `/admin/lms/questions`

Recommended extra admin operation:
- `POST /admin/lms/migrations/v2` (idempotent migration runner)

### 3.7 Migration Plan (Migrate All Previous Contents)

Run once in a transactional/idempotent migration:

1. **Schema additions**
   - add `lessons.body`, `lessons.video_iframe`
   - create `video_segments`
   - add `questions.video_segment_id`, `questions.placement`, `questions.prompt`, `questions.answer_key`, `questions.order`

2. **Backfill lesson text**
   - set `lessons.body = lessons.content` for all existing rows

3. **Backfill question fields**
   - `questions.prompt = questions.content`
   - `questions.placement = 'lesson'`
   - `questions.video_segment_id = NULL`
   - convert old `correct_answer` and `options` to `answer_key`/JSON form

4. **Preserve solves**
   - keep all existing `question_solves` rows unchanged
   - old solves remain valid because question IDs are preserved

5. **Compatibility period**
   - for one release, keep reading legacy fields (`content`, `correct_answer`) if new fields are empty
   - then remove legacy fields in a later cleanup migration

### 3.8 Validation Rules
- Reject `video_iframe` if not an `<iframe ...>` with a YouTube embed URL.
- If `placement = segment`, require non-null `video_segment_id` belonging to same lesson.
- If `placement = lesson`, enforce `video_segment_id = null`.
- Enforce `points >= 0`.
- Enforce `start_seconds < end_seconds` for segments.

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
LMS v2 endpoints support full RESTful CRUD operations (GET, POST, PUT, DELETE).
Base paths:
- `/admin/lms/modules`
- `/admin/lms/lessons`
- `/admin/lms/video-segments`
- `/admin/lms/questions`

- `[GET] /admin/lms/<resource>` - List all instances of `<resource>`
- `[GET] /admin/lms/<resource>/{id}` - Fetch single instance
- `[POST] /admin/lms/<resource>` - Create a new instance (JSON payload)
  - For `lessons`, include optional `video_iframe`.
  - For `video-segments`, include `lesson_id`, `start_seconds`, `end_seconds`, and metadata.
  - For `questions`, include `placement` (`lesson`/`segment`) and question `type` (`single_choice`, `multi_choice`, `true_false`, `short_text`, `long_text`, `numeric`, `code`).
- `[PUT] /admin/lms/<resource>/{id}` - Update an instance (JSON payload)
- `[DELETE] /admin/lms/<resource>/{id}` - Soft-delete an instance
- `[POST] /admin/lms/migrations/v2` - Run idempotent migration from legacy LMS schema

---

## 5. Scoreboards 🏆

These endpoints aggregate all correctly solved tasks across the user base and return an ordered array of user scores, automatically filtering out admin accounts.

### Fetch CTF Scoreboard
- **URL**: `/scoreboard/ctf`
- **Method**: `GET`
- **Access**: Public
- **Response**: `200 OK`
  ```json
  {
      "status": "success",
      "scoreboard": [
          {
              "user_id": 2,
              "username": "hacker123",
              "name": "Jane Doe",
              "score": 1400
          }
      ]
  }
  ```

### Fetch LMS Scoreboard
- **URL**: `/scoreboard/lms`
- **Method**: `GET`
- **Access**: Public
- **Response**: `200 OK` (Same schema as above, calculating from the theoretical LMS Question solves instead)

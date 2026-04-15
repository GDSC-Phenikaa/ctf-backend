# Notes API

This document defines the backend contract for per-user notes keyed by page URL.

## Scope

- Notes belong to the authenticated user who created them.
- Each note is associated with one `href` value.
- Frontend can store and render any additional UX behavior as needed.
- All routes below are `[Protected]` and require authentication.

## Data Shape

A note contains exactly these writable fields:

- `href`: string
- `content`: string

Recommended database fields:

- `id`: integer primary key
- `user_id`: authenticated user owner
- `href`: string
- `content`: text
- `created_at`
- `updated_at`

## Endpoints

### List Notes
- **URL**: `/user/notes`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Query Params**:
  - `href` optional, filter notes for a specific page URL
  - `page` optional, integer >= 1, default `1`
  - `limit` optional, integer between `1` and `100`, default `20`
- **Response**: `200 OK`

Example response:
```json
{
  "status": "success",
  "page": 1,
  "limit": 20,
  "total": 42,
  "total_pages": 3,
  "notes": [
    {
      "id": 1,
      "href": "https://example.com/article/1",
      "content": "Remember to check the cookie flags.",
      "created_at": "2026-04-12T10:00:00Z",
      "updated_at": "2026-04-12T10:05:00Z"
    }
  ]
}
```

Pagination behavior:

- The backend should return only the notes for the requested page.
- Results should be ordered by newest first unless you choose a different stable sort, but the sort must be documented and consistent.
- `total` is the total number of matching notes before pagination.
- `total_pages = ceil(total / limit)`.
- If `page` exceeds `total_pages`, return an empty `notes` array with the same metadata.
- If `limit` is missing or invalid, use the default of `20`.

### Get Note
- **URL**: `/user/notes/{id}`
- **Method**: `GET`
- **Access**: `[Protected]`
- **Response**: `200 OK`

Example response:
```json
{
  "status": "success",
  "note": {
    "id": 1,
    "href": "https://example.com/article/1",
    "content": "Remember to check the cookie flags.",
    "created_at": "2026-04-12T10:00:00Z",
    "updated_at": "2026-04-12T10:05:00Z"
  }
}
```

### Create Note
- **URL**: `/user/notes`
- **Method**: `POST`
- **Access**: `[Protected]`
- **Request Body**:
  - `href`: required, string
  - `content`: required, string

Example body:
```json
{
  "href": "https://example.com/article/1",
  "content": "Remember to check the cookie flags."
}
```

- **Response**: `201 Created`

Example response:
```json
{
  "status": "success",
  "note": {
    "id": 1,
    "href": "https://example.com/article/1",
    "content": "Remember to check the cookie flags."
  }
}
```

### Update Note
- **URL**: `/user/notes/{id}`
- **Method**: `PUT`
- **Access**: `[Protected]`
- **Request Body**:
  - `href`: required, string
  - `content`: required, string

Example body:
```json
{
  "href": "https://example.com/article/1",
  "content": "Updated note text."
}
```

- **Response**: `200 OK`

Example response:
```json
{
  "status": "success",
  "note": {
    "id": 1,
    "href": "https://example.com/article/1",
    "content": "Updated note text."
  }
}
```

### Delete Note
- **URL**: `/user/notes/{id}`
- **Method**: `DELETE`
- **Access**: `[Protected]`
- **Response**: `200 OK`

Example response:
```json
{
  "status": "success",
  "message": "Note deleted"
}
```

## Validation Rules

- `href` must be present and must not be empty.
- `content` must be present and must not be empty.
- `href` should be treated as an opaque page key or full URL by the backend.
- Notes should be scoped to the current user only.
- A user should not be able to read, update, or delete another user’s notes.

## Suggested Behavior

- Backend implementation policy: allow multiple notes per page for the same user.
- Notes are returned newest first.
- If you want one note per page later, enforce uniqueness on `(user_id, href)`.

This policy keeps the model flexible for personal page annotations and avoids a uniqueness constraint on `href`.

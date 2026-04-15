# Database ERD

This diagram reflects the current GORM-managed schema used by the backend.
`Container` and `Settings` exist as models in code, but they are not part of the current auto-migration set.

```mermaid
erDiagram
    users ||--o{ notes : owns
    users ||--o{ challanges : authors
    users ||--o{ solves : submits
    challanges ||--o{ solves : receives
    users ||--o| workspaces : allocates

    modules ||--o{ lessons : contains
    lessons ||--o{ video_segments : contains
    lessons ||--o{ questions : contains
    video_segments ||--o{ questions : anchors
    users ||--o{ question_solves : attempts
    questions ||--o{ question_solves : records

    users {
        uint id PK
        string name
        string email
        string username
        string password
        bool is_admin
    }

    notes {
        uint id PK
        uint user_id FK
        string href
        string content
        datetime created_at
        datetime updated_at
    }

    challanges {
        uint id PK
        string title
        string description
        string difficulty
        string type
        int points
        string flag
        string created_at
        string updated_at
        uint author_id FK
        bool docker
        string docker_image
        int solves
        bool hidden
    }

    solves {
        uint id PK
        uint challenge_id FK
        uint user_id FK
        string flag
        bool correct
    }

    workspaces {
        uint id PK
        uint user_id FK
        string container_id
        string status
        string target_url
        datetime created_at
        datetime expires_at
    }

    modules {
        uint id PK
        string title
        string description
        int order
    }

    lessons {
        uint id PK
        uint module_id FK
        string title
        string content
        string body
        string video_iframe
        int order
    }

    video_segments {
        uint id PK
        uint lesson_id FK
        string title
        string description
        int start_seconds
        int end_seconds
        int order
    }

    questions {
        uint id PK
        uint lesson_id FK
        uint video_segment_id FK
        string placement
        string content
        string prompt
        string type
        string options
        string correct_answer
        string answer_key
        int points
        int order
    }

    question_solves {
        uint id PK
        uint question_id FK
        uint user_id FK
        string submitted_answer
        string normalized_answer
        int attempt_no
        bool correct
    }
```

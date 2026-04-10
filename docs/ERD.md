# Database Entity-Relationship Diagram (ERD)

This diagram outlines the complete SQLite/PostgreSQL schema architecture utilized by the GDSC CTF Backend, illustrating both the core hacking mechanics and the LMS structural relationships.

```mermaid
erDiagram
    User ||--o{ Challanges : "creates (Author)"
    User ||--o{ Solves : "submits"
    Challanges ||--o{ Solves : "has"
    User ||--o{ QuestionSolve : "answers"
    
    Module ||--o{ Lesson : "contains"
    Lesson ||--o{ Question : "contains"
    Question ||--o{ QuestionSolve : "has"

    User {
        uint ID PK
        string Name
        string Email
        string Username
        string Password
        bool IsAdmin
    }

    Challanges {
        uint ID PK
        string Title
        string Description
        string Difficulty
        string Type
        int Points
        string Flag
        uint AuthorID FK
        bool Docker
        bool Hidden
        int Solves
    }

    Solves {
        uint ID PK
        uint ChallengeID FK
        uint UserID FK
        string Flag
        bool Correct
    }

    Module {
        uint ID PK
        string Title
        string Description
        int Order
    }

    Lesson {
        uint ID PK
        uint ModuleID FK
        string Title
        string Content
        int Order
    }

    Question {
        uint ID PK
        uint LessonID FK
        string Content
        string Type
        string Options
        string CorrectAnswer
        int Points
    }

    QuestionSolve {
        uint ID PK
        uint QuestionID FK
        uint UserID FK
        string SubmittedAnswer
        bool Correct
    }
```

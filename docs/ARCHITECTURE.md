# Architecture Diagram

This diagram summarizes the current request flow and the main internal subsystems of the backend.

```mermaid
flowchart TD
    Client[Browser / Frontend / API Client]
    Swagger[Swagger UI\n(debug only)]

    Main[main.go\nbootstraps app]
    Chi[chi.Router]
    Middleware[Global middleware\nRequestID, Logger, Recoverer, CORS]
    Auth[AuthMiddleware\nBearer token, query token, or cookie]
    Admin[AdminMiddleware\nchecks users.is_admin]

    subgraph Routes[Route groups]
        UserRoutes[/user\nlogin, register, profile, notes/]
        Profile[/profile]
        AdminRoutes[/admin\nchallenge admin CRUD/]
        LMSAdmin[/admin/lms\nmodules, lessons, segments, questions, migration/]
        Scoreboard[/scoreboard\nCTF and LMS rankings/]
        UserChallenges[/user/challenges\nlist and submit flags/]
        UserLMS[/user/lms\nmodules, lessons, submits, progress/]
        Secret[/secret\nsecret flag endpoint/]
        Certificate[/certificate\nPDF generation/]
        Workspace[/workspace\nstart, stop, status, proxy/]
        Websockify[/websockify*\nroot proxy bridge for VNC/]
    end

    subgraph Core[Internal layers]
        Helpers[helpers\nCORS, response, flags, Docker, PDF helpers]
        Sessions[sessions\nCookie session store + JWT generation]
        Models[models\nGORM entities]
        DB[db.Connect\nSQLite or PostgreSQL via GORM]
        Env[env\nconfig from environment]
    end

    subgraph External[External systems]
        Database[(SQLite / PostgreSQL)]
        Docker[(Docker daemon\nworkspace containers)]
        PDF[(go-pdf/fpdf)]
        JWT[(JWT secret)]
    end

    Client --> Main
    Swagger --> Main
    Main --> Chi
    Main --> DB
    Main --> Sessions
    Main --> Env

    Chi --> Middleware
    Chi --> Routes
    Routes --> Auth
    Routes --> Admin

    UserRoutes --> Models
    UserRoutes --> DB
    Profile --> Models
    AdminRoutes --> Admin
    LMSAdmin --> Admin
    Scoreboard --> DB
    UserChallenges --> DB
    UserLMS --> DB
    Secret --> Env
    Certificate --> PDF
    Workspace --> Helpers
    Workspace --> DB
    Websockify --> Auth

    Middleware --> Sessions
    Auth --> JWT
    Admin --> DB
    DB --> Models
    DB --> Database
    Helpers --> Docker
    Helpers --> PDF
    Sessions --> JWT
    Workspace --> Docker
```
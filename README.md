# MuslimReferralsBE

## Project Overview

This backend application facilitates a referral system, likely connecting job candidates with referrers within companies. It's built using Go and utilizes several packages including GORM for database interaction, Gorilla Mux for routing, Resend for email sending, and Google OAuth2 for authentication.

## Core Components

### 1. Database (`database/`)

*   **Technology:** Uses GORM with a SQLite database (`database.go`). Includes mutexes (`sync.RWMutex`) for managing concurrent access.
*   **Models (`models.go`):** Defines the core data structures:
    *   `User`: Basic user information (name, email, contact details, social links).
    *   `Company`: Represents companies, including their domains and whether they are supported.
    *   `Referrer`: A user associated with a specific company, identified by their corporate email (which needs verification).
    *   `Candidate`: A user seeking referrals, including work experience and resume URL.
    *   `ReferralRequest`: The central object linking a `Candidate` to a `Company` for a specific job/role type, potentially assigned to a `Referrer`. Includes status tracking (Requested, Referred, Accepted, Rejected, Issue).
    *   `EmailVerification`: Tracks email verification requests (code, expiry, status).
*   **Operations:** Each model has associated Go files (e.g., `user.go`, `company.go`) containing CRUD (Create, Read, Update, Delete) functions using the `DbDriver`. Operations often include preloading related data (e.g., `Preload("User")`).

### 2. Service (`service/`)

*   **Purpose:** Encapsulates the core business logic, acting as an intermediary between the API and Database layers.
*   **Structure (`service.go`):**
    *   Defines interfaces (`DatabaseOperations`, `EmailSender`) for dependencies (database driver, email client), enabling dependency injection and testability.
    *   The `Service` struct holds the OAuth configuration, a TTL cache (`userToIdCache`) for mapping token digests to user IDs, and the injected dependencies.
    *   `NewService` constructor initializes the service with these dependencies.
*   **Authentication (`service.go`, `user.go`):**
    *   Uses Google OAuth2.
    *   `GetTokenFromCode`: Exchanges an authorization code (received from the frontend callback) for an OAuth2 token using `oauthConfig.Exchange`.
    *   `queryGoogleForEmail`: Uses the obtained token to fetch user information (email, name, etc.) from Google's userinfo endpoint.
    *   `GetUserIdFromTokenDigest`:
        *   Takes a base64 encoded token digest.
        *   Checks a local TTL cache first for the user ID associated with the token digest.
        *   On cache miss, decodes the token, queries Google for user info.
        *   Checks if the user exists in the database via `dbDriver.GetUserByEmail`.
        *   If the user is new, calls `HandleNewUser`.
        *   Caches the mapping from token digest to user ID.
        *   Returns the user ID and a flag indicating if it was a new user.
    *   `HandleNewUser`: Creates a new `User` record in the database using the information from Google.
*   **Email Verification (`email_verification.go`):**
    *   Manages the process of verifying a user's (specifically a Referrer's) corporate email.
    *   `RequestEmailVerification`:
        *   Performs preconditions checks (rate limiting via `maxActiveVerificationsPerUser`, checks for existing active requests for the email).
        *   Creates a `database.EmailVerification` record with a unique UUID code, TTL (`emailVerificationTTL`), and initial status (`Claimed`).
        *   Uses the injected `EmailSender` (Resend client) to send a verification email containing a unique link (`verificationBaseURL` + code).
        *   Updates the verification record status to `Sent` on success or `SendFailed` on error. Handles cases where the email sender might be disabled (e.g., missing API key).
    *   `VerifyEmail`:
        *   Retrieves the `EmailVerification` record by code.
        *   Checks if the code is valid (exists, not expired, status is `Sent`).
        *   If valid, updates the verification status to `Verified`.
        *   Updates the associated `Referrer`'s `CorporateEmail` field with the verified email.
    *   Uses specific error types (e.g., `ErrVerificationNotFound`, `ErrVerificationExpired`).
*   **Testing (`email_verification_test.go`):** Includes comprehensive unit tests using mocks for the database (`MockDatabaseDriver`) and the email sender (`MockResendEmailsAPI`), demonstrating good testing practices.

### 3. API (`api/`)

*   **Framework:** Uses Gorilla Mux (`mux.Router`) for routing (`endpoints.go`).
*   **Structure (`endpoints.go`):**
    *   `HttpServer` struct holds the router, database driver, and service instances.
    *   `NewHttpServer` initializes the server and sets up routes.
    *   Middleware: Includes CORS (`corsMiddleware`) and request logging (`loggingMiddleware`).
    *   Routes are organized into sub-routers based on user roles/entities (User, Candidate, Referrer) and functionality (Login, Email Verification).
    *   `GetUserIDFromContext`: Helper function to extract the user ID from the request context, likely populated by an authentication middleware (details not fully shown, but it uses the `auth` cookie and `service.GetUserIdFromTokenDigest`).
*   **Key Routes:**
    *   `/login` (`login_routes.go`): Handles the OAuth callback. Receives the `code`, exchanges it for a token via the service, retrieves/creates the user, sets an `auth` cookie containing the base64 encoded token, and redirects the user (to a new user path or default path).
    *   `/api/email-verification` (`email_verification_routes.go`):
        *   `POST /`: Authenticated users (Referrers) request verification for an email address. Calls `service.RequestEmailVerification`.
        *   `GET /verify/{verification_code}`: Handles the link clicked from the verification email. Calls `service.VerifyEmail`. No authentication needed for this endpoint itself, as the code provides the verification context.
    *   User Routes (`user_routes.go`): CRUD operations for User profile, Company (creation/listing), Referrer profile, Candidate profile. Requires authentication.
    *   Candidate Routes (`candidate_routes.go`): CRUD operations for `ReferralRequest` from the candidate's perspective. Requires authentication as a candidate.
    *   Referrer Routes (`referrer_routes.go`): Read operations for `ReferralRequest` relevant to the referrer (e.g., requests for their company). Requires authentication as a referrer.

### 4. API Objects (`api_objects/`)

*   **Purpose:** Defines Data Transfer Objects (DTOs) or "View Models" used specifically for API request/response payloads. This decouples the API structure from the internal database models.
*   **Structure:** Organizes views based on the perspective:
    *   `user_view.go`: Objects for general user profile management (User, Company, Referrer, Candidate details editable by the user).
    *   `candidate_view.go`: Objects tailored for what a candidate sees (e.g., `CandidateViewReferralRequest` includes limited referrer info).
    *   `referrer_view.go`: Objects tailored for what a referrer sees (e.g., `ReferrerViewReferralRequest` includes candidate resume details).
    *   `general_view.go`: Common, simplified views (e.g., `GeneralViewCompany` with just ID and Name).
*   **Conversion:** Contains explicit functions to convert between database models and these API view objects (e.g., `ConvertDbReferralRequestToCandidateViewReferralRequest`, `ConvertUserViewUserToUser`). This ensures only necessary/allowed data is exposed via the API.

## Workflow Summary

1.  **Login:** User initiates Google OAuth flow (frontend). Google redirects to `/login` callback with an authorization `code`. Backend exchanges code for token, fetches/creates user, sets `auth` cookie, redirects frontend.
2.  **Authenticated Requests:** Frontend sends subsequent requests with the `auth` cookie. Backend middleware/handlers use `GetUserIDFromContext` to verify the token (via cache or Google) and get the `userID`.
3.  **Referrer Setup:** A user registers as a referrer for a company (`/api/user/referrer/create`).
4.  **Email Verification:** The referrer needs to verify their corporate email. They trigger a request (`POST /api/email-verification`) providing the email. Backend sends a verification link via Resend. User clicks the link (`GET /api/email-verification/verify/{code}`), backend verifies the code and updates the referrer's record.
5.  **Candidate Setup:** A user registers as a candidate (`/api/user/candidate/create`).
6.  **Referral Request:** Candidate creates a referral request for a specific company (`/api/candidate/referral_request/create`).
7.  **Referrer View:** Referrer views pending requests for their company (`/api/referrer/referral_requests/all` or `/api/referrer/referral_requests/company/{id}`).
8.  **(Future/Implied):** Referrer acts on the request (accepts/rejects/submits), updating the `ReferralRequest` status.

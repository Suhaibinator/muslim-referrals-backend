# Referral Management API Documentation

## Introduction

The Referral Management API provides endpoints to manage referrals for MuslimReferrals. The API is versioned and requires authentication tokens to be passed in the headers for security. Below is the detailed documentation of the available endpoints, their required parameters, and the expected responses.

## Base URL

```
https://api.muslimreferrals.com/v1
```

## Authentication

All API requests must include an authentication token in the header. The token should be passed as follows:

```
Authorization: Bearer YOUR_AUTH_TOKEN
```

## Endpoints

### 1. User Profile

#### 1.1 Create User Profile

**Endpoint:** `POST /users`

**Description:** Create a new user profile.

**Request Body:**

```json
{
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890", // Optional
    "work_experience": "5 years", // Optional, enum type: University Student, New Grad, 1-3 years, 3-5 years, 5-10 years, 10-15 years, 15+ years
    "current_company": "Company A", // Optional
    "current_position": "Software Engineer", // Optional
    "current_title": "Senior Developer", // Optional
    "linkedin_link": "https://linkedin.com/in/johndoe", // Optional
    "website": "http://johndoe.com" // Optional
}
```

**Response:**

```json
{
    "id": "user_123",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "work_experience": "5 years",
    "current_company": "Company A",
    "current_position": "Software Engineer",
    "current_title": "Senior Developer",
    "linkedin_link": "https://linkedin.com/in/johndoe",
    "website": "http://johndoe.com"
}
```

#### 1.2 Get User Profile

**Endpoint:** `GET /users/{user_id}`

**Description:** Retrieve a user profile by user ID.

**Path Parameters:**

- `user_id` (string): The ID of the user to retrieve.

**Response:**

```json
{
    "id": "user_123",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890", // Optional
    "work_experience": "5 years", // Optional
    "current_company": "Company A", // Optional
    "current_position": "Software Engineer", // Optional
    "current_title": "Senior Developer", // Optional
    "linkedin_link": "https://linkedin.com/in/johndoe", // Optional
    "website": "http://johndoe.com" // Optional
}
```

#### 1.3 Update User Profile

**Endpoint:** `PUT /users/{user_id}`

**Description:** Update a user profile by user ID.

**Path Parameters:**

- `user_id` (string): The ID of the user to update.

**Request Body:**

```json
{
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890", // Optional
    "work_experience": "5 years", // Optional
    "current_company": "Company A", // Optional
    "current_position": "Software Engineer", // Optional
    "current_title": "Senior Developer", // Optional
    "linkedin_link": "https://linkedin.com/in/johndoe", // Optional
    "website": "http://johndoe.com" // Optional
}
```

**Response:**

```json
{
    "id": "user_123",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "phone_number": "+1234567890",
    "work_experience": "5 years",
    "current_company": "Company A",
    "current_position": "Software Engineer",
    "current_title": "Senior Developer",
    "linkedin_link": "https://linkedin.com/in/johndoe",
    "website": "http://johndoe.com"
}
```

#### 1.4 Delete User Profile

**Endpoint:** `DELETE /users/{user_id}`

**Description:** Delete a user profile by user ID.

**Path Parameters:**

- `user_id` (string): The ID of the user to delete.

**Response:**

```json
{
    "message": "User profile deleted successfully."
}
```

### 2. Referrer View

#### 2.1 Create Referrer

**Endpoint:** `POST /referrers`

**Description:** Create a new referrer profile.

**Request Body:**

```json
{
  "corporate_email": "referrer@example.com",
  "company": "Example Corp"
}
```

**Response:**

```json
{
  "id": "referrer_123", // corresponding user_id
  "corporate_email": "referrer@example.com",
  "company": "Example Corp"
}
```

#### 2.2 Get Referrer

**Endpoint:** `GET /referrers/{user_id}`

**Description:** Retrieve a referrer profile by referrer ID.

**Path Parameters:**

- `user_id` (string): The ID of the referrer to retrieve.

**Response:**

```json
{
  "id": "referrer_123",
  "corporate_email": "referrer@example.com",
  "company": "Example Corp"
}
```

#### 2.3 Update Referrer

**Endpoint:** `PUT /referrers/{referrer_id}`

**Description:** Update a referrer profile by referrer ID.

**Path Parameters:**

- `user_id` (string): The ID of the referrer to update.

**Request Body:**

```json
{
  "corporate_email": "newreferrer@example.com",
  "company": "New Example Corp"
}
```

**Response:**

```json
{
  "id": "referrer_123",
  "corporate_email": "newreferrer@example.com",
  "company": "New Example Corp"
}
```

#### 2.4 Delete Referrer

**Endpoint:** `DELETE /referrers/{user_id}`

**Description:** Delete a referrer profile by referrer ID.

**Path Parameters:**

- `user_id` (string): The ID of the referrer to delete.

**Response:**

```json
{
  "message": "Referrer profile deleted successfully."
}
```

### 3. Requester View

#### 3.1 Create Requester

**Endpoint:** `POST /requesters`

**Description:** Create a new requester profile.

**Request Body:**

```json
{
        "referral_request_ids": ["referral_1234", "referral_5678"],
        "resume_id": "resume_1234"
}
```

**Response:**

```json
{
        "id": "requester_123",
        "referral_request_ids": ["referral_1234", "referral_5678"],
        "resume_id": "resume_1234"
}
```

#### 3.2 Get Requester

**Endpoint:** `GET /requesters/{requester_id}`

**Description:** Retrieve a requester profile by requester ID.

**Path Parameters:**

- `requester_id` (string): The ID of the requester to retrieve.

**Response:**

```json
{
    "id": "requester_123",
    "work_experience": "5 years",
    "referral_request_ids": ["referral_1234", "referral_5678"],
    "resume_id": "resume_1234"
}
```

#### 3.3 Update Requester

**Endpoint:** `PUT /requesters/{requester_id}`

**Description:** Update a requester profile by requester ID.

**Path Parameters:**

- `requester_id` (string): The ID of the requester to update.

**Request Body:**

```json
{
    "work_experience": "6 years",
    "referral_request_ids": ["referral_1234", "referral_5678", "referral_9101"]
}
```

**Response:**

```json
{
    "id": "requester_123",
    "work_experience": "6 years",
    "referral_request_ids": ["referral_1234", "referral_5678", "referral_9101"],
    "resume_id": "resume_1234"
}
```

#### 3.4 Delete Requester

**Endpoint:** `DELETE /requesters/{requester_id}`

**Description:** Delete a requester profile by requester ID.

**Path Parameters:**

- `requester_id` (string): The ID of the requester to delete.

**Response:**

```json
{
    "message": "Requester profile deleted successfully."
}
```

#### 3.5 Upload Resume

**Endpoint:** `POST /resumes`

**Description:** Upload a resume. The resume should be in PDF format and should not exceed 5MB. The resume will be stored in an object storage and a unique resume ID will be generated and returned.

**Request Body:**

- `resume` (file): The PDF file of the resume. The file size should not exceed 5MB.

**Response:**

```json
{
    "message": "Resume uploaded successfully.",
    "resume_id": "resume_1234"
}
```

## Error Handling

All error responses will follow the structure below:

**Response:**

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Detailed error message"
  }
}
```

### 4. Referral Request

#### 4.1 Create Referral Request

**Endpoint:** `POST /referral-requests`

**Description:** Create a new referral request.

**Request Body:**

```json
{
    "referral_type": "Internship", // Internship or Full-time
    "job_title": "Software Engineer",
    "company_id": "company_123",
    "locations": ["San Francisco, CA"], 
    "description": "I am seeking a software engineering internship. I am skilled at Java and Python and I am qualified for this position because...", // Description that requester can give
    "job_links": ["https://example.com/job1", "https://example.com/job2"] // Optional
}
```

**Response:**

```json
{
    "id": "referral_request_123",
    "referral_type": "Internship",
    "job_title": "Software Engineer",
    "company_id": "company_123",
    "locations": ["San Francisco, CA"],
    "description": "I am seeking a software engineering internship. I am skilled at Java and Python and I am qualified for this position because...",
    "job_links": ["https://example.com/job1", "https://example.com/job2"],
    "status": "referral requested", // Enumerable type: referral requested, referred for job, job offer received, job offer accepted, job offer declined
    "requester_id": "requester_123", // Automatically assigned based on requester auth code
    "referrer_id": "" // Assigned when a referrer accepts the request, initially empty
}
```

#### 4.2 Get Referral Request

**Endpoint:** `GET /referral-requests/{referral_request_id}`

**Description:** Retrieve a referral request by referral request ID.

**Path Parameters:**

- `referral_request_id` (string): The ID of the referral request to retrieve.

**Response:**

```json
{
    "id": "referral_request_123",
    "referral_type": "Internship",
    "job_title": "Software Engineer",
    "company_id": "company_123",
    "locations": ["San Francisco, CA"],
    "description": "I am seeking a software engineering internship. I am skilled at Java and Python and I am qualified for this position because...",
    "job_links": ["https://example.com/job1", "https://example.com/job2"],
    "status": "referral requested",
    "requester_id": "requester_123",
    "referrer_id": ""
}
```

#### 4.3 Update Referral Request

**Endpoint:** `PUT /referral-requests/{referral_request_id}`

**Description:** Update a referral request by referral request ID.

**Path Parameters:**

- `referral_request_id` (string): The ID of the referral request to update.

**Request Body:**

```json
{
    "referral_type": "Full-time",
    "job_title": "Senior Software Engineer",
    "company_id": "company_123",
    "locations": ["San Francisco, CA", "New York, NY"],
    "description": "I am seeking a full-time position as a senior software engineer...",
    "job_links": ["https://example.com/job1", "https://example.com/job2", "https://example.com/job3"],
    "status": "referred for job",
    "referrer_id": "referrer_123"
}
```

**Response:**

```json
{
    "id": "referral_request_123",
    "referral_type": "Full-time",
    "job_title": "Senior Software Engineer",
    "company_id": "company_123",
    "locations": ["San Francisco, CA", "New York, NY"],
    "description": "I am seeking a full-time position as a senior software engineer...",
    "job_links": ["https://example.com/job1", "https://example.com/job2", "https://example.com/job3"],
    "status": "referred for job",
    "requester_id": "requester_123",
    "referrer_id": "referrer_123"
}
```

#### 4.4 Delete Referral Request

**Endpoint:** `DELETE /referral-requests/{referral_request_id}`

**Description:** Delete a referral request by referral request ID.

**Path Parameters:**

- `referral_request_id` (string): The ID of the referral request to delete.

**Response:**

```json
{
  "message": "Referral request deleted successfully."
}
```

#### 4.5 Get Referral Requests by Requester

**Endpoint:** `GET /referral-requests/requester/{requester_id}`

**Description:** Retrieve all referral requests associated with a specific requester.

**Path Parameters:**

- `requester_id` (string): The ID of the requester whose referral requests to retrieve.

**Response:**

```json
[
    {
        "id": "referral_request_123",
        "referral_type": "Internship",
        "job_title": "Software Engineer",
        "company_id": "company_123",
        "locations": ["San Francisco, CA"],
        "description": "I am seeking a software engineering internship. I am skilled at Java and Python and I am qualified for this position because...",
        "job_links": ["https://example.com/job1", "https://example.com/job2"],
        "status": "referral requested",
        "requester_id": "requester_123",
        "referrer_id": ""
    },
    {
        "id": "referral_request_456",
        "referral_type": "Full-time",
        "job_title": "Senior Software Engineer",
        "company_id": "company_456",
        "locations": ["New York, NY"],
        "description": "I am seeking a full-time position as a senior software engineer...",
        "job_links": ["https://example.com/job3", "https://example.com/job4"],
        "status": "referred for job",
        "requester_id": "requester_123",
        "referrer_id": "referrer_456"
    }
]
```

#### 4.6 Get Referral Requests by Referrer

**Endpoint:** `GET /referral-requests/referrer/{referrer_id}`

**Description:** Retrieve all referral requests associated with a specific referrer.

**Path Parameters:**

- `referrer_id` (string): The ID of the referrer whose referral requests to retrieve.

**Response:**

```json
[
    {
        "id": "referral_request_789",
        "referral_type": "Internship",
        "job_title": "Software Engineer",
        "company_id": "company_789",
        "locations": ["San Francisco, CA"],
        "description": "I am seeking a software engineering internship. I am skilled at Java and Python and I am qualified for this position because...",
        "job_links": ["https://example.com/job5", "https://example.com/job6"],
        "status": "referral requested",
        "requester_id": "requester_789",
        "referrer_id": "referrer_123"
    },
    {
        "id": "referral_request_012",
        "referral_type": "Full-time",
        "job_title": "Senior Software Engineer",
        "company_id": "company_012",
        "locations": ["New York, NY"],
        "description": "I am seeking a full-time position as a senior software engineer...",
        "job_links": ["https://example.com/job7", "https://example.com/job8"],
        "status": "referred for job",
        "requester_id": "requester_012",
        "referrer_id": "referrer_123"
    }
]
```

#### 4.7 Get Referral Requests by Company

**Endpoint:** `GET /referral-requests/company/{company_id}`

**Description:** Retrieve all referral requests associated with a specific company.

**Path Parameters:**

- `company_id` (string): The ID of the company whose referral requests to retrieve.

**Response:**

```json
[
    {
        "id": "referral_request_345",
        "referral_type": "Internship",
        "job_title": "Software Engineer",
        "company_id": "company_123",
        "locations": ["San Francisco, CA"],
        "description": "I am seeking a software engineering internship. I am skilled at Java and Python and I am qualified for this position because...",
        "job_links": ["https://example.com/job9", "https://example.com/job10"],
        "status": "referral requested",
        "requester_id": "requester_345",
        "referrer_id": ""
    },
    {
        "id": "referral_request_678",
        "referral_type": "Full-time",
        "job_title": "Senior Software Engineer",
        "company_id": "company_123",
        "locations": ["New York, NY"],
        "description": "I am seeking a full-time position as a senior software engineer...",
        "job_links": ["https://example.com/job11", "https://example.com/job12"],
        "status": "referred for job",
        "requester_id": "requester_678",
        "referrer_id": "referrer_678"
    }
]
```

### 5. Company

#### 5.1 Create Company

**Endpoint:** `POST /companies`

**Description:** Create a new company.

**Request Body:**

```json
{
    "id": "company_123",
    "name": "Example Company",
    "corporate_email_domain": ["example.com"]
}
```

**Response:**

```json
{
    "id": "company_123",
    "name": "Example Company",
    "corporate_email_domain": ["example.com"]
}
```

#### 5.2 Get Company

**Endpoint:** `GET /companies/{company_id}`

**Description:** Retrieve a company by company ID.

**Path Parameters:**

- `company_id` (string): The ID of the company to retrieve.

**Response:**

```json
{
    "id": "company_123",
    "name": "Example Company",
    "corporate_email_domain": ["example.com"]
}
```

#### 5.3 Update Company

**Endpoint:** `PUT /companies/{company_id}`

**Description:** Update a company by company ID.

**Path Parameters:**

- `company_id` (string): The ID of the company to update.

**Request Body:**

```json
{
    "name": "Updated Company",
    "corporate_email_domain": ["updated.com"]
}
```

**Response:**

```json
{
    "id": "company_123",
    "name": "Updated Company",
    "corporate_email_domain": ["updated.com"]
}
```

#### 5.4 Delete Company

**Endpoint:** `DELETE /companies/{company_id}`

**Description:** Delete a company by company ID.

**Path Parameters:**

- `company_id` (string): The ID of the company to delete.

**Response:**

```json
{
    "message": "Company successfully deleted."
}
```
### Common Error Codes

- `401 Unauthorized`: Authentication token is missing or invalid.
- `404 Not Found`: The specified resource could not be found.
- `400 Bad Request`: The request is malformed or contains invalid data.
- `500 Internal Server Error`: An unexpected error occurred on the server.

## Conclusion

This API documentation provides all the necessary endpoints for managing user profiles, referrer views, and requester views in the referral management web application. Ensure to include the authentication token in all requests to securely interact with the API.


### Muslim Referrals API Documentation

---

#### **1. User Management**

- **Create User**
  - **Endpoint:** `/user/create`
  - **Method:** POST
  - **Description:** Registers a new user in the system.
  - **Request Body:**
    ```json
    {
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "phoneNumber": "+1234567890",
      "phoneExt": "123",
      "linkedIn": "https://linkedin.com/in/johndoe",
      "github": "https://github.com/johndoe",
      "website": "https://johndoe.com"
    }
    ```
  - **Response:**
    - **Success:** HTTP 200 OK with user details.
    - **Error:** HTTP 400 Bad Request, HTTP 401 Unauthorized, or HTTP 500 Internal Server Error.

- **Get User**
  - **Endpoint:** `/user`
  - **Method:** GET
  - **Description:** Retrieves details of the authenticated user, includes `Id` field
  - **Response:**
    - **Success:** HTTP 200 OK with user details.
    - **Error:** HTTP 401 Unauthorized or HTTP 404 Not Found.

#### **2. Company Management**

- **Create Company**
  - **Endpoint:** `/user/company/create`
  - **Method:** POST
  - **Description:** Allows a user to register a new company.
  - **Request Body:**
    ```json
    {
      "name": "NewCo",
      "domains": ["newco.com"]
    }
    ```
  - **Response:**
    - **Success:** HTTP 200 OK with company details.
    - **Error:** HTTP 400 Bad Request, HTTP 401 Unauthorized, or HTTP 500 Internal Server Error.

- **Get All Companies**
  - **Endpoint:** `/user/company/get/all`
  - **Method:** GET
  - **Description:** Retrieves all companies associated with the authenticated user.
  - **Response:**
    - **Success:** HTTP 200 OK with list of companies.
    - **Error:** HTTP 401 Unauthorized or HTTP 500 Internal Server Error.

- **Get Specific Company**
  - **Endpoint:** `/user/company/get/{company_id}`
  - **Method:** GET
  - **Description:** Fetches details of a specific company by ID.
  - **URL Parameters:** `company_id` (integer)
  - **Response:**
    - **Success:** HTTP 200 OK with company details.
    - **Error:** HTTP 401 Unauthorized, HTTP 404 Not Found, or HTTP 500 Internal Server Error.

#### **3. Referrer Management**

- **Create Referrer**
  - **Endpoint:** `/user/referrer/create`
  - **Method:** POST
  - **Description:** Registers a user as a referrer.
  - **Request Body:**
    ```json
    {
      "userId": 123,
      "companyId": 456,
      "corporateEmail": "john.doe@newco.com"
    }
    ```
  - **Response:**
    - **Success:** HTTP 200 OK with referrer details.
    - **Error:** HTTP 400 Bad Request, HTTP 401 Unauthorized, or HTTP 500 Internal Server Error.

- **Update Referrer Profile**
  - **Endpoint:** `/user/referrer/update`
  - **Method:** PUT
  - **Description:** Updates details of an existing referrer.
  - **Request Body:**
    ```json
    {
      "referrerId": 789,
      "corporateEmail": "new.email@newco.com"
    }
    ```
  - **Response:**
    - **Success:** HTTP 200 OK with updated referrer details.
    - **Error:** HTTP 400 Bad Request, HTTP 401 Unauthorized, or HTTP 500 Internal Server Error.

- **Get Referrer**
  - **Endpoint:** `/user/referrer/get`
  - **Method:** GET
  - **Description:** Retrieves details of the authenticated referrer.
  - **Response:**
    - **Success:** HTTP 200 OK with referrer details.
    - **Error:** HTTP 401 Unauthorized or HTTP 404 Not Found.

- **Delete Referrer**
  - **Endpoint:** `/user/referrer/delete`
  - **Method:** DELETE
  - **Description:** Removes a referrer from the system.
  - **Response:**
    - **Success:** HTTP 204 No Content.
    - **Error:** HTTP 401 Unauthorized or HTTP 500 Internal Server Error.

#### **4. Candidate Management**

- **Create Candidate**
  - **Endpoint:** `/user/candidate/create`
  - **Method:** POST
  - **Description:** Registers a user as a candidate.
  - **Request Body:**
    ```json
    {
      "userId": 123,
      "workExperience": 5,
      "resumeUrl": "https://resumes.com/johndoe.pdf"
    }
    ```
  - **Response:**
    - **Success:** HTTP 200 OK with candidate details.
    - **Error:** HTTP 400 Bad Request, HTTP 401 Unauthorized, or HTTP 500 Internal Server Error.

- **Update Candidate Profile**
  - **Endpoint:** `/user/candidate/update`
  - **Method:** PUT
  - **Description:** Updates details of an existing candidate.
  - **Request Body:**
    ```json
    {
      "candidateId": 123,
      "workExperience": 7,
      "resumeUrl": "https://resumes.com/updated_johndoe.pdf"
    }
    ```
  - **Response:**
    - **Success:** HTTP 200 OK with updated candidate details.
    - **Error:** HTTP 400 Bad Request, HTTP 401 Unauthorized, or HTTP 500 Internal Server Error.

- **Get Candidate**
  - **Endpoint:** `/user/candidate/get`
  - **Method:** GET
  - **Description:** Retrieves details of the authenticated candidate.
  - **Response:**
    - **Success:** HTTP 200 OK with candidate details.
    - **Error:** HTTP 401 Unauthorized or HTTP 404 Not Found.

- **Delete Candidate**
  - **Endpoint:** `/user/candidate/delete`
  - **Method:** DELETE
  - **Description:** Removes a candidate from the system.
  - **Response:**
    - **Success:** HTTP 204 No Content.
    - **Error:** HTTP 401 Unauthorized or HTTP 500 Internal Server Error.

#### **5. Referral Request Management**

- **Get All Referral Requests for Referrer**
	- **Endpoint:** `/referrer/referral_requests/all`
	- **Method:** GET
	- **Description:** Retrieves all referral requests associated with the authenticated referrer.
	- **Response:**
	  - **Success:** HTTP 200 OK with a list of referral requests.
	  - **Error:**
	    - HTTP 401 Unauthorized: Authentication failed or user not authorized.
	    - HTTP 500 Internal Server Error: An unexpected error occurred on the server.
  - **Response Body:**
    ```json
    [
      {
        "id": 101,
        "candidate": {
          "firstName": "Jane",
          "lastName": "Doe",
          "workExperience": 5,
          "resumeUrl": "https://resumes.com/janedoe.pdf"
        },
        "company_id": 303,
        "company": {
          "id": 303,
          "name": "TechCorp",
          "domains": ["techcorp.com"]
        },
        "job_title": "Software Engineer",
        "job_links": [
          "https://techcorp.com/careers/software-engineer"
        ],
        "description": "Looking for a backend engineering role.",
        "locations": ["Remote", "New York, NY"],
        "referral_type": "EmployeeReferral",
        "referrer": {
          "firstName": "John",
          "lastName": "Smith",
          "company": {
            "id": 303,
            "name": "TechCorp",
            "domains": ["techcorp.com"]
          }
        },
        "status": "Pending"
      },
      {
        "id": 102,
        "candidate": {
          "firstName": "Alice",
          "lastName": "Johnson",
          "workExperience": 3,
          "resumeUrl": "https://resumes.com/alicejohnson.pdf"
        },
        "company_id": 304,
        "company": {
          "id": 304,
          "name": "InnovateX",
          "domains": ["innovatex.com"]
        },
        "job_title": "Product Manager",
        "job_links": [
          "https://innovatex.com/careers/product-manager"
        ],
        "description": "Passionate about product development.",
        "locations": ["San Francisco, CA"],
        "referral_type": "EmployeeReferral",
        "referrer": {
          "firstName": "John",
          "lastName": "Smith",
          "company": {
            "id": 303,
            "name": "TechCorp",
            "domains": ["techcorp.com"]
          }
        },
        "status": "Approved"
      }
    ]
    ```

- **Get Referral Requests by Company**

  - **Endpoint:** `/referrer/referral_requests/company/{company_id}`
  - **Method:** `GET`
  - **Description:** Retrieves all referral requests for a specific company associated with the authenticated referrer.
  - **URL Parameters:**
    - `company_id` (integer): The ID of the company.
  - **Response:**
    - **Success:** HTTP 200 OK with a list of referral requests for the specified company.
    - **Error:**
      - **HTTP 401 Unauthorized:** Authentication failed or user not authorized.
      - **HTTP 404 Not Found:** The specified company does not exist or is not associated with the referrer.
      - **HTTP 500 Internal Server Error:** An unexpected error occurred on the server.
  - **Response Body Example:**

    ```json
    [
      {
        "id": 103,
        "candidate": {
          "firstName": "Michael",
          "lastName": "Brown",
          "workExperience": 4,
          "resumeUrl": "https://resumes.com/michaelbrown.pdf"
        },
        "company_id": 303,
        "company": {
          "id": 303,
          "name": "TechCorp",
          "domains": ["techcorp.com"]
        },
        "job_title": "Data Scientist",
        "job_links": [
          "https://techcorp.com/careers/data-scientist"
        ],
        "description": "Experienced in machine learning and data analysis.",
        "locations": ["Boston, MA"],
        "referral_type": "EmployeeReferral",
        "referrer": {
          "firstName": "John",
          "lastName": "Smith",
          "company": {
            "id": 303,
            "name": "TechCorp",
            "domains": ["techcorp.com"]
          }
        },
        "status": "Pending"
      }
    ]

- **Get Specific Referral Request**

  - **Endpoint:** `/referrer/referral_requests/{request_id}`
  - **Method:** `GET`
  - **Description:** Fetches details of a specific referral request by ID for the authenticated referrer.
  - **URL Parameters:**
    - `request_id` (integer): The ID of the referral request.
  - **Response:**
    - **Success:** HTTP 200 OK with referral request details.
    - **Error:**
      - **HTTP 401 Unauthorized:** Authentication failed or user not authorized.
      - **HTTP 404 Not Found:** The specified referral request does not exist or is not associated with the referrer.
      - **HTTP 500 Internal Server Error:** An unexpected error occurred on the server.
  - **Response Body Example:**

    ```json
    {
      "id": 101,
      "candidate": {
        "firstName": "Jane",
        "lastName": "Doe",
        "workExperience": 5,
        "resumeUrl": "https://resumes.com/janedoe.pdf"
      },
      "company_id": 303,
      "company": {
        "id": 303,
        "name": "TechCorp",
        "domains": ["techcorp.com"]
      },
      "job_title": "Software Engineer",
      "job_links": [
        "https://techcorp.com/careers/software-engineer"
      ],
      "description": "Looking for a backend engineering role.",
      "locations": ["Remote", "New York, NY"],
      "referral_type": "EmployeeReferral",
      "referrer": {
        "firstName": "John",
        "lastName": "Smith",
        "company": {
          "id": 303,
          "name": "TechCorp",
          "domains": ["techcorp.com"]
        }
      },
      "status": "Pending"
    }
    ```

### CandidateViewReferralRequest Data Structure

The `CandidateViewReferralRequest` object represents a referral request from the candidate's perspective.

#### Fields:

- `id` (uint64): The ID of the referral request.
- `candidate` (`CandidateViewCandidate`): Information about the candidate.
- `company_id` (uint64): The ID of the company.
- `company` (`GeneralViewCompany`): Basic information about the company.
- `job_title` (string): The primary job title the candidate is seeking.
- `job_links` (array of strings): URLs to specific job postings.
- `description` (string): A summary or description provided by the candidate.
- `locations` (array of strings): Preferred job locations.
- `referral_type` (string): The type of referral (e.g., `"EmployeeReferral"`).
- `referrer` (`CandidateViewCandidate`): Information about the referrer.
- `status` (string): The current status of the referral request (e.g., `"Pending"`, `"Approved"`, `"Rejected"`).

- **Create Referral Request**

  - **Endpoint:** `/candidate/referral_request/create`
  - **Method:** `POST`
  - **Description:** Allows a candidate to create a new referral request.
  - **Request Body:**

    ```json
    {
      "company_id": 303,
      "job_title": "Software Engineer",
      "job_links": [
        "https://techcorp.com/careers/software-engineer"
      ],
      "description": "Looking for a backend engineering role.",
      "locations": ["Remote", "New York, NY"],
      "referral_type": "EmployeeReferral"
    }
    ```

- **Update Referral Request**

  - **Endpoint:** `/candidate/referral_request/update`
  - **Method:** `PUT`
  - **Description:** Allows a candidate to update an existing referral request.
  - **Response:**
    - **Success:** HTTP 200 OK with the updated referral request details.
    - **Error:**
      - **HTTP 400 Bad Request:** Invalid input data.
      - **HTTP 401 Unauthorized:** Authentication failed or user not authorized.
      - **HTTP 404 Not Found:** The referral request does not exist or is not associated with the candidate.
      - **HTTP 500 Internal Server Error:** An unexpected error occurred on the server.

- **Delete Referral Request**

  - **Endpoint:** `/candidate/referral_request/delete/{referral_request_id}`
  - **Method:** `DELETE`
  - **Description:** Allows a candidate to delete an existing referral request.
  - **URL Parameters:**
    - `referral_request_id` (integer): The ID of the referral request to delete.
  - **Response:**
    - **Success:** HTTP 204 No Content.
    - **Error:**
      - **HTTP 401 Unauthorized:** Authentication failed or user not authorized.
      - **HTTP 404 Not Found:** The referral request does not exist or is not associated with the candidate.
      - **HTTP 500 Internal Server Error:** An unexpected error occurred on the server.

- **Get All Referral Requests**

  - **Endpoint:** `/candidate/referral_request/get/all`
  - **Method:** `GET`
  - **Description:** Retrieves all referral requests associated with the authenticated candidate.
  - **Response:**
    - **Success:** HTTP 200 OK with a list of referral requests.
    - **Error:**
      - **HTTP 401 Unauthorized:** Authentication failed or user not authorized.
      - **HTTP 500 Internal Server Error:** An unexpected error occurred on the server.

- **Get Specific Referral Request**

  - **Endpoint:** `/candidate/referral_request/get/{referral_request_id}`
  - **Method:** `GET`
  - **Description:** Fetches details of a specific referral request by ID for the authenticated candidate.
  - **URL Parameters:**
    - `referral_request_id` (integer): The ID of the referral request.
  - **Response:**
    - **Success:** HTTP 200 OK with referral request details.
    - **Error:**
      - **HTTP 401 Unauthorized:** Authentication failed or user not authorized.
      - **HTTP 404 Not Found:** The referral request does not exist or is not associated with the candidate.
      - **HTTP 500 Internal Server Error:** An unexpected error occurred on the server.

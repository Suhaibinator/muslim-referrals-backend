

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

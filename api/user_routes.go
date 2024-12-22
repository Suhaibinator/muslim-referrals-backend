package api

import (
	"encoding/json"
	"fmt"
	"log"
	"muslim-referrals-backend/api_objects"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// UserUpdateUserHandler handles user creation
func (hs *HttpServer) UserUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserUpdateUserHandler")
	userId, err := hs.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	userDbModel := hs.dbDriver.GetUser(userId)
	if userDbModel == nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	requestUser := api_objects.UserViewUser{}
	// Convert the request body to a UserViewUser object
	unmarshalErr := json.NewDecoder(r.Body).Decode(&requestUser)
	if unmarshalErr != nil {
		http.Error(w, unmarshalErr.Error(), http.StatusBadRequest)
		return
	}
	requestUser.Id = userDbModel.Id
	requestUser.Email = userDbModel.Email

	// Convert the UserViewUser object to a User object
	user := api_objects.ConvertUserViewUserToUser(requestUser, time.Now(), time.Now(), nil)
	// Create the user in the database
	userUpdateErr := hs.dbDriver.UpdateUser(&user)
	if userUpdateErr != nil {
		http.Error(w, userUpdateErr.Error(), http.StatusInternalServerError)
		return
	}

	resultUser := api_objects.ConvertUserToUserViewUser(user)

	response, marshalErr := json.Marshal(resultUser)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

// UserGetUserHandler handles fetching the user details
func (hs *HttpServer) UserGetUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Called UserGetUserHandler")

	userID, userIdRetrievalError := hs.GetUserIDFromContext(r)
	if userIdRetrievalError != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Fetch the user details from the database
	user := hs.dbDriver.GetUser(userID)
	if user == nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	userView := api_objects.ConvertUserToUserViewUser(*user)

	response, marshalErr := json.Marshal(userView)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

// UserCreateCompanyHandler handles company creation for a user
func (hs *HttpServer) UserCreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserCreateCompanyHandler")

	var requestCompany api_objects.UserViewCompany
	// Convert the request body to a CompanyView object
	if err := json.NewDecoder(r.Body).Decode(&requestCompany); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	requestCompany.Id = 0

	// Assuming session or context provides user_id
	userID, userIdRetrievalError := hs.GetUserIDFromContext(r)
	if userIdRetrievalError != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Convert the CompanyView object to a Company domain object
	company := api_objects.ConvertUserViewCompanyToCompany(requestCompany, userID, time.Now(), time.Now(), nil)

	// Create the company in the database
	createdCompany, creationErr := hs.dbDriver.CreateCompany(&company)
	if creationErr != nil {
		http.Error(w, creationErr.Error(), http.StatusInternalServerError)
		return
	}

	resultCompany := api_objects.ConvertCompanyToUserViewCompany(*createdCompany)

	response, marshalErr := json.Marshal(resultCompany)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UserGetAllCompaniesHandler handles fetching all companies for a user
func (hs *HttpServer) UserGetAllCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserGetAllCompaniesHandler")

	_, err := hs.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	companies := hs.dbDriver.GetAllCompanies()

	response, marshalErr := json.Marshal(companies)
	if marshalErr != nil {
		http.Error(w, "Error marshaling companies: "+marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UserGetCompanyHandler handles fetching a specific company for a user
func (hs *HttpServer) UserGetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserGetCompanyHandler")

	vars := mux.Vars(r)
	companyIdString := vars["company_id"]
	_, err := hs.GetUserIDFromContext(r)
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	companyID, parseUint64Err := parseUint64FromString(companyIdString)
	if parseUint64Err != nil {
		http.Error(w, "Invalid company ID: "+parseUint64Err.Error(), http.StatusBadRequest)
		return
	}

	company := hs.dbDriver.GetCompanyById(companyID)

	response, err := json.Marshal(company)
	if err != nil {
		http.Error(w, "Error marshaling company data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UserCreateReferrerHandler handles the creation of a referrer
func (hs *HttpServer) UserCreateReferrerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserCreateReferrerHandler")

	// Decode the request body into UserViewReferrer struct
	var requestReferrer api_objects.UserViewReferrer
	if err := json.NewDecoder(r.Body).Decode(&requestReferrer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	requestReferrer.ReferrerId = 0

	// Authenticate user
	userID, authErr := hs.GetUserIDFromContext(r)
	if authErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Convert the UserViewReferrer object to a Referrer domain object
	referrer := api_objects.ConvertUserViewReferrerToReferrer(requestReferrer, userID, time.Now(), time.Now(), nil)

	// Create the referrer in the database
	createdReferrer, createErr := hs.dbDriver.CreateReferrer(&referrer)
	if createErr != nil {
		http.Error(w, createErr.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the domain model back to a view model
	resultReferrer := api_objects.ConvertReferrerToUserViewReferrer(*createdReferrer)

	// Marshal the result to JSON
	response, marshalErr := json.Marshal(resultReferrer)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UserUpdateReferrerHandler handles updating a referrer
func (hs *HttpServer) UserUpdateReferrerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserUpdateReferrerHandler")

	// Decode the request body
	var updateReferrer api_objects.UserViewReferrer
	if err := json.NewDecoder(r.Body).Decode(&updateReferrer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authenticate and retrieve user ID
	userID, authErr := hs.GetUserIDFromContext(r)
	if authErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	updateReferrer.UserId = userID

	referrerDbObject := api_objects.ConvertUserViewReferrerToReferrer(updateReferrer, userID, time.Now(), time.Now(), nil)

	// Update the referrer in the database
	updatedReferrer, updateErr := hs.dbDriver.UpdateReferrer(userID, &referrerDbObject)
	if updateErr != nil {
		http.Error(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	// Convert updated referrer to the view model
	resultReferrer := api_objects.ConvertReferrerToUserViewReferrer(*updatedReferrer)

	// Marshal the updated referrer
	response, marshalErr := json.Marshal(resultReferrer)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	// Set header and write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UserGetReferrerHandler handles fetching referrer details
func (hs *HttpServer) UserGetReferrerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserGetReferrerHandler")

	// Authenticate and retrieve user ID
	userID, authErr := hs.GetUserIDFromContext(r)
	if authErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Fetch the referrer details from the database
	referrer := hs.dbDriver.GetReferrerByUserId(userID)
	if referrer == nil {
		http.Error(w, "Referrer not found", http.StatusInternalServerError)
		return
	}

	// Convert the domain model to the view model
	resultReferrer := api_objects.ConvertReferrerToUserViewReferrer(*referrer)

	// Marshal the referrer
	response, marshalErr := json.Marshal(resultReferrer)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	// Set header and write the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UserDeleteReferrerHandler handles deleting a referrer
func (hs *HttpServer) UserDeleteReferrerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserDeleteReferrerHandler")

	// Authenticate and retrieve user ID
	userID, authErr := hs.GetUserIDFromContext(r)
	if authErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Fetch the referrer details from the database
	referrer := hs.dbDriver.GetReferrerByUserId(userID)
	if referrer == nil {
		http.Error(w, "Referrer not found", http.StatusInternalServerError)
		return
	}

	// Delete the referrer from the database
	deleteErr := hs.dbDriver.DeleteReferrer(userID, referrer)
	if deleteErr != nil {
		http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Referrer deleted successfully"))
}

// UserCreateCandidateHandler handles creating a candidate
func (hs *HttpServer) UserCreateCandidateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserCreateCandidateHandler")

	var requestCandidate api_objects.UserViewCandidate
	if err := json.NewDecoder(r.Body).Decode(&requestCandidate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, authErr := hs.GetUserIDFromContext(r)
	if authErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	candidate := api_objects.ConvertUserViewCandidateToCandidate(requestCandidate, userID, time.Now(), time.Now(), nil)

	createdCandidate, createErr := hs.dbDriver.CreateCandidate(&candidate)
	if createErr != nil {
		http.Error(w, createErr.Error(), http.StatusInternalServerError)
		return
	}

	resultCandidate := api_objects.ConvertCandidateToUserViewCandidate(*createdCandidate)

	response, marshalErr := json.Marshal(resultCandidate)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UserUpdateCandidateHandler handles updating candidate information
func (hs *HttpServer) UserUpdateCandidateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserUpdateCandidateHandler")

	var updateCandidate api_objects.UserViewCandidate
	if err := json.NewDecoder(r.Body).Decode(&updateCandidate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, authErr := hs.GetUserIDFromContext(r)
	if authErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	candidateDbObject := api_objects.ConvertUserViewCandidateToCandidate(updateCandidate, userID, time.Now(), time.Now(), nil)

	updatedCandidate, updateErr := hs.dbDriver.UpdateCandidate(userID, &candidateDbObject)
	if updateErr != nil {
		http.Error(w, updateErr.Error(), http.StatusInternalServerError)
		return
	}

	resultCandidate := api_objects.ConvertCandidateToUserViewCandidate(*updatedCandidate)

	response, marshalErr := json.Marshal(resultCandidate)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UserGetCandidateHandler handles fetching candidate details
func (hs *HttpServer) UserGetCandidateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserGetCandidateHandler")

	userID, authErr := hs.GetUserIDFromContext(r)
	if authErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	candidate := hs.dbDriver.GetCandidateByUserId(userID)
	if candidate == nil {
		http.Error(w, "Candidate not found", http.StatusInternalServerError)
		return
	}

	resultCandidate := api_objects.ConvertCandidateToUserViewCandidate(*candidate)

	response, marshalErr := json.Marshal(resultCandidate)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UserDeleteCandidateHandler handles deleting a candidate
func (hs *HttpServer) UserDeleteCandidateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Called UserDeleteCandidateHandler")

	userID, authErr := hs.GetUserIDFromContext(r)
	if authErr != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	candidate := hs.dbDriver.GetCandidateByUserId(userID)
	if candidate == nil {
		http.Error(w, "Candidate not found", http.StatusInternalServerError)
		return
	}

	deleteErr := hs.dbDriver.DeleteCandidate(userID, candidate)
	if deleteErr != nil {
		http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Indicates successful deletion with no content to return
}

func parseUint64FromString(str string) (uint64, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

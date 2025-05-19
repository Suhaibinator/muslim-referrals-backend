package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Suhaibinator/muslim-referrals-backend/database"
	"github.com/Suhaibinator/muslim-referrals-backend/service"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

// helper to create http server with in-memory db and prepopulated cache

func setupTestServer(userID uint64, token string) *HttpServer {
	db := database.NewDbDriver(":memory:")
	svc := service.NewService(&oauth2.Config{}, db, nil)

	// Seed the service cache directly using the exported helper
	svc.SetUserIDForToken(token, userID)

	return NewHttpServer(svc, db)
}

func TestUserGetCandidateHandler_NotFound(t *testing.T) {
	token := "tok1"
	hs := setupTestServer(1, token)

	req := httptest.NewRequest(http.MethodGet, "/api/user/candidate/get", nil)
	req.AddCookie(&http.Cookie{Name: "auth", Value: token})
	rr := httptest.NewRecorder()

	hs.UserGetCandidateHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d got %d", http.StatusNotFound, rr.Code)
	}
}

func TestUserGetReferrerHandler_NotFound(t *testing.T) {
	token := "tok2"
	hs := setupTestServer(2, token)

	req := httptest.NewRequest(http.MethodGet, "/api/user/referrer/get", nil)
	req.AddCookie(&http.Cookie{Name: "auth", Value: token})
	rr := httptest.NewRecorder()

	hs.UserGetReferrerHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d got %d", http.StatusNotFound, rr.Code)
	}
}

func TestUserGetCompanyHandler_NotFound(t *testing.T) {
	token := "tok3"
	hs := setupTestServer(3, token)

	req := httptest.NewRequest(http.MethodGet, "/api/user/company/get/1", nil)
	req.AddCookie(&http.Cookie{Name: "auth", Value: token})
	rr := httptest.NewRecorder()

	// We call handler directly with mux vars set
	req = mux.SetURLVars(req, map[string]string{"company_id": "1"})
	hs.UserGetCompanyHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d got %d", http.StatusNotFound, rr.Code)
	}
}

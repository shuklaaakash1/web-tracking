package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	reqBody := `{"username": "testuser", "email": "test@example.com"}`

	req, err := http.NewRequest("POST", "/login", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(loginHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Here, you can decode the JWT token and verify its claims
	// For brevity, we skip the JWT token validation in this example
	// You can add more tests for token validation separately
	assert.NotEmpty(t, response.Token)
}

func TestProtectedHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/protected", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Provide a valid JWT token in the Authorization header
	tokenString := "valid-token"
	req.Header.Add("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(protectedHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "user authorized", rr.Body.String())
}

func TestCreateUserHandler(t *testing.T) {

	userPayload := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
		Role:     "user",
	}
	payloadBytes, _ := json.Marshal(userPayload)

	req, err := http.NewRequest("POST", "/create-user", strings.NewReader(string(payloadBytes)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createUserHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "User created successfully")
}

// Add more unit tests for other functions as needed

func TestGenerateToken(t *testing.T) {
	// Create a test user
	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	token, err := GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

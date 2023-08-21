package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetUserCount(t *testing.T) {
	mockResponseData := `{"total": 5}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponseData))
	}))
	defer server.Close()

	oldBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = oldBaseURL }()

	count := getUserCount()
	if count != 5 {
		t.Errorf("Expected user count: 5, but got: %d", count)
	}
}

func TestGenerateHTMLTable(t *testing.T) {
	users := []User{
		{
			ID:        1,
			Email:     "user1@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Avatar:    "avatar1.jpg",
		},
		{
			ID:        2,
			Email:     "user2@example.com",
			FirstName: "Jane",
			LastName:  "Smith",
			Avatar:    "avatar2.jpg",
		},
	}

	filename := "test_table.html"
	err := generateHTMLTable(users, filename)
	if err != nil {
		t.Errorf("Error generating HTML table: %s", err)
	}

	// Check if the file was generated
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		t.Errorf("Generated HTML table file does not exist")
	}
}

func TestFetchUsersFromAPI(t *testing.T) {
	mockResponseData := struct {
		Data []User `json:"data"`
	}{
		Data: []User{
			{
				ID:        1,
				Email:     "user1@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Avatar:    "avatar1.jpg",
			},
			{
				ID:        2,
				Email:     "user2@example.com",
				FirstName: "Jane",
				LastName:  "Smith",
				Avatar:    "avatar2.jpg",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		responseData, _ := json.Marshal(mockResponseData)
		w.Write(responseData)
	}))
	defer server.Close()

	oldBaseURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = oldBaseURL }()

	users, err := fetchUsersFromAPI()
	if err != nil {
		t.Errorf("Error fetching users from API: %s", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, but got %d", len(users))
	}
}

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/raulaguila/go-api/internal/core/dto"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	client := &http.Client{}

	// 1. Login
	loginPayload := map[string]string{
		"login":    "testuser",
		"password": "password",
	}
	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", baseURL+"/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp dto.AuthOutput
	json.NewDecoder(resp.Body).Decode(&authResp)
	resp.Body.Close()

	assert.NotEmpty(t, authResp.AccessToken)
	token := authResp.AccessToken

	// 2. Me (Check Token Validity)
	req, _ = http.NewRequest("GET", baseURL+"/auth", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 3. Logout (Revoke Token)
	req, _ = http.NewRequest("POST", baseURL+"/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// 4. Me Again (Should fail)
	req, _ = http.NewRequest("GET", baseURL+"/auth", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	resp.Body.Close()
}

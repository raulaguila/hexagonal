package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raulaguila/go-api/internal/core/dto"
)

func TestAuthFlow(t *testing.T) {
	client := &http.Client{}

	// 1. Login
	loginPayload := map[string]string{
		"login":    "testuser",
		"password": "password",
	}
	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", baseURL+"/v1/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp dto.AuthOutput
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	assert.NoError(t, err)
	_ = resp.Body.Close()

	assert.NotEmpty(t, authResp.AccessToken)
	token := authResp.AccessToken

	// 2. Me (Check Token Validity)
	req, _ = http.NewRequest("GET", baseURL+"/v1/auth", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	// 3. Logout (Revoke Token)
	req, _ = http.NewRequest("POST", baseURL+"/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	_ = resp.Body.Close()

	// 4. Me Again (Should fail)
	req, _ = http.NewRequest("GET", baseURL+"/v1/auth", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	_ = resp.Body.Close()
}

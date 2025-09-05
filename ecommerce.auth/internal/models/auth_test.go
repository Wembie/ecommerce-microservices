package models

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTokenRequestBinding_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	form := url.Values{}
	form.Add("username", "juan")
	form.Add("password", "secret")

	c, _ := gin.CreateTestContext(nil)
	c.Request = &http.Request{
		Method: "POST",
		Header: map[string][]string{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		Body:          io.NopCloser(strings.NewReader(form.Encode())),
		ContentLength: int64(len(form.Encode())),
	}

	var req TokenRequest
	err := c.ShouldBind(&req)
	assert.NoError(t, err)
	assert.Equal(t, "juan", req.Username)
	assert.Equal(t, "secret", req.Password)
}

func TestTokenRequestBinding_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	form := url.Values{}
	form.Add("username", "juan")

	c, _ := gin.CreateTestContext(nil)
	c.Request = &http.Request{
		Method: "POST",
		Header: map[string][]string{
			"Content-Type": {"application/x-www-form-urlencoded"},
		},
		Body:          io.NopCloser(strings.NewReader(form.Encode())),
		ContentLength: int64(len(form.Encode())),
	}

	var req TokenRequest
	err := c.ShouldBind(&req)
	assert.Error(t, err)
}

func TestTokenResponseJSON(t *testing.T) {
	resp := TokenResponse{
		AccessToken: "abc123",
		TokenType:   "bearer",
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded TokenResponse
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, "abc123", decoded.AccessToken)
	assert.Equal(t, "bearer", decoded.TokenType)
}

package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	middleware "ecommerce.products.manager/internal/api/routes/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimeOutResponse(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	middleware.TimeOutResponse(c)

	assert.Equal(t, http.StatusRequestTimeout, w.Code)
	assert.Equal(t, "timeout", w.Body.String())
}

func TestCorsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.CorsMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "http://example.com")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
}

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.RequestIDMiddleware())
	r.GET("/test", func(c *gin.Context) {
		traceID, exists := c.Get("trace_id")
		if !exists {
			c.String(http.StatusInternalServerError, "trace_id missing")
			return
		}
		c.String(http.StatusOK, traceID.(string))
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Body.String())
	assert.Equal(t, w.Body.String(), w.Header().Get("X-Trace-Id"))
}

func TestTimeoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.TimeoutMiddleware())
	r.GET("/slow", func(c *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		c.String(http.StatusOK, "done")
	})

	req, _ := http.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "done", w.Body.String())
}

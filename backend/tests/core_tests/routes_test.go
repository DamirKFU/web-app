package core_tests

import (
	"app/tests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	tests.TestMain(m)
}

func TestCoreServerInitialization(t *testing.T) {
	server, _ := tests.GetTestServerWithTx(t)

	assert.NotNil(t, server, "server should not be nil")
	assert.Equal(t, "testdb", server.Cfg.DB.Name, "expected test database name")

	req := httptest.NewRequest(http.MethodGet, server.Reverse("healf", nil), nil)
	resp := httptest.NewRecorder()
	server.Eng.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "expected 200 OK from /healf")
	assert.Contains(t, resp.Body.String(), "ok", "response should contain 'ok'")
}

package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/config"
	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/engine"
	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/httpapi"
)

func newTestServer() http.Handler {
	policies := config.DefaultPolicySet()
	service := engine.NewService(policies)
	server := httpapi.NewServer(config.Load(), policies, service)
	return server.Routes()
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	newTestServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestDashboardSummary(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)

	newTestServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestEvaluateBot(t *testing.T) {
	body := `{"path":"/launch-kit","botScore":95,"suspiciousSignals":3,"userAgentClass":"automation"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/evaluate/bot", strings.NewReader(body))

	newTestServer().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

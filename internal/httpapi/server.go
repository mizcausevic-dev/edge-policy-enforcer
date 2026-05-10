package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/config"
	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/domain"
	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/engine"
)

type Server struct {
	config   config.Config
	policies domain.PolicySet
	service  *engine.Service
	bootedAt time.Time
}

func NewServer(cfg config.Config, policies domain.PolicySet, service *engine.Service) *Server {
	return &Server{
		config:   cfg,
		policies: policies,
		service:  service,
		bootedAt: time.Now().UTC(),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/origins", s.handleOrigins)
	mux.HandleFunc("/api/policies", s.handlePolicies)
	mux.HandleFunc("/api/redirect-rules", s.handleRedirectRules)
	mux.HandleFunc("/api/dashboard/summary", s.handleDashboard)
	mux.HandleFunc("/api/evaluate/request", s.handleEvaluateRequest)
	mux.HandleFunc("/api/evaluate/rate-pressure", s.handleRatePressure)
	mux.HandleFunc("/api/evaluate/bot", s.handleEvaluateBot)
	return withJSONHeaders(mux)
}

func (s *Server) handleRoot(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service": "edge-policy-enforcer",
		"status":  "ready",
		"docs": []string{
			"/health",
			"/api/dashboard/summary",
			"/api/evaluate/request",
			"/api/evaluate/rate-pressure",
			"/api/evaluate/bot",
		},
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "ok",
		"service":  "edge-policy-enforcer",
		"bootedAt": s.bootedAt.Format(time.RFC3339),
	})
}

func (s *Server) handleOrigins(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.policies.Origins)
}

func (s *Server) handlePolicies(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.policies.Policies)
}

func (s *Server) handleRedirectRules(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.policies.RedirectRules)
}

func (s *Server) handleDashboard(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.service.DashboardSummary())
}

func (s *Server) handleEvaluateRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var input domain.RequestEvaluationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	writeJSON(w, http.StatusOK, s.service.EvaluateRequest(input))
}

func (s *Server) handleRatePressure(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var input domain.RatePressureInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	writeJSON(w, http.StatusOK, s.service.EvaluateRatePressure(input))
}

func (s *Server) handleEvaluateBot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var input domain.BotEvaluationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	writeJSON(w, http.StatusOK, s.service.EvaluateBot(input))
}

func withJSONHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Service", "edge-policy-enforcer")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

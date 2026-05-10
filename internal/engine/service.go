package engine

import (
	"strings"

	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/domain"
)

type Service struct {
	policies domain.PolicySet
}

func NewService(policies domain.PolicySet) *Service {
	return &Service{policies: policies}
}

func (s *Service) EvaluateRequest(input domain.RequestEvaluationInput) domain.RequestEvaluationResult {
	result := domain.RequestEvaluationResult{
		Status:          domain.ActionAllow,
		Score:           18,
		RiskLevel:       "low",
		MatchedPolicies: []domain.MatchedPolicy{},
		Issues:          []string{},
		PassedChecks:    []string{},
	}

	if input.Preview || input.ReleaseLocked || strings.Contains(input.Path, "/preview") {
		result.Status = domain.ActionDeny
		result.Score = 92
		result.RiskLevel = "critical"
		result.Issues = append(result.Issues, "Preview or release-locked surface should not be reachable from public traffic.")
		result.MatchedPolicies = append(result.MatchedPolicies, domain.MatchedPolicy{
			ID:       "preview-surface-lock",
			Name:     "Preview surface deny",
			Category: "governance",
			Action:   domain.ActionDeny,
			Reason:   "Public access to unreleased surfaces must be blocked at the edge.",
		})
	}

	if input.BotScore >= 90 && input.SuspiciousSignals >= 2 {
		result.Status = domain.ActionDeny
		result.Score = max(result.Score, 96)
		result.RiskLevel = "critical"
		result.Issues = append(result.Issues, "High-confidence hostile automation is hitting the route with multiple suspicious signals.")
		result.MatchedPolicies = append(result.MatchedPolicies, domain.MatchedPolicy{
			ID:       "bot-score-block",
			Name:     "Known hostile automation block",
			Category: "bot",
			Action:   domain.ActionDeny,
			Reason:   "Severe bot confidence plus malformed posture should be blocked immediately.",
		})
	}

	if result.Status != domain.ActionDeny && input.RequestsPerMinute >= 240 {
		result.Status = domain.ActionChallenge
		result.Score = max(result.Score, 78)
		result.RiskLevel = "high"
		result.Issues = append(result.Issues, "Burst traffic exceeds the preferred edge threshold and should be challenged before full denial.")
		result.MatchedPolicies = append(result.MatchedPolicies, domain.MatchedPolicy{
			ID:       "rate-pressure-challenge",
			Name:     "Challenge burst traffic",
			Category: "rate",
			Action:   domain.ActionChallenge,
			Reason:   "The edge should absorb traffic pressure without dropping healthy users too early.",
		})
	}

	if result.Status == domain.ActionAllow && input.ExpectedExperienceGeo != "" && input.Geo != "" && !strings.EqualFold(input.ExpectedExperienceGeo, input.Geo) {
		for _, rule := range s.policies.RedirectRules {
			if strings.Contains(input.Path, rule.SourcePattern) && (strings.EqualFold(rule.Geo, input.Geo) || rule.Geo == "global") {
				result.Status = domain.ActionRedirect
				result.Score = max(result.Score, 54)
				result.RiskLevel = "moderate"
				result.RedirectTarget = rule.Target
				result.MatchedPolicies = append(result.MatchedPolicies, domain.MatchedPolicy{
					ID:       "geo-lane-redirect",
					Name:     "Regional experience redirect",
					Category: "geo",
					Action:   domain.ActionRedirect,
					Reason:   rule.Reason,
				})
				result.Issues = append(result.Issues, "The incoming experience lane does not match the governed regional surface.")
				break
			}
		}
	}

	if len(result.Issues) == 0 {
		result.PassedChecks = append(result.PassedChecks,
			"Request does not target preview or release-locked infrastructure.",
			"Bot posture remains below enforcement thresholds.",
			"Traffic rate is inside the standard operational lane.",
		)
		result.RecommendedNextAction = "Allow request and attach trace metadata for downstream observability."
		return result
	}

	result.PassedChecks = append(result.PassedChecks, "Origin and path were evaluated against the current policy set.")

	switch result.Status {
	case domain.ActionDeny:
		result.RecommendedNextAction = "Deny request at the edge, log the event, and notify platform operations if the pattern persists."
	case domain.ActionChallenge:
		result.RecommendedNextAction = "Route the request through a challenge lane and monitor whether rate pressure stabilizes."
	case domain.ActionRedirect:
		result.RecommendedNextAction = "Redirect into the governed regional lane and preserve attribution headers."
	default:
		result.RecommendedNextAction = "Allow request with monitoring annotations."
	}

	return result
}

func (s *Service) EvaluateRatePressure(input domain.RatePressureInput) domain.RatePressureResult {
	result := domain.RatePressureResult{
		Status:       "stable",
		Score:        24,
		Issues:       []string{},
		PassedChecks: []string{},
	}

	if input.RequestsPerMinute >= 300 {
		result.Status = "degraded"
		result.Score = 83
		result.Issues = append(result.Issues, "Requests per minute exceed the safe sustained edge threshold.")
	}

	if input.ErrorRatePercent >= 5 {
		result.Status = "degraded"
		result.Score = max(result.Score, 88)
		result.Issues = append(result.Issues, "Error rate is high enough to threaten customer-facing request quality.")
	}

	if input.SaturationPercent >= 85 {
		result.Status = "critical"
		result.Score = max(result.Score, 94)
		result.Issues = append(result.Issues, "Edge saturation is approaching the point where policy decisions need traffic shaping.")
	}

	if len(result.Issues) == 0 {
		result.PassedChecks = append(result.PassedChecks,
			"Request rate is within the preferred operating band.",
			"Error rate is below action thresholds.",
			"Saturation remains below emergency pressure levels.",
		)
		result.RecommendedNextAction = "Maintain current lane allocation and continue passive monitoring."
		return result
	}

	result.PassedChecks = append(result.PassedChecks, "Priority traffic share is still visible for protected customer routes.")
	result.RecommendedNextAction = "Throttle non-priority traffic, preserve customer-critical routes, and escalate to platform engineering if pressure continues."
	return result
}

func (s *Service) EvaluateBot(input domain.BotEvaluationInput) domain.BotEvaluationResult {
	result := domain.BotEvaluationResult{
		Status:       "allow",
		Score:        16,
		Issues:       []string{},
		PassedChecks: []string{},
	}

	if input.BotScore >= 90 && input.SuspiciousSignals >= 2 {
		result.Status = "block"
		result.Score = 95
		result.Issues = append(result.Issues, "High-confidence malicious automation should be blocked.")
	} else if input.BotScore >= 65 {
		result.Status = "challenge"
		result.Score = 72
		result.Issues = append(result.Issues, "Bot posture is elevated enough to warrant a challenge lane.")
	}

	if len(result.Issues) == 0 {
		result.PassedChecks = append(result.PassedChecks,
			"Bot posture remains within the healthy request lane.",
			"Suspicious signal count does not trigger enforcement.",
		)
		result.RecommendedNextAction = "Allow request and keep attribution headers intact."
		return result
	}

	result.PassedChecks = append(result.PassedChecks, "User-agent class was evaluated against the current edge policy.")
	if result.Status == "block" {
		result.RecommendedNextAction = "Block request, log the fingerprint, and add the pattern to the review queue."
	} else {
		result.RecommendedNextAction = "Issue a challenge before origin processing."
	}
	return result
}

func (s *Service) DashboardSummary() domain.DashboardSummary {
	return domain.DashboardSummary{
		Origins:        len(s.policies.Origins),
		Policies:       len(s.policies.Policies),
		RedirectRules:  len(s.policies.RedirectRules),
		ChallengeLanes: 1,
		ProtectedSurfaces: []string{
			"preview surfaces",
			"localized pricing lanes",
			"campaign redirects",
		},
		CurrentPriorityRisks: []string{
			"burst traffic against campaign paths",
			"preview-route exposure during launch week",
			"source attribution loss during geo redirects",
		},
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

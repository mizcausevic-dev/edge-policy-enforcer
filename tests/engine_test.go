package tests

import (
	"testing"

	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/config"
	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/domain"
	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/engine"
)

func TestEvaluateRequestDenyPreview(t *testing.T) {
	service := engine.NewService(config.DefaultPolicySet())

	result := service.EvaluateRequest(domain.RequestEvaluationInput{
		OriginID:      "growth-site",
		Path:          "/preview/launch",
		Method:        "GET",
		Geo:           "US",
		Preview:       true,
		BotScore:      12,
		RequestsPerMinute: 12,
	})

	if result.Status != domain.ActionDeny {
		t.Fatalf("expected deny, got %s", result.Status)
	}
}

func TestEvaluateRequestRedirect(t *testing.T) {
	service := engine.NewService(config.DefaultPolicySet())

	result := service.EvaluateRequest(domain.RequestEvaluationInput{
		OriginID:              "growth-site",
		Path:                  "/pricing",
		Method:                "GET",
		Geo:                   "GB",
		ExpectedExperienceGeo: "US",
		BotScore:              18,
		RequestsPerMinute:     20,
	})

	if result.Status != domain.ActionRedirect {
		t.Fatalf("expected redirect, got %s", result.Status)
	}
}

func TestEvaluateRatePressureCritical(t *testing.T) {
	service := engine.NewService(config.DefaultPolicySet())

	result := service.EvaluateRatePressure(domain.RatePressureInput{
		OriginID:          "campaign-site",
		RequestsPerMinute: 320,
		ErrorRatePercent:  7,
		SaturationPercent: 91,
	})

	if result.Status != "critical" {
		t.Fatalf("expected critical, got %s", result.Status)
	}
}

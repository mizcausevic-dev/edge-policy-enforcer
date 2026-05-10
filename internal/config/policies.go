package config

import "github.com/mizcausevic-dev/edge-policy-enforcer/internal/domain"

func DefaultPolicySet() domain.PolicySet {
	return domain.PolicySet{
		Origins: []domain.Origin{
			{
				ID:          "growth-site",
				Name:        "Growth Platform",
				Environment: "production",
				Region:      "global",
			},
			{
				ID:          "docs-site",
				Name:        "Documentation Surface",
				Environment: "production",
				Region:      "global",
			},
			{
				ID:          "campaign-site",
				Name:        "Campaign Surface",
				Environment: "production",
				Region:      "north-america",
			},
		},
		Policies: []domain.Policy{
			{
				ID:          "bot-score-block",
				Name:        "Known hostile automation block",
				Category:    "bot",
				Description: "Block requests with severe bot confidence and malformed request posture.",
				Threshold:   "botScore >= 90 and suspiciousSignals >= 2",
				Action:      domain.ActionDeny,
			},
			{
				ID:          "rate-pressure-challenge",
				Name:        "Challenge burst traffic",
				Category:    "rate",
				Description: "Challenge requests that exceed operational burst thresholds before hard denial.",
				Threshold:   "requestsPerMinute >= 240",
				Action:      domain.ActionChallenge,
			},
			{
				ID:          "geo-lane-redirect",
				Name:        "Regional experience redirect",
				Category:    "geo",
				Description: "Redirect eligible traffic to the preferred regional lane when geo and language are mismatched.",
				Threshold:   "geo mismatch on localized experience",
				Action:      domain.ActionRedirect,
			},
			{
				ID:          "preview-surface-lock",
				Name:        "Preview surface deny",
				Category:    "governance",
				Description: "Deny crawl or public access to preview surfaces and unreleased campaign paths.",
				Threshold:   "preview or releaseLocked path",
				Action:      domain.ActionDeny,
			},
			{
				ID:          "safe-allow",
				Name:        "Normal customer passage",
				Category:    "baseline",
				Description: "Allow healthy human traffic through the edge lane with trace context intact.",
				Threshold:   "all checks clear",
				Action:      domain.ActionAllow,
			},
		},
		RedirectRules: []domain.RedirectRule{
			{
				ID:            "uk-pricing-lane",
				SourcePattern: "/pricing",
				Target:        "/uk/pricing",
				Geo:           "GB",
				Reason:        "Route UK visitors into the localized conversion path.",
			},
			{
				ID:            "campaign-docs-hardening",
				SourcePattern: "/launch-kit",
				Target:        "/resources/launch-kit",
				Geo:           "global",
				Reason:        "Consolidate legacy campaign routes into governed content surfaces.",
			},
		},
	}
}

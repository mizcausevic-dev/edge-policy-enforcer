package domain

type Action string

const (
	ActionAllow     Action = "allow"
	ActionChallenge Action = "challenge"
	ActionRedirect  Action = "redirect"
	ActionDeny      Action = "deny"
)

type Origin struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Region      string `json:"region"`
}

type Policy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Threshold   string `json:"threshold"`
	Action      Action `json:"action"`
}

type RedirectRule struct {
	ID            string `json:"id"`
	SourcePattern string `json:"sourcePattern"`
	Target        string `json:"target"`
	Geo           string `json:"geo"`
	Reason        string `json:"reason"`
}

type PolicySet struct {
	Origins       []Origin       `json:"origins"`
	Policies      []Policy       `json:"policies"`
	RedirectRules []RedirectRule `json:"redirectRules"`
}

type RequestEvaluationInput struct {
	OriginID              string            `json:"originId"`
	Path                  string            `json:"path"`
	Method                string            `json:"method"`
	Geo                   string            `json:"geo"`
	Language              string            `json:"language"`
	UserAgentClass        string            `json:"userAgentClass"`
	BotScore              int               `json:"botScore"`
	RequestsPerMinute     int               `json:"requestsPerMinute"`
	SuspiciousSignals     int               `json:"suspiciousSignals"`
	Preview               bool              `json:"preview"`
	ReleaseLocked         bool              `json:"releaseLocked"`
	Headers               map[string]string `json:"headers,omitempty"`
	ExpectedExperienceGeo string            `json:"expectedExperienceGeo,omitempty"`
}

type MatchedPolicy struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Action   Action `json:"action"`
	Reason   string `json:"reason"`
}

type RequestEvaluationResult struct {
	Status                Action          `json:"status"`
	Score                 int             `json:"score"`
	RiskLevel             string          `json:"riskLevel"`
	MatchedPolicies       []MatchedPolicy `json:"matchedPolicies"`
	Issues                []string        `json:"issues"`
	PassedChecks          []string        `json:"passedChecks"`
	RecommendedNextAction string          `json:"recommendedNextAction"`
	RedirectTarget        string          `json:"redirectTarget,omitempty"`
}

type RatePressureInput struct {
	OriginID             string `json:"originId"`
	RequestsPerMinute    int    `json:"requestsPerMinute"`
	ErrorRatePercent     int    `json:"errorRatePercent"`
	SaturationPercent    int    `json:"saturationPercent"`
	PriorityTrafficShare int    `json:"priorityTrafficShare"`
}

type RatePressureResult struct {
	Status                string   `json:"status"`
	Score                 int      `json:"score"`
	Issues                []string `json:"issues"`
	PassedChecks          []string `json:"passedChecks"`
	RecommendedNextAction string   `json:"recommendedNextAction"`
}

type BotEvaluationInput struct {
	Path              string `json:"path"`
	BotScore          int    `json:"botScore"`
	SuspiciousSignals int    `json:"suspiciousSignals"`
	UserAgentClass    string `json:"userAgentClass"`
}

type BotEvaluationResult struct {
	Status                string   `json:"status"`
	Score                 int      `json:"score"`
	Issues                []string `json:"issues"`
	PassedChecks          []string `json:"passedChecks"`
	RecommendedNextAction string   `json:"recommendedNextAction"`
}

type DashboardSummary struct {
	Origins              int      `json:"origins"`
	Policies             int      `json:"policies"`
	RedirectRules        int      `json:"redirectRules"`
	ChallengeLanes       int      `json:"challengeLanes"`
	ProtectedSurfaces    []string `json:"protectedSurfaces"`
	CurrentPriorityRisks []string `json:"currentPriorityRisks"`
}

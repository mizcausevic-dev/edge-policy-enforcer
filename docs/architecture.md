# Edge Policy Enforcer Architecture

## Service Overview

Edge Policy Enforcer models a request-governance service that sits in front of origin applications and decides whether traffic should be allowed, challenged, redirected, or denied.

It is designed to represent the kind of policy layer platform engineering teams use for:

- regional request routing
- bot handling
- burst-traffic shaping
- preview surface protection
- governed redirect control

## Request Flow

1. A request payload enters the evaluation layer.
2. The engine classifies path, geo, request posture, and traffic pressure.
3. Policy rules are matched against request attributes.
4. The engine returns an operational action:
   - `allow`
   - `challenge`
   - `redirect`
   - `deny`
5. Edge operators consume the result through request, bot, and rate-pressure endpoints plus the dashboard summary.

## Endpoint Map

- `GET /`
- `GET /health`
- `GET /api/origins`
- `GET /api/policies`
- `GET /api/redirect-rules`
- `GET /api/dashboard/summary`
- `POST /api/evaluate/request`
- `POST /api/evaluate/rate-pressure`
- `POST /api/evaluate/bot`

## Policy Themes

### Governance

- preview surface lockout
- release gating for unreleased campaign paths

### Bot Handling

- hostile automation blocking
- challenge lanes for suspicious but non-terminal traffic

### Traffic Pressure

- burst protection
- saturation-aware response shaping

### Experience Routing

- geo-targeted redirect decisions
- attribution-preserving regional lane selection

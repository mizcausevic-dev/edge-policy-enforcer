package main

import (
	"log"
	"net/http"

	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/config"
	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/engine"
	"github.com/mizcausevic-dev/edge-policy-enforcer/internal/httpapi"
)

func main() {
	cfg := config.Load()
	policies := config.DefaultPolicySet()
	service := engine.NewService(policies)

	server := httpapi.NewServer(cfg, policies, service)

	log.Printf("edge-policy-enforcer listening on %s", cfg.Address())
	if err := http.ListenAndServe(cfg.Address(), server.Routes()); err != nil {
		log.Fatal(err)
	}
}

package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouterRoutes(t *testing.T) {
	router, _ := setupRouterWithMemoryService()
	routes := router.Routes()

	expected := []string{
		"GET /loans",
		"GET /loans/:id",
		"POST /loans",
		"POST /loans/:id/approve",
		"POST /loans/:id/invest",
		"POST /loans/:id/disburse",
	}

	for _, route := range expected {
		found := false
		for _, r := range routes {
			if r.Method+" "+r.Path == route {
				found = true
				break
			}
		}
		assert.True(t, found, "Route missing: "+route)
	}
}

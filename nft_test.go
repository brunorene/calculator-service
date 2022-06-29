//go:build nft

package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestEndpoints(t *testing.T) {
	assert.Eventually(t, func() bool {
		response, err := http.Get("http://localhost:8090/add/2/3")
		if err == nil {
			return response.StatusCode == 200
		}

		return false
	}, 10*time.Second, time.Second, "waiting for server")

	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := time.Minute
	targeter := vegeta.NewStaticTargeter(
		vegeta.Target{
			Method: "GET",
			URL:    "http://localhost:8090/add/2/3",
		},
		vegeta.Target{
			Method: "GET",
			URL:    "http://localhost:8090/sub/2/3",
		},
		vegeta.Target{
			Method: "GET",
			URL:    "http://localhost:8090/mul/2/3",
		},
		vegeta.Target{
			Method: "GET",
			URL:    "http://localhost:8090/div/2/3",
		},
	)
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	defer metrics.Close()
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}

	t.Logf("Metrics: %#v", metrics)

	assert.Empty(t, metrics.Errors)
	assert.Equal(t, 60000, metrics.StatusCodes["200"])
}

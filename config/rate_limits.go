package config

import "github.com/pocketbase/pocketbase/core"

func InitRateLimitRule(settings *core.Settings) {
	settings.RateLimits.Enabled = true
	settings.RateLimits.Rules = []core.RateLimitRule{
		{
			Label:       "*:auth",
			MaxRequests: 2,
			Duration:    3,
			Audience:    "",
		},
		{
			Label:       "/api/",
			MaxRequests: 300,
			Duration:    10,
			Audience:    "",
		},
		{
			Label:       "/api/webhook/stripe",
			MaxRequests: 500,
			Duration:    5,
			Audience:    "",
		},
		{
			Label:       "/api/checkout/subscription",
			MaxRequests: 5,
			Duration:    10,
			Audience:    "@auth",
		},
		{
			Label:       "/api/checkout/subscription",
			MaxRequests: 2,
			Duration:    10,
			Audience:    "@guest",
		},
	}
}

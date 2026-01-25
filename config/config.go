package config

import (
	"os"
)

type Config struct {
	StripeKey     string
	StripeSignKey string
	PlanToPrice   map[string]string
}

func New() *Config {
	return &Config{
		StripeKey:     getEnv("STRIPE_KEY", ""),
		StripeSignKey: getEnv("STRIPE_SIGN_KEY", ""),
		// 动态从环境变量读取价格 ID
		PlanToPrice: map[string]string{
			"starter": os.Getenv("STRIPE_PLAN_STARTER"),
			"pro":     os.Getenv("STRIPE_PLAN_PRO"),
			"plus":    os.Getenv("STRIPE_PLAN_PLUS"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

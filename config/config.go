package config

import (
	"os"
)

type Config struct {
	StripeKey     string
	StripeSignKey string
	PlanToPrice   map[string]string

	PriceToPlan map[string]string
}

func New() *Config {
	// 原始映射
	planToPrice := map[string]string{
		"starter": getEnv("STRIPE_PLAN_STARTER", ""),
		"pro":     getEnv("STRIPE_PLAN_PRO", ""),
		"plus":    getEnv("STRIPE_PLAN_PLUS", ""),
	}

	// 自动生成反向映射
	priceToPlan := make(map[string]string)
	for plan, price := range planToPrice {
		if price != "" {
			priceToPlan[price] = plan
		}
	}

	return &Config{
		StripeKey:     getEnv("STRIPE_KEY", ""),
		StripeSignKey: getEnv("STRIPE_SIGN_KEY", ""),
		PlanToPrice:   planToPrice,
		PriceToPlan:   priceToPlan,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if defaultValue == "" {
		panic("environment variable '" + key + "' has not been set, and is required")
	}

	return defaultValue
}

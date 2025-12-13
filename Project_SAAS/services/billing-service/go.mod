module project_saas/services/billing-service

go 1.22

require (
	github.com/go-chi/chi/v5 v5.0.10
	go.uber.org/zap v1.27.0
	project_saas/shared v0.0.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
)

replace project_saas/shared => ../../shared

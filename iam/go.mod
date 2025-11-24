module github.com/max-kriv0s/go-microservices-edu/iam

go 1.24.5

replace github.com/max-kriv0s/go-microservices-edu/shared => ../shared

replace github.com/max-kriv0s/go-microservices-edu/platform => ../platform

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/jackc/pgx/v5 v5.7.6
	github.com/joho/godotenv v1.5.1
	github.com/max-kriv0s/go-microservices-edu/platform v0.0.0-00010101000000-000000000000
	github.com/max-kriv0s/go-microservices-edu/shared v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.76.0
)

require (
	github.com/alexedwards/argon2id v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pressly/goose/v3 v3.26.0 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.44.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250825161204-c5933d9347a5 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

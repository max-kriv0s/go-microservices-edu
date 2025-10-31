module github.com/max-kriv0s/go-microservices-edu/payment

go 1.24.5

replace github.com/max-kriv0s/go-microservices-edu/shared => ../shared

require (
	github.com/brianvoe/gofakeit/v7 v7.8.0
	github.com/google/uuid v1.6.0
	github.com/max-kriv0s/go-microservices-edu/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.76.0
)

require (
	github.com/caarlos0/env/v11 v11.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

up:
	docker compose up -d

down:
	docker compose down

test: clear vet build-server build-agent metrics-test

gen-keys:
	openssl genrsa -out ./private.pem 2048
	openssl rsa -pubout -in ./private.pem -out ./public.pem

clear:
	clear

fmt:
	goimports -local "github.com/psfpro/metrics" -w ./

doc:
	godoc -http=:8080 -play && http://localhost:8080/pkg/github.com/psfpro/metrics/internal/?m=all

vet:
	go vet ./...

staticlint:
	go run ./cmd/staticlint ./...

build-server:
	cd cmd/server && go build -buildvcs=false -o server

build-agent:
	cd cmd/agent && go build -buildvcs=false  -o agent

run-server:
	cd cmd/server && go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$$(git rev-parse HEAD)'" main.go

run-agent:
	cd cmd/agent && go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildDate=$$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$$(git rev-parse HEAD)'" main.go

metrics-test:
	metricstest -test.v -test.run=^TestIteration14$$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -database-dsn='postgres://app:pass@localhost:5432/app?sslmode=disable' \
                -key=123 \
                -file-storage-path=tmp \
                -server-port=8080

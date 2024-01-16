up:
	docker compose up -d

down:
	docker compose down

test: clear vet build-server build-agent metrics-test

clear:
	clear

vet:
	go vet ./...

build-server:
	cd cmd/server && go build -buildvcs=false -o server

build-agent:
	cd cmd/agent && go build -buildvcs=false  -o agent

metrics-test:
	metricstest -test.v -test.run=^TestIteration10*$$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -file-storage-path=tmp \
                -server-port=8888

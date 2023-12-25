test: build-server build-agent metrics-test

build-server:
	cd cmd/server && go build -buildvcs=false -o server

build-agent:
	cd cmd/agent && go build -buildvcs=false  -o agent

metrics-test:
	metricstest -test.v -test.run=^TestIteration6*$$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server \
                -server-port=8888

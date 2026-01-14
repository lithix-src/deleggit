module catalyst/docker-watcher

go 1.24.0

toolchain go1.24.6

require (
	github.com/docker/docker v24.0.7+incompatible
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/prometheus/client_golang v1.19.0
)

require (
	github.com/Microsoft/go-winio v0.4.21 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/docker/distribution v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/go-connections v0.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/morikuni/aec v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gotest.tools/v3 v3.5.2 // indirect
)

// Fix for Docker dependency hell
replace github.com/docker/distribution => github.com/docker/distribution v2.8.2+incompatible

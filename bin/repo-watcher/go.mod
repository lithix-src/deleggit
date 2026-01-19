module catalyst/repo-watcher

go 1.24.6

require (
	github.com/eclipse/paho.mqtt.golang v1.5.1
	github.com/point-unknown/catalyst/pkg v0.0.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
)

replace github.com/point-unknown/catalyst/pkg => ../../pkg

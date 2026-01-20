module github.com/datacraft/catalyst/bin/test-messenger

go 1.24.6

replace (
	github.com/datacraft/catalyst/core => ../../core
	github.com/point-unknown/catalyst/pkg => ../../pkg
)

require github.com/eclipse/paho.mqtt.golang v1.5.1

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
)

package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/datacraft/deleggit/core/internal/domain"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Adapter struct {
	client mqtt.Client
}

// NewAdapter creates a connected MQTT adapter
func NewAdapter(brokerURL, clientID string) (*Adapter, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	opts.SetKeepAlive(60 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Adapter{client: client}, nil
}

// Publish sends a CloudEvent to the specified topic
func (a *Adapter) Publish(topic string, event domain.CloudEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	token := a.client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

// Subscribe listens for messages on a topic and triggers the handler
func (a *Adapter) Subscribe(topic string, handler func(domain.CloudEvent)) error {
	cb := func(client mqtt.Client, msg mqtt.Message) {
		var event domain.CloudEvent
		if err := json.Unmarshal(msg.Payload(), &event); err != nil {
			fmt.Printf("Error unmarshalling event on %s: %v\n", topic, err)
			return
		}
		handler(event)
	}

	token := a.client.Subscribe(topic, 0, cb)
	token.Wait()
	return token.Error()
}

func (a *Adapter) Close() {
	a.client.Disconnect(250)
}

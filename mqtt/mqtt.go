package mqtt

import (
	"errors"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Message contains mqtt message information
type Message struct {
	Topic   string
	Message string
}

type ConnectionOptions struct {
	Broker       string
	ClientID     string
	Username     string
	Password     string
	CleanSession bool
}

type Client struct {
	client mqtt.Client
}

// Connect creates the mqtt connection
func Connect(o ConnectionOptions) (*Client, error) {
	log.Println("Connecting to mqtt")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(o.Broker)
	opts.SetClientID(o.ClientID)
	opts.SetUsername(o.Username)
	opts.SetPassword(o.Password)
	opts.SetCleanSession(o.CleanSession)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("Error connecting to mqtt: %s", token.Error())
	}

	c := &Client{client: client}

	return c, nil
}

// Disconnect from the mqtt broker
func (c *Client) Disconnect() {
	log.Println("Disconnecting from mqtt")
	c.client.Disconnect(250)
}

// GetMessagesOnTopic creates a mqtt subscription and sends messages back on the channel
func (c *Client) GetMessagesOnTopic(topic string, cMsg chan Message, cIn chan string) {
	if c.client.IsConnected() == false {
		log.Println("No mqtt connection created yet")
		return
	}

	cMqtt := make(chan Message)

	token := c.client.Subscribe(topic, byte(0), func(client mqtt.Client, msg mqtt.Message) {
		cMqtt <- Message{msg.Topic(), string(msg.Payload())}
	})

	if token.Wait() && token.Error() != nil {
		log.Printf("Error subscribing to mqtt: %s", token.Error())
		return
	}

	log.Printf("Monitoring message on topic %s", topic)

	for {
		select {
		case msg := <-cMqtt:
			cMsg <- msg

		case command := <-cIn:
			if command == "shutdown" {
				return
			}
		}
	}
}

func (c *Client) SendMessage(topic, message string) error {
	if c.client.IsConnected() == false {
		return errors.New("No mqtt connection created yet")
	}

	log.Printf("Sending message %s to %s\n", message, topic)
	if t := c.client.Publish(topic, 0, false, message); t.WaitTimeout(5*time.Second) && t.Error() != nil {
		return t.Error()
	}

	return nil
}

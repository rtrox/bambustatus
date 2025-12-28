package printer

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTConfig holds the configuration for connecting to the Bambu Lab printer
type MQTTConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Serial   string // Optional - will be auto-discovered if empty
}

// MQTTClient manages the MQTT connection and status updates
type MQTTClient struct {
	config  MQTTConfig
	client  mqtt.Client
	status  *Status
	mu      sync.RWMutex
	onUpdate func(*Status)
}

// NewMQTTClient creates a new MQTT client
func NewMQTTClient(config MQTTConfig) *MQTTClient {
	return &MQTTClient{
		config:  config,
		status:  NewStatus(),
	}
}

// SetUpdateCallback sets a callback function that will be called when status updates
func (m *MQTTClient) SetUpdateCallback(callback func(*Status)) {
	m.onUpdate = callback
}

// GetStatus returns the current status (thread-safe)
func (m *MQTTClient) GetStatus() *Status {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid race conditions
	statusCopy := *m.status
	return &statusCopy
}

// updateStatus updates the internal status (thread-safe)
func (m *MQTTClient) updateStatus(status *Status) {
	m.mu.Lock()
	m.status = status
	m.mu.Unlock()

	// Call update callback if set
	if m.onUpdate != nil {
		m.onUpdate(status)
	}
}

// messageHandler processes incoming MQTT messages
func (m *MQTTClient) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic: %s", msg.Topic())

	var bambuMsg BambuMessage
	if err := json.Unmarshal(msg.Payload(), &bambuMsg); err != nil {
		log.Printf("Error parsing MQTT message: %v", err)
		return
	}

	// Convert to our status format
	status := bambuMsg.ToStatus()
	m.updateStatus(status)
}

// connectionLostHandler handles connection loss
func (m *MQTTClient) connectionLostHandler(client mqtt.Client, err error) {
	log.Printf("MQTT connection lost: %v", err)
}

// onConnectHandler handles successful connection and subscription
func (m *MQTTClient) onConnectHandler(client mqtt.Client) {
	log.Println("Connected to MQTT broker")

	var topic string
	if m.config.Serial != "" {
		// Use configured serial number
		topic = fmt.Sprintf("device/%s/report", m.config.Serial)
	} else {
		// Subscribe to all devices and auto-discover
		topic = "device/+/report"
		log.Println("Auto-discovering printer serial number...")
	}

	token := client.Subscribe(topic, 0, m.messageHandler)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Error subscribing to topic %s: %v", topic, token.Error())
	} else {
		log.Printf("Subscribed to topic: %s", topic)
	}
}

// Connect establishes the MQTT connection
func (m *MQTTClient) Connect() error {
	// Configure TLS with insecure skip verify
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Build broker URL
	broker := fmt.Sprintf("ssl://%s:%d", m.config.Host, m.config.Port)

	// Configure MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(fmt.Sprintf("bambustatus-%d", time.Now().Unix()))
	opts.SetUsername(m.config.Username)
	opts.SetPassword(m.config.Password)
	opts.SetTLSConfig(tlsConfig)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	opts.SetConnectionLostHandler(m.connectionLostHandler)
	opts.SetOnConnectHandler(m.onConnectHandler)

	// Create and connect the client
	m.client = mqtt.NewClient(opts)

	log.Printf("Connecting to MQTT broker at %s...", broker)
	token := m.client.Connect()
	token.Wait()

	if token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	return nil
}

// Disconnect closes the MQTT connection
func (m *MQTTClient) Disconnect() {
	if m.client != nil && m.client.IsConnected() {
		m.client.Disconnect(250)
		log.Println("Disconnected from MQTT broker")
	}
}

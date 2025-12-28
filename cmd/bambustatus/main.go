package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/rtrox/bambustatus/pkg/printer"
)

type Server struct {
	templates  *template.Template
	mqttClient *printer.MQTTClient
}

func NewServer(mqttClient *printer.MQTTClient) (*Server, error) {
	tmplPath := filepath.Join("web", "templates", "*.html")
	tmpl, err := template.ParseGlob(tmplPath)
	if err != nil {
		return nil, err
	}

	return &Server{
		templates:  tmpl,
		mqttClient: mqttClient,
	}, nil
}

func (s *Server) handleOverlay(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	status := s.mqttClient.GetStatus()
	if err := s.templates.ExecuteTemplate(w, "overlay.html", status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
	}
}

func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := s.mqttClient.GetStatus()
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error encoding JSON: %v", err)
	}
}

func main() {
	// Define command-line flags
	mqttHost := flag.String("host", "", "MQTT broker hostname/IP (required)")
	mqttPort := flag.Int("port", 8883, "MQTT broker port")
	mqttUsername := flag.String("username", "bblp", "MQTT username")
	mqttPassword := flag.String("password", "", "MQTT password (required)")
	mqttSerial := flag.String("serial", "", "Printer serial number (auto-discovered if empty)")
	httpPort := flag.String("http-port", "8080", "HTTP server port")

	flag.Parse()

	// Validate required flags
	if *mqttHost == "" {
		log.Fatal("Error: -host is required")
	}
	if *mqttPassword == "" {
		log.Fatal("Error: -password is required")
	}

	// Create MQTT client configuration
	config := printer.MQTTConfig{
		Host:     *mqttHost,
		Port:     *mqttPort,
		Username: *mqttUsername,
		Password: *mqttPassword,
		Serial:   *mqttSerial,
	}

	// Create MQTT client
	mqttClient := printer.NewMQTTClient(config)

	// Connect to MQTT broker
	if err := mqttClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer mqttClient.Disconnect()

	// Create HTTP server
	server, err := NewServer(mqttClient)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Set up HTTP routes
	http.HandleFunc("/", server.handleOverlay)
	http.HandleFunc("/api/status", server.handleAPI)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Start HTTP server in a goroutine
	go func() {
		addr := ":" + *httpPort
		log.Printf("Starting HTTP server on http://localhost%s", addr)
		log.Printf("OBS Browser Source URL: http://localhost%s/", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}

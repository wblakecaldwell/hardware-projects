package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ListenHostPort    string
	PanicWaterLevel   float64
	ServerSecret      string
	EmailRecipients   []string
	EmailSMTPLogin    string
	EmailSMTPPassword string
	EmailSMTPServer   string
	EmailSMTPPort     int
	EmailSender       string
}

func NewConfigFromEnv() (*Config, error) {
	var envVar string
	c := Config{}

	envVar = os.Getenv("LISTEN_HOST_PORT")
	if envVar == "" {
		return nil, fmt.Errorf("Missing environment variable LISTEN_HOST_PORT")
	}
	c.ListenHostPort = envVar

	// water value
	envVar = os.Getenv("PANIC_WATER_LEVEL")
	if envVar == "" {
		return nil, fmt.Errorf("Missing environment variable PANIC_WATER_LEVEL")
	}
	panicWaterLevel, err := strconv.ParseFloat(envVar, 64)
	if panicWaterLevel <= 0.0 {
		return nil, fmt.Errorf("Invalid water level: %f", panicWaterLevel)
	}
	c.PanicWaterLevel = panicWaterLevel

	// server secret
	envVar = os.Getenv("SERVER_SECRET")
	if envVar == "" {
		return nil, fmt.Errorf("Missing SERVER_SECRET")
	}
	c.ServerSecret = envVar

	// email recepients
	envVar = os.Getenv("EMAIL_RECIPIENTS")
	if len(envVar) == 0 {
		return nil, fmt.Errorf("Missing EMAIL_RECIPIENTS")
	}
	emailRecipients := strings.Split(envVar, ",")
	if len(emailRecipients) == 0 {
		return nil, fmt.Errorf("Missing email recipients")
	}
	c.EmailRecipients = emailRecipients

	// email smtp login
	envVar = os.Getenv("SMTP_LOGIN")
	if envVar == "" {
		return nil, fmt.Errorf("Missing SMTP_LOGIN")
	}
	c.EmailSMTPLogin = envVar

	// email smtp password
	envVar = os.Getenv("SMTP_PASSWORD")
	if envVar == "" {
		return nil, fmt.Errorf("Missing SMTP_PASSWORD")
	}
	c.EmailSMTPPassword = envVar

	// email smtp server
	envVar = os.Getenv("SMTP_SERVER")
	if envVar == "" {
		return nil, fmt.Errorf("Missing SMTP_SERVER")
	}
	c.EmailSMTPServer = envVar

	// email smtp port
	envVar = os.Getenv("SMTP_PORT")
	if envVar == "" {
		return nil, fmt.Errorf("Missing SMTP_PORT")
	}
	intVal, err := strconv.Atoi(envVar)
	if err != nil {
		return nil, fmt.Errorf("SMTP_PORT parse error: %s", err)
	}
	if intVal <= 0 {
		return nil, fmt.Errorf("SMTP_PORT less than zero: %d", intVal)
	}
	c.EmailSMTPPort = intVal

	// email sender
	envVar = os.Getenv("EMAIL_SENDER")
	if envVar == "" {
		return nil, fmt.Errorf("Missing EMAIL_SENDER")
	}
	c.EmailSender = envVar

	return &c, nil
}

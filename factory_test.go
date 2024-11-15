package main

import (
	"testing"
)

func TestFactory(t *testing.T) {
	factory := NotifyServiceFactory{}
	services := []struct {
		name        string
		serviceType string
		expectedErr string
	}{
		{"CreateTelegramService", "telegram", ""},
		{"CreateOSService", "os", ""},
		{"InvalidSMSService", "sms", "Service sms isn't valid\n"},
	}
	for _, service := range services {
		t.Run(service.name, func(t *testing.T) {
			_, err := factory.CreateService(service.serviceType)
			if err == nil && service.expectedErr != "" {
				t.Errorf("Expected '%s', got '%s'", service.expectedErr, err)
			}
			if err != nil && err.Error() != service.expectedErr {
				t.Errorf("Expected '%s', got '%s'", service.expectedErr, err)
			}
		})
	}
}

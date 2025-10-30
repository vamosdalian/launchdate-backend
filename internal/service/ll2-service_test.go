package service

import (
	"testing"

	"github.com/vamosdalian/launchdate-backend/internal/config"
)

func TestLoadLaunches(t *testing.T) {
	s := NewLL2Service(&config.Config{LL2URLPrefix: "https://lldev.thespacedevs.com"}, nil)
	launches, err := s.LoadLaunches(1, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(launches.Results) != 1 {
		t.Fatalf("Expected 1 launch, got %d", len(launches.Results))
	}
	t.Logf("%+v", launches)
}

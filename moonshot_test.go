package moonshot_test

import (
	"testing"
	
	"github.com/rizome-dev/go-moonshot"
)

func TestVersion(t *testing.T) {
	v := moonshot.Version()
	if v == "" {
		t.Error("Version() returned empty string")
	}
	if v != "0.1.0" {
		t.Errorf("Version() = %v, want 0.1.0", v)
	}
}

func TestNew(t *testing.T) {
	// Test with environment variable
	t.Run("with env var", func(t *testing.T) {
		t.Setenv("MOONSHOT_API_KEY", "test-key")
		sdk := moonshot.New()
		if sdk == nil {
			t.Fatal("New() returned nil")
		}
		if sdk.Client == nil {
			t.Error("SDK.Client is nil")
		}
		if sdk.Chat == nil {
			t.Error("SDK.Chat is nil")
		}
		if sdk.Files == nil {
			t.Error("SDK.Files is nil")
		}
	})
	
	// Test with API key parameter
	t.Run("with api key", func(t *testing.T) {
		sdk := moonshot.New("test-key")
		if sdk == nil {
			t.Fatal("New() returned nil")
		}
	})
}
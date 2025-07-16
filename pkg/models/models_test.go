package models_test

import (
	"testing"

	"github.com/rizome-dev/go-moonshot/pkg/models"
)

func TestModel_String(t *testing.T) {
	tests := []struct {
		name  string
		model models.Model
		want  string
	}{
		{"MoonshotV18K", models.MoonshotV18K, "moonshot-v1-8k"},
		{"MoonshotV132K", models.MoonshotV132K, "moonshot-v1-32k"},
		{"MoonshotV1128K", models.MoonshotV1128K, "moonshot-v1-128k"},
		{"KimiK2", models.KimiK2, "kimi-k2"},
		{"KimiK2Base", models.KimiK2Base, "kimi-k2-base"},
		{"KimiK2Instruct", models.KimiK2Instruct, "kimi-k2-instruct"},
		{"Legacy Moonshot8K", models.Moonshot8K, "moonshot-v1-8k"},
		{"Legacy Moonshot32K", models.Moonshot32K, "moonshot-v1-32k"},
		{"Legacy Moonshot128K", models.Moonshot128K, "moonshot-v1-128k"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.String()
			if got != tt.want {
				t.Errorf("Model.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		model models.Model
		want  bool
	}{
		{"MoonshotV18K", models.MoonshotV18K, true},
		{"MoonshotV132K", models.MoonshotV132K, true},
		{"MoonshotV1128K", models.MoonshotV1128K, true},
		{"KimiK2", models.KimiK2, true},
		{"KimiK2Base", models.KimiK2Base, true},
		{"KimiK2Instruct", models.KimiK2Instruct, true},
		{"Invalid model", models.Model("invalid-model"), false},
		{"Empty model", models.Model(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.IsValid()
			if got != tt.want {
				t.Errorf("Model.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_MaxTokens(t *testing.T) {
	tests := []struct {
		name  string
		model models.Model
		want  int
	}{
		{"MoonshotV18K", models.MoonshotV18K, 8192},
		{"MoonshotV132K", models.MoonshotV132K, 32768},
		{"MoonshotV1128K", models.MoonshotV1128K, 131072},
		{"KimiK2", models.KimiK2, 131072},
		{"KimiK2Base", models.KimiK2Base, 131072},
		{"KimiK2Instruct", models.KimiK2Instruct, 131072},
		{"Invalid model", models.Model("invalid-model"), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.MaxTokens()
			if got != tt.want {
				t.Errorf("Model.MaxTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_SupportsTools(t *testing.T) {
	tests := []struct {
		name  string
		model models.Model
		want  bool
	}{
		{"MoonshotV18K", models.MoonshotV18K, true},
		{"MoonshotV132K", models.MoonshotV132K, true},
		{"MoonshotV1128K", models.MoonshotV1128K, true},
		{"KimiK2", models.KimiK2, true},
		{"KimiK2Base", models.KimiK2Base, true},
		{"KimiK2Instruct", models.KimiK2Instruct, true},
		{"Invalid model", models.Model("invalid-model"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.SupportsTools()
			if got != tt.want {
				t.Errorf("Model.SupportsTools() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_SupportsVision(t *testing.T) {
	tests := []struct {
		name  string
		model models.Model
		want  bool
	}{
		{"MoonshotV18K", models.MoonshotV18K, false},
		{"MoonshotV132K", models.MoonshotV132K, false},
		{"MoonshotV1128K", models.MoonshotV1128K, false},
		{"KimiK2", models.KimiK2, true},
		{"KimiK2Base", models.KimiK2Base, false},
		{"KimiK2Instruct", models.KimiK2Instruct, true},
		{"Invalid model", models.Model("invalid-model"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.model.SupportsVision()
			if got != tt.want {
				t.Errorf("Model.SupportsVision() = %v, want %v", got, tt.want)
			}
		})
	}
}
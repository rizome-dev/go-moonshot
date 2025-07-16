package models

// Model represents a Moonshot AI model
type Model string

// Available Moonshot models
const (
	// Moonshot V1 models with different context lengths
	MoonshotV18K   Model = "moonshot-v1-8k"   // 8K context
	MoonshotV132K  Model = "moonshot-v1-32k"  // 32K context
	MoonshotV1128K Model = "moonshot-v1-128k" // 128K context
	
	// Kimi K2 models
	KimiK2         Model = "kimi-k2"          // Latest Kimi K2 model
	KimiK2Base     Model = "kimi-k2-base"     // Base model for fine-tuning
	KimiK2Instruct Model = "kimi-k2-instruct" // Instruction-tuned model
	
	// Legacy model aliases (for compatibility)
	Moonshot8K  = MoonshotV18K
	Moonshot32K = MoonshotV132K
	Moonshot128K = MoonshotV1128K
)

// String returns the string representation of a model
func (m Model) String() string {
	return string(m)
}

// IsValid checks if a model is valid
func (m Model) IsValid() bool {
	switch m {
	case MoonshotV18K, MoonshotV132K, MoonshotV1128K,
	     KimiK2, KimiK2Base, KimiK2Instruct:
		return true
	}
	return false
}

// MaxTokens returns the maximum tokens for a model
func (m Model) MaxTokens() int {
	switch m {
	case MoonshotV18K:
		return 8192
	case MoonshotV132K:
		return 32768
	case MoonshotV1128K:
		return 131072
	case KimiK2, KimiK2Base, KimiK2Instruct:
		return 131072 // K2 models support long context
	default:
		return 0
	}
}

// SupportsTools returns whether a model supports tool/function calling
func (m Model) SupportsTools() bool {
	// All current Moonshot models support tools
	return m.IsValid()
}

// SupportsVision returns whether a model supports vision/image inputs
func (m Model) SupportsVision() bool {
	// Based on documentation, Kimi models support vision
	switch m {
	case KimiK2, KimiK2Instruct:
		return true
	default:
		return false
	}
}
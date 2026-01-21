package ai

// Provider representa un proveedor de IA.
type Provider string

const (
	// ProviderCLI usa Claude CLI instalado localmente.
	ProviderCLI Provider = "cli"

	// ProviderClaudeAPI usa Anthropic Claude API.
	ProviderClaudeAPI Provider = "claude-api"

	// ProviderOpenAI usa OpenAI API.
	ProviderOpenAI Provider = "openai"

	// ProviderZAI usa z.ai API.
	ProviderZAI Provider = "zai"

	// ProviderGemini usa Google Gemini API.
	ProviderGemini Provider = "gemini"

	// ProviderGroq usa Groq API.
	ProviderGroq Provider = "groq"
)

// String retorna el nombre del provider.
func (p Provider) String() string {
	switch p {
	case ProviderCLI:
		return "Claude CLI"
	case ProviderClaudeAPI:
		return "Claude API"
	case ProviderOpenAI:
		return "OpenAI"
	case ProviderZAI:
		return "Z.AI"
	case ProviderGemini:
		return "Gemini"
	case ProviderGroq:
		return "Groq"
	default:
		return "Unknown"
	}
}

// Description retorna una descripción del provider para mostrar en el selector.
func (p Provider) Description() string {
	switch p {
	case ProviderCLI:
		return "Claude CLI instalado localmente (gratis con Claude Code PRO, más lento)"
	case ProviderClaudeAPI:
		return "Anthropic Claude API (requiere API key de pago, más rápido)"
	case ProviderOpenAI:
		return "OpenAI API (requiere API key, más rápido)"
	case ProviderZAI:
		return "Z.AI API (requiere API key, más rápido)"
	case ProviderGemini:
		return "Google Gemini API (free tier disponible, muy rápido)"
	case ProviderGroq:
		return "Groq API (free tier disponible, extremadamente rápido)"
	default:
		return "Proveedor desconocido"
	}
}

// RequiresAPIKey retorna true si el provider requiere API key.
func (p Provider) RequiresAPIKey() bool {
	return p != ProviderCLI
}

// AllProviders retorna todos los providers disponibles.
func AllProviders() []Provider {
	return []Provider{
		ProviderCLI,
		ProviderClaudeAPI,
		ProviderOpenAI,
		ProviderZAI,
		ProviderGemini,
		ProviderGroq,
	}
}

// IsValid retorna true si el provider es válido.
func (p Provider) IsValid() bool {
	switch p {
	case ProviderCLI, ProviderClaudeAPI, ProviderOpenAI, ProviderZAI, ProviderGemini, ProviderGroq:
		return true
	default:
		return false
	}
}

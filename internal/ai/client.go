package ai

// Message representa un mensaje enviado al provider de IA.
type Message struct {
	Role    string // "system" o "user"
	Content string
}

// Client define la interfaz común para todos los clientes de IA.
type Client interface {
	// SendMessage envía un mensaje y retorna la respuesta.
	SendMessage(systemPrompt, userMessage string) (string, error)

	// SendSimpleMessage envía un mensaje sin system prompt.
	SendSimpleMessage(message string) (string, error)

	// Provider retorna el tipo de provider.
	Provider() Provider

	// IsAvailable retorna true si el provider está disponible (CLI instalado o API key configurada).
	IsAvailable() (bool, error)

	// Close libera recursos del cliente (si aplica).
	Close() error
}

// ValidationResult contiene el resultado de validar las respuestas del usuario.
type ValidationResult struct {
	IsValid     bool     // true si las respuestas son válidas
	MissingInfo []string // información faltante detectada
	Suggestions []string // sugerencias de mejora
	Questions   []string // preguntas adicionales a hacer al usuario
}

// Recommendation contiene la recomendación de estructura del proyecto.
type Recommendation struct {
	Agents      []string // agents a generar
	Commands    []string // commands a generar
	Skills      []string // skills a incluir
	Description string   // descripción de la estructura recomendada
}

// NewValidationResult crea un nuevo ValidationResult.
func NewValidationResult(isValid bool, missingInfo, suggestions, questions []string) *ValidationResult {
	return &ValidationResult{
		IsValid:     isValid,
		MissingInfo: missingInfo,
		Suggestions: suggestions,
		Questions:   questions,
	}
}

// NewRecommendation crea una nueva Recommendation.
func NewRecommendation(agents, commands, skills []string, description string) *Recommendation {
	return &Recommendation{
		Agents:      agents,
		Commands:    commands,
		Skills:      skills,
		Description: description,
	}
}

package claude

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/drossan/claude-init/internal/ai"
	"github.com/drossan/claude-init/internal/logger"
	"github.com/drossan/claude-init/internal/survey"
)

//go:embed embeds/agent_guide.md
var agentGuideContent string

//go:embed embeds/command_guide.md
var commandGuideContent string

//go:embed embeds/skill_guide.md
var skillGuideContent string

// Recommendation contiene la recomendación de estructura del proyecto.
type Recommendation struct {
	Agents      []string `json:"agents"`
	Commands    []string `json:"commands"`
	Skills      []string `json:"skills"`
	Description string   `json:"description"`
}

// AgentInfo contiene información extraída de un archivo de agente.
type AgentInfo struct {
	Name        string   // Nombre del agente (ej: "architect")
	Description string   // Descripción del rol
	Color       string   // Color del agente
	Model       string   // Modelo preferido
	Tools       []string // Lista de herramientas disponibles
	Skills      []string // Skills inyectadas en el agente
	FilePath    string   // Ruta al archivo .md
}

// SkillInfo contiene información extraída de un archivo de skill.
type SkillInfo struct {
	Name        string // Nombre de la skill (ej: "go-expert")
	Category    string // Categoría: "language", "framework", "base"
	Description string // Descripción de la skill
	Purpose     string // Propósito de la skill
	FilePath    string // Ruta al archivo .md
}

// CommandInfo contiene información extraída de un archivo de comando.
type CommandInfo struct {
	Name        string // Nombre del comando (ej: "test", "lint")
	Description string // Descripción del comando
	Usage       string // Uso del comando
	FilePath    string // Ruta al archivo .md
}

// Generator es un wrapper para generar configuraciones usando un Client de IA.
type Generator struct {
	projectPath    string
	answers        *survey.Answers
	logger         *logger.Logger
	promptBuilder  *PromptBuilder
	templateLoader *TemplateLoader
	client         ai.Client
}

// NewGenerator crea una nueva instancia de Generator.
func NewGenerator(projectPath string, answers *survey.Answers, client ai.Client) *Generator {
	return &Generator{
		projectPath:    projectPath,
		answers:        answers,
		logger:         logger.New(nil, logger.WARNLevel),
		promptBuilder:  NewPromptBuilder(answers),
		templateLoader: NewTemplateLoader(),
		client:         client,
	}
}

// SetLogger establece el logger para el generador.
func (g *Generator) SetLogger(l *logger.Logger) {
	g.logger = l
}

// GenerateAgent genera un archivo de agente usando templates base o Claude CLI.
//
// Primero intenta usar un template base de claude_examples/ adaptado al proyecto.
// Si no existe, usa templates incrustados o Claude CLI para generarlo.
func (g *Generator) GenerateAgent(agentType string) error {
	g.logger.Debug("Generando agent %s para %s", agentType, g.answers.ProjectName)

	// Crear directorio de agentes
	agentsDir := filepath.Join(g.projectPath, ".claude", "agents")
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio agents: %w", err)
	}

	// Sanitizar el nombre del agente a kebab-case
	safeAgentName := sanitizeFilename(agentType)
	outputPath := filepath.Join(agentsDir, safeAgentName+".md")
	var content string
	var err error

	// 1. Intentar usar template base de claude_examples/ primero
	if g.templateLoader.HasTemplate("agent", agentType) {
		g.logger.Debug("Usando template base para agent %s", agentType)
		template, templateErr := g.templateLoader.LoadTemplate("agent", agentType)
		if templateErr == nil {
			// Adaptar template al proyecto actual
			content = g.templateLoader.AdaptTemplate(template, g.answers)
			err = nil
		} else {
			err = templateErr
		}
	}

	// 2. Si no hay template externo, usar template incrustado
	if content == "" {
		g.logger.Debug("Usando template incrustado para agent %s", agentType)
		content = g.getEmbeddedAgentTemplate(agentType)
	}

	// 3. Si todo falla, usar Claude CLI para generar
	if content == "" {
		g.logger.Debug("Generando agent %s con AI", agentType)
		prompt := g.buildAgentYAMLTemplate(agentType)
		content, err = g.generateWithClaude(prompt, nil)
		if err != nil {
			return fmt.Errorf("error generando agent %s: %w", agentType, err)
		}
	}

	// Escribir archivo
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo agent %s: %w", agentType, err)
	}

	g.logger.Info("Agent %s generado en %s", agentType, outputPath)
	return nil
}

// GenerateSkill genera un archivo de skill usando templates base o Claude CLI.
//
// Primero intenta usar un template base de claude_examples/ adaptado al proyecto.
// Si no existe, usa templates incrustados o Claude CLI para generarlo.
func (g *Generator) GenerateSkill(skillType, skillName string) error {
	g.logger.Debug("Generando skill %s:%s", skillType, skillName)

	// Crear directorio de skills
	skillsDir := filepath.Join(g.projectPath, ".claude", "skills", skillType)
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio skills/%s: %w", skillType, err)
	}

	// Sanitizar el nombre del archivo para que sea válido
	safeFileName := sanitizeFilename(skillName)
	outputPath := filepath.Join(skillsDir, safeFileName+".md")
	var content string
	var err error

	// 1. Intentar usar template base primero
	templateKey := skillType + "/" + skillName
	if g.templateLoader.HasTemplate("skill", templateKey) {
		g.logger.Debug("Usando template base para skill %s", templateKey)
		template, templateErr := g.templateLoader.LoadTemplate("skill", templateKey)
		if templateErr == nil {
			// Adaptar template al proyecto actual
			content = g.templateLoader.AdaptTemplate(template, g.answers)
			err = nil
		} else {
			err = templateErr
		}
	}

	// 2. Si no hay template externo, usar template incrustado
	if content == "" {
		g.logger.Debug("Usando template incrustado para skill %s", skillName)
		content = g.getEmbeddedSkillTemplate(skillName, skillType)
	}

	// 3. Si todo falla, usar Claude CLI para generar
	if content == "" {
		g.logger.Debug("Generando skill %s con AI", skillName)
		prompt := g.buildSkillTemplate(skillType, skillName)
		content, err = g.generateWithClaude(prompt, nil)
		if err != nil {
			return fmt.Errorf("error generando skill %s: %w", skillName, err)
		}
	}

	// Escribir archivo
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo skill %s: %w", skillName, err)
	}

	g.logger.Info("Skill %s generado en %s", skillName, outputPath)
	return nil
}

// GenerateCommand genera un archivo de comando usando templates base o Claude CLI.
//
// Primero intenta usar un template base de claude_examples/ adaptado al proyecto.
// Si no existe, usa templates incrustados o Claude CLI para generarlo.
// Esta versión es un wrapper que llama a GenerateCommandWithContext sin contexto.
func (g *Generator) GenerateCommand(commandType string) error {
	return g.GenerateCommandWithContext(commandType, "", "")
}

// GenerateCommandWithContext genera un archivo de comando usando templates base o Claude CLI.
//
// Similar a GenerateCommand pero acepta contexto adicional sobre agents y skills disponibles.
// Este contexto se usa cuando se genera con IA para informar qué agentes y skills referenciar.
func (g *Generator) GenerateCommandWithContext(commandType, agentsContext, skillsContext string) error {
	g.logger.Debug("Generando command %s", commandType)

	// Crear directorio de comandos
	commandsDir := filepath.Join(g.projectPath, ".claude", "commands")
	if err := os.MkdirAll(commandsDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio commands: %w", err)
	}

	// Sanitizar el nombre del comando a kebab-case
	safeCommandName := sanitizeFilename(commandType)
	outputPath := filepath.Join(commandsDir, safeCommandName+".md")
	var content string
	var err error

	// IMPORTANTE: Si se pasó contexto de agents/skills, usar SIEMPRE IA para que pueda
	// hacer matching dinámico con los agents/skills reales del proyecto.
	// Los templates base tienen agentes/skills hardcoded que pueden no existir.
	hasContext := agentsContext != "" || skillsContext != ""

	if !hasContext {
		// 1. Sin contexto: Intentar usar template base primero
		if g.templateLoader.HasTemplate("command", commandType) {
			g.logger.Debug("Usando template base para command %s", commandType)
			template, templateErr := g.templateLoader.LoadTemplate("command", commandType)
			if templateErr == nil {
				// Adaptar template al proyecto actual
				content = g.templateLoader.AdaptTemplate(template, g.answers)
				err = nil
			} else {
				err = templateErr
			}
		}

		// 2. Si no hay template externo, usar template incrustado
		if content == "" {
			g.logger.Debug("Usando template incrustado para command %s", commandType)
			content = g.getEmbeddedCommandTemplate(commandType)
		}
	}

	// 3. Si hay contexto O no hay templates, usar Claude CLI para generar CON CONTEXTO
	if content == "" {
		g.logger.Debug("Generando command %s con AI y contexto (hasContext=%v)", commandType, hasContext)
		prompt := g.buildCommandTemplateWithContext(commandType, agentsContext, skillsContext)
		content, err = g.generateWithClaude(prompt, nil)
		if err != nil {
			return fmt.Errorf("error generando command %s: %w", commandType, err)
		}
	}

	// Escribir archivo
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo command %s: %w", commandType, err)
	}

	g.logger.Info("Command %s generado en %s", commandType, outputPath)
	return nil
}

// GetRecommendation obtiene una recomendación de estructura usando Claude CLI.
//
// Usa claude -p con un prompt específico para obtener recomendaciones.
func (g *Generator) GetRecommendation() (*Recommendation, error) {
	g.logger.Debug("Obteniendo recomendación para %s", g.answers.ProjectName)

	prompt := g.promptBuilder.buildRecommendationPrompt()

	// Generar recomendación usando Claude CLI
	content, err := g.generateWithClaude(prompt, map[string]string{
		"output-format": "json",
	})
	if err != nil {
		return nil, fmt.Errorf("error obteniendo recomendación: %w", err)
	}

	// Intentar extraer JSON de la respuesta
	jsonStr := g.extractJSON(content)
	if jsonStr == "" {
		// Fallback a recomendación por defecto
		return g.getDefaultRecommendation(), nil
	}

	// Parsear JSON
	var rec Recommendation
	if err := json.Unmarshal([]byte(jsonStr), &rec); err != nil {
		g.logger.Warn("Error parseando JSON de recomendación, usando defaults: %v", err)
		return g.getDefaultRecommendation(), nil
	}

	// Normalizar todos los nombres a kebab-case
	rec = g.normalizeRecommendationNames(rec)

	return &rec, nil
}

// normalizeRecommendationNames normaliza todos los nombres de una recomendación a kebab-case.
func (g *Generator) normalizeRecommendationNames(rec Recommendation) Recommendation {
	normalized := Recommendation{
		Agents:      make([]string, len(rec.Agents)),
		Commands:    make([]string, len(rec.Commands)),
		Skills:      make([]string, len(rec.Skills)),
		Description: rec.Description,
	}

	// Normalizar nombres de agentes
	for i, agent := range rec.Agents {
		normalized.Agents[i] = sanitizeFilename(agent)
	}

	// Normalizar nombres de comandos
	for i, cmd := range rec.Commands {
		normalized.Commands[i] = sanitizeFilename(cmd)
	}

	// Normalizar nombres de skills
	for i, skill := range rec.Skills {
		normalized.Skills[i] = sanitizeFilename(skill)
	}

	return normalized
}

// GenerateAll genera toda la estructura .claude/ basada en una recomendación.
// Combina las recomendaciones de la IA con los items base obligatorios.
func (g *Generator) GenerateAll(rec *Recommendation) error {
	g.logger.Info("Generando estructura completa para %s", g.answers.ProjectName)

	// PASO 1: Generar CLAUDE.md PRIMERO para proporcionar contexto a las generaciones posteriores
	if err := g.GenerateClaudeMD(); err != nil {
		g.logger.Warn("Error generando CLAUDE.md: %v", err)
	}

	// Crear directorio base .claude
	configDir := filepath.Join(g.projectPath, ".claude")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio .claude: %w", err)
	}

	// Obtener items base que siempre deben estar presentes
	baseItems := GetBaseItems()

	// Combinar recomendaciones de IA con items base (sin duplicados)
	agents := g.combineUnique(rec.Agents, baseItems.Agents)
	commands := g.combineUnique(rec.Commands, baseItems.Commands)
	skills := g.combineUnique(rec.Skills, baseItems.Skills)

	// Generar agentes
	g.logger.Info("Generando %d agentes...", len(agents))
	for _, agent := range agents {
		if err := g.GenerateAgent(agent); err != nil {
			g.logger.Warn("Error generando agent %s: %v", agent, err)
		}
	}

	// PASO 3: Generar agents/README.md con la lista de agentes
	if err := g.GenerateAgentsReadme(agents); err != nil {
		g.logger.Warn("Error generando agents/README.md: %v", err)
		// Continuar aunque falle la generación del README
	}

	// Generar skills
	g.logger.Info("Generando %d skills...", len(skills))
	for _, skill := range skills {
		// Determinar tipo de skill basado en el contexto
		skillType := g.determineSkillType(skill)
		if err := g.GenerateSkill(skillType, skill); err != nil {
			g.logger.Warn("Error generando skill %s: %v", skill, err)
		}
	}

	// PASO 4: Generar skills/README.md con la lista de skills
	if err := g.GenerateSkillsReadme(skills); err != nil {
		g.logger.Warn("Error generando skills/README.md: %v", err)
		// Continuar aunque falle la generación del README
	}

	// Leer el contenido de los READMEs para pasarlo como contexto a los commands
	agentsReadme := g.getReadmeContent("agents")
	skillsReadme := g.getReadmeContent("skills")

	// PASO 5: Generar comandos con contexto de agents y skills
	g.logger.Info("Generando %d comandos...", len(commands))
	for _, cmd := range commands {
		if err := g.GenerateCommandWithContext(cmd, agentsReadme, skillsReadme); err != nil {
			g.logger.Warn("Error generando command %s: %v", cmd, err)
		}
	}

	// PASO 5.5: Generar commands/README.md con la lista de comandos
	if err := g.GenerateCommandsReadme(commands); err != nil {
		g.logger.Warn("Error generando commands/README.md: %v", err)
		// Continuar aunque falle la generación del README
	}

	// PASO 6: Generar development_guide.md CON CONTEXTO COMPLETO
	// Ahora tenemos toda la estructura creada, podemos pasar contexto al development guide
	if err := g.GenerateDevelopmentGuideWithContext(agents, commands, skills); err != nil {
		g.logger.Warn("Error generando development_guide.md: %v", err)
	}

	g.logger.Info("Estructura .claude/ generada exitosamente")
	return nil
}

// combineUnique combina dos slices eliminando duplicados.
func (g *Generator) combineUnique(recommended, base []string) []string {
	seen := make(map[string]bool)
	var result []string

	// Primero agregar items base
	for _, item := range base {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	// Luego agregar recomendaciones (si no están ya)
	for _, item := range recommended {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// determineSkillType determina el tipo de skill basado en el nombre.
// Skills base se clasifican en una carpeta especial "base".
func (g *Generator) determineSkillType(skillName string) string {
	// Verificar si es un skill base
	baseItems := GetBaseItems()
	for _, baseSkill := range baseItems.Skills {
		if skillName == baseSkill {
			return "base"
		}
	}

	// Si coincide con el lenguaje del proyecto
	if strings.EqualFold(skillName, g.answers.Language) {
		return "language"
	}

	// Si coincide con el framework
	if g.answers.Framework != "" && strings.Contains(strings.ToLower(g.answers.Framework), strings.ToLower(skillName)) {
		return "framework"
	}

	// Por defecto, language
	return "language"
}

// generateWithClaude ejecuta el cliente de IA con el prompt dado y retorna la salida.
// Incluye el contexto del CLAUDE.md si existe.
func (g *Generator) generateWithClaude(prompt string, extraFlags map[string]string) (string, error) {
	g.logger.Debug("Enviando prompt a AI client")

	// Construir system prompt con contexto del proyecto si está disponible
	systemPrompt := g.buildSystemPrompt()

	// Intentar leer CLAUDE.md para agregar contexto adicional
	if claudeContext := g.readClaudeMDContext(); claudeContext != "" {
		systemPrompt += `

` + claudeContext
	}

	// Enviar mensaje al cliente de IA
	response, err := g.client.SendMessage(systemPrompt, prompt)
	if err != nil {
		return "", fmt.Errorf("AI client error: %w", err)
	}

	return response, nil
}

// readClaudeMDContext lee el archivo CLAUDE.md si existe y retorna su contenido como contexto.
func (g *Generator) readClaudeMDContext() string {
	claudeMDPath := filepath.Join(g.projectPath, "CLAUDE.md")
	content, err := os.ReadFile(claudeMDPath)
	if err != nil {
		// El archivo no existe o no se puede leer, no hay problema
		return ""
	}

	// Retornar el contenido como contexto adicional
	return fmt.Sprintf(`# CONTEXTO DEL PROYECTO (desde CLAUDE.md)

A continuación se proporciona el contexto actual del proyecto que debes tener en cuenta
para generar contenido coherente y consistente con el proyecto existente:

%s
---
FIN DEL CONTEXTO DEL PROYECTO
`, string(content))
}

// getReadmeContent lee el contenido de un README (agents o skills) para pasarlo como contexto.
// Si el README no existe, retorna un string vacío (no es error).
func (g *Generator) getReadmeContent(readmeType string) string {
	var readmePath string
	if readmeType == "agents" {
		readmePath = filepath.Join(g.projectPath, ".claude", "agents", "README.md")
	} else if readmeType == "skills" {
		readmePath = filepath.Join(g.projectPath, ".claude", "skills", "README.md")
	} else {
		return ""
	}

	content, err := os.ReadFile(readmePath)
	if err != nil {
		// El README no existe todavía, no es error
		return ""
	}

	return string(content)
}

// buildSystemPrompt construye el system prompt para Claude.
func (g *Generator) buildSystemPrompt() string {
	return `Eres un experto en desarrollo de software y generacion de configuraciones para proyectos.

Tu tarea es generar archivos de configuracion para agentes, comandos y skills de desarrollo asistido por IA.

Debes ser preciso y generar contenido que sea directamente utilizable sin necesidad de edicion posterior.`
}

// extractJSON extrae un objeto JSON de un string que puede contener texto adicional.
func (g *Generator) extractJSON(s string) string {
	s = strings.TrimSpace(s)

	// Buscar el primer { y el último }
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")

	if start == -1 || end == -1 || end <= start {
		return ""
	}

	return s[start : end+1]
}

// getDefaultRecommendation retorna una recomendación por defecto basada en el proyecto.
func (g *Generator) getDefaultRecommendation() *Recommendation {
	rec := &Recommendation{
		Agents:      []string{"architect", "developer", "tester", "reviewer"},
		Commands:    []string{"test", "lint", "build"},
		Skills:      []string{sanitizeFilename(g.answers.Language)},
		Description: fmt.Sprintf("Default structure for %s (%s)", g.answers.ProjectName, g.answers.Language),
	}

	// Agregar agent de debugger si es un proyecto complejo
	if g.answers.Architecture != "" && g.answers.Architecture != "monolith" {
		rec.Agents = append(rec.Agents, "debugger")
	}

	// Agregar skill de framework si existe (normalizado a kebab-case)
	if g.answers.Framework != "" {
		rec.Skills = append(rec.Skills, sanitizeFilename(g.answers.Framework))
	}

	return rec
}

// Constants para los agents/commands/skills base que siempre deben estar presentes.
const (
	// Agents base obligatorios
	baseAgents = "architect,writer,debugger,planning-agent,orchestrator-agent"
	// Commands base obligatorios
	baseCommands = "plan-manage,orchestrator,pre-flight"
	// Skills base obligatorios
	baseSkills = "technical-writer,code-reviewer,debug-master"
)

// BaseItems contiene los items base que siempre deben estar presentes.
type BaseItems struct {
	Agents   []string
	Commands []string
	Skills   []string
}

// GetBaseItems retorna los items base que siempre deben estar presentes.
func GetBaseItems() *BaseItems {
	return &BaseItems{
		Agents:   []string{"architect", "writer", "debugger", "planning-agent", "orchestrator-agent"},
		Commands: []string{"plan-manage", "orchestrator", "pre-flight"},
		Skills:   []string{"technical-writer", "code-reviewer", "debug-master"},
	}
}

// sanitizeFilename sanitiza un nombre y lo convierte a kebab-case.
// Convierte a minúsculas, reemplaza espacios y caracteres inválidos por guiones.
// Ejemplo: "PlanningAgent" -> "planning-agent", "codeReviewer" -> "code-reviewer"
func sanitizeFilename(name string) string {
	// Primero, manejar casos especiales de camelCase y PascalCase
	// Insertar guiones antes de mayúsculas (excepto al inicio)
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			// Verificar si el carácter anterior es minúscula o número
			prev := name[i-1]
			if (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9') {
				result.WriteRune('-')
			}
		}
		result.WriteRune(r)
	}
	name = result.String()

	// Convertir a minúsculas
	name = strings.ToLower(name)

	// Reemplazar caracteres inválidos por guiones (excepto guiones y guiones bajos existentes)
	reg := regexp.MustCompile(`[^\w\-]`)
	name = reg.ReplaceAllString(name, "-")

	// Eliminar múltiples guiones consecutivos
	reg = regexp.MustCompile(`-+`)
	name = reg.ReplaceAllString(name, "-")

	// Eliminar guiones al inicio y al final
	name = strings.Trim(name, "-")

	// Si queda vacío, usar un nombre por defecto
	if name == "" {
		name = "unnamed"
	}

	return name
}

// parseAgentFrontmatter extrae metadata del frontmatter YAML de un archivo de agente.
// Si el archivo no tiene frontmatter válido, retorna un AgentInfo con valores por defecto.
func (g *Generator) parseAgentFrontmatter(filePath string) (*AgentInfo, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Extraer nombre del archivo
	fileName := filepath.Base(filePath)
	name := strings.TrimSuffix(fileName, ".md")

	info := &AgentInfo{
		Name:     name,
		FilePath: filePath,
		Color:    "gray",
		Model:    "sonnet",
		Tools:    []string{},
		Skills:   []string{},
	}

	// Buscar frontmatter YAML entre ---
	strContent := string(content)
	startIdx := strings.Index(strContent, "---")
	if startIdx == -1 {
		// No hay frontmatter, usar valores por defecto
		info.Description = fmt.Sprintf("Agent %s", name)
		return info, nil
	}

	endIdx := strings.Index(strContent[startIdx+3:], "---")
	if endIdx == -1 {
		// No hay cierre del frontmatter
		info.Description = fmt.Sprintf("Agent %s", name)
		return info, nil
	}

	frontmatter := strContent[startIdx+3 : startIdx+3+endIdx]
	lines := strings.Split(frontmatter, "\n")

	// Parsear líneas del frontmatter
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Formato "key: value"
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])

			switch key {
			case "name":
				info.Name = value
			case "description":
				info.Description = value
			case "color":
				info.Color = value
			case "model":
				info.Model = value
			case "tools":
				// Parsear lista de herramientas
				info.Tools = parseListValue(value)
			case "skills":
				// Parsear lista de skills
				info.Skills = parseListValue(value)
			}
		}
	}

	// Si no se encontró descripción en frontmatter, buscar en el body
	if info.Description == "" {
		// Buscar primera línea de título después del frontmatter
		bodyStart := startIdx + 3 + endIdx + 3
		if bodyStart < len(strContent) {
			body := strContent[bodyStart:]
			// Buscar primera línea # o texto no vacío
			bodyLines := strings.Split(body, "\n")
			for _, line := range bodyLines {
				line = strings.TrimSpace(line)
				if line != "" {
					// Remover # si existe
					info.Description = strings.TrimPrefix(line, "#")
					info.Description = strings.TrimSpace(info.Description)
					if len(info.Description) > 200 {
						info.Description = info.Description[:200] + "..."
					}
					break
				}
			}
		}
	}

	// Fallback final si aún no hay descripción
	if info.Description == "" {
		info.Description = fmt.Sprintf("Agent %s", name)
	}

	return info, nil
}

// parseSkillFrontmatter extrae metadata del frontmatter YAML de un archivo de skill.
// Si el archivo no tiene frontmatter válido, retorna un SkillInfo con valores por defecto.
func (g *Generator) parseSkillFrontmatter(filePath string) (*SkillInfo, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Extraer nombre del archivo y categoría del directorio
	fileName := filepath.Base(filePath)
	name := strings.TrimSuffix(fileName, ".md")

	// Determinar categoría del directorio padre
	category := "unknown"
	parentDir := filepath.Dir(filePath)
	if strings.Contains(parentDir, "/language/") || strings.HasSuffix(parentDir, "/language") {
		category = "language"
	} else if strings.Contains(parentDir, "/framework/") || strings.HasSuffix(parentDir, "/framework") {
		category = "framework"
	} else if strings.Contains(parentDir, "/base/") || strings.HasSuffix(parentDir, "/base") {
		category = "base"
	}

	info := &SkillInfo{
		Name:     name,
		Category: category,
		FilePath: filePath,
	}

	// Buscar frontmatter YAML entre ---
	strContent := string(content)
	startIdx := strings.Index(strContent, "---")
	if startIdx == -1 {
		// No hay frontmatter, usar valores por defecto
		info.Description = fmt.Sprintf("Skill %s", name)
		info.Purpose = fmt.Sprintf("Provides %s capabilities", name)
		return info, nil
	}

	endIdx := strings.Index(strContent[startIdx+3:], "---")
	if endIdx == -1 {
		info.Description = fmt.Sprintf("Skill %s", name)
		info.Purpose = fmt.Sprintf("Provides %s capabilities", name)
		return info, nil
	}

	frontmatter := strContent[startIdx+3 : startIdx+3+endIdx]
	lines := strings.Split(frontmatter, "\n")

	// Parsear líneas del frontmatter
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Formato "key: value"
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			value := strings.TrimSpace(line[idx+1:])

			switch key {
			case "name":
				info.Name = value
			case "description":
				info.Description = value
			case "category":
				info.Category = value
			case "purpose":
				info.Purpose = value
			}
		}
	}

	// Si no se encontró descripción en frontmatter, buscar en el body
	if info.Description == "" {
		bodyStart := startIdx + 3 + endIdx + 3
		if bodyStart < len(strContent) {
			body := strContent[bodyStart:]
			bodyLines := strings.Split(body, "\n")
			for _, line := range bodyLines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					info.Description = line
					if len(info.Description) > 200 {
						info.Description = info.Description[:200] + "..."
					}
					break
				}
			}
		}
	}

	// Fallbacks
	if info.Description == "" {
		info.Description = fmt.Sprintf("Skill %s", name)
	}
	if info.Purpose == "" {
		info.Purpose = fmt.Sprintf("Provides %s-related expertise and capabilities", name)
	}

	return info, nil
}

// parseListValue parsea un valor que puede ser una lista YAML o string simple.
func parseListValue(value string) []string {
	value = strings.TrimSpace(value)

	// Si está vacío, retornar lista vacía
	if value == "" {
		return []string{}
	}

	// Si contiene corchetes, es un array YAML
	if strings.HasPrefix(value, "[") {
		value = strings.TrimPrefix(value, "[")
		value = strings.TrimSuffix(value, "]")
		parts := strings.Split(value, ",")
		var result []string
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
		return result
	}

	// Si es un string simple, retornar como lista de un elemento
	return []string{value}
}

// GenerateAgentsReadme genera un README.md en el directorio agents/ que lista todos los agentes.
// El README incluye nombre, descripción, color, herramientas y skills de cada agente.
func (g *Generator) GenerateAgentsReadme(agentNames []string) error {
	g.logger.Debug("Generando agents/README.md")

	agentsDir := filepath.Join(g.projectPath, ".claude", "agents")
	outputPath := filepath.Join(agentsDir, "README.md")

	// Parsear todos los agentes para obtener su metadata
	var agents []*AgentInfo
	for _, agentName := range agentNames {
		agentPath := filepath.Join(agentsDir, sanitizeFilename(agentName)+".md")
		info, err := g.parseAgentFrontmatter(agentPath)
		if err != nil {
			g.logger.Warn("Error parsing agent %s: %v", agentName, err)
			continue
		}
		agents = append(agents, info)
	}

	// Generar contenido del README
	content := g.buildAgentsReadmeContent(agents)

	// Escribir archivo
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo agents/README.md: %w", err)
	}

	g.logger.Info("Agents README generado en %s", outputPath)
	return nil
}

// buildAgentsReadmeContent construye el contenido markdown del README de agents.
func (g *Generator) buildAgentsReadmeContent(agents []*AgentInfo) string {
	var sb strings.Builder

	sb.WriteString("# Agents Disponibles\n\n")
	sb.WriteString("Este directorio contiene los agentes especializados del proyecto ")
	sb.WriteString(g.answers.ProjectName)
	sb.WriteString(".\n\n")

	if len(agents) == 0 {
		sb.WriteString("*No hay agentes configurados aún.*\n")
		return sb.String()
	}

	sb.WriteString("## Agentes Configurados\n\n")

	for _, agent := range agents {
		sb.WriteString(fmt.Sprintf("### %s\n\n", agent.Name))

		if agent.Description != "" {
			sb.WriteString(fmt.Sprintf("**Descripción**: %s\n\n", agent.Description))
		}

		// Metadata en una lista compacta
		sb.WriteString(fmt.Sprintf("**Color**: %s\n", agent.Color))
		if agent.Model != "" {
			sb.WriteString(fmt.Sprintf("\n**Modelo**: %s", agent.Model))
		}

		if len(agent.Tools) > 0 {
			sb.WriteString("\n\n**Herramientas**: ")
			sb.WriteString(strings.Join(agent.Tools, ", "))
		}

		if len(agent.Skills) > 0 {
			sb.WriteString("\n\n**Skills Inyectadas**:\n")
			for _, skill := range agent.Skills {
				sb.WriteString(fmt.Sprintf("- %s\n", skill))
			}
		}

		sb.WriteString("\n---\n\n")
	}

	sb.WriteString("## Uso de los Agentes\n\n")
	sb.WriteString("Los agentes se definen en los comandos del directorio `commands/`. ")
	sb.WriteString("Cada comando orquesta uno o más agentes para ejecutar una tarea específica.\n\n")
	sb.WriteString("Para más información sobre cómo se crean estos agentes, ")
	sb.WriteString("consulta las guías en `.claude/embeds/agent_guide.md`.\n")

	return sb.String()
}

// GenerateSkillsReadme genera un README.md en el directorio skills/ que lista todas las skills.
// El README agrupa skills por categoría (language, framework, base) con descripción y propósito.
func (g *Generator) GenerateSkillsReadme(skillNames []string) error {
	g.logger.Debug("Generando skills/README.md")

	// Primero necesitamos determinar la categoría de cada skill
	// Para esto, necesitamos recorrer los subdirectorios
	skillsDir := filepath.Join(g.projectPath, ".claude", "skills")
	outputPath := filepath.Join(skillsDir, "README.md")

	// Parsear todas las skills para obtener su metadata
	var skills []*SkillInfo
	for _, skillName := range skillNames {
		// Buscar en diferentes subdirectorios
		categories := []string{"base", "language", "framework"}
		var foundPath string

		for _, cat := range categories {
			catDir := filepath.Join(skillsDir, cat)
			testPath := filepath.Join(catDir, sanitizeFilename(skillName)+".md")
			if _, err := os.Stat(testPath); err == nil {
				foundPath = testPath
				break
			}
		}

		// Si no se encontró en subdirectorios, buscar directamente en skills/
		if foundPath == "" {
			testPath := filepath.Join(skillsDir, sanitizeFilename(skillName)+".md")
			if _, err := os.Stat(testPath); err == nil {
				foundPath = testPath
			}
		}

		if foundPath == "" {
			g.logger.Warn("No se encontró archivo para skill %s", skillName)
			continue
		}

		info, err := g.parseSkillFrontmatter(foundPath)
		if err != nil {
			g.logger.Warn("Error parsing skill %s: %v", skillName, err)
			continue
		}
		skills = append(skills, info)
	}

	// Generar contenido del README
	content := g.buildSkillsReadmeContent(skills)

	// Escribir archivo
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo skills/README.md: %w", err)
	}

	g.logger.Info("Skills README generado en %s", outputPath)
	return nil
}

// buildSkillsReadmeContent construye el contenido markdown del README de skills.
func (g *Generator) buildSkillsReadmeContent(skills []*SkillInfo) string {
	var sb strings.Builder

	sb.WriteString("# Skills Disponibles\n\n")
	sb.WriteString("Este directorio contiene las habilidades técnicas inyectables en los agentes del proyecto ")
	sb.WriteString(g.answers.ProjectName)
	sb.WriteString(".\n\n")

	if len(skills) == 0 {
		sb.WriteString("*No hay skills configuradas aún.*\n")
		return sb.String()
	}

	// Agrupar skills por categoría
	byCategory := make(map[string][]*SkillInfo)
	for _, skill := range skills {
		cat := skill.Category
		if cat == "" || cat == "unknown" {
			cat = "other"
		}
		byCategory[cat] = append(byCategory[cat], skill)
	}

	sb.WriteString("## Skills por Categoría\n\n")

	// Orden de categorías
	categories := []string{"language", "framework", "base", "other"}
	for _, cat := range categories {
		skillsInCat := byCategory[cat]
		if len(skillsInCat) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("### %s/\n\n", strings.Title(cat)))

		for _, skill := range skillsInCat {
			sb.WriteString(fmt.Sprintf("#### %s\n\n", skill.Name))

			if skill.Description != "" {
				sb.WriteString(fmt.Sprintf("**Descripción**: %s\n\n", skill.Description))
			}

			if skill.Purpose != "" {
				sb.WriteString(fmt.Sprintf("**Propósito**: %s\n\n", skill.Purpose))
			}

			sb.WriteString("---\n\n")
		}
	}

	sb.WriteString("## Cómo Funcionan las Skills\n\n")
	sb.WriteString("Las skills se inyectan en los agentes mediante el frontmatter YAML:\n\n")
	sb.WriteString("```yaml\n")
	sb.WriteString("---\n")
	sb.WriteString("skills:\n")
	sb.WriteString("  - go-expert\n")
	sb.WriteString("  - cobra-cli\n")
	sb.WriteString("  - http-client\n")
	sb.WriteString("---\n")
	sb.WriteString("```\n\n")
	sb.WriteString("Para más información sobre cómo se crean estas skills, ")
	sb.WriteString("consulta las guías en `.claude/embeds/skill_guide.md`.\n")

	return sb.String()
}

// getEmbeddedAgentTemplate retorna un template incrustado para un agente.
func (g *Generator) getEmbeddedAgentTemplate(agentType string) string {
	roleDesc := g.getAgentRoleDescription(agentType)
	agentColor := g.getAgentColor(agentType)
	agentTools := g.getAgentToolsList(agentType)

	// Skills según tipo de agente
	skills := g.getAgentSkills(agentType)

	return fmt.Sprintf(`---
name: %s
description: %s
tools: %s
model: sonnet
color: %s
---

# Agente %s (%s) - %s

## Rol
%s

## Tu Especialidad
Tu maestría técnica se fundamenta en las **skills** que se te proporcionan en cada intervención:
%s

## Proceso de Trabajo
1. **Análisis**: Comprender el contexto y los requisitos de la tarea.
2. **Planificación**: Desglosar la tarea en pasos accionables.
3. **Implementación**: Ejecutar siguiendo las mejores prácticas de %s.
4. **Validación**: Asegurar que el resultado cumple con los estándares de calidad.
5. **Documentación**: Registrar decisiones y patrones aplicados.

## Convenciones de Código
- **Tipado**: Seguir las convenciones de tipo de %s.
- **Naming**: Usar convenciones estándar del lenguaje.
- **Estructura**: Mantener una organización clara y modular.
- **Estilo**: Seguir las guías de estilo del proyecto.

## Reglas de Oro
- **Calidad**: Priorizar código limpio y mantenible.
- **Testing**: Asegurar cobertura de pruebas adecuada.
- **Documentación**: Documentar código complejo y decisiones arquitectónicas.
- **Colaboración**: Coordinar con otros agentes cuando sea necesario.`,
		agentType,
		roleDesc,
		agentTools,
		agentColor,
		capitalize(agentType),
		strings.ToUpper(agentType[:1])+agentType[1:],
		g.answers.ProjectName,
		roleDesc,
		skills,
		g.answers.Language,
		g.answers.Language,
	)
}

// getAgentColor retorna el color para un tipo de agente.
func (g *Generator) getAgentColor(agentType string) string {
	colors := map[string]string{
		"architect":    "cyan",
		"developer":    "pink",
		"tester":       "green",
		"reviewer":     "yellow",
		"debugger":     "red",
		"writer":       "blue",
		"planner":      "purple",
		"orchestrator": "orange",
	}

	if color, ok := colors[agentType]; ok {
		return color
	}
	return "gray"
}

// getAgentToolsList retorna la lista de herramientas para un agente.
func (g *Generator) getAgentToolsList(agentType string) string {
	// Herramientas base disponibles
	baseTools := "Read, Write, Edit, Bash, Glob, Grep"

	// Herramientas adicionales según tipo
	extraTools := map[string]string{
		"architect": ", Test, WebSearch",
		"developer": ", Test",
		"tester":    ", Test, RunTests",
		"reviewer":  "",
		"debugger":  ", Bash, RunTests, Browser",
		"writer":    "",
	}

	if extra, ok := extraTools[agentType]; ok {
		return baseTools + extra
	}
	return baseTools
}

// getAgentSkills retorna las skills para un tipo de agente.
func (g *Generator) getAgentSkills(agentType string) string {
	langSkill := strings.ToLower(g.answers.Language)

	skills := map[string]string{
		"architect": fmt.Sprintf(`- **%s-expert**: Para el dominio del lenguaje y patrones arquitectónicos.
- **system-architect**: Para el diseño de arquitectura y estructura de módulos.
- **code-reviewer**: Para validar la calidad del diseño arquitectónico.`,
			langSkill,
		),
		"developer": fmt.Sprintf(`- **%s-expert**: Para el dominio avanzado de %s y optimización.
- **domain-expert**: Para la creación de entidades robustas y DTOs.
- **usecase-developer**: Para la implementación de la lógica de negocio.`,
			langSkill, g.answers.Language,
		),
		"tester": fmt.Sprintf(`- **qa-engineer**: Para estrategias de testing y cobertura.
- **%s-expert**: Para tests específicos del lenguaje.
- **tdd-champion**: Para metodología TDD.`,
			langSkill,
		),
		"reviewer": fmt.Sprintf(`- **code-reviewer**: Para revisión de calidad y mejores prácticas.
- **%s-expert**: Para validar convenciones del lenguaje.`,
			langSkill,
		),
		"debugger": fmt.Sprintf(`- **debug-master**: Para análisis y resolución de bugs.
- **%s-expert**: Para entender problemas específicos del lenguaje.`,
			langSkill,
		),
		"writer": `- **technical-writer**: Para documentación técnica.
- **code-reviewer**: Para entender el código a documentar.`,
	}

	if skill, ok := skills[agentType]; ok {
		return skill
	}

	return fmt.Sprintf(`- **%s-expert**: Para dominio del lenguaje.
- **code-reviewer**: Para revisión de código.`,
		langSkill,
	)
}

// getAgentRoleDescription retorna la descripción del rol para un tipo de agente.
func (g *Generator) getAgentRoleDescription(agentType string) string {
	descriptions := map[string]string{
		"architect": fmt.Sprintf("The Architect Agent for %s is responsible for designing system architecture following best practices, defining package structure and interfaces, and ensuring scalability and maintainability.", g.answers.ProjectName),
		"developer": fmt.Sprintf("The Developer Agent for %s is responsible for developing and maintaining the codebase. This includes writing clean, maintainable %s code, following project architecture, writing comprehensive tests, and ensuring high code quality.", g.answers.ProjectName, g.answers.Language),
		"tester":    fmt.Sprintf("The QA/Testing Agent for %s is responsible for writing comprehensive tests, finding edge cases and potential bugs, ensuring test coverage meets quality standards, and validating requirements against test results.", g.answers.ProjectName),
		"reviewer":  fmt.Sprintf("The Code Review Agent for %s is responsible for reviewing code for quality, maintainability, and best practices, ensuring adherence to %s coding standards, identifying potential bugs and performance problems, and providing constructive feedback.", g.answers.ProjectName, g.answers.Language),
		"debugger":  fmt.Sprintf("The Debugging Specialist Agent for %s is responsible for investigating and resolving complex bugs, analyzing stack traces and error messages, identifying race conditions and memory leaks, and proposing fixes.", g.answers.ProjectName),
		"writer":    fmt.Sprintf("The Technical Writer Agent for %s is responsible for creating and maintaining documentation, writing clear and concise technical guides, ensuring documentation is up-to-date, and following documentation best practices.", g.answers.ProjectName),
	}

	if desc, ok := descriptions[agentType]; ok {
		return desc
	}
	return fmt.Sprintf("The %s Agent for %s assists with development tasks.", capitalize(agentType), g.answers.ProjectName)
}

// getAgentResponsibilities retorna las responsabilidades para un tipo de agente.
func (g *Generator) getAgentResponsibilities(agentType string) string {
	resps := map[string]string{
		"architect": `      - Design system architecture following best practices
      - Define package structure and interfaces
      - Ensure scalability and maintainability
      - Follow SOLID principles and %s conventions`,
		"developer": `      - Write clean, maintainable code following %s best practices
      - Follow project conventions and architecture
      - Write comprehensive tests for new features
      - Document complex logic and decisions
      - Ensure code quality and performance`,
		"tester": `      - Write comprehensive tests (unit, integration, e2e)
      - Find edge cases and potential bugs
      - Ensure test coverage meets quality standards
      - Validate requirements against test results
      - Create testing strategies and plans`,
		"reviewer": `      - Review code for quality, maintainability, and best practices
      - Ensure adherence to %s coding standards
      - Identify potential bugs, security issues, and performance problems
      - Provide constructive feedback and suggestions
      - Validate architectural decisions`,
		"debugger": `      - Investigate and resolve complex bugs
      - Analyze stack traces and error messages
      - Identify race conditions and memory leaks
      - Root cause analysis for production issues
      - Propose and implement fixes`,
		"writer": `      - Create and maintain project documentation
      - Write clear API documentation
      - Create technical guides and tutorials
      - Keep documentation synchronized with code changes
      - Follow documentation best practices`,
	}

	if resp, ok := resps[agentType]; ok {
		return fmt.Sprintf(resp, g.answers.Language)
	}
	return fmt.Sprintf(`      - Assist with %s-related tasks
      - Follow %s best practices
      - Ensure code quality`, agentType, g.answers.Language)
}

// getAgentGuidelines retorna las guías de desarrollo para un tipo de agente.
func (g *Generator) getAgentGuidelines(agentType string) string {
	baseGuidelines := []string{
		fmt.Sprintf("    - Adhere to %s best practices and coding standards", g.answers.Language),
		"    - Write modular and reusable code structures",
		"    - Follow SOLID principles for software design",
		"    - Use descriptive and consistent naming conventions",
		"    - Document complex logic with inline comments",
		"    - Maintain project structure as per existing conventions",
	}

	// Agregar guías específicas según lenguaje
	languageExtras := map[string]string{
		"TypeScript": "    - Use strict typing and interfaces for type safety",
		"Go":         "    - Follow Go conventions and effective Go guidelines",
		"Python":     "    - Follow PEP 8 style guide and Python conventions",
		"Rust":       "    - Follow Rust idioms and ownership patterns",
	}

	if extra, ok := languageExtras[g.answers.Language]; ok {
		baseGuidelines = append(baseGuidelines, extra)
	}

	return strings.Join(baseGuidelines, "\n")
}

// getAgentTools retorna las herramientas para un tipo de agente.
func (g *Generator) getAgentTools(agentType string) string {
	baseTools := []string{
		"  - name: Git",
		"    description: Version control system",
	}

	// Herramientas según tipo
	extraTools := map[string][]string{
		"architect": {
			"  - name: Architecture Diagram Tools",
			"    description: Tools for creating architecture diagrams",
		},
		"developer": {
			"  - name: IDE",
			fmt.Sprintf("    description: IDE for %s development", g.answers.Language),
			"  - name: Linter",
			"    description: Code quality tool",
			"  - name: Testing Framework",
			"    description: Framework for running tests",
		},
		"tester": {
			"  - name: Testing Framework",
			"    description: Framework for writing and running tests",
			"  - name: Coverage Tool",
			"    description: Tool for measuring test coverage",
		},
		"reviewer": {
			"  - name: Static Analysis Tools",
			"    description: Tools for code analysis",
			"  - name: Linter",
			"    description: Code quality tool",
		},
		"debugger": {
			"  - name: Debugger",
			"    description: Debugging tool",
			"  - name: Profiler",
			"    description: Performance profiling tool",
		},
	}

	if tools, ok := extraTools[agentType]; ok {
		baseTools = append(baseTools, tools...)
	}

	return strings.Join(baseTools, "\n")
}

// buildAgentYAMLTemplate construye un prompt para generar un agente en formato YAML.
// Incluye la guía completa de creación de agentes como contexto.
func (g *Generator) buildAgentYAMLTemplate(agentType string) string {
	// Leer la guía completa de agentes
	agentGuide := g.getAgentGuide()

	roleDesc := g.getAgentRoleDescription(agentType)
	responsibilities := g.getAgentResponsibilities(agentType)
	guidelines := g.getAgentGuidelines(agentType)
	tools := g.getAgentTools(agentType)

	prompt := fmt.Sprintf(`Generate a comprehensive agent configuration file for a %s agent for a project called %s.

## PROJECT DETAILS
- Language: %s
- Framework: %s
- Architecture: %s
- Description: %s

## AGENT SPECIFICATION
- Role: %s
- Responsibilities:
%s
- Coding Guidelines:
%s
- Tools:
%s

## CRITICAL: AGENT CREATION GUIDE

You MUST follow this guide to create the agent. Read it carefully and apply ALL principles:

%s

## OUTPUT FORMAT

Generate the agent file following the template structure from the guide above. The output should be a complete, production-ready agent configuration that:

1. Follows the "Template Oficial de Agentes" structure from the guide
2. Implements the "Bucle Operativo (Agent Loop)" with all 4 phases
3. Defines "Capacidades Inyectadas" (Skills and Tools) properly
4. Includes "Estrategia de Toma de Decisiones" with examples
5. Documents "Reglas de Oro" and "Restricciones y Políticas"
6. Provides an "Invocación de Ejemplo" with expected output

Remember the core principle from the guide:
🧠 **Agente = Razonamiento puro** - Sin conocimiento técnico hardcodeado
📚 **Skills = Conocimiento inyectado** - Convenciones, frameworks, lenguajes
🛠️ **Tools = Capacidad de acción** - "Darle un ordenador al agente"

CRITICAL: The agent MUST be agnostic to specific technologies. All technical knowledge comes from injected skills, NOT from the agent definition itself.`,
		agentType,
		g.answers.ProjectName,
		g.answers.Language,
		g.answers.Framework,
		g.answers.Architecture,
		g.answers.Description,
		roleDesc,
		responsibilities,
		guidelines,
		tools,
		agentGuide,
	)

	return prompt
}

// getAgentGuide retorna la guía completa de creación de agentes.
// El contenido está embebido en el binario compilado usando go:embed.
func (g *Generator) getAgentGuide() string {
	// El contenido está embebido en la variable agentGuideContent
	// Si por alguna razón está vacía, retornar una guía básica
	if agentGuideContent == "" {
		g.logger.Warn("El contenido de la guía de agentes está vacío")
		return `# Guía Básica de Agentes (Contenido de respaldo)

## Principios Fundamentales
- El agente NO posee conocimiento técnico hardcodeado
- Todo conocimiento técnico viene de skills inyectadas
- El agente aplica razonamiento estructurado sobre el contexto

## Estructura Mínima Requerida
1. Perfil de Razonamiento (rol, principios, objetivo)
2. Bucle Operativo (4 fases: Recopilar, Actuar, Verificar, Iterar)
3. Capacidades Inyectadas (skills y tools)
4. Estrategia de Toma de Decisiones
5. Reglas de Oro (No alucinar, Verificación empírica, Trazabilidad)
`
	}

	return agentGuideContent
}

// buildCommandTemplate construye un prompt para generar un comando.
// Incluye la guía completa de creación de comandos como contexto.
func (g *Generator) buildCommandTemplate(commandType string) string {
	// Leer la guía completa de comandos
	commandGuide := g.getCommandGuide()

	cmdDesc := g.getCommandDescription(commandType)
	cmdUsage := g.getCommandUsage(commandType)
	cmdFlow := g.getCommandFlow(commandType)

	prompt := fmt.Sprintf(`Generate a comprehensive command configuration file for a %s command for a project called %s.

## PROJECT DETAILS
- Language: %s
- Framework: %s
- Architecture: %s
- Description: %s

## COMMAND SPECIFICATION
- Description: %s
- Usage: %s
- Flow:
%s

## CRITICAL: COMMAND CREATION GUIDE

You MUST follow this guide to create the command. Read it carefully and apply ALL principles:

%s

## OUTPUT FORMAT

Generate the command file following the template structure from the guide above. The output should be a complete, production-ready command configuration that:

1. Follows the "Template Oficial de Comandos" structure from the guide
2. Implements the "Flujo de Implementación Orquestado" with all agents and phases
3. Includes "Reglas Críticas" and validation steps
4. Documents the orquestation process clearly
5. Specifies which agents are involved and their responsibilities

Remember the core principles from the guide:
- **Planificación**: Analizar y desglosar la tarea
- **Ejecución**: Coordinar los agentes apropiados
- **Verificación**: Asegurar calidad y no regresiones
- **Documentación**: Registrar cambios y decisiones

CRITICAL: The command must orchestrate agents properly, defining clear responsibilities and handoff points between agents.`,
		commandType,
		g.answers.ProjectName,
		g.answers.Language,
		g.answers.Framework,
		g.answers.Architecture,
		g.answers.Description,
		cmdDesc,
		cmdUsage,
		cmdFlow,
		commandGuide,
	)

	return prompt
}

// buildCommandTemplateWithContext construye un prompt para generar un comando con contexto.
// Similar a buildCommandTemplate pero incluye información sobre agents y skills disponibles.
// Este contexto permite que la IA referencie agentes y skills que realmente existen.
func (g *Generator) buildCommandTemplateWithContext(commandType, agentsContext, skillsContext string) string {
	commandGuide := g.getCommandGuide()

	cmdDesc := g.getCommandDescription(commandType)
	cmdUsage := g.getCommandUsage(commandType)
	cmdFlow := g.getCommandFlow(commandType)

	prompt := fmt.Sprintf(`Generate a comprehensive command configuration file for a %s command for a project called %s.

## PROJECT DETAILS
- Language: %s
- Framework: %s
- Architecture: %s
- Description: %s

## COMMAND SPECIFICATION
- Description: %s
- Usage: %s
- Flow:
%s

%s

## AVAILABLE AGENTS AND SKILLS

Below are the agents and skills that have been created for this project.
You MUST reference ONLY these agents and skills when creating this command.

### AVAILABLE AGENTS
%s

### AVAILABLE SKILLS
%s

## HOW TO SELECT THE RIGHT AGENTS AND SKILLS

**IMPORTANT**: You must ANALYZE the "AVAILABLE AGENTS" and "AVAILABLE SKILLS" lists above.
Look at each agent's DESCRIPTION and each skill's PURPOSE to determine which ones are best suited for this command.

### Selection Process

1. **READ the descriptions** of all available agents
2. **IDENTIFY which agent's role matches** this command's primary goal
3. **READ the purposes** of all available skills
4. **SELECT skills that complement** the chosen agent for this specific task
5. **DO NOT default to the same agent/skill** - each command is unique

### Matching Examples (Apply this pattern to YOUR specific agents/skills)

**For a "test" command**: Look for an agent whose description mentions testing, quality, verification, or validation. Select skills that provide testing expertise.

**For a "lint" or "review" command**: Look for an agent whose description mentions code review, quality, validation, or standards. Select skills related to linting or static analysis.

**For a "bug-fix" command**: Look for an agent whose description mentions debugging, investigation, or problem-solving. Select skills that provide debugging techniques.

**For a "new-feature" command**: Look for agents that handle planning/design AND implementation. You may need multiple agents. Select skills for architecture, coding, and integration.

**For a "refactor" command**: Look for an agent whose description mentions code improvement, cleanup, or optimization. Select skills for code quality and design patterns.

**For a "config" command**: Look for agents that handle configuration or setup. Select skills for configuration management.

## CRITICAL REQUIREMENTS

When creating this command:

1. **ANALYZE THE AVAILABLE AGENTS**: Read each agent's description in the list above. Find the agent whose DESCRIPTION best matches this command's purpose.

2. **ANALYZE THE AVAILABLE SKILLS**: Read each skill's purpose in the list above. Select skills that PROVIDE what this command needs.

3. **MATCH BY DESCRIPTION/PURPOSE**:
   - This command is about: [describe the command's goal]
   - The best agent for this is: [agent whose description matches]
   - Because its description says: [quote relevant part of description]
   - The relevant skills are: [skills whose purposes match]

4. **VARY YOUR SELECTION**: Each command should use DIFFERENT agents and skills based on what that command specifically does.

5. **REFERENCE ONLY THE LIST**: Use ONLY agents and skills from "AVAILABLE AGENTS" and "AVAILABLE SKILLS" above.

Example format when creating your command:

---
name: [command-name]
description: [your description]

# Command: [Command Name]

## Implementation Flow

### 1. [Phase Name]
- **Agent**: [agent-name-from-above] - selected because its description mentions "[relevant-description]"
- **Skills**:
    - [skill-from-above]: provides "[relevant-purpose]"
    - [skill-from-above]: provides "[relevant-purpose]"
- [Action to perform]

Remember the core principles from the guide:
- **Planificación**: Analizar y desglosar la tarea
- **Ejecución**: Coordinar los agentes apropiados
- **Verificación**: Asegurar calidad y no regresiones
- **Documentación**: Registrar cambios y decisiones

**CRITICAL**: You MUST READ the agent descriptions and skill purposes above. Do NOT pick randomly - pick based on WHAT THE AGENT/SKILL DOES according to its description/purpose!`,
		commandType,
		g.answers.ProjectName,
		g.answers.Language,
		g.answers.Framework,
		g.answers.Architecture,
		g.answers.Description,
		cmdDesc,
		cmdUsage,
		cmdFlow,
		commandGuide,
		func() string {
			if agentsContext != "" {
				return "\n## AVAILABLE AGENTS\n\n" + agentsContext
			}
			return ""
		}(),
		func() string {
			if skillsContext != "" {
				return "\n### AVAILABLE SKILLS\n\n" + skillsContext
			}
			return ""
		}(),
	)

	// Agregar instrucciones específicas para OpenAI
	if g.isOpenAI() {
		prompt += `

## OPENAI SPECIFIC INSTRUCTIONS

You MUST generate a comprehensive, detailed command configuration:

1. **BE CREATIVE**: Don't default to the same agent/skill combinations
2. **BE SPECIFIC**: Include concrete examples in the Usage section
3. **BE THOROUGH**: Generate at least 4-6 workflow phases
4. **USE ALL AVAILABLE AGENTS**: Don't limit yourself to 2-3 agents
5. **INCLUDE DIAGRAMS**: Add Mermaid diagrams when appropriate
6. **ADD ROLLBACK PROCEDURES**: Include detailed rollback steps

The user expects a RICH, DETAILED configuration similar to what Claude would generate.

Minimum requirements:
- At least 4 different agents
- At least 4 distinct workflow phases
- Concrete usage examples
- Mermaid workflow diagram
- Detailed rollback procedures
`
	}

	return prompt
}

// isOpenAI retorna true si el cliente actual es OpenAI.
func (g *Generator) isOpenAI() bool {
	return g.client.Provider() == ai.ProviderOpenAI
}

// getCommandGuide retorna la guía completa de creación de comandos.
// El contenido está embebido en el binario compilado usando go:embed.
func (g *Generator) getCommandGuide() string {
	// El contenido está embebido en la variable commandGuideContent
	// Si por alguna razón está vacía, retornar una guía básica
	if commandGuideContent == "" {
		g.logger.Warn("El contenido de la guía de comandos está vacío")
		return `# Guía Básica de Comandos (Contenido de respaldo)

## Principios Fundamentales
- Los comandos orquestan uno o más agentes
- Definen un flujo claro de implementación
- Incluyen reglas críticas y verificación

## Estructura Mínima Requerida
1. Descripción del comando
2. Flujo de implementación orquestado
3. Reglas críticas
4. Agentes involucrados
`
	}

	return commandGuideContent
}

// buildSkillTemplate construye un prompt para generar una skill.
// Incluye la guía completa de creación de skills como contexto.
func (g *Generator) buildSkillTemplate(skillType, skillName string) string {
	// Leer la guía completa de skills
	skillGuide := g.getSkillGuide()

	skillDesc := g.getSkillDescription(skillName)
	skillTitle := g.getSkillTitle(skillName)

	prompt := fmt.Sprintf(`Generate a comprehensive skill configuration file for a %s skill called %s for a project called %s.

## PROJECT DETAILS
- Language: %s
- Framework: %s
- Architecture: %s
- Description: %s

## SKILL SPECIFICATION
- Skill Type: %s
- Description: %s
- Title: %s

## CRITICAL: SKILL CREATION GUIDE

You MUST follow this guide to create the skill. Read it carefully and apply ALL principles:

%s

## OUTPUT FORMAT

Generate the skill file following the template structure from the guide above. The output should be a complete, production-ready skill configuration that:

1. Follows the "Template Oficial de Skills" structure from the guide
2. Includes proper frontmatter with name and description
3. Documents "How It Works" with clear phases
4. Provides "Usage" examples and trigger phrases
5. Lists "Capabilities" with best practices and patterns
6. Includes "Output Examples" and "Troubleshooting" sections

Remember the core principles from the guide:
- **Skills = Conocimiento Inyectado**: Technical expertise externalized from agents
- **Declarativo, no Imperativo**: Define WHAT and WHY, not HOW
- **Específico al Dominio**: Focused on particular technologies or patterns
- **Autocontenido**: Complete and independently usable

CRITICAL: The skill must be domain-specific knowledge that agents can inject, not procedural instructions.`,
		skillType,
		skillName,
		g.answers.ProjectName,
		g.answers.Language,
		g.answers.Framework,
		g.answers.Architecture,
		g.answers.Description,
		skillType,
		skillDesc,
		skillTitle,
		skillGuide,
	)

	return prompt
}

// getSkillGuide retorna la guía completa de creación de skills.
// El contenido está embebido en el binario compilado usando go:embed.
func (g *Generator) getSkillGuide() string {
	// El contenido está embebido en la variable skillGuideContent
	// Si por alguna razón está vacía, retornar una guía básica
	if skillGuideContent == "" {
		g.logger.Warn("El contenido de la guía de skills está vacío")
		return `# Guía Básica de Skills (Contenido de respaldo)

## Principios Fundamentales
- Skills = Conocimiento técnico inyectado en agentes
- Declarativo, no imperativo (define qué, no cómo)
- Específico al dominio (lenguajes, frameworks, patrones)
- Autocontenido y reutilizable

## Estructura Mínima Requerida
1. Frontmatter con nombre y descripción
2. Descripción detallada del skill
3. Cómo funciona (How It Works)
4. Ejemplos de uso y trigger phrases
5. Capacidades y mejores prácticas
`
	}

	return skillGuideContent
}

// capitalize capitaliza la primera letra de un string.
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// getEmbeddedSkillTemplate retorna un template markdown incrustado para un skill.
func (g *Generator) getEmbeddedSkillTemplate(skillName, skillType string) string {
	skillDesc := g.getSkillDescription(skillName)
	skillTitle := g.getSkillTitle(skillName)
	triggers := g.getSkillTriggers(skillName)

	return fmt.Sprintf(`---
name: %s
description: %s
---

# %s

%s

## How It Works

1. The agent identifies the need to use this %s skill.
2. The specific context and requirements are analyzed.
3. Appropriate %s patterns and best practices are applied.

## Usage

This skill is automatically activated when working with %s in the %s project.

**Trigger phrases:**
%s

## Capabilities

- Language-specific best practices for %s
- Framework-specific patterns and conventions
- Code optimization techniques
- Testing strategies
- Common pitfalls and solutions

## Output Examples

%s

## Present Results to User

When applying this skill, present results in a clear, structured format:
- Brief summary of what was done
- Key changes or recommendations
- Any relevant code snippets or examples
- Next steps or considerations

## Troubleshooting

**Common issues:**
- **Incorrect syntax**: Follow %s conventions and best practices
- **Missing dependencies**: Ensure all required packages are installed
- **Type errors**: Check type annotations and interfaces
- **Performance issues**: Review for optimization opportunities`,
		skillName,
		skillDesc,
		skillTitle,
		skillDesc,
		skillName,
		skillName,
		skillName,
		g.answers.ProjectName,
		triggers,
		skillName,
		g.getSkillOutputExamples(skillName),
		g.answers.Language,
	)
}

// getSkillDescription retorna la descripción para un skill.
func (g *Generator) getSkillDescription(skillName string) string {
	descriptions := map[string]string{
		"technical-writer": "Assist with technical writing, documentation, and code comments",
		"code-reviewer":    "Review code for quality, best practices, and potential issues",
		"debug-master":     "Debug complex issues, analyze errors, and propose solutions",
		"typescript":       "Optimize TypeScript performance and configure TypeScript projects",
		"go":               "Optimize Go performance and configure Go projects",
		"python":           "Optimize Python performance and configure Python projects",
		"javascript":       "Optimize JavaScript performance and configure JavaScript projects",
	}

	if desc, ok := descriptions[skillName]; ok {
		return desc
	}

	// Descripción genérica basada en el lenguaje/framework
	if strings.EqualFold(skillName, g.answers.Language) {
		return fmt.Sprintf("Optimize %s performance and configure %s projects", g.answers.Language, g.answers.Language)
	}

	if strings.EqualFold(skillName, g.answers.Framework) {
		return fmt.Sprintf("Configure and optimize %s framework components", g.answers.Framework)
	}

	return fmt.Sprintf("Assist with %s-related tasks and configurations", skillName)
}

// getSkillTitle retorna el título para un skill.
func (g *Generator) getSkillTitle(skillName string) string {
	titles := map[string]string{
		"technical-writer": "Technical Writing Skill",
		"code-reviewer":    "Code Review Skill",
		"debug-master":     "Debugging Skill",
	}

	if title, ok := titles[skillName]; ok {
		return title
	}

	return fmt.Sprintf("%s Skill", capitalize(skillName))
}

// getSkillTriggers retorna las frases trigger para un skill.
func (g *Generator) getSkillTriggers(skillName string) string {
	triggers := map[string]string{
		"technical-writer": `- "Write documentation for..."
- "Add comments to..."
- "Explain how this works..."
- "Create a guide for..."`,
		"code-reviewer": `- "Review this code..."
- "Check for best practices..."
- "Identify potential issues..."
- "Suggest improvements..."`,
		"debug-master": `- "Debug this error..."
- "Fix this bug..."
- "Why is this not working..."
- "Analyze this stack trace..."`,
	}

	if trigger, ok := triggers[skillName]; ok {
		return trigger
	}

	// Triggers genéricos basados en el skill
	return fmt.Sprintf(`- "Optimize %s..."
- "Configure %s..."
- "Fix %s issue..."
- "Improve %s..."`,
		skillName, skillName, skillName, skillName,
	)
}

// getSkillOutputExamples retorna ejemplos de output para un skill.
func (g *Generator) getSkillOutputExamples(skillName string) string {
	examples := map[string]string{
		"technical-writer": `**Documentation Example:**

%smarkdown
## Function Name

Brief description of what the function does.

### Parameters
- param1: Description
- param2: Description

### Returns
Description of return value.

### Example
%s%s
// Example code
%s
`,
		"code-reviewer": `**Review Summary:**
- Code follows %s best practices
- Consider handling edge cases in...
- Suggestion: Extract this logic into a separate function`,
		"debug-master": `**Analysis:**
The issue is caused by...
**Solution:**
1. Fix the root cause by...
2. Add error handling for...
3. Test with...`,
	}

	if example, ok := examples[skillName]; ok {
		if skillName == "technical-writer" {
			return fmt.Sprintf(example, "```", "```", g.answers.Language, "```", "```")
		}
		return fmt.Sprintf(example, g.answers.Language)
	}

	return fmt.Sprintf(`**Example Output:**
Applied %s best practices to improve code quality and performance.`, skillName)
}

// getEmbeddedCommandTemplate retorna un template incrustado para un comando.
func (g *Generator) getEmbeddedCommandTemplate(commandType string) string {
	cmdDesc := g.getCommandDescription(commandType)
	cmdUsage := g.getCommandUsage(commandType)
	cmdFlow := g.getCommandFlow(commandType)

	return fmt.Sprintf(`---
name: %s
description: %s
usage: %s
---

# Comando: %s

%s

## Flujo de Implementación Orquestado

%s

## Reglas Críticas
- **Calidad**: Mantener altos estándares de calidad en el código.
- **Testing**: Asegurar cobertura de pruebas adecuada.
- **Documentación**: Documentar cambios y decisiones.
- **Colaboración**: Coordinar con otros agentes cuando sea necesario.

---

¿Qué %s deseas realizar en %s?`,
		commandType,
		cmdDesc,
		cmdUsage,
		capitalize(commandType),
		cmdDesc,
		cmdFlow,
		commandType,
		g.answers.ProjectName,
	)
}

// getCommandDescription retorna la descripción para un comando.
func (g *Generator) getCommandDescription(commandType string) string {
	descriptions := map[string]string{
		"test":        fmt.Sprintf("Ejecuta tests para %s y reporta resultados", g.answers.ProjectName),
		"lint":        fmt.Sprintf("Ejecuta linters para verificar calidad de código en %s", g.answers.ProjectName),
		"build":       fmt.Sprintf("Compila %s y genera artefactos de deploy", g.answers.ProjectName),
		"new-feature": fmt.Sprintf("Planifica e implementa una nueva funcionalidad en %s", g.answers.ProjectName),
		"refactor":    fmt.Sprintf("Refactoriza código en %s mejorando su estructura", g.answers.ProjectName),
		"bug-fix":     fmt.Sprintf("Investiga y corrige bugs en %s", g.answers.ProjectName),
	}

	if desc, ok := descriptions[commandType]; ok {
		return desc
	}
	return fmt.Sprintf("Comando para %s en %s", commandType, g.answers.ProjectName)
}

// getCommandUsage retorna el usage para un comando.
func (g *Generator) getCommandUsage(commandType string) string {
	usages := map[string]string{
		"test":        "test [nombre-test]",
		"lint":        "lint [archivo-o-directorio]",
		"build":       "build [entorno]",
		"new-feature": "new-feature [nombre-funcionalidad] [descripcion]",
		"refactor":    "refactor [archivo-o-componente]",
		"bug-fix":     "bug-fix [descripcion-del-error]",
	}

	if usage, ok := usages[commandType]; ok {
		return usage
	}
	return commandType
}

// getCommandFlow retorna el flujo para un comando.
func (g *Generator) getCommandFlow(commandType string) string {
	langSkill := strings.ToLower(g.answers.Language)

	switch commandType {
	case "test":
		return fmt.Sprintf(`### 1. Ejecutar Tests
- **Agente**: tester
- **Skills**: %s-expert, qa-engineer
- Ejecutar suite de tests completa o específica

### 2. Análisis de Resultados
- Revisar tests fallidos
- Identificar patrones de errores
- Proponer correcciones

### 3. Reporte
- Generar reporte de cobertura
- Listar tests fallidos con errores
- Sugerir mejoras`, langSkill)

	case "lint":
		return fmt.Sprintf(`### 1. Ejecutar Linters
- **Agente**: reviewer
- **Skills**: code-reviewer, %s-expert
- Ejecutar linters del proyecto

### 2. Análisis de Problemas
- Revisar advertencias y errores
- Clasificar por severidad
- Priorizar correcciones

### 3. Corrección Automática
- Aplicar auto-fix cuando sea posible
- Generar reporte con sugerencias`, langSkill)

	case "build":
		return `### 1. Preparación
- **Agente**: developer
- Verificar dependencias
- Limpiar artefactos previos

### 2. Compilación
- Compilar para el entorno objetivo
- Verificar que no haya errores de compilación
- Generar artefactos

### 3. Validación
- Ejecutar tests básicos
- Verificar tamaño de artefactos
- Generar reporte de build`

	case "new-feature":
		return fmt.Sprintf(`### 1. Planificación
- **Agente**: planning-agent
- **Skills**: system-architect, %s-expert
- Analizar requerimientos
- Crear plan de implementación

### 2. Desarrollo (TDD)
- **Agente**: developer
- **Skills**: %s-expert, usecase-developer
- Escribir tests primero (TDD)
- Implementar funcionalidad mínima

### 3. Integración
- **Agente**: tester, reviewer
- **Skills**: qa-engineer, code-reviewer
- Tests de integración
- Revisión de código

### 4. Documentación
- **Agente**: writer
- **Skills**: technical-writer
- Actualizar documentación
- Registrar cambios`, langSkill, langSkill)

	case "refactor":
		return fmt.Sprintf(`### 1. Análisis
- **Agente**: reviewer
- **Skills**: code-reviewer, %s-expert
- Identificar code smells
- Proponer mejoras arquitectónicas

### 2. Refactorización
- **Agente**: developer
- **Skills**: %s-expert
- Aplicar refactorizaciones paso a paso
- Mantener tests pasando

### 3. Validación
- **Agente**: tester
- Ejecutar suite completa de tests
- Verificar que no haya regresiones`, langSkill, langSkill)

	case "bug-fix":
		return fmt.Sprintf(`### 1. Investigación
- **Agente**: debugger
- **Skills**: debug-master, %s-expert
- Reproducir el error
- Analizar stack traces y logs

### 2. Diagnóstico
- Identificar causa raíz
- Proponer hipótesis
- Validar con tests

### 3. Solución
- **Agente**: developer
- Implementar fix
- Verificar que tests pasen
- Prevenir regresiones`, langSkill)

	default:
		return fmt.Sprintf(`### 1. Ejecución
- **Agente**: developer
- **Skills**: %s-expert
- Ejecutar comando %s

### 2. Validación
- Verificar resultados
- Reportar estado`, langSkill, commandType)
	}
}

// GenerateDevelopmentGuide genera el archivo development_guide.md.
func (g *Generator) GenerateDevelopmentGuide() error {
	outputPath := filepath.Join(g.projectPath, ".claude", "development_guide.md")

	// Verificar si ya existe y no estamos en modo force
	if _, err := os.Stat(outputPath); err == nil {
		// Ya existe, no sobrescribir
		return nil
	}

	content, err := g.getDevelopmentGuideTemplate()
	if err != nil {
		return fmt.Errorf("error generando development_guide.md: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo development_guide.md: %w", err)
	}

	g.logger.Info("Development guide generado en %s", outputPath)
	return nil
}

// GenerateDevelopmentGuideWithContext genera el archivo development_guide.md con contexto de la estructura completa.
// A diferencia de GenerateDevelopmentGuide, esta versión acepta información sobre agents, commands y skills
// para incluir en la guía de desarrollo.
func (g *Generator) GenerateDevelopmentGuideWithContext(agentNames, commandNames, skillNames []string) error {
	outputPath := filepath.Join(g.projectPath, ".claude", "development_guide.md")

	// Generar contenido con contexto
	content, err := g.getDevelopmentGuideTemplateWithContext(agentNames, commandNames, skillNames)
	if err != nil {
		return fmt.Errorf("error generando development_guide.md: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo development_guide.md: %w", err)
	}

	g.logger.Info("Development guide con contexto generado en %s", outputPath)
	return nil
}

// GenerateClaudeMD genera el archivo CLAUDE.md con contexto del proyecto.
// Este archivo se genera PRIMERO para proporcionar contexto a las generaciones posteriores.
func (g *Generator) GenerateClaudeMD() error {
	outputPath := filepath.Join(g.projectPath, "CLAUDE.md")

	// Verificar si ya existe
	if _, err := os.Stat(outputPath); err == nil {
		g.logger.Debug("CLAUDE.md ya existe, no se sobrescribe")
		return nil
	}

	content, err := g.getClaudeMDTemplate()
	if err != nil {
		return fmt.Errorf("error generando CLAUDE.md: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo CLAUDE.md: %w", err)
	}

	g.logger.Info("CLAUDE.md generado en %s", outputPath)
	return nil
}

// getClaudeMDTemplate genera el contenido del CLAUDE.md usando IA.
// Analiza el proyecto real para obtener información precisa.
func (g *Generator) getClaudeMDTemplate() (string, error) {
	// Analizar el proyecto para obtener información real
	projectContext := g.analyzeProjectContext()

	prompt := fmt.Sprintf(`Genera un archivo CLAUDE.md completo y detallado para el siguiente proyecto.

Este archivo será usado por Claude Code (o similar) para entender el contexto del proyecto y proporcionar asistencia precisa. El formato debe ser similar al que Claude Code genera nativamente.

**Información básica del proyecto:**
- **Nombre:** %s
- **Descripción:** %s
- **Lenguaje principal:** %s
- **Framework:** %s

**Información adicional del análisis del proyecto:**
%s

IMPORTANTE: Genera un CLAUDE.md que siga este formato y nivel de detalle:

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
[Descripción concisa pero completa del proyecto, su propósito y objetivos principales]

## Tech Stack
- **Framework**: [Framework y versión específicas]
- **Language**: [Lenguaje y versión]
- **Package Manager**: [npm/yarn/pnpm - MUY IMPORTANTE especificar cuál se debe usar]
- [Otras tecnologías importantes con versiones si es posible]

## Essential Commands

### Development
[Listado de comandos para desarrollo con bloques de código]

### Building
[Listado de comandos de build con bloques de código]

### Testing
[Comandos para testing con bloques de código]

### Other
[Otros comandos útiles]

## Architecture

### Directory Structure
[Estructura de directorios REAL del proyecto, no genérica]

### Key Architectural Concepts
[Conceptos arquitectónicos importantes específicos del proyecto]

## Import Path Aliases
[Si aplica, listar los alias de importación de tsconfig.json o similares]

## Code Quality
[Información sobre ESLint, Prettier, y otras herramientas de calidad]

## Environment Setup
[Instrucciones de configuración del entorno]

## Component/Module Guidelines
[Guías específicas para el desarrollo de componentes/módulos]

Genera el contenido completo en markdown, específico y detallado basado en la información del proyecto. NO uses placeholders como "..." o comandos genéricos. Si hay información específica disponible (como los scripts de package.json), ÚSALA.`,
		g.answers.ProjectName,
		g.answers.Description,
		g.answers.Language,
		g.answers.Framework,
		projectContext,
	)

	return g.generateWithClaude(prompt, nil)
}

// analyzeProjectContext analiza el proyecto real para extraer información precisa.
func (g *Generator) analyzeProjectContext() string {
	var context strings.Builder

	// Analizar package.json si existe
	if pkgInfo := g.analyzePackageJSON(); pkgInfo != "" {
		context.WriteString("**Package.json info:**\n")
		context.WriteString(pkgInfo)
		context.WriteString("\n\n")
	}

	// Analizar tsconfig.json si existe
	if tsconfigInfo := g.analyzeTsConfig(); tsconfigInfo != "" {
		context.WriteString("**TypeScript config:**\n")
		context.WriteString(tsconfigInfo)
		context.WriteString("\n\n")
	}

	// Analizar estructura de directorios
	if dirInfo := g.analyzeDirectoryStructure(); dirInfo != "" {
		context.WriteString("**Directory structure:**\n")
		context.WriteString(dirInfo)
		context.WriteString("\n\n")
	}

	// Analizar archivos de configuración de code quality
	if qualityInfo := g.analyzeCodeQualityConfigs(); qualityInfo != "" {
		context.WriteString("**Code quality tools:**\n")
		context.WriteString(qualityInfo)
		context.WriteString("\n\n")
	}

	// Analizar directorios de documentación
	if docInfo := g.analyzeDocumentation(); docInfo != "" {
		context.WriteString("**Documentation:**\n")
		context.WriteString(docInfo)
		context.WriteString("\n\n")
	}

	result := context.String()
	if result == "" {
		return "No se pudo extraer información adicional del proyecto."
	}
	return result
}

// analyzePackageJSON analiza el package.json si existe.
func (g *Generator) analyzePackageJSON() string {
	pkgPath := filepath.Join(g.projectPath, "package.json")
	content, err := os.ReadFile(pkgPath)
	if err != nil {
		return ""
	}

	var pkg struct {
		Scripts         map[string]string `json:"scripts"`
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		Name            string            `json:"name"`
		Version         string            `json:"version"`
		Description     string            `json:"description"`
	}

	if err := json.Unmarshal(content, &pkg); err != nil {
		return ""
	}

	var info strings.Builder

	// Scripts más importantes
	importantScripts := []string{"dev", "start", "build", "test", "lint", "type-check", "watch"}
	info.WriteString("- **Scripts:**\n")
	for _, key := range importantScripts {
		if script, ok := pkg.Scripts[key]; ok {
			info.WriteString(fmt.Sprintf("  - `%s`: %s\n", key, script))
		}
	}

	// Detectar package manager
	if _, hasYarnLock := os.ReadFile(filepath.Join(g.projectPath, "yarn.lock")); hasYarnLock == nil {
		info.WriteString("\n- **Package Manager**: yarn (SIEMPRE usar yarn, no npm)\n")
	} else if _, hasPNPMLock := os.ReadFile(filepath.Join(g.projectPath, "pnpm-lock.yaml")); hasPNPMLock == nil {
		info.WriteString("\n- **Package Manager**: pnpm\n")
	} else {
		info.WriteString("\n- **Package Manager**: npm\n")
	}

	// Dependencias principales
	if len(pkg.Dependencies) > 0 {
		info.WriteString("\n- **Main dependencies:**\n")
		for name, version := range pkg.Dependencies {
			// Mostrar solo las más importantes
			if isImportantDependency(name) {
				info.WriteString(fmt.Sprintf("  - %s: %s\n", name, version))
			}
		}
	}

	// DevDependencies importantes
	if len(pkg.DevDependencies) > 0 {
		info.WriteString("\n- **Dev dependencies:**\n")
		for name, version := range pkg.DevDependencies {
			if isImportantDevDependency(name) {
				info.WriteString(fmt.Sprintf("  - %s: %s\n", name, version))
			}
		}
	}

	return info.String()
}

// analyzeTsConfig analiza el tsconfig.json si existe.
func (g *Generator) analyzeTsConfig() string {
	tsconfigPath := filepath.Join(g.projectPath, "tsconfig.json")
	content, err := os.ReadFile(tsconfigPath)
	if err != nil {
		return ""
	}

	var tsconfig map[string]interface{}
	if err := json.Unmarshal(content, &tsconfig); err != nil {
		return ""
	}

	var info strings.Builder

	// Compiler options
	if compilerOpts, ok := tsconfig["compilerOptions"].(map[string]interface{}); ok {
		info.WriteString("- **TypeScript Configuration:**\n")

		if target, ok := compilerOpts["target"].(string); ok {
			info.WriteString(fmt.Sprintf("  - target: %s\n", target))
		}

		if strict, ok := compilerOpts["strict"].(bool); ok && strict {
			info.WriteString("  - strict mode: enabled\n")
		}

		// Paths (aliases)
		if paths, ok := compilerOpts["paths"].(map[string]interface{}); ok && len(paths) > 0 {
			info.WriteString("\n- **Import Path Aliases:**\n")
			for alias, path := range paths {
				info.WriteString(fmt.Sprintf("  - %s: %v\n", alias, path))
			}
		}
	}

	return info.String()
}

// analyzeDirectoryStructure analiza la estructura de directorios del proyecto.
func (g *Generator) analyzeDirectoryStructure() string {
	// Directorios comunes a analizar
	commonDirs := []string{"src", "app", "components", "lib", "utils", "hooks", "types", "tests", "__tests__", "test", "spec", "styles", "assets", "public", "dist", "build", "config", "scripts"}

	var info strings.Builder
	info.WriteString("```\n")

	for _, dir := range commonDirs {
		fullPath := filepath.Join(g.projectPath, dir)
		if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
			info.WriteString(dir + "/\n")

			// Listar subdirectorios principales
			if entries, err := os.ReadDir(fullPath); err == nil {
				for _, entry := range entries {
					if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
						info.WriteString(fmt.Sprintf("  └── %s/\n", entry.Name()))
					}
				}
			}
		}
	}

	info.WriteString("```")
	return info.String()
}

// analyzeCodeQualityConfigs analiza configuraciones de calidad de código.
func (g *Generator) analyzeCodeQualityConfigs() string {
	var info strings.Builder

	// ESLint
	eslintPath := filepath.Join(g.projectPath, "eslint.config.mjs")
	if _, err := os.Stat(eslintPath); err == nil {
		info.WriteString("- ESLint: Modern flat config (eslint.config.mjs)\n")
	} else {
		if _, err := os.Stat(filepath.Join(g.projectPath, ".eslintrc.json")); err == nil {
			info.WriteString("- ESLint: .eslintrc.json\n")
		} else if _, err := os.Stat(filepath.Join(g.projectPath, ".eslintrc.js")); err == nil {
			info.WriteString("- ESLint: .eslintrc.js\n")
		}
	}

	// Prettier
	if _, err := os.Stat(filepath.Join(g.projectPath, ".prettierrc")); err == nil {
		info.WriteString("- Prettier: Configurado\n")
	} else if _, err := os.Stat(filepath.Join(g.projectPath, ".prettierrc.json")); err == nil {
		info.WriteString("- Prettier: Configurado (.prettierrc.json)\n")
	}

	// Husky
	if _, err := os.Stat(filepath.Join(g.projectPath, ".husky")); err == nil {
		info.WriteString("- Husky: Git hooks configurados\n")
	}

	return info.String()
}

// analyzeDocumentation analiza directorios de documentación del proyecto.
func (g *Generator) analyzeDocumentation() string {
	var info strings.Builder

	// Directorios comunes de documentación a buscar
	commonDocDirs := []string{"docs", "documentation", "guide", "guides", "wiki", "help"}

	// Recopilar todos los directorios de documentación a analizar
	dirsToAnalyze := make([]string, 0, len(commonDocDirs)+len(g.answers.DocumentationDirs))
	dirsToAnalyze = append(dirsToAnalyze, commonDocDirs...)
	dirsToAnalyze = append(dirsToAnalyze, g.answers.DocumentationDirs...)

	foundAny := false

	for _, docDir := range dirsToAnalyze {
		if docDir == "" {
			continue
		}

		docPath := filepath.Join(g.projectPath, docDir)
		if stat, err := os.Stat(docPath); err != nil || !stat.IsDir() {
			continue
		}

		if !foundAny {
			foundAny = true
			info.WriteString(fmt.Sprintf("- Directorio '%s/' encontrado:\n", docDir))
		} else {
			info.WriteString(fmt.Sprintf("\n- Directorio '%s/' encontrado:\n", docDir))
		}

		// Listar archivos de documentación importantes en este directorio
		if entries, err := os.ReadDir(docPath); err == nil {
			var mdFiles []string
			for _, entry := range entries {
				name := entry.Name()
				if !entry.IsDir() && (strings.HasSuffix(name, ".md") || strings.HasSuffix(name, ".txt")) {
					mdFiles = append(mdFiles, name)
				}
			}

			if len(mdFiles) > 0 {
				info.WriteString("  Archivos de documentación:\n")
				for _, file := range mdFiles {
					// Leer las primeras líneas del archivo para obtener un resumen
					filePath := filepath.Join(docPath, file)
					if summary := g.getDocFileSummary(filePath); summary != "" {
						info.WriteString(fmt.Sprintf("  - %s: %s\n", file, summary))
					} else {
						info.WriteString(fmt.Sprintf("  - %s\n", file))
					}
				}
			}
		}
	}

	// También buscar README.md en la raíz
	if _, err := os.Stat(filepath.Join(g.projectPath, "README.md")); err == nil {
		if !foundAny {
			info.WriteString("- README.md encontrado en la raíz del proyecto\n")
		} else {
			info.WriteString("\n- README.md encontrado en la raíz del proyecto\n")
		}
	}

	return info.String()
}

// getDocFileSummary obtiene un resumen breve de un archivo de documentación.
func (g *Generator) getDocFileSummary(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	// Leer primeras líneas para obtener título o descripción
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Buscar línea que parezca un título (empieza con #)
		if strings.HasPrefix(line, "#") {
			title := strings.TrimPrefix(line, "#")
			title = strings.TrimSpace(title)
			if title != "" && len(title) < 100 {
				return title
			}
		}
		// Si encontramos una línea de texto sustancioso
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "---") && len(line) > 10 && len(line) < 100 {
			return line
		}
	}

	return ""
}

// isImportantDependency determina si una dependencia es importante de mencionar.
func isImportantDependency(name string) bool {
	important := []string{
		"react", "vue", "angular", "svelte", "next", "nuxt", "gatsby",
		"express", "fastify", "koa", "nest",
		"@griddo", "axios", "lodash", "ramda",
		"zod", "typeorm", "prisma", "mongoose", "sequelize",
		"joi", "yup",
	}
	for _, imp := range important {
		if strings.Contains(name, imp) {
			return true
		}
	}
	return false
}

// isImportantDevDependency determina si una devDependency es importante de mencionar.
func isImportantDevDependency(name string) bool {
	important := []string{
		"typescript", "vite", "webpack", "rollup", "esbuild", "parcel",
		"jest", "vitest", "mocha", "jasmine", "cypress", "playwright", "@testing-library",
		"eslint", "prettier", "@typescript-eslint",
		"babel", "postcss", "tailwindcss", "sass", "less",
		"storybook", "@storybook",
	}
	for _, imp := range important {
		if strings.Contains(name, imp) {
			return true
		}
	}
	return false
}

// getDevelopmentGuideTemplate genera el contenido del development_guide.md usando IA.
// Analiza el proyecto real para obtener información precisa.
func (g *Generator) getDevelopmentGuideTemplate() (string, error) {
	// Analizar el proyecto para obtener información real
	projectContext := g.analyzeProjectContext()

	prompt := fmt.Sprintf(`Genera una guía de desarrollo completa y detallada para el siguiente proyecto.

**Información básica del proyecto:**
- **Nombre:** %s
- **Descripción:** %s
- **Lenguaje principal:** %s
- **Framework:** %s
- **Arquitectura:** %s

**Información adicional del análisis del proyecto:**
%s

La guía debe ser un documento markdown completo que incluya:

1. **Estructura del Proyecto**:
   - Describe la estructura de directorios REAL del proyecto (usa la información del análisis)
   - Explica la organización del código según la arquitectura especificada
   - Describe los directorios y subdirectorios principales encontrados

2. **Convenciones de Código**:
   - Estilo y formato específico (usa info de ESLint/Prettier si está disponible)
   - Nomenclatura de archivos (kebab-case, camelCase, etc.)
   - Nomenclatura de variables, funciones, clases, constantes
   - Comentarios y documentación (JSDoc, TSDoc, etc.)
   - Patrones específicos del proyecto

3. **Sistema de Builds y Scripts**:
   - Lista los scripts IMPORTANTES del package.json (dev, start, build, test, lint, watch, etc.)
   - Explica qué hace cada script importante
   - Menciona comandos específicos del proyecto (como build:themes, build:icons, etc.)
   - ESPECIFICA claramente qué package manager usar (yarn/npm/pnpm)

4. **Configuración Específica**:
   - Si hay alias de importación (de tsconfig.json), inclúyelos
   - Si hay sistemas específicos (como Griddo, Next.js, etc.), explícalos
   - Si hay procesos de generación auto-mática (iconos, temas, tipos), descríbelos

5. **Testing**:
   - Framework de testing usado (Jest, Vitest, Mocha, etc.)
   - Estrategia de pruebas (unitarias, integración, e2e)
   - Comandos para ejecutar tests
   - Cobertura objetivo

6. **Git y Commits**:
   - Convención de mensajes (conventional commits)
   - Estrategia de ramas
   - Proceso de Pull Request

7. **Code Review**:
   - Checklist de revisión específico
   - Criterios de calidad

8. **Despliegue**:
   - Proceso de despliegue específico del proyecto
   - Entornos (development, staging, production)
   - Comandos de build para producción

IMPORTANTE:
- GENERA contenido basado en la información REAL del proyecto proporcionada
- NO uses placeholders genéricos como "..."
- Si hay información específica disponible (scripts reales, alias reales), ÚSALA
- Incluye comandos reales que se ejecutarían
- Para proyectos TypeScript/React, incluye detalles sobre tsc, type-checking, etc.`,
		g.answers.ProjectName,
		g.answers.Description,
		g.answers.Language,
		g.answers.Framework,
		g.answers.Architecture,
		projectContext,
	)

	return g.generateWithClaude(prompt, nil)
}

// getDevelopmentGuideTemplateWithContext genera el contenido del development_guide.md usando IA con contexto de la estructura .claude/.
// A diferencia de getDevelopmentGuideTemplate, esta versión incluye información sobre agents, commands y skills creados.
func (g *Generator) getDevelopmentGuideTemplateWithContext(agentNames, commandNames, skillNames []string) (string, error) {
	// Analizar el proyecto para obtener información real
	projectContext := g.analyzeProjectContext()

	// Leer los READMEs de agents y skills para incluir contexto
	agentsReadme := g.getReadmeContent("agents")
	skillsReadme := g.getReadmeContent("skills")

	// Construir lista de commands para mostrar en la guía
	commandsList := strings.Join(commandNames, ", ")

	prompt := fmt.Sprintf(`Genera una guía de desarrollo completa y detallada para el siguiente proyecto.

**Información básica del proyecto:**
- **Nombre:** %s
- **Descripción:** %s
- **Lenguaje principal:** %s
- **Framework:** %s
- **Arquitectura:** %s

**Información adicional del análisis del proyecto:**
%s

## ESTRUCTURA .CLAUDE/ DEL PROYECTO

Este proyecto tiene una estructura .claude/ con agentes, skills y comandos personalizados.

### Agents Configurados
%s

### Skills Disponibles
%s

### Commands Disponibles
%s

La guía debe ser un documento markdown completo que incluya:

1. **Estructura del Proyecto**:
   - Describe la estructura de directorios REAL del proyecto (usa la información del análisis)
   - Explica la organización del código según la arquitectura especificada
   - Describe los directorios y subdirectorios principales encontrados

2. **Convenciones de Código**:
   - Estilo y formato específico (usa info de ESLint/Prettier si está disponible)
   - Nomenclatura de archivos (kebab-case, camelCase, etc.)
   - Nomenclatura de variables, funciones, clases, constantes
   - Comentarios y documentación (JSDoc, TSDoc, etc.)
   - Patrones específicos del proyecto

3. **Sistema de Builds y Scripts**:
   - Lista los scripts IMPORTANTES del package.json (dev, start, build, test, lint, watch, etc.)
   - Explica qué hace cada script importante
   - Menciona comandos específicos del proyecto (como build:themes, build:icons, etc.)
   - ESPECIFICA claramente qué package manager usar (yarn/npm/pnpm)

4. **Configuración Específica**:
   - Si hay alias de importación (de tsconfig.json), inclúyelos
   - Si hay sistemas específicos (como Griddo, Next.js, etc.), explícalos
   - Si hay procesos de generación auto-mática (iconos, temas, tipos), descríbelos

5. **Testing**:
   - Framework de testing usado (Jest, Vitest, Mocha, etc.)
   - Estrategia de pruebas (unitarias, integración, e2e)
   - Comandos para ejecutar tests
   - Cobertura objetivo

6. **Git y Commits**:
   - Convención de mensajes (conventional commits)
   - Estrategia de ramas
   - Proceso de Pull Request

7. **Code Review**:
   - Checklist de revisión específico
   - Criterios de calidad

8. **Despliegue**:
   - Proceso de despliegue específico del proyecto
   - Entornos (development, staging, production)
   - Comandos de build para producción

9. **Uso de la Estructura .claude/**:
   - Cómo usar los commands disponibles
   - Cuándo invocar agentes específicos
   - Skills recomendadas para cada tarea

IMPORTANTE:
- GENERA contenido basado en la información REAL del proyecto proporcionada
- NO uses placeholders genéricos como "..."
- Si hay información específica disponible (scripts reales, alias reales), ÚSALA
- Incluye comandos reales que se ejecutarían
- Para proyectos TypeScript/React, incluye detalles sobre tsc, type-checking, etc.`,
		g.answers.ProjectName,
		g.answers.Description,
		g.answers.Language,
		g.answers.Framework,
		g.answers.Architecture,
		projectContext,
		agentsReadme,
		skillsReadme,
		commandsList,
	)

	return g.generateWithClaude(prompt, nil)
}

// getProjectStructure retorna la estructura del proyecto según la arquitectura.
func (g *Generator) getProjectStructure() string {
	switch g.answers.Architecture {
	case "Hexagonal":
		return `### Capas

1. **Domain (src/)**: Lógica de negocio, entidades y reglas
2. **Application (src/application/)**: Casos de uso y orquestación
3. **Infrastructure (src/infrastructure/)**: Implementaciones técnicas (DB, API, etc.)

### Directorios Principales
- entities/: Entidades de dominio
- usecases/: Casos de uso o servicios
- repositories/: Acceso a datos
- controllers/: Manejo de peticiones HTTP
`
	case "Microservicios":
		return `### Servicios

Cada servicio es autónomo y se comunica mediante APIs.

### Directorios Principales
- services/: Cada microservicio
- api/: Gateway/Router
- shared/: Código compartido
`
	default:
		return `### Directorios

- src/: Código fuente principal
- tests/: Tests del proyecto
- docs/: Documentación
- config/: Configuraciones
`
	}
}

// getLanguageConventions retorna las convenciones específicas del lenguaje.
func (g *Generator) getLanguageConventions() string {
	switch strings.ToLower(g.answers.Language) {
	case "typescript", "javascript":
		return `- Usar **strict mode**
- Evitar "any" tanto como sea posible
- Usar **interfaces** para contratos
- **async/await** para código asíncrono
- **Arrow functions** para callbacks cortos
`
	case "go":
		return `- Usar **go fmt** para formateo
- **Error handling** explícito (if err != nil)
- **Interfaces** para contratos
- **Goroutines** para concurrencia
- **Context** para cancelación y timeouts
`
	case "python":
		return `- **PEP 8** para estilo
- **Type hints** para funciones
- **Docstrings** para documentación
- **List comprehensions** cuando sea apropiado
- **Context managers** para recursos
`
	default:
		return fmt.Sprintf(`- Seguir convenciones estándar de %s
- Consultar guías de estilo oficiales
- Mantener consistencia en todo el proyecto
`, g.answers.Language)
	}
}

// getNamingConventions retorna las convenciones de nomenclatura.
func (g *Generator) getNamingConventions() string {
	switch strings.ToLower(g.answers.Language) {
	case "typescript", "javascript":
		return `camelCase, PascalCase`
	case "go":
		return `camelCase, PascalCase`
	case "python":
		return `snake_case, PascalCase`
	case "rust":
		return `snake_case, PascalCase`
	default:
		return `camelCase, PascalCase`
	}
}

// getTestingFramework retorna el framework de testing.
func (g *Generator) getTestingFramework() string {
	switch strings.ToLower(g.answers.Language) {
	case "typescript", "javascript":
		if g.answers.Framework != "" {
			return fmt.Sprintf("Jest/Vitest (%s)", g.answers.Framework)
		}
		return "Jest/Vitest"
	case "go":
		return "testing package"
	case "python":
		return "pytest"
	case "rust":
		return "cargo test"
	default:
		return "Framework de testing del lenguaje"
	}
}

// getDeploymentSection retorna la sección de despliegue.
func (g *Generator) getDeploymentSection() string {
	return `### Proceso de Despliegue

1. **Build**: Compilar/empaquetar la aplicación
2. **Test**: Ejecutar suite de tests
3. **Deploy**: Desplegar al entorno correspondiente
4. **Verify**: Verificar que el despliegue fue exitoso

### Monitoreo
- Revisar logs de aplicación
- Monitorear métricas clave
- Configurar alertas`
}

// getReferences retorna las referencias del proyecto.
func (g *Generator) getReferences() string {
	baseRefs := "- Documentación oficial del lenguaje\n"

	switch strings.ToLower(g.answers.Language) {
	case "typescript":
		baseRefs += `- TypeScript Handbook: https://www.typescriptlang.org/docs/\n`
		if g.answers.Framework != "" {
			baseRefs += fmt.Sprintf("- %s Documentation: [link]\n", g.answers.Framework)
		}
	case "go":
		baseRefs += "- Effective Go: https://go.dev/doc/effective_go\n"
		baseRefs += "- Go Reference: https://pkg.go.dev/\n"
	case "python":
		baseRefs += "- PEP 8: https://peps.python.org/pep-0008/\n"
		baseRefs += "- Python Documentation: https://docs.python.org/\n"
	}

	return baseRefs
}

// parseCommandFrontmatter extrae metadata del frontmatter de un archivo de comando.
func (g *Generator) parseCommandFrontmatter(filePath string) (*CommandInfo, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Extraer nombre del archivo
	fileName := filepath.Base(filePath)
	name := strings.TrimSuffix(fileName, ".md")

	info := &CommandInfo{
		Name:     name,
		FilePath: filePath,
	}

	// Buscar frontmatter YAML entre ---
	strContent := string(content)
	startIdx := strings.Index(strContent, "---")
	if startIdx == -1 {
		// No hay frontmatter, usar valores por defecto
		info.Description = fmt.Sprintf("Comando %s", name)
		info.Usage = name
		return info, nil
	}

	endIdx := strings.Index(strContent[startIdx+3:], "---")
	if endIdx == -1 {
		info.Description = fmt.Sprintf("Comando %s", name)
		info.Usage = name
		return info, nil
	}

	frontmatter := strContent[startIdx+3 : startIdx+3+endIdx]
	lines := strings.Split(frontmatter, "\n")

	// Parsear líneas del frontmatter
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "name":
			// Usar el nombre del archivo en lugar del del frontmatter
		case "description":
			info.Description = strings.Trim(value, `"`)
		case "usage":
			info.Usage = strings.Trim(value, `"`)
		}
	}

	// Fallback si no se encontraron valores
	if info.Description == "" {
		info.Description = fmt.Sprintf("Comando %s para %s", name, g.answers.ProjectName)
	}
	if info.Usage == "" {
		info.Usage = name
	}

	return info, nil
}

// GenerateCommandsReadme genera el README.md del directorio commands con la lista de comandos.
func (g *Generator) GenerateCommandsReadme(commandNames []string) error {
	g.logger.Debug("Generando commands/README.md")

	commandsDir := filepath.Join(g.projectPath, ".claude", "commands")
	outputPath := filepath.Join(commandsDir, "README.md")

	// Parsear todos los comandos para obtener su metadata
	var commands []*CommandInfo
	for _, cmdName := range commandNames {
		cmdPath := filepath.Join(commandsDir, sanitizeFilename(cmdName)+".md")
		info, err := g.parseCommandFrontmatter(cmdPath)
		if err != nil {
			g.logger.Warn("Error parsing command %s: %v", cmdName, err)
			continue
		}
		commands = append(commands, info)
	}

	// Generar contenido del README
	content := g.buildCommandsReadmeContent(commands)

	// Escribir archivo
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error escribiendo commands/README.md: %w", err)
	}

	g.logger.Info("Commands README generado en %s", outputPath)
	return nil
}

// buildCommandsReadmeContent construye el contenido markdown del README de commands.
func (g *Generator) buildCommandsReadmeContent(commands []*CommandInfo) string {
	var sb strings.Builder

	sb.WriteString("# Commands Disponibles\n\n")
	sb.WriteString("Este directorio contiene los comandos del proyecto ")
	sb.WriteString(g.answers.ProjectName)
	sb.WriteString(".\n\n")

	if len(commands) == 0 {
		sb.WriteString("*No hay comandos configurados aún.*\n")
		return sb.String()
	}

	sb.WriteString("## Comandos Configurados\n\n")

	for _, cmd := range commands {
		sb.WriteString(fmt.Sprintf("### %s\n\n", cmd.Name))

		if cmd.Description != "" {
			sb.WriteString(fmt.Sprintf("**Descripción**: %s\n\n", cmd.Description))
		}

		if cmd.Usage != "" {
			sb.WriteString(fmt.Sprintf("**Uso**: `%s`\n\n", cmd.Usage))
		}

		sb.WriteString("---\n\n")
	}

	sb.WriteString("## Uso de los Comandos\n\n")
	sb.WriteString("Los comandos definen flujos de trabajo orquestados que utilizan uno o más agentes.\n\n")
	sb.WriteString(fmt.Sprintf("Para más información sobre cómo se crean estos comandos, consulta las guías en `.claude/embeds/command_guide.md`.\n"))

	return sb.String()
}

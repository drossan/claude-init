package survey

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// Runner ejecuta un survey de forma interactiva.
type Runner struct {
	Survey *Survey
}

// NewRunner crea un nuevo Runner con las preguntas especificadas.
func NewRunner(questions []*Question) *Runner {
	return &Runner{
		Survey: &Survey{
			Questions: questions,
		},
	}
}

// Run ejecuta el survey de forma interactiva y recolecta las respuestas.
func (r *Runner) Run() (*Answers, error) {
	return r.RunWithPrefill(nil)
}

// RunWithPrefill ejecuta el survey con valores pre-llenados que el usuario puede editar.
// Si prefill es nil, se ejecuta el survey normal sin valores pre-llenados.
func (r *Runner) RunWithPrefill(prefill *Answers) (*Answers, error) {
	answers := &Answers{}

	// Si hay prefill, copiar esos valores iniciales
	if prefill != nil {
		*answers = *prefill
	}

	for _, q := range r.Survey.Questions {
		// Obtener valor pre-llenado para esta pregunta
		defaultValue := r.getPrefillValue(prefill, q.ID)
		if defaultValue == "" {
			// Si no hay prefill, usar el default de la pregunta
			defaultValue = q.Default
		}

		var answer string
		var err error

		switch q.Type {
		case QuestionTypeInput:
			prompt := &survey.Input{
				Message: q.Text,
				Default: defaultValue,
				Help:    q.Placeholder,
			}
			if q.Required {
				err = survey.AskOne(prompt, &answer, survey.WithValidator(survey.Required))
			} else {
				err = survey.AskOne(prompt, &answer)
			}

		case QuestionTypeMultiline:
			// Usar input estándar con validación de longitud mínima
			// El usuario puede pegar texto o escribir una descripción corta
			prompt := &survey.Input{
				Message: q.Text,
				Default: defaultValue,
				Help:    "Describe el contexto del negocio, objetivos y requisitos. Puedes pegar texto o escribir una descripción corta.",
			}
			if q.Required {
				err = survey.AskOne(prompt, &answer, survey.WithValidator(func(ans interface{}) error {
					str, ok := ans.(string)
					if !ok || strings.TrimSpace(str) == "" {
						return fmt.Errorf("se requiere una descripción del contexto del negocio (mínimo 20 caracteres)")
					}
					if len(strings.TrimSpace(str)) < 20 {
						return fmt.Errorf("la descripción es muy corta. Por favor, proporciona más detalles (mínimo 20 caracteres)")
					}
					return nil
				}))
			} else {
				err = survey.AskOne(prompt, &answer)
			}

		case QuestionTypeSelect:
			prompt := &survey.Select{
				Message: q.Text,
				Options: q.Options,
				Default: defaultValue,
			}
			if q.Required {
				err = survey.AskOne(prompt, &answer, survey.WithValidator(survey.Required))
			} else {
				err = survey.AskOne(prompt, &answer)
			}

		case QuestionTypeMultiSelect:
			var selected []string
			prompt := &survey.MultiSelect{
				Message: q.Text,
				Options: q.Options,
			}
			if q.Required {
				err = survey.AskOne(prompt, &selected, survey.WithValidator(survey.MinItems(1)))
			} else {
				err = survey.AskOne(prompt, &selected)
			}
			answer = fmt.Sprintf("%v", selected)

		case QuestionTypeConfirm:
			var confirmed bool
			prompt := &survey.Confirm{
				Message: q.Text,
				Default: defaultValue == "true" || defaultValue == "yes",
			}
			err = survey.AskOne(prompt, &confirmed)
			if confirmed {
				answer = "yes"
			} else {
				answer = "no"
			}
		}

		if err != nil {
			return nil, fmt.Errorf("failed to ask question %s: %w", q.ID, err)
		}

		// Asignar respuesta al campo correspondiente
		r.setAnswer(answers, q.ID, answer)
	}

	return answers, nil
}

// setAnswer asigna la respuesta al campo correspondiente en Answers.
func (r *Runner) setAnswer(answers *Answers, id, value string) {
	switch id {
	case "project_origin":
		answers.ProjectOrigin = value
	case "project_name":
		answers.ProjectName = value
	case "description":
		answers.Description = value
	case "language":
		answers.Language = value
	case "framework":
		answers.Framework = value
	case "architecture":
		answers.Architecture = value
	case "database":
		answers.Database = value
	case "project_category":
		answers.ProjectCategory = value
	case "business_context":
		answers.BusinessContext = value
	}
}

// getPrefillValue obtiene el valor pre-llenado para una pregunta específica.
func (r *Runner) getPrefillValue(prefill *Answers, id string) string {
	if prefill == nil {
		return ""
	}

	switch id {
	case "project_origin":
		return prefill.ProjectOrigin
	case "project_name":
		return prefill.ProjectName
	case "description":
		return prefill.Description
	case "language":
		return prefill.Language
	case "framework":
		return prefill.Framework
	case "architecture":
		return prefill.Architecture
	case "database":
		return prefill.Database
	case "project_category":
		return prefill.ProjectCategory
	case "business_context":
		return prefill.BusinessContext
	default:
		return ""
	}
}

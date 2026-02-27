package prompts

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

type PromptManager struct {
	tmpl *template.Template
}

func New() (*PromptManager, error) {
	t, err := template.ParseFS(templatesFS, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return &PromptManager{tmpl: t}, nil
}

func (p *PromptManager) RenderSystemPrompt(ctx SystemPromptContext) (string, error) {
	var buf bytes.Buffer
	if err := p.tmpl.ExecuteTemplate(&buf, "system_prompt.tmpl", ctx); err != nil {
		return "", fmt.Errorf("failed to execute system_prompt template: %w", err)
	}
	return buf.String(), nil
}

func (p *PromptManager) RenderAdditionalInfo(ctx AdditionalInfoContext) (string, error) {
	var buf bytes.Buffer
	if err := p.tmpl.ExecuteTemplate(&buf, "additional_info.tmpl", ctx); err != nil {
		return "", fmt.Errorf("failed to execute additional_info template: %w", err)
	}
	return buf.String(), nil
}

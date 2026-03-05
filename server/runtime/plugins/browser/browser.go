package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"openhands-go/server/runtime"
	"sync"

	"github.com/playwright-community/playwright-go"
	"github.com/tmc/langchaingo/llms"
)

type BrowserPlugin struct {
	mu            sync.Mutex
	pw            *playwright.Playwright
	browser       playwright.Browser
	page          playwright.Page
	isInitialized bool
}

func NewBrowserPlugin() *BrowserPlugin {
	return &BrowserPlugin{}
}

func (p *BrowserPlugin) Name() string {
	return "browser"
}

func (p *BrowserPlugin) Init(ctx context.Context, rt runtime.Runtime) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isInitialized {
		return nil
	}

	// Install drivers if missing (needed for playwright to run if not already installed globally)
	err := playwright.Install()
	if err != nil {
		fmt.Printf("Warning: failed to install playwright drivers: %v\n", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not start playwright: %v", err)
	}
	p.pw = pw

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("could not launch browser: %v", err)
	}
	p.browser = browser

	page, err := browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not create page: %v", err)
	}
	p.page = page
	p.isInitialized = true
	return nil
}

func (p *BrowserPlugin) Tools() []llms.Tool {
	return []llms.Tool{
		{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        "browser_action",
				Description: "Interact with a web browser. Supports: goto, scrape, screenshot.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"action": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"goto", "scrape", "screenshot"},
							"description": "The action to perform.",
						},
						"url": map[string]interface{}{
							"type":        "string",
							"description": "URL to navigate to (for goto).",
						},
					},
					"required": []string{"action"},
				},
			},
		},
	}
}

func (p *BrowserPlugin) HandleToolCall(ctx context.Context, name string, args string) (string, bool, error) {
	if name != "browser_action" {
		return "", false, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isInitialized {
		return "Browser not initialized", true, nil
	}

	var params struct {
		Action string `json:"action"`
		URL    string `json:"url"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error unmarshalling args: %v", err), true, nil
	}

	switch params.Action {
	case "goto":
		if _, err := p.page.Goto(params.URL); err != nil {
			return fmt.Sprintf("Error navigating to %s: %v", params.URL, err), true, nil
		}
		title, _ := p.page.Title()
		return fmt.Sprintf("Navigated to %s. Title: %s", params.URL, title), true, nil

	case "scrape":
		content, err := p.page.Content()
		if err != nil {
			return fmt.Sprintf("Error getting content: %v", err), true, nil
		}
		// Truncate content for sanity
		if len(content) > 5000 {
			content = content[:5000] + "... (truncated)"
		}
		return content, true, nil

	case "screenshot":
		// Return base64 or path?
		// For now just say executed.
		path := fmt.Sprintf("screenshot_%d.png", 1)
		if _, err := p.page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String(path),
		}); err != nil {
			return fmt.Sprintf("Error taking screenshot: %v", err), true, nil
		}
		return fmt.Sprintf("Screenshot saved to %s", path), true, nil
	}

	return fmt.Sprintf("Unknown action: %s", params.Action), true, nil
}

func (p *BrowserPlugin) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.browser != nil {
		p.browser.Close()
	}
	if p.pw != nil {
		p.pw.Stop()
	}
	return nil
}

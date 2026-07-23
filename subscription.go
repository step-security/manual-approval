package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func checkSubscription() {
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	repoPrivate := ""
	if eventPath != "" {
		data, err := os.ReadFile(eventPath)
		if err == nil {
			var event map[string]interface{}
			if json.Unmarshal(data, &event) == nil {
				if repo, ok := event["repository"].(map[string]interface{}); ok {
					if private, ok := repo["private"].(bool); ok {
						if private {
							repoPrivate = "true"
						} else {
							repoPrivate = "false"
						}
					}
				}
			}
		}
	}

	upstream := "trstringer/manual-approval"
	docsURL := "https://docs.stepsecurity.io/actions/stepsecurity-maintained-actions"

	fmt.Println()
	fmt.Println("\033[1;36mStepSecurity Maintained Action\033[0m")
	fmt.Printf("Secure drop-in replacement for %s\n", upstream)
	if repoPrivate == "false" {
		fmt.Println("\033[32m✓ Free for public repositories\033[0m")
	}
	fmt.Printf("\033[36mLearn more:\033[0m %s\n", docsURL)
	fmt.Println()

	if repoPrivate == "false" {
		return
	}

	serverURL := os.Getenv("GITHUB_SERVER_URL")
	if serverURL == "" {
		serverURL = "https://github.com"
	}
	actionRepo := os.Getenv("GITHUB_ACTION_REPOSITORY")
	githubRepo := os.Getenv("GITHUB_REPOSITORY")

	type requestBody struct {
		Action     string `json:"action"`
		GhesServer string `json:"ghes_server,omitempty"`
	}

	reqBody := requestBody{Action: actionRepo}
	if serverURL != "https://github.com" {
		reqBody.GhesServer = serverURL
	}

	bodyBytes, _ := json.Marshal(reqBody)
	apiURL := fmt.Sprintf("https://agent.api.stepsecurity.io/v1/github/%s/actions/maintained-actions-subscription", githubRepo)

	httpClient := &http.Client{Timeout: 3 * time.Second}
	resp, err := httpClient.Post(apiURL, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		fmt.Println("Timeout or API not reachable. Continuing to next step.")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusForbidden {
		fmt.Println("::error::\033[1;31mThis action requires a StepSecurity subscription for private repositories.\033[0m")
		fmt.Printf("::error::\033[31mLearn how to enable a subscription: %s\033[0m\n", docsURL)
		os.Exit(1)
	}
}

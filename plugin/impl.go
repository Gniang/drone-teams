// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Settings for the plugin.
type Settings struct {
	Webhook string
	Status  string
	Card    string
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	// Verify the webhook endpoint
	if p.settings.Webhook == "" {
		// If webhook is undefined, check if the ${DRONE_BRANCH}_teams_webhook env var is defined.
		branchWebhook := fmt.Sprintf("%s_teams_webhook", os.Getenv("DRONE_BRANCH"))
		if os.Getenv(branchWebhook) == "" {
			return fmt.Errorf("no webhook endpoint provided")
		}
		// Set webhook setting to ${DRONE_BRANCH}_teams_webhook
		p.settings.Webhook = os.Getenv(branchWebhook)
	}

	// If the plugin status setting is defined, use that as the build status
	if p.settings.Status == "" {
		p.settings.Status = p.pipeline.Build.Status
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {

	isAdaptiveCard := strings.ToLower(p.settings.Card) == "adaptive"

	var card interface{}
	if isAdaptiveCard {
		card = CreateAcaptiveCard(p)

	} else {
		card = CreateMessageCard(p)
	}

	log.Info("Generated card: ", card)

	// MS teams webhook post
	jsonValue, _ := json.Marshal(card)
	_, err := http.Post(p.settings.Webhook, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Error("Failed to send request to teams webhook")
		return err
	}
	return nil
}

// Create post data for AdaptiveCard
func CreateAcaptiveCard(p *Plugin) WebhookContent {
	auther := fmt.Sprintf("%s (%s)", p.pipeline.Commit.Author, p.pipeline.Commit.Author.Email)
	style := "default"
	if p.pipeline.Build.Status == "failure" {
		style = "attention"
	}

	// Create rich message card body
	card := WebhookContent{
		Attachments: []AdaptiveCard{{
			ContentType: "application/vnd.microsoft.card.adaptive",
			Content: AdaptiveCardContent{
				Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
				Type:    "AdaptiveCard",
				Version: "1.4",
				Body: []CardContainer{{
					Type:  "Container",
					Style: style,
					Items: []CardColumnSet{
						CreateNameValueLabel("Build Status", p.settings.Status),
						CreateNameValueLabel("URL", ToUrlMarkdown(p.pipeline.Build.Link)),
						CreateNameValueLabel("Branch", p.pipeline.Commit.Branch),
						CreateNameValueLabel("Auther", auther),
					},
				}},
			},
		}},
	}
	return card
}

func ToUrlMarkdown(url string) string {
	return fmt.Sprintf("[%s](%s)", url, url)
}

func CreateNameValueLabel(name string, value string) CardColumnSet {
	return CardColumnSet{
		Type: "ColumnSet",
		Columns: []CardColumn{
			{
				Type:                     "Column",
				Width:                    "auto",
				VerticalContentAlignment: "Center",
				HorizontalAlignment:      "Left",
				Items: []CardTextBlock{{
					Type: "TextBlock",
					Text: name,
					Wrap: true,
				}},
			},
			{
				Type:                     "Column",
				Width:                    "stretch",
				VerticalContentAlignment: "Center",
				HorizontalAlignment:      "Left",
				Items: []CardTextBlock{{
					Type: "TextBlock",
					Text: value,
					Wrap: true,
				}},
			}},
	}
}

// If commit link is not null add commit link
func GetCommitLink(p *Plugin) string {
	if p.pipeline.Commit.Link != "" {
		return p.pipeline.Commit.Link
	} else if cl, present := os.LookupEnv("DRONE_COMMIT_LINK"); present && cl != "" {
		return cl
	}
	return ""
}

// Create post data for MessageCard
func CreateMessageCard(p *Plugin) MessageCard {

	// Default card color is green
	themeColor := "96FF33"

	// Create list of card facts
	facts := []MessageCardSectionFact{
		{
			Name:  "Build Number",
			Value: fmt.Sprintf("%d", p.pipeline.Build.Number),
		},
		{
			Name:  "Time",
			Value: p.pipeline.Build.Started.String(),
		},
		{
			Name:  "Repo Link",
			Value: p.pipeline.Repo.Link,
		},
		{
			Name:  "Branch",
			Value: p.pipeline.Build.Branch,
		},
		{
			Name:  "Git Author",
			Value: fmt.Sprintf("%s (%s)", p.pipeline.Commit.Author, p.pipeline.Commit.Author.Email),
		},
		{
			Name:  "Commit Message Title",
			Value: p.pipeline.Commit.Message.Title,
		},
		{
			Name:  "Commit Message Body",
			Value: p.pipeline.Commit.Message.Body,
		}}

	// If commit link is not null add commit link fact to card
	if p.pipeline.Commit.Link != "" {
		facts = append(facts, MessageCardSectionFact{
			Name:  "Commit Link",
			Value: p.pipeline.Commit.Link,
		})
	} else if commitLink, present := os.LookupEnv("DRONE_COMMIT_LINK"); present && commitLink != "" {
		facts = append(facts, MessageCardSectionFact{
			Name:  "Commit Link",
			Value: commitLink,
		})
	}

	// If build link is not null add build link fact to card
	if p.pipeline.Build.Link != "" && p.pipeline.Stage.Number > 0 {
		facts = append(facts, MessageCardSectionFact{
			Name:  "Build Link",
			Value: "[" + p.pipeline.Build.Link + "/" + strconv.Itoa(p.pipeline.Stage.Number) + "](" + p.pipeline.Build.Link + "/" + strconv.Itoa(p.pipeline.Stage.Number) + ")",
		})
	} else {
		buildLink, presentLink := os.LookupEnv("DRONE_BUILD_LINK")
		buildStage, presentStage := os.LookupEnv("DRONE_STAGE_NUMBER")
		if presentLink && presentStage && buildLink != "" && buildStage != "" {
			facts = append(facts, MessageCardSectionFact{
				Name:  "Build Link",
				Value: "[" + buildLink + "/" + buildStage + "](" + buildLink + "/" + buildStage + ")",
			})
		}
	}

	// If build has failed, change color to red and add failed step fact
	if p.settings.Status == "failure" {
		themeColor = "FF5733"
		facts = append(facts, MessageCardSectionFact{
			Name:  "Failed Build Steps",
			Value: strings.Join(p.pipeline.Build.FailedSteps, " "),
		})
		// If the plugin status setting is defined and is "building", set the color to blue
	} else if p.settings.Status == "building" {
		themeColor = "002BFF"
	}

	// Create rich message card body
	card := MessageCard{
		Type:       "MessageCard",
		Context:    "http://schema.org/extensions",
		ThemeColor: themeColor,
		Summary:    p.pipeline.Repo.Slug,
		Sections: []MessageCardSection{{
			ActivityTitle:    p.pipeline.Repo.Slug,
			ActivitySubtitle: strings.ToUpper(p.settings.Status),
			ActivityImage:    "https://github.com/jdamata/drone-teams/raw/master/drone.png",
			Markdown:         true,
			Facts:            facts,
		}},
	}
	return card
}

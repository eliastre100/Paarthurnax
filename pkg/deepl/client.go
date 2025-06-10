package deepl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	DeepLDomainFree = "api-free.deepl.com"
	DeepLDomainPro  = "api.deepl.com"

	DefaultMaxRetry = 100
)

type Client struct {
	MaxRetry  int
	apiKey    string
	domain    string
	languages []Language
}

type Language struct {
	Language          string `json:"language"`
	Name              string `json:"name"`
	SupportsFormality bool   `json:"supports_formality"`
}

type translationRequest struct {
	Text           []string `json:"text"`
	SourceLanguage string   `json:"source_lang"`
	SargetLanguage string   `json:"target_lang"`
	TagHandling    string   `json:"tag_handling"`
}

type translationResponse struct {
	Translations []struct {
		Text string `json:"text"`
	} `json:"translations"`
}

func NewClient(apiKey string, domain string) (*Client, error) {
	client := Client{apiKey: apiKey, domain: domain, MaxRetry: DefaultMaxRetry}
	if err := client.refreshLanguages(); err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *Client) refreshLanguages() error {
	payload, err := c.fetch("v2/languages?type=source", nil, 0)
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, &c.languages)
}

func (c *Client) Translate(text string, sourceLocale string, destinationLocale string) (string, error) {
	request, err := json.Marshal(&translationRequest{Text: []string{text}, SourceLanguage: sourceLocale, SargetLanguage: destinationLocale, TagHandling: "html"})
	if err != nil {
		return "", errors.New("Unable to create translation request: " + err.Error())
	}
	payload, err := c.fetch("/v2/translate", request, 0)

	var translations translationResponse
	if err = json.Unmarshal(payload, &translations); err != nil {
		return "", err
	}
	return translations.Translations[0].Text, nil
}

func (c *Client) fetch(endpoint string, body []byte, try int) ([]byte, error) {
	endpoint, _ = strings.CutPrefix(endpoint, "/")
	url := fmt.Sprintf("https://%s/%s", c.domain, endpoint)
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", c.apiKey))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case 200:
		return io.ReadAll(resp.Body)
	case 500:
		slog.Warn("Request failed with code 500", "endpoint", endpoint)
	case 429:
		if try+1 > c.MaxRetry {
			return nil, fmt.Errorf("deepl responded with status %d (%s) (too many retries)", resp.StatusCode, resp.Status)
		}
		time.Sleep(time.Duration(math.Min(math.Ceil(float64(try)*math.Pow(1.6, float64(try))), 1000)) * time.Second)
		return c.fetch(endpoint, body, try+1)
	}
	return nil, fmt.Errorf("unexpected response code %d", resp.StatusCode)
}

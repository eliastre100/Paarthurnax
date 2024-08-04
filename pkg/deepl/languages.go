package deepl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type language struct {
	Language          string `json:"language"`
	Name              string `json:"name"`
	SupportsFormality bool   `json:"supports_formality"`
}

func (client *Deepl) Languages() ([]string, error) {
	url := fmt.Sprintf("https://%s/v2/languages?type=source", client.domain)
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", client.apiKey))
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("DeepL responded with status " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var languages []language
	err = json.Unmarshal(body, &languages)
	if err != nil {
		return nil, err
	}
	return extractLanguages(languages), nil
}

func extractLanguages(languages []language) []string {
	result := make([]string, 0, len(languages))

	for _, l := range languages {
		result = append(result, l.Language)
	}

	return result
}

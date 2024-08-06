package deepl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

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

func (client *Deepl) Translate(text string, source string, destination string) (string, error) {
	url := fmt.Sprintf("https://%s/v2/translate", client.domain)
	request, err := json.Marshal(&translationRequest{Text: []string{text}, SourceLanguage: source, SargetLanguage: destination, TagHandling: "html"})
	if err != nil {
		return "", errors.New("Unable to create translation request: " + err.Error())
	}
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request))
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", client.apiKey))
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("DeepL responded with status " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var translations translationResponse
	err = json.Unmarshal(body, &translations)
	if err != nil {
		return "", err
	}
	return translations.Translations[0].Text, nil
}

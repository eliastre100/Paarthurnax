package deepl

func New(apiKey string, domain string) (*Deepl, error) {
	client := Deepl{apiKey: apiKey, domain: domain}
	if _, err := client.Languages(); err != nil {
		return nil, err
	}

	return &client, nil
}

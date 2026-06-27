package v1

type Locales struct {
	Source       string   `toml:"source"`
	Destinations []string `toml:"destinations"`
}

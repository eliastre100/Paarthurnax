package translation

func LoadOrCreate(filename string) (*TranslationFile, error) {
	file, err := Load(filename)
	if err == nil {
		return file, nil
	}
	return Create(filename)
}

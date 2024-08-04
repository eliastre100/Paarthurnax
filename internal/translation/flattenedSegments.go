package translation

func (translation *TranslationFile) FlattenedSegments() map[string]string {
	return extractFlattenedSegments(translation.segments)
}

func extractFlattenedSegments(m map[string]interface{}) map[string]string {
	result := make(map[string]string)

	for key, value := range m {
		switch value.(type) {
		case string:
			result[key] = value.(string)
		case map[string]interface{}:
			for subKey, subValue := range extractFlattenedSegments(value.(map[string]interface{})) {
				result[key+"."+subKey] = subValue
			}
		}
	}

	return result
}

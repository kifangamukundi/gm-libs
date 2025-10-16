package transformations

func Transform[T any](
	models []T,
	fieldNames []string,
	fieldExtractors ...func(T) any,
) []map[string]any {
	if len(fieldNames) != len(fieldExtractors) {
		panic("fieldNames and fieldExtractors must have the same length")
	}

	transformed := make([]map[string]any, len(models))
	for i, model := range models {
		entry := make(map[string]any)

		for j, extractor := range fieldExtractors {
			entry[fieldNames[j]] = extractor(model)
		}

		transformed[i] = entry
	}

	return transformed
}

func TransformSingle[T any](
	model T,
	fieldNames []string,
	fieldExtractors ...func(T) any,
) map[string]any {
	if len(fieldNames) != len(fieldExtractors) {
		panic("fieldNames and fieldExtractors must have the same length")
	}

	result := make(map[string]any)
	for i, extractor := range fieldExtractors {
		result[fieldNames[i]] = extractor(model)
	}

	return result
}

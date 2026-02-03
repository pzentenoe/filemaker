package filemaker

// validateRequired validates that a required field is not empty.
// Returns a ValidationError if the field is empty, nil otherwise.
func validateRequired(field, value, message string) error {
	if value == "" {
		return &ValidationError{
			Field:   field,
			Message: message,
		}
	}
	return nil
}

// validateDatabase validates that a database name is provided.
func validateDatabase(database string) error {
	return validateRequired("database", database, "database name is required")
}

// validateLayout validates that a layout name is provided.
func validateLayout(layout string) error {
	return validateRequired("layout", layout, "layout name is required")
}

// validateToken validates that a session token is provided.
func validateToken(token string) error {
	return validateRequired("token", token, "session token is required")
}

// validateRecordID validates that a record ID is provided.
func validateRecordID(recordID string) error {
	return validateRequired("recordID", recordID, "record ID is required")
}

// validateScript validates that a script name is provided.
func validateScript(script string) error {
	return validateRequired("script", script, "script name is required")
}

// validateFieldName validates that a field name is provided.
func validateFieldName(fieldName string) error {
	return validateRequired("fieldName", fieldName, "field name is required")
}

// validateFilePath validates that a file path is provided.
func validateFilePath(filePath string) error {
	return validateRequired("filePath", filePath, "file path is required")
}

// validateFilename validates that a filename is provided.
func validateFilename(filename string) error {
	return validateRequired("filename", filename, "filename is required")
}

// validateRepetition validates that a repetition number is valid (>= 1).
func validateRepetition(repetition int) error {
	if repetition < 1 {
		return &ValidationError{
			Field:   "repetition",
			Message: "repetition must be >= 1",
		}
	}
	return nil
}

// validateGlobalFields validates that at least one global field is provided.
func validateGlobalFields(fields map[string]any) error {
	if len(fields) == 0 {
		return &ValidationError{
			Field:   "globalFields",
			Message: "at least one global field must be specified",
		}
	}
	return nil
}

// validateFileData validates that file data is not empty.
func validateFileData(data []byte) error {
	if len(data) == 0 {
		return &ValidationError{
			Field:   "data",
			Message: "file data cannot be empty",
		}
	}
	return nil
}

// validateURL validates that a URL is provided.
func validateURL(url string) error {
	return validateRequired("url", url, "url is required")
}

// validateAction validates that an action is one of the allowed values.
func validateAction(action string) error {
	if action == "" {
		return &ValidationError{
			Field:   "action",
			Message: "action is required (create, edit, delete)",
		}
	}

	switch action {
	case "create", "edit", "delete":
		return nil
	default:
		return &ValidationError{
			Field:   "action",
			Message: "invalid action, must be create, edit, or delete",
		}
	}
}

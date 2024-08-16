package main

type Language struct {
	Key   string
	Value string
}

// Language represents a programming language with a key and a value.
func getLanguages() []Language {
	pairs := []Language{
		{Key: "plaintext", Value: "Plain Text"},
		{Key: "html", Value: "HTML"},
		{Key: "css", Value: "CSS"},
		{Key: "scss", Value: "SCSS"},
		{Key: "javascript", Value: "JavaScript"},
		{Key: "typescript", Value: "TypeScript"},
		{Key: "json", Value: "JSON"},
		{Key: "php", Value: "PHP"},
		{Key: "python", Value: "Python"},
		{Key: "go", Value: "GO"},
		{Key: "sql", Value: "SQL"},
		{Key: "bash", Value: "Bash"},
		{Key: "xml", Value: "XML"},
		{Key: "c", Value: "C"},
		{Key: "cpp", Value: "C++"},
		{Key: "csharp", Value: "C#"},
		{Key: "java", Value: "Java"},
		{Key: "swift", Value: "Swift"},
		{Key: "rust", Value: "Rust"},
		{Key: "ruby", Value: "Ruby"},
		{Key: "perl", Value: "Perl"},
		{Key: "lua", Value: "Lua"},
		{Key: "shell", Value: "Shell"},
	}
	return pairs
}

// getLanguages returns a slice of Language structs containing supported languages.
// Each Language struct contains a key (used internally) and a value (display name).
func getLanguageKeys() []string {
	// Retrieve the list of languages.
	languages := getLanguages()

	// Initialize a slice to hold the keys.
	keys := make([]string, len(languages))

	// Populate the slice with the keys of the languages.
	for i, lang := range languages {
		keys[i] = lang.Key
	}
	return keys
}

// getLanguageLabel returns the display name of a language given its key.
// If the key does not match any language, it returns "Plain Text" as the default value.
func getLanguageLabel(s string) string {
	// Retrieve the list of languages.
	languages := getLanguages()

	// Iterate through the languages to find a match for the key.
	for _, pair := range languages {
		if pair.Key == s {
			return pair.Value
		}
	}

	// If no match is found, return "Plain Text" as the default value.
	return "Plain Text"
}

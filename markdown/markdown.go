package markdown

import (
	"fmt"
	"strings"
)

var (
	// is thread safe/goroutine safe
	markdownReplacer = strings.NewReplacer(
		"\\", "\\\\",
		"`", "\\`",
		"*", "\\*",
		"_", "\\_",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		".", "\\.",
		"!", "\\!",
	)
)

// Escape user input outside of inline code blocks
func Escape(userInput string) string {
	return markdownReplacer.Replace(userInput)
}

func CodeHighlight(language string, code string) string {
	hasLeadingNewline := strings.HasPrefix(code, "\n")
	hasTrailingNewline := strings.HasSuffix(code, "\n")

	if hasLeadingNewline && hasTrailingNewline {
		return WrapInMultilineCodeBlock(fmt.Sprintf("%s%s", language, code))
	} else if hasLeadingNewline {
		return WrapInMultilineCodeBlock(fmt.Sprintf("%s%s\n", language, code))
	} else if hasTrailingNewline {
		return WrapInMultilineCodeBlock(fmt.Sprintf("%s\n%s", language, code))
	}
	return WrapInMultilineCodeBlock(fmt.Sprintf("%s\n%s\n", language, code))
}

func WrapInMultilineCodeBlock(text string) string {
	return WrapInCustom(text, "```")
}

// WrapInInlineCodeBlock puts the user input into a inline codeblock that is properly escaped.
func WrapInInlineCodeBlock(text string) string {
	return WrapInCustom(text, "`")
}

func WrapInFat(text string) string {
	return WrapInCustom(text, "**")
}

func WrapInCustom(text, wrap string) (result string) {
	if text == "" {
		return ""
	}

	numWraps := strings.Count(text, wrap) + 1
	result = text
	for idx := 0; idx < numWraps; idx++ {
		result = fmt.Sprintf("%s%s%s", wrap, result, wrap)
	}
	return result
}

package wpimport

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// NestedListsToMarkdown converts nested HTML lists to Markdown
func NestedListsToMarkdown(content string) string {
	// First sanitize the content
	content = SanitizeWordPressContent(content)

	// Replace any remaining HTML tags with appropriate markdown syntax

	// First handle ordered lists (ol)
	olRegex := regexp.MustCompile(`(?s)<ol[^>]*>(.*?)</ol>`)
	content = olRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		liRegex := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`)
		items := liRegex.FindAllStringSubmatch(match, -1)

		result := "\n"
		for i, item := range items {
			if len(item) > 1 {
				// Handle nested lists within this item
				itemContent := item[1]

				// Extract text before any nested list
				mainText := itemContent
				if strings.Contains(itemContent, "<ul") {
					idx := strings.Index(itemContent, "<ul")
					if idx > 0 {
						mainText = itemContent[:idx]
					} else {
						mainText = ""
					}
				} else if strings.Contains(itemContent, "<ol") {
					idx := strings.Index(itemContent, "<ol")
					if idx > 0 {
						mainText = itemContent[:idx]
					} else {
						mainText = ""
					}
				}

				// Convert inline formatting to markdown
				mainText = regexp.MustCompile(`<strong[^>]*>(.*?)</strong>`).ReplaceAllString(mainText, "**$1**")
				mainText = regexp.MustCompile(`<b[^>]*>(.*?)</b>`).ReplaceAllString(mainText, "**$1**")
				mainText = regexp.MustCompile(`<em[^>]*>(.*?)</em>`).ReplaceAllString(mainText, "*$1*")
				mainText = regexp.MustCompile(`<i[^>]*>(.*?)</i>`).ReplaceAllString(mainText, "*$1*")
				mainText = regexp.MustCompile(`<a[^>]*href=["']([^"']*)["'][^>]*>(.*?)</a>`).ReplaceAllString(mainText, "[$2]($1)")

				// Clean any remaining tags
				mainText = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(mainText, "")
				mainText = html.UnescapeString(mainText)
				mainText = strings.TrimSpace(mainText)

				// Add the item number and text
				if mainText != "" {
					result += fmt.Sprintf("%d. %s\n", i+1, mainText)
				} else {
					result += fmt.Sprintf("%d. \n", i+1)
				}

				// Process nested lists with indentation
				if strings.Contains(itemContent, "<ol") {
					nestedOlRegex := regexp.MustCompile(`(?s)<ol[^>]*>(.*?)</ol>`)
					if matches := nestedOlRegex.FindStringSubmatch(itemContent); len(matches) > 1 {
						// Convert nested ordered list
						nestedItems := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(matches[1], -1)
						for j, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<strong[^>]*>(.*?)</strong>`).ReplaceAllString(nestedItem[1], "**$1**")
								nestedText = regexp.MustCompile(`<b[^>]*>(.*?)</b>`).ReplaceAllString(nestedText, "**$1**")
								nestedText = regexp.MustCompile(`<em[^>]*>(.*?)</em>`).ReplaceAllString(nestedText, "*$1*")
								nestedText = regexp.MustCompile(`<i[^>]*>(.*?)</i>`).ReplaceAllString(nestedText, "*$1*")
								nestedText = regexp.MustCompile(`<a[^>]*href=["']([^"']*)["'][^>]*>(.*?)</a>`).ReplaceAllString(nestedText, "[$2]($1)")
								nestedText = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedText, "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								result += fmt.Sprintf("    %d. %s\n", j+1, nestedText)
							}
						}
					}
				} else if strings.Contains(itemContent, "<ul") {
					nestedUlRegex := regexp.MustCompile(`(?s)<ul[^>]*>(.*?)</ul>`)
					if matches := nestedUlRegex.FindStringSubmatch(itemContent); len(matches) > 1 {
						// Convert nested unordered list
						nestedItems := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(matches[1], -1)
						for _, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<strong[^>]*>(.*?)</strong>`).ReplaceAllString(nestedItem[1], "**$1**")
								nestedText = regexp.MustCompile(`<b[^>]*>(.*?)</b>`).ReplaceAllString(nestedText, "**$1**")
								nestedText = regexp.MustCompile(`<em[^>]*>(.*?)</em>`).ReplaceAllString(nestedText, "*$1*")
								nestedText = regexp.MustCompile(`<i[^>]*>(.*?)</i>`).ReplaceAllString(nestedText, "*$1*")
								nestedText = regexp.MustCompile(`<a[^>]*href=["']([^"']*)["'][^>]*>(.*?)</a>`).ReplaceAllString(nestedText, "[$2]($1)")
								nestedText = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedText, "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								result += fmt.Sprintf("    - %s\n", nestedText)
							}
						}
					}
				}
			}
		}
		return result
	})

	// Then handle unordered lists (ul)
	ulRegex := regexp.MustCompile(`(?s)<ul[^>]*>(.*?)</ul>`)
	content = ulRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		liRegex := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`)
		items := liRegex.FindAllStringSubmatch(match, -1)

		result := "\n"
		for _, item := range items {
			if len(item) > 1 {
				// Handle nested lists within this item
				itemContent := item[1]

				// Extract text before any nested list
				mainText := itemContent
				if strings.Contains(itemContent, "<ul") {
					idx := strings.Index(itemContent, "<ul")
					if idx > 0 {
						mainText = itemContent[:idx]
					} else {
						mainText = ""
					}
				} else if strings.Contains(itemContent, "<ol") {
					idx := strings.Index(itemContent, "<ol")
					if idx > 0 {
						mainText = itemContent[:idx]
					} else {
						mainText = ""
					}
				}

				// Convert inline formatting to markdown
				mainText = regexp.MustCompile(`<strong[^>]*>(.*?)</strong>`).ReplaceAllString(mainText, "**$1**")
				mainText = regexp.MustCompile(`<b[^>]*>(.*?)</b>`).ReplaceAllString(mainText, "**$1**")
				mainText = regexp.MustCompile(`<em[^>]*>(.*?)</em>`).ReplaceAllString(mainText, "*$1*")
				mainText = regexp.MustCompile(`<i[^>]*>(.*?)</i>`).ReplaceAllString(mainText, "*$1*")
				mainText = regexp.MustCompile(`<a[^>]*href=["']([^"']*)["'][^>]*>(.*?)</a>`).ReplaceAllString(mainText, "[$2]($1)")

				// Clean any remaining tags
				mainText = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(mainText, "")
				mainText = html.UnescapeString(mainText)
				mainText = strings.TrimSpace(mainText)

				// Add bullet point and text
				if mainText != "" {
					result += fmt.Sprintf("- %s\n", mainText)
				} else {
					result += "- \n"
				}

				// Process nested lists with indentation
				if strings.Contains(itemContent, "<ol") {
					nestedOlRegex := regexp.MustCompile(`(?s)<ol[^>]*>(.*?)</ol>`)
					if matches := nestedOlRegex.FindStringSubmatch(itemContent); len(matches) > 1 {
						// Convert nested ordered list
						nestedItems := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(matches[1], -1)
						for j, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<strong[^>]*>(.*?)</strong>`).ReplaceAllString(nestedItem[1], "**$1**")
								nestedText = regexp.MustCompile(`<b[^>]*>(.*?)</b>`).ReplaceAllString(nestedText, "**$1**")
								nestedText = regexp.MustCompile(`<em[^>]*>(.*?)</em>`).ReplaceAllString(nestedText, "*$1*")
								nestedText = regexp.MustCompile(`<i[^>]*>(.*?)</i>`).ReplaceAllString(nestedText, "*$1*")
								nestedText = regexp.MustCompile(`<a[^>]*href=["']([^"']*)["'][^>]*>(.*?)</a>`).ReplaceAllString(nestedText, "[$2]($1)")
								nestedText = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedText, "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								result += fmt.Sprintf("    %d. %s\n", j+1, nestedText)
							}
						}
					}
				} else if strings.Contains(itemContent, "<ul") {
					nestedUlRegex := regexp.MustCompile(`(?s)<ul[^>]*>(.*?)</ul>`)
					if matches := nestedUlRegex.FindStringSubmatch(itemContent); len(matches) > 1 {
						// Convert nested unordered list
						nestedItems := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(matches[1], -1)
						for _, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<strong[^>]*>(.*?)</strong>`).ReplaceAllString(nestedItem[1], "**$1**")
								nestedText = regexp.MustCompile(`<b[^>]*>(.*?)</b>`).ReplaceAllString(nestedText, "**$1**")
								nestedText = regexp.MustCompile(`<em[^>]*>(.*?)</em>`).ReplaceAllString(nestedText, "*$1*")
								nestedText = regexp.MustCompile(`<i[^>]*>(.*?)</i>`).ReplaceAllString(nestedText, "*$1*")
								nestedText = regexp.MustCompile(`<a[^>]*href=["']([^"']*)["'][^>]*>(.*?)</a>`).ReplaceAllString(nestedText, "[$2]($1)")
								nestedText = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedText, "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								result += fmt.Sprintf("    - %s\n", nestedText)
							}
						}
					}
				}
			}
		}
		return result
	})

	// Clean up any remaining HTML tags
	content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")

	// Clean up any excessive newlines
	content = regexp.MustCompile(`\n{3,}`).ReplaceAllString(content, "\n\n")

	return strings.TrimSpace(content)
}

// NestedListsToPlainText converts nested HTML lists to plain text
func NestedListsToPlainText(content string) string {
	// First sanitize the content
	content = SanitizeWordPressContent(content)

	// Replace any remaining HTML tags with appropriate text and structure

	// First handle ordered lists (ol)
	olRegex := regexp.MustCompile(`(?s)<ol[^>]*>(.*?)</ol>`)
	content = olRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		liRegex := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`)
		items := liRegex.FindAllStringSubmatch(match, -1)

		result := "\n"
		for i, item := range items {
			if len(item) > 1 {
				// Handle nested lists within this item
				itemContent := item[1]

				// Extract text before any nested list
				mainText := itemContent
				if strings.Contains(itemContent, "<ul") {
					idx := strings.Index(itemContent, "<ul")
					if idx > 0 {
						mainText = itemContent[:idx]
					} else {
						mainText = ""
					}
				} else if strings.Contains(itemContent, "<ol") {
					idx := strings.Index(itemContent, "<ol")
					if idx > 0 {
						mainText = itemContent[:idx]
					} else {
						mainText = ""
					}
				}

				// Clean the text
				mainText = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(mainText, "")
				mainText = html.UnescapeString(mainText)
				mainText = strings.TrimSpace(mainText)

				// Add the item number and text
				if mainText != "" {
					result += fmt.Sprintf("%d. %s\n", i+1, mainText)
				} else {
					result += fmt.Sprintf("%d. \n", i+1)
				}

				// Process nested lists with indentation
				if strings.Contains(itemContent, "<ol") {
					nestedOlRegex := regexp.MustCompile(`(?s)<ol[^>]*>(.*?)</ol>`)
					if matches := nestedOlRegex.FindStringSubmatch(itemContent); len(matches) > 1 {
						// Convert nested ordered list
						nestedItems := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(matches[1], -1)
						for j, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedItem[1], "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								result += fmt.Sprintf("    %d.%d. %s\n", i+1, j+1, nestedText)
							}
						}
					}
				} else if strings.Contains(itemContent, "<ul") {
					nestedUlRegex := regexp.MustCompile(`(?s)<ul[^>]*>(.*?)</ul>`)
					if matches := nestedUlRegex.FindStringSubmatch(itemContent); len(matches) > 1 {
						// Convert nested unordered list
						nestedItems := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(matches[1], -1)
						for _, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedItem[1], "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								result += fmt.Sprintf("    • %s\n", nestedText)
							}
						}
					}
				}
			}
		}
		return result
	})

	// Then handle unordered lists (ul)
	ulRegex := regexp.MustCompile(`(?s)<ul[^>]*>(.*?)</ul>`)
	content = ulRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		liRegex := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`)
		items := liRegex.FindAllStringSubmatch(match, -1)

		result := "\n"
		for _, item := range items {
			if len(item) > 1 {
				// Handle nested lists within this item
				itemContent := item[1]

				// Extract text before any nested list
				mainText := itemContent
				if strings.Contains(itemContent, "<ul") {
					idx := strings.Index(itemContent, "<ul")
					if idx > 0 {
						mainText = itemContent[:idx]
					} else {
						mainText = ""
					}
				} else if strings.Contains(itemContent, "<ol") {
					idx := strings.Index(itemContent, "<ol")
					if idx > 0 {
						mainText = itemContent[:idx]
					} else {
						mainText = ""
					}
				}

				// Clean the text
				mainText = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(mainText, "")
				mainText = html.UnescapeString(mainText)
				mainText = strings.TrimSpace(mainText)

				// Add bullet point and text
				if mainText != "" {
					result += fmt.Sprintf("• %s\n", mainText)
				} else {
					result += "• \n"
				}

				// Process nested lists with indentation
				if strings.Contains(itemContent, "<ol") {
					nestedOlRegex := regexp.MustCompile(`(?s)<ol[^>]*>(.*?)</ol>`)
					if matches := nestedOlRegex.FindStringSubmatch(itemContent); len(matches) > 1 {
						// Convert nested ordered list
						nestedItems := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(matches[1], -1)
						for j, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedItem[1], "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								result += fmt.Sprintf("    %d. %s\n", j+1, nestedText)
							}
						}
					}
				} else if strings.Contains(itemContent, "<ul") {
					nestedUlRegex := regexp.MustCompile(`(?s)<ul[^>]*>(.*?)</ul>`)
					if matches := nestedUlRegex.FindStringSubmatch(itemContent); len(matches) > 1 {
						// Convert nested unordered list
						nestedItems := regexp.MustCompile(`(?s)<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(matches[1], -1)
						for _, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedItem[1], "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								result += fmt.Sprintf("    • %s\n", nestedText)
							}
						}
					}
				}
			}
		}
		return result
	})

	// Clean up any remaining HTML tags
	content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")

	// Clean up any excessive newlines
	content = regexp.MustCompile(`\n{3,}`).ReplaceAllString(content, "\n\n")

	return strings.TrimSpace(content)
}

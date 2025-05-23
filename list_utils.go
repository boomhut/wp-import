package wpimport

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// improvedConvertUnorderedLists converts HTML unordered lists to markdown with nested list support
func improvedConvertUnorderedLists(content string) string {
	// First, mark nested lists for special handling
	content = regexp.MustCompile(`<li[^>]*>(.*?)<ul[^>]*>(.*?)</ul>(.*?)</li>`).
		ReplaceAllString(content, `<li data-has-child="true">$1<ul>$2</ul>$3</li>`)

	// Now handle outer lists
	re := regexp.MustCompile(`<ul[^>]*>(.*?)</ul>`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		itemRe := regexp.MustCompile(`<li[^>]*?(?:data-has-child="true")?(.*?)</li>`)
		items := itemRe.FindAllStringSubmatch(match, -1)

		result := "\n" // Start with a newline for proper Markdown spacing
		for _, item := range items {
			if len(item) > 1 {
				itemContent := item[1]

				// Check if this item has nested lists
				if strings.Contains(itemContent, "<ul>") {
					// Process content before the nested list
					beforeList := itemContent
					if idx := strings.Index(itemContent, "<ul>"); idx != -1 {
						beforeList = itemContent[:idx]
					}

					// Clean the content before the nested list
					beforeList = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(beforeList, "")
					beforeList = html.UnescapeString(beforeList)
					beforeList = strings.TrimSpace(beforeList)

					// Start the item
					result += fmt.Sprintf("- %s\n", beforeList)

					// Process nested lists with indentation
					nestedContent := ""
					if nestedMatch := regexp.MustCompile(`<ul>(.*?)</ul>`).FindStringSubmatch(itemContent); len(nestedMatch) > 1 {
						// Convert the nested list items with indentation
						nestedItems := regexp.MustCompile(`<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(nestedMatch[1], -1)
						for _, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedItem[1], "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								// Add indentation for nested items
								nestedContent += fmt.Sprintf("  - %s\n", nestedText)
							}
						}
						result += nestedContent
					}

					// Process content after the nested list if any
					afterList := ""
					if idx := strings.LastIndex(itemContent, "</ul>"); idx != -1 && idx+5 < len(itemContent) {
						afterList = itemContent[idx+5:]
						afterList = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(afterList, "")
						afterList = html.UnescapeString(afterList)
						afterList = strings.TrimSpace(afterList)
						if afterList != "" {
							result += fmt.Sprintf("- %s\n", afterList)
						}
					}
				} else {
					// Regular item without nested content
					itemContent = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(itemContent, "")
					itemContent = html.UnescapeString(itemContent)
					itemContent = strings.TrimSpace(itemContent)
					result += fmt.Sprintf("- %s\n", itemContent)
				}
			}
		}
		result += "\n" // Add newline after the list
		return result
	})

	return content
}

// improvedConvertOrderedLists converts HTML ordered lists to markdown with nested list support
func improvedConvertOrderedLists(content string) string {
	// First, mark nested lists for special handling
	content = regexp.MustCompile(`<li[^>]*>(.*?)<ol[^>]*>(.*?)</ol>(.*?)</li>`).
		ReplaceAllString(content, `<li data-has-child="true">$1<ol>$2</ol>$3</li>`)

	re := regexp.MustCompile(`<ol[^>]*>(.*?)</ol>`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		itemRe := regexp.MustCompile(`<li[^>]*?(?:data-has-child="true")?(.*?)</li>`)
		items := itemRe.FindAllStringSubmatch(match, -1)

		result := "\n" // Start with a newline for proper Markdown spacing
		for i, item := range items {
			if len(item) > 1 {
				itemContent := item[1]

				// Check if this item has nested lists
				if strings.Contains(itemContent, "<ol>") {
					// Process content before the nested list
					beforeList := itemContent
					if idx := strings.Index(itemContent, "<ol>"); idx != -1 {
						beforeList = itemContent[:idx]
					}

					// Clean the content before the nested list
					beforeList = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(beforeList, "")
					beforeList = html.UnescapeString(beforeList)
					beforeList = strings.TrimSpace(beforeList)

					// Start the item
					result += fmt.Sprintf("%d. %s\n", i+1, beforeList)

					// Process nested lists with indentation
					if nestedMatch := regexp.MustCompile(`<ol>(.*?)</ol>`).FindStringSubmatch(itemContent); len(nestedMatch) > 1 {
						// Convert the nested list items with indentation
						nestedItems := regexp.MustCompile(`<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(nestedMatch[1], -1)
						for j, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedItem[1], "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								// Add indentation for nested items
								result += fmt.Sprintf("    %d.%d. %s\n", i+1, j+1, nestedText)
							}
						}
					}

					// Process content after the nested list if any
					afterList := ""
					if idx := strings.LastIndex(itemContent, "</ol>"); idx != -1 && idx+5 < len(itemContent) {
						afterList = itemContent[idx+5:]
						afterList = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(afterList, "")
						afterList = html.UnescapeString(afterList)
						afterList = strings.TrimSpace(afterList)
						if afterList != "" {
							result += fmt.Sprintf("%d. %s\n", i+2, afterList)
						}
					}
				} else {
					// Regular item without nested content
					itemContent = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(itemContent, "")
					itemContent = html.UnescapeString(itemContent)
					itemContent = strings.TrimSpace(itemContent)
					result += fmt.Sprintf("%d. %s\n", i+1, itemContent)
				}
			}
		}
		result += "\n" // Add newline after the list
		return result
	})

	return content
}

// improvedConvertOrderedListsToPlainText converts HTML ordered lists to plain text with nested list support
func improvedConvertOrderedListsToPlainText(content string) string {
	// First, mark nested lists for special handling
	content = regexp.MustCompile(`<li[^>]*>(.*?)<ol[^>]*>(.*?)</ol>(.*?)</li>`).
		ReplaceAllString(content, `<li data-has-child="true">$1<ol>$2</ol>$3</li>`)

	re := regexp.MustCompile(`<ol[^>]*>(.*?)</ol>`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		itemRe := regexp.MustCompile(`<li[^>]*?(?:data-has-child="true")?(.*?)</li>`)
		items := itemRe.FindAllStringSubmatch(match, -1)

		result := ""
		for i, item := range items {
			if len(item) > 1 {
				itemContent := item[1]

				// Check if this item has nested lists
				if strings.Contains(itemContent, "<ol>") {
					// Process content before the nested list
					beforeList := itemContent
					if idx := strings.Index(itemContent, "<ol>"); idx != -1 {
						beforeList = itemContent[:idx]
					}

					// Clean the content before the nested list
					beforeList = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(beforeList, "")
					beforeList = html.UnescapeString(beforeList)
					beforeList = strings.TrimSpace(beforeList)

					// Start the item
					result += fmt.Sprintf("%d. %s\n", i+1, beforeList)

					// Process nested lists with indentation
					if nestedMatch := regexp.MustCompile(`<ol>(.*?)</ol>`).FindStringSubmatch(itemContent); len(nestedMatch) > 1 {
						// Convert the nested list items with indentation
						nestedItems := regexp.MustCompile(`<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(nestedMatch[1], -1)
						for j, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedItem[1], "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								// Add indentation for nested items
								result += fmt.Sprintf("    %d.%d. %s\n", i+1, j+1, nestedText)
							}
						}
					}

					// Process content after the nested list if any
					afterList := ""
					if idx := strings.LastIndex(itemContent, "</ol>"); idx != -1 && idx+5 < len(itemContent) {
						afterList = itemContent[idx+5:]
						afterList = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(afterList, "")
						afterList = html.UnescapeString(afterList)
						afterList = strings.TrimSpace(afterList)
						if afterList != "" {
							result += fmt.Sprintf("%d. %s\n", i+2, afterList)
						}
					}
				} else {
					// Regular item without nested content
					itemContent = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(itemContent, "")
					itemContent = html.UnescapeString(strings.TrimSpace(itemContent))
					result += fmt.Sprintf("%d. %s\n", i+1, itemContent)
				}
			}
		}
		return result
	})
	return content
}

// improvedConvertUnorderedListsToPlainText converts HTML unordered lists to plain text with nested list support
func improvedConvertUnorderedListsToPlainText(content string) string {
	// First, mark nested lists for special handling
	content = regexp.MustCompile(`<li[^>]*>(.*?)<ul[^>]*>(.*?)</ul>(.*?)</li>`).
		ReplaceAllString(content, `<li data-has-child="true">$1<ul>$2</ul>$3</li>`)

	re := regexp.MustCompile(`<ul[^>]*>(.*?)</ul>`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		itemRe := regexp.MustCompile(`<li[^>]*?(?:data-has-child="true")?(.*?)</li>`)
		items := itemRe.FindAllStringSubmatch(match, -1)

		result := ""
		for _, item := range items {
			if len(item) > 1 {
				itemContent := item[1]

				// Check if this item has nested lists
				if strings.Contains(itemContent, "<ul>") {
					// Process content before the nested list
					beforeList := itemContent
					if idx := strings.Index(itemContent, "<ul>"); idx != -1 {
						beforeList = itemContent[:idx]
					}

					// Clean the content before the nested list
					beforeList = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(beforeList, "")
					beforeList = html.UnescapeString(beforeList)
					beforeList = strings.TrimSpace(beforeList)

					// Start the item
					result += fmt.Sprintf("• %s\n", beforeList)

					// Process nested lists with indentation
					if nestedMatch := regexp.MustCompile(`<ul>(.*?)</ul>`).FindStringSubmatch(itemContent); len(nestedMatch) > 1 {
						// Convert the nested list items with indentation
						nestedItems := regexp.MustCompile(`<li[^>]*>(.*?)</li>`).FindAllStringSubmatch(nestedMatch[1], -1)
						for _, nestedItem := range nestedItems {
							if len(nestedItem) > 1 {
								nestedText := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(nestedItem[1], "")
								nestedText = html.UnescapeString(nestedText)
								nestedText = strings.TrimSpace(nestedText)
								// Add indentation for nested items
								result += fmt.Sprintf("  • %s\n", nestedText)
							}
						}
					}

					// Process content after the nested list if any
					afterList := ""
					if idx := strings.LastIndex(itemContent, "</ul>"); idx != -1 && idx+5 < len(itemContent) {
						afterList = itemContent[idx+5:]
						afterList = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(afterList, "")
						afterList = html.UnescapeString(afterList)
						afterList = strings.TrimSpace(afterList)
						if afterList != "" {
							result += fmt.Sprintf("• %s\n", afterList)
						}
					}
				} else {
					// Regular item without nested content
					itemContent = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(itemContent, "")
					itemContent = html.UnescapeString(strings.TrimSpace(itemContent))
					result += fmt.Sprintf("• %s\n", itemContent)
				}
			}
		}
		return result
	})
	return content
}

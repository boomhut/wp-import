package main

import (
	"fmt"
	"log"

	wpimport "github.com/boomhut/wp-import" // Replace with your actual import path
)

func main() {
	// Example usage
	filename := "wordpress-export.xml" // Replace with your actual file path

	site, err := wpimport.ParseWordPressXML(filename)
	if err != nil {
		log.Fatalf("Error parsing WordPress XML: %v", err)
	}

	// Display some basic information
	fmt.Printf("Site Title: %s\n", site.Channel.Title)
	fmt.Printf("Site URL: %s\n", site.Channel.Link)
	fmt.Printf("Description: %s\n", site.Channel.Description)
	fmt.Printf("WordPress Version: %s\n", site.Channel.WXRVersion)
	fmt.Printf("\n")

	// Count different content types
	posts := 0
	pages := 0
	attachments := 0
	other := 0

	for _, item := range site.Channel.Items {
		switch item.PostType {
		case "post":
			posts++
		case "page":
			pages++
		case "attachment":
			attachments++
		default:
			other++
		}
	}

	fmt.Printf("Content Summary:\n")
	fmt.Printf("- Posts: %d\n", posts)
	fmt.Printf("- Pages: %d\n", pages)
	fmt.Printf("- Attachments: %d\n", attachments)
	fmt.Printf("- Other: %d\n", other)
	fmt.Printf("- Authors: %d\n", len(site.Channel.Authors))
	fmt.Printf("- Categories: %d\n", len(site.Channel.Categories))
	fmt.Printf("- Tags: %d\n", len(site.Channel.Tags))
	fmt.Printf("- Custom Terms: %d\n", len(site.Channel.Terms))
	fmt.Printf("\n")

	// Analyze plugin data
	pluginData := site.AnalyzePluginData()
	if len(pluginData) > 0 {
		fmt.Println("Plugin Data Detected:")
		for plugin, data := range pluginData {
			fmt.Printf("- %s: %v\n", plugin, data)
		}
		fmt.Printf("\n")
	}
	// Example of content sanitization
	fmt.Println("Content Sanitization Example:")
	for _, item := range site.Channel.Items {
		if item.PostType == "post" && item.Status == "publish" && item.Content != "" {
			fmt.Printf("Original content length: %d characters\n", len(item.Content))

			// Sanitize WordPress content
			sanitized := wpimport.SanitizeWordPressContent(item.Content)
			fmt.Printf("Sanitized content length: %d characters\n", len(sanitized))

			// Convert to plain text
			plainText := wpimport.ConvertToPlainText(item.Content)
			fmt.Printf("Plain text length: %d characters\n", len(plainText))

			// Convert to Markdown
			markdown := wpimport.ConvertToMarkdown(item.Content)
			fmt.Printf("Markdown length: %d characters\n", len(markdown))

			// Show first 200 characters of each
			if len(sanitized) > 0 {
				fmt.Printf("Sanitized preview: %s...\n", truncateString(sanitized, 500))
			}
			if len(plainText) > 0 {
				fmt.Printf("Plain text preview: %s...\n", truncateString(plainText, 500))
			}
			if len(markdown) > 0 {
				fmt.Printf("Markdown preview: %s...\n", truncateString(markdown, 500))
			}
			break // Just show first post as example
		}
	}

	// Analyze styles and theme data
	customStyles := site.GetCustomStyles()
	if len(customStyles) > 0 {
		fmt.Println("Custom Styles Found:")
		for styleType, content := range customStyles {
			fmt.Printf("- %s: %d characters\n", styleType, len(content))
		}
		fmt.Printf("\n")
	}

	// Analyze page builder usage
	builderData := site.GetPageBuilderData()
	if len(builderData) > 0 {
		fmt.Println("Page Builder Data:")
		for builder, data := range builderData {
			fmt.Printf("- %s: %v\n", builder, data)
		}
		fmt.Printf("\n")
	}

	// Theme information
	themeInfo := site.GetThemeInfo()
	if len(themeInfo) > 0 {
		fmt.Println("Theme Information:")
		for key, value := range themeInfo {
			if len(value) > 100 {
				fmt.Printf("- %s: %d characters of data\n", key, len(value))
			} else {
				fmt.Printf("- %s: %s\n", key, value)
			}
		}
		fmt.Printf("\n")
	}

	// Display first few published posts
	fmt.Println("Recent Published Posts:")
	count := 0
	for _, item := range site.Channel.Items {
		if item.PostType == "post" && item.Status == "publish" && count < 5 {
			fmt.Printf("- %s (ID: %d)\n", item.Title, item.PostID)
			if item.PostDate != "" {
				if date, err := wpimport.ParseWordPressDate(item.PostDate); err == nil {
					fmt.Printf("  Date: %s\n", date.Format("2006-01-02 15:04:05"))
				}
			}
			fmt.Printf("  Categories: ")
			for i, cat := range item.Categories {
				if cat.Domain == "category" {
					if i > 0 {
						fmt.Print(", ")
					}
					fmt.Print(cat.Name)
				}
			}
			fmt.Println()
			count++
		}
	}

	// Display specific post by ID
	postID := 3144 // Replace with actual post ID
	post := site.GetPostByID(postID)
	if post != nil {
		fmt.Printf("\nPost Details (ID: %d):\n", post.PostID)
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Link: %s\n", post.Link)
		fmt.Printf("Published Date: %s\n", post.PubDate)
		fmt.Printf("Creator: %s\n", post.Creator)
		fmt.Printf("Content Length: %d characters\n", len(post.Content))

		// display content preview
		if len(post.Content) > 0 {
			fmt.Printf("Content Preview: %s...\n\n\n", truncateString(post.Content, 1700))
			// clean content
			cleanContent := wpimport.ConvertToPlainText(post.Content)
			fmt.Printf("Cleaned Content Preview:\n%s\n\n\n", truncateString(cleanContent, 1700))
			// Markdown content
			markdownContent := wpimport.ConvertToMarkdown(post.Content)
			fmt.Printf("Markdown Content Preview:\n%s\n\n\n", truncateString(markdownContent, 1700))

			// clean HTML content
			cleanHTMLContent := wpimport.SanitizeWordPressContent(post.Content)
			fmt.Printf("Cleaned HTML Content Preview:\n%s\n\n\n", truncateString(cleanHTMLContent, 1700))

		}

	} else {
		fmt.Printf("\nPost with ID %d not found.\n", postID)
	}
}

// Helper function to truncate string for preview
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

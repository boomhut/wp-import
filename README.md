- [WordPress Import Package (wp-import)](#wordpress-import-package-wp-import)
  - [Features](#features)
  - [Installation](#installation)
    - [Requirements](#requirements)
    - [Compatibility](#compatibility)
  - [Quick Start](#quick-start)
  - [Core Functionality](#core-functionality)
    - [Data Parsing](#data-parsing)
    - [Content Processing](#content-processing)
    - [Data Retrieval](#data-retrieval)
    - [Data Analysis](#data-analysis)
  - [Data Types](#data-types)
  - [Content Conversion Examples](#content-conversion-examples)
    - [Converting HTML to Plain Text](#converting-html-to-plain-text)
    - [Converting HTML to Markdown](#converting-html-to-markdown)
    - [Memory Usage](#memory-usage)
    - [Optimization Tips](#optimization-tips)
  - [Enhanced List Formatting](#enhanced-list-formatting)
    - [For Plain Text Conversion](#for-plain-text-conversion)
    - [For Markdown Conversion](#for-markdown-conversion)
    - [For Table Conversion](#for-table-conversion)
  - [How It Works](#how-it-works)
    - [HTML to Plain Text Conversion](#html-to-plain-text-conversion)
    - [HTML to Markdown Conversion](#html-to-markdown-conversion)
    - [Regular Expression Patterns](#regular-expression-patterns)
  - [Advanced Usage](#advanced-usage)
    - [Custom Conversion Options](#custom-conversion-options)
    - [Processing Large Exports](#processing-large-exports)
  - [Troubleshooting](#troubleshooting)
    - [Common Issues](#common-issues)
      - [Parsing Errors with Large XML Files](#parsing-errors-with-large-xml-files)
      - [Memory Usage Considerations](#memory-usage-considerations)
      - [Handling Complex WordPress Shortcodes](#handling-complex-wordpress-shortcodes)
  - [Contributing](#contributing)
    - [Reporting Issues](#reporting-issues)
  - [License](#license)

# WordPress Import Package (wp-import)

[![Go Reference](https://pkg.go.dev/badge/github.com/boomhut/wp-import.svg)](https://pkg.go.dev/github.com/boomhut/wp-import)
[![Go Report Card](https://goreportcard.com/badge/github.com/boomhut/wp-import)](https://goreportcard.com/report/github.com/boomhut/wp-import)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/boomhut/wp-import)](https://github.com/boomhut/wp-import)

A Go package for parsing, analyzing, and converting WordPress export XML files. This library provides robust functionality for accessing WordPress content programmatically and transforming it into various formats.

## Features

- **Complete WordPress XML Parsing**: Parse WordPress export files into structured Go types
- **Content Format Conversion**:
  - Convert WordPress HTML content to clean plain text
  - Convert WordPress HTML content to properly formatted Markdown
  - Sanitize WordPress content (remove Gutenberg blocks, clean HTML)
- **WordPress Data Analysis**:
  - Extract metadata from posts and pages
  - Analyze custom post types, plugins, and themes
  - Gather information about media attachments
- **Advanced Content Processing**:
  - Handle lists (ordered and unordered) with proper formatting:
    - Ordered lists convert to numbered points (1., 2., etc.)
    - Unordered lists convert to bullet points (•)
    - Proper spacing before and after lists
  - Convert HTML tables to Markdown tables with header separators
  - Proper formatting of headings, code blocks, blockquotes, and images
  - Support for inline formatting (bold, italic, strikethrough)

## Installation

```bash
go get github.com/boomhut/wp-import
```

### Requirements

- Go 1.24 or higher
- No external dependencies outside the Go standard library

### Compatibility

This package is compatible with:
- Standard WordPress export files (WXR format)
- WordPress exports from version 4.0 and newer
- Both single-site and multisite exports

## Quick Start

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/boomhut/wp-import"
)

func main() {
	// Parse WordPress export file
	site, err := wpimport.ParseWordPressXML("wordpress-export.xml")
	if err != nil {
		log.Fatalf("Error parsing WordPress XML: %v", err)
	}
	
	// Display basic site info
	fmt.Printf("Site Title: %s\n", site.Channel.Title)
	fmt.Printf("Site URL: %s\n", site.Channel.Link)
	
	// Get posts by type
	posts := site.GetPostsByType("post")
	fmt.Printf("Found %d posts\n", len(posts))
	
	// Convert content to Markdown
	if len(posts) > 0 {
		markdown := wpimport.ConvertToMarkdown(posts[0].Content)
		fmt.Printf("Markdown preview: %s...\n", truncateString(markdown, 200))
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
```

## Core Functionality

### Data Parsing

- `ParseWordPressXML(filename string) (*WordPressSite, error)` - Parse WordPress export XML file
- `ParseWordPressDate(dateStr string) (time.Time, error)` - Parse WordPress date format

### Content Processing

- `SanitizeWordPressContent(content string) string` - Clean up WordPress content
- `CleanHTML(content string) string` - Sanitize HTML while preserving HTML structure
- `ConvertToPlainText(content string) string` - Convert HTML content to plain text
- `ConvertToMarkdown(content string) string` - Convert HTML content to Markdown
- `convertOrderedListsToPlainText(content string) string` - Convert `<ol>` lists to numbered plain text
- `convertUnorderedListsToPlainText(content string) string` - Convert `<ul>` lists to bullet points
- `convertOrderedLists(content string) string` - Convert `<ol>` lists to Markdown format
- `convertUnorderedLists(content string) string` - Convert `<ul>` lists to Markdown format
- `convertTables(content string) string` - Convert HTML tables to Markdown table format

### Data Retrieval

- `GetPostsByType(postType string) []Item` - Get posts of a specific type
- `GetPublishedPosts() []Item` - Get only published posts
- `GetPostByID(id int) *Item` - Find a post by ID
- `GetAuthors() []Author` - Get all authors
- `GetCustomTerms() []Term` - Get custom taxonomy terms
- `GetAttachmentURLs() []string` - Get all attachment URLs

### Data Analysis

- `AnalyzePluginData() map[string]interface{}` - Analyze plugin usage
- `GetCustomStyles() map[string]string` - Extract custom CSS and styles
- `GetPageBuilderData() map[string]interface{}` - Analyze page builder usage
- `GetThemeInfo() map[string]string` - Extract theme information

## Data Types

The package provides comprehensive types that map to WordPress export structures:

- `WordPressSite` - Root structure for the WordPress export
- `Channel` - Contains site information and all content items
- `Author` - WordPress user account information
- `Item` - Post, page, or other content type
- `Category` - WordPress category
- `Tag` - WordPress tag
- `Term` - Custom taxonomy term
- `PostMeta` - Custom fields and metadata
- `Comment` - Post comment
- `CommentMeta` - Comment metadata

## Content Conversion Examples

### Converting HTML to Plain Text

```go
htmlContent := `<p>This is a <strong>paragraph</strong> with <em>formatting</em>.</p>
<ul>
  <li>Bullet point 1</li>
  <li>Bullet point 2</li>
</ul>
<ol>
  <li>First ordered item</li>
  <li>Second ordered item</li>
</ol>`

plainText := wpimport.ConvertToPlainText(htmlContent)
fmt.Println(plainText)
```

Output:
```
This is a paragraph with formatting.

• Bullet point 1
• Bullet point 2

1. First ordered item
2. Second ordered item
```

### Converting HTML to Markdown

```go
htmlContent := `<h1>Heading</h1>
<p>This is a <strong>paragraph</strong> with <em>formatting</em>.</p>
<ul>
  <li>Bullet point</li>
  <li>Another bullet point</li>
</ul>
<ol>
  <li>First item</li>
  <li>Second item</li>
</ol>
<blockquote>This is a blockquote</blockquote>
<pre><code>function example() {
  return "This is a code block";
}</code></pre>
<table>
  <tr>
    <th>Header 1</th>
    <th>Header 2</th>
  </tr>
  <tr>
    <td>Cell 1</td>
    <td>Cell 2</td>
  </tr>
</table>`

markdown := wpimport.ConvertToMarkdown(htmlContent)
fmt.Println(markdown)
```

Output (after conversion to Markdown):
```markdown
# Heading

This is a **paragraph** with *formatting*.

- Bullet point
- Another bullet point

1. First item
2. Second item

> This is a blockquote

<!-- Code block converted from HTML -->
```
function example() {
  return "This is a code block";
}
```

| Header 1 | Header 2 |
| --- | --- |
| Cell 1 | Cell 2 |
```

### Sanitizing HTML while Preserving Structure

```go
htmlContent := `<p>This is <b>bold</b> text with a <font color="red">colored font</font> tag.</p>
<center>This text is centered</center>
<p class="wp-block-paragraph aligncenter" style="">Paragraph with WordPress classes</p>
<div data-wp-block="true" data-align="wide">Block with data attributes</div>
<ul><li>Bullet</li><li>Another <strong>bullet</strong> with <em>formatting</em></li></ul>`

cleanHTML := wpimport.CleanHTML(htmlContent)
fmt.Println(cleanHTML)
```

Output (after cleaning):
```html
<p>This is <strong>bold</strong> text with a <span>colored font</span> tag.</p>
<div style="text-align: center;">This text is centered</div>
<p>Paragraph with WordPress classes</p>
<div>Block with data attributes</div>
<ul>
  <li>Bullet</li>
  <li>Another <strong>bullet</strong> with <em>formatting</em></li>
</ul>
```

### Memory Usage

Memory usage is optimized for processing large WordPress exports:

- Parsing a 100MB WordPress export uses approximately 200-300MB RAM
- Converting 1000 posts to Markdown (average 10KB each) uses approximately 50-100MB RAM

### Optimization Tips

For optimal performance:

1. Process posts in batches or use goroutines for parallel processing
2. For extremely large exports, consider splitting the XML file
3. Use the conversion functions directly on individual posts rather than processing the entire content at once

## Enhanced List Formatting

The package includes specialized helper functions for properly converting different types of HTML lists:

### For Plain Text Conversion

```go
// Helper function to convert ordered lists to numbered plain text
func convertOrderedListsToPlainText(content string) string {
    // Extracts each list item and numbers them (1., 2., etc.)
    // Preserves list item contents while stripping HTML
}

// Helper function to convert unordered lists to bullet point plain text
func convertUnorderedListsToPlainText(content string) string {
    // Extracts each list item and adds bullet points (•)
    // Preserves list item contents while stripping HTML
}
```

### For Markdown Conversion

```go
// Helper function to convert ordered lists to markdown
func convertOrderedLists(content string) string {
    // Extracts and converts ordered lists to markdown format
    // Adds proper spacing before and after lists
    // Ensures correct numbering (1., 2., etc.)
}

// Helper function to convert unordered lists to markdown
func convertUnorderedLists(content string) string {
    // Extracts and converts unordered lists to markdown format
    // Uses proper markdown bullet point style (-)
    // Adds proper spacing before and after lists
}
```

### For Table Conversion

```go
// Helper function to convert HTML tables to markdown format
func convertTables(content string) string {
    // Processes HTML tables and converts them to markdown tables
    // Creates header row with separator
    // Handles cell content and escapes pipe characters
    // Maintains proper alignment and spacing
}
```

These functions work by using regular expressions to locate structured elements in the HTML content, extract their components, and format them according to the target format. They handle proper spacing and ensure that nested content is correctly processed.

## How It Works

### HTML Cleaning and Sanitization

The `CleanHTML` function provides a way to sanitize WordPress HTML content while preserving its structure:

1. **Remove WordPress-specific Elements**: 
   - Removes Gutenberg block comments
   - Cleans out empty paragraphs and unnecessary whitespace
   
2. **Fix HTML Structure Issues**:
   - Repairs unclosed or improperly nested tags
   - Fixes malformed list structures
   - Ensures proper HTML structure is maintained
   
3. **Modernize HTML**:
   - Updates deprecated tags to modern HTML5 equivalents
   - Converts `<center>` to styled divs
   - Converts `<font>` to spans
   - Converts `<b>` to `<strong>` and `<i>` to `<em>`
   
4. **Clean Up Attributes**:
   - Removes empty and WordPress-specific attributes
   - Cleans up unnecessary styling attributes
   - Removes data attributes that are WordPress-specific

### HTML to Plain Text Conversion

The `ConvertToPlainText` function works through these steps:

1. **Sanitize**: First, the WordPress content is cleaned up to remove non-standard HTML
2. **List Processing**: 
   - Ordered lists (`<ol>`) are converted to numbered text (1., 2., etc.)
   - Unordered lists (`<ul>`) are converted to bullet points (•)
3. **Tag Processing**:
   - `<br>` tags are replaced with newlines
   - `</p>` tags are replaced with double newlines
   - All other HTML tags are removed
4. **Entity Decoding**: All HTML entities are decoded (e.g., `&amp;` to `&`)
5. **Whitespace Cleanup**: Excessive whitespace is normalized

### HTML to Markdown Conversion

The `ConvertToMarkdown` function follows a more comprehensive process:

1. **Sanitize Content**: Remove WordPress-specific HTML and clean up
2. **Process Block Elements**:
   - Lists are converted first to avoid interference with other processing
   - Tables are processed into properly formatted Markdown tables
   - Headers, blockquotes, and code blocks are converted
3. **Process Inline Elements**:
   - Process links, images, and inline formatting
   - Handle text formatting (bold, italic, strikethrough)
4. **Final Cleanup**:
   - Remove any remaining HTML tags
   - Decode HTML entities
   - Normalize whitespace
   - Ensure proper spacing between elements

### Regular Expression Patterns

The conversion relies on carefully crafted regular expressions:

- **Multiline matching**: Uses `(?s)` flag to match across newlines
- **Non-greedy matching**: Uses `.*?` to avoid over-capturing
- **Attribute-aware matching**: Handles variations in HTML tag attributes

## Advanced Usage

### Custom Conversion Options

You can combine the package's functions for specialized conversion needs:

```go
// First sanitize, then perform custom transformations, then convert
content := wpimport.SanitizeWordPressContent(htmlContent)
// Apply your custom transformations to content
markdown := wpimport.ConvertToMarkdown(content)
```

### Processing Large Exports

For performance when processing many posts:

```go
site, err := wpimport.ParseWordPressXML("wordpress-export.xml")
if err != nil {
    log.Fatal(err)
}

// Process posts in parallel with goroutines
var wg sync.WaitGroup
posts := site.GetPostsByType("post")

for _, post := range posts {
    wg.Add(1)
    go func(p wpimport.Item) {
        defer wg.Done()
        
        // Process the post content
        markdown := wpimport.ConvertToMarkdown(p.Content)
        
        // Do something with the markdown...
    }(post)
}

wg.Wait()
```

## Troubleshooting

### Common Issues

#### Parsing Errors with Large XML Files

If you encounter errors parsing very large WordPress exports:

```go
// For large files, try increasing buffers
site, err := wpimport.ParseWordPressXML("large-wordpress-export.xml")
if err != nil {
    if strings.Contains(err.Error(), "token too large") {
        // Large file handling - split the XML or process in chunks
        log.Println("XML file too large, consider splitting")
    }
    log.Fatal(err)
}
```

#### Memory Usage Considerations

For very large WordPress sites:

```go
// Process posts in batches to manage memory
posts := site.GetPostsByType("post")
batchSize := 50
totalPosts := len(posts)

for i := 0; i < totalPosts; i += batchSize {
    end := i + batchSize
    if end > totalPosts {
        end = totalPosts
    }
    
    batch := posts[i:end]
    // Process this batch of posts
    processPostBatch(batch)
}
```

#### Handling Complex WordPress Shortcodes

If your WordPress content has complex shortcodes:

```go
// Apply more aggressive shortcode removal before conversion
content = wpimport.SanitizeWordPressContent(post.Content)
// Additional shortcode handling if needed
content = regexp.MustCompile(`\[.*?\]`).ReplaceAllString(content, "")
markdown := wpimport.ConvertToMarkdown(content)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Reporting Issues

If you encounter any bugs or have feature requests, please open an issue on the GitHub repository.

## License

MIT License. See the [LICENSE](LICENSE) file for details.
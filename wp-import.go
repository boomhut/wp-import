package wpimport

import (
	"encoding/xml"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"
	"time"
)

// WordPressSite represents the root structure of a WordPress export
type WordPressSite struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel contains the main site information and all posts
type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	PubDate       string `xml:"pubDate"`
	LastBuildDate string `xml:"lastBuildDate"`
	Generator     string `xml:"generator"`

	// WordPress specific fields
	WXRVersion  string `xml:"http://wordpress.org/export/1.2/ wxr_version"`
	BaseSiteURL string `xml:"http://wordpress.org/export/1.2/ base_site_url"`
	BaseBlogURL string `xml:"http://wordpress.org/export/1.2/ base_blog_url"`

	// Authors
	Authors []Author `xml:"http://wordpress.org/export/1.2/ author"`

	// Categories and Tags
	Categories []Category `xml:"http://wordpress.org/export/1.2/ category"`
	Tags       []Tag      `xml:"http://wordpress.org/export/1.2/ tag"`
	Terms      []Term     `xml:"http://wordpress.org/export/1.2/ term"`

	// Posts and Pages
	Items []Item `xml:"item"`
}

// Author represents WordPress user accounts
type Author struct {
	ID          int    `xml:"http://wordpress.org/export/1.2/ author_id"`
	Login       string `xml:"http://wordpress.org/export/1.2/ author_login"`
	Email       string `xml:"http://wordpress.org/export/1.2/ author_email"`
	DisplayName string `xml:"http://wordpress.org/export/1.2/ author_display_name"`
	FirstName   string `xml:"http://wordpress.org/export/1.2/ author_first_name"`
	LastName    string `xml:"http://wordpress.org/export/1.2/ author_last_name"`
}

// Term represents custom taxonomy terms
type Term struct {
	TermID      int    `xml:"http://wordpress.org/export/1.2/ term_id"`
	Taxonomy    string `xml:"http://wordpress.org/export/1.2/ term_taxonomy"`
	Slug        string `xml:"http://wordpress.org/export/1.2/ term_slug"`
	Parent      string `xml:"http://wordpress.org/export/1.2/ term_parent"`
	Name        string `xml:"http://wordpress.org/export/1.2/ term_name"`
	Description string `xml:"http://wordpress.org/export/1.2/ term_description"`
}

// Category represents WordPress categories
type Category struct {
	TermID   int    `xml:"http://wordpress.org/export/1.2/ term_id"`
	NiceName string `xml:"http://wordpress.org/export/1.2/ category_nicename"`
	Parent   string `xml:"http://wordpress.org/export/1.2/ category_parent"`
	Name     string `xml:"http://wordpress.org/export/1.2/ cat_name"`
}

// Tag represents WordPress tags
type Tag struct {
	TermID int    `xml:"http://wordpress.org/export/1.2/ term_id"`
	Slug   string `xml:"http://wordpress.org/export/1.2/ tag_slug"`
	Name   string `xml:"http://wordpress.org/export/1.2/ tag_name"`
}

// Item represents posts, pages, and other content types
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Creator     string `xml:"http://purl.org/dc/elements/1.1/ creator"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	Content     string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	Excerpt     string `xml:"http://wordpress.org/export/1.2/ post_excerpt"`

	// WordPress specific fields
	PostID       int    `xml:"http://wordpress.org/export/1.2/ post_id"`
	PostDate     string `xml:"http://wordpress.org/export/1.2/ post_date"`
	PostDateGMT  string `xml:"http://wordpress.org/export/1.2/ post_date_gmt"`
	PostName     string `xml:"http://wordpress.org/export/1.2/ post_name"`
	PostType     string `xml:"http://wordpress.org/export/1.2/ post_type"`
	Status       string `xml:"http://wordpress.org/export/1.2/ status"`
	PostParent   int    `xml:"http://wordpress.org/export/1.2/ post_parent"`
	MenuOrder    int    `xml:"http://wordpress.org/export/1.2/ menu_order"`
	PostPassword string `xml:"http://wordpress.org/export/1.2/ post_password"`
	IsSticky     int    `xml:"http://wordpress.org/export/1.2/ is_sticky"`

	// Categories and Tags for this post
	Categories []ItemCategory `xml:"category"`

	// Post metadata
	PostMeta []PostMeta `xml:"http://wordpress.org/export/1.2/ postmeta"`

	// Comments
	Comments []Comment `xml:"http://wordpress.org/export/1.2/ comment"`
}

// ItemCategory represents category/tag assignments for posts
type ItemCategory struct {
	Domain   string `xml:"domain,attr"`
	NiceName string `xml:"nicename,attr"`
	Name     string `xml:",chardata"`
}

// PostMeta represents custom fields and metadata
type PostMeta struct {
	Key   string `xml:"http://wordpress.org/export/1.2/ meta_key"`
	Value string `xml:"http://wordpress.org/export/1.2/ meta_value"`
}

// Comment represents post comments
type Comment struct {
	ID          int           `xml:"http://wordpress.org/export/1.2/ comment_id"`
	Author      string        `xml:"http://wordpress.org/export/1.2/ comment_author"`
	AuthorEmail string        `xml:"http://wordpress.org/export/1.2/ comment_author_email"`
	AuthorURL   string        `xml:"http://wordpress.org/export/1.2/ comment_author_url"`
	AuthorIP    string        `xml:"http://wordpress.org/export/1.2/ comment_author_IP"`
	Date        string        `xml:"http://wordpress.org/export/1.2/ comment_date"`
	DateGMT     string        `xml:"http://wordpress.org/export/1.2/ comment_date_gmt"`
	Content     string        `xml:"http://wordpress.org/export/1.2/ comment_content"`
	Approved    string        `xml:"http://wordpress.org/export/1.2/ comment_approved"`
	Type        string        `xml:"http://wordpress.org/export/1.2/ comment_type"`
	Parent      int           `xml:"http://wordpress.org/export/1.2/ comment_parent"`
	UserID      int           `xml:"http://wordpress.org/export/1.2/ comment_user_id"`
	CommentMeta []CommentMeta `xml:"http://wordpress.org/export/1.2/ commentmeta"`
}

// CommentMeta represents comment metadata
type CommentMeta struct {
	Key   string `xml:"http://wordpress.org/export/1.2/ meta_key"`
	Value string `xml:"http://wordpress.org/export/1.2/ meta_value"`
}

// ParseWordPressXML reads and parses a WordPress export XML file
func ParseWordPressXML(filename string) (*WordPressSite, error) {
	// Read the XML file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the XML
	var site WordPressSite
	err = xml.Unmarshal(data, &site)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &site, nil
}

// ParseWordPressDate parses WordPress date format
func ParseWordPressDate(dateStr string) (time.Time, error) {
	// WordPress typically uses "2006-01-02 15:04:05" format
	return time.Parse("2006-01-02 15:04:05", dateStr)
}

// Helper function to get all posts of a specific type
func (site *WordPressSite) GetPostsByType(postType string) []Item {
	var posts []Item
	for _, item := range site.Channel.Items {
		if item.PostType == postType {
			posts = append(posts, item)
		}
	}
	return posts
}

// Helper function to get published posts only
func (site *WordPressSite) GetPublishedPosts() []Item {
	var posts []Item
	for _, item := range site.Channel.Items {
		if item.Status == "publish" {
			posts = append(posts, item)
		}
	}
	return posts
}

// Helper function to find a post by ID
func (site *WordPressSite) GetPostByID(id int) *Item {
	for _, item := range site.Channel.Items {
		if item.PostID == id {
			return &item
		}
	}
	return nil
}

// Helper function to get all authors
func (site *WordPressSite) GetAuthors() []Author {
	return site.Channel.Authors
}

// Helper function to get author by ID
func (site *WordPressSite) GetAuthorByID(id int) *Author {
	for _, author := range site.Channel.Authors {
		if author.ID == id {
			return &author
		}
	}
	return nil
}

// Helper function to get posts by author
func (site *WordPressSite) GetPostsByAuthor(authorLogin string) []Item {
	var posts []Item
	for _, item := range site.Channel.Items {
		if item.Creator == authorLogin {
			posts = append(posts, item)
		}
	}
	return posts
}

// Helper function to get all custom taxonomy terms
func (site *WordPressSite) GetCustomTerms() []Term {
	return site.Channel.Terms
}

// Helper function to get terms by taxonomy
func (site *WordPressSite) GetTermsByTaxonomy(taxonomy string) []Term {
	var terms []Term
	for _, term := range site.Channel.Terms {
		if term.Taxonomy == taxonomy {
			terms = append(terms, term)
		}
	}
	return terms
}

// Helper function to get all attachment URLs from post content
func (site *WordPressSite) GetAttachmentURLs() []string {
	var urls []string
	for _, item := range site.Channel.Items {
		if item.PostType == "attachment" {
			if item.GUID != "" {
				urls = append(urls, item.GUID)
			}
		}
	}
	return urls
}

// Helper function to get post meta by key
func (item *Item) GetMetaValue(key string) string {
	for _, meta := range item.PostMeta {
		if meta.Key == key {
			return meta.Value
		}
	}
	return ""
}

// Helper function to check if post has specific meta key
func (item *Item) HasMeta(key string) bool {
	for _, meta := range item.PostMeta {
		if meta.Key == key {
			return true
		}
	}
	return false
}

// Helper function to extract WooCommerce order data (if present in meta)
func (item *Item) GetWooCommerceOrderData() map[string]string {
	orderData := make(map[string]string)
	if item.PostType != "shop_order" {
		return orderData
	}

	// Common WooCommerce order meta keys
	wooKeys := []string{
		"_order_total", "_order_tax", "_order_shipping", "_order_discount",
		"_billing_first_name", "_billing_last_name", "_billing_email",
		"_shipping_first_name", "_shipping_last_name", "_payment_method",
		"_order_currency", "_customer_user", "_order_key",
	}

	for _, key := range wooKeys {
		if value := item.GetMetaValue(key); value != "" {
			orderData[key] = value
		}
	}

	return orderData
}

// Helper function to get all plugin-related meta data
func (item *Item) GetPluginMeta() map[string]string {
	pluginMeta := make(map[string]string)

	// Common plugin meta prefixes
	pluginPrefixes := []string{
		"_yoast_", "_aioseop_", "_elementor_", "_wpb_", "_vc_",
		"_wc_", "_edd_", "_bbp_", "_bp_", "_tribe_",
	}

	for _, meta := range item.PostMeta {
		for _, prefix := range pluginPrefixes {
			if len(meta.Key) > len(prefix) && meta.Key[:len(prefix)] == prefix {
				pluginMeta[meta.Key] = meta.Value
				break
			}
		}
	}

	return pluginMeta
}

// Helper function to detect custom post types that might be from plugins
func (site *WordPressSite) GetCustomPostTypes() map[string]int {
	postTypes := make(map[string]int)

	// WordPress core post types
	coreTypes := map[string]bool{
		"post": true, "page": true, "attachment": true, "revision": true,
		"nav_menu_item": true, "custom_css": true, "customize_changeset": true,
	}

	for _, item := range site.Channel.Items {
		if !coreTypes[item.PostType] {
			postTypes[item.PostType]++
		}
	}

	return postTypes
}

// Helper function to analyze what plugin data might be present
func (site *WordPressSite) AnalyzePluginData() map[string]interface{} {
	analysis := make(map[string]interface{})

	// Check for WooCommerce
	wooOrders := 0
	wooProducts := 0
	for _, item := range site.Channel.Items {
		switch item.PostType {
		case "shop_order":
			wooOrders++
		case "product":
			wooProducts++
		}
	}
	if wooOrders > 0 || wooProducts > 0 {
		analysis["woocommerce"] = map[string]int{
			"orders":   wooOrders,
			"products": wooProducts,
		}
	}

	// Check for Easy Digital Downloads
	eddDownloads := 0
	eddPayments := 0
	for _, item := range site.Channel.Items {
		switch item.PostType {
		case "download":
			eddDownloads++
		case "edd_payment":
			eddPayments++
		}
	}
	if eddDownloads > 0 || eddPayments > 0 {
		analysis["easy_digital_downloads"] = map[string]int{
			"downloads": eddDownloads,
			"payments":  eddPayments,
		}
	}

	// Check for bbPress
	bbpForums := 0
	bbpTopics := 0
	bbpReplies := 0
	for _, item := range site.Channel.Items {
		switch item.PostType {
		case "forum":
			bbpForums++
		case "topic":
			bbpTopics++
		case "reply":
			bbpReplies++
		}
	}
	if bbpForums > 0 || bbpTopics > 0 || bbpReplies > 0 {
		analysis["bbpress"] = map[string]int{
			"forums":  bbpForums,
			"topics":  bbpTopics,
			"replies": bbpReplies,
		}
	}

	// Get all custom post types
	customTypes := site.GetCustomPostTypes()
	if len(customTypes) > 0 {
		analysis["custom_post_types"] = customTypes
	}

	return analysis
}

// Helper function to extract custom CSS and styles
func (site *WordPressSite) GetCustomStyles() map[string]string {
	styles := make(map[string]string)

	for _, item := range site.Channel.Items {
		// WordPress Custom CSS (Appearance > Customize > Additional CSS)
		if item.PostType == "custom_css" {
			styles["custom_css"] = item.Content
		}

		// Check for inline styles in content
		if item.PostType == "page" || item.PostType == "post" {
			// Look for style attributes in content
			if hasInlineStyles(item.Content) {
				styles[fmt.Sprintf("inline_styles_%d", item.PostID)] = extractInlineStyles(item.Content)
			}
		}

		// Check for Elementor styles (if using Elementor page builder)
		if elementorCSS := item.GetMetaValue("_elementor_css"); elementorCSS != "" {
			styles[fmt.Sprintf("elementor_css_%d", item.PostID)] = elementorCSS
		}

		// Check for theme customizer settings
		if item.PostType == "customize_changeset" {
			styles[fmt.Sprintf("customizer_%d", item.PostID)] = item.Content
		}
	}

	return styles
}

// Helper function to extract theme information from meta
func (site *WordPressSite) GetThemeInfo() map[string]string {
	themeInfo := make(map[string]string)

	for _, item := range site.Channel.Items {
		// Check for theme-related meta
		if template := item.GetMetaValue("_wp_page_template"); template != "" {
			themeInfo[fmt.Sprintf("page_template_%d", item.PostID)] = template
		}

		// Check for theme customizer data
		if item.PostType == "customize_changeset" && item.Status == "publish" {
			themeInfo["active_customizer_settings"] = item.Content
		}
	}

	return themeInfo
}

// Helper function to get page builder data
func (site *WordPressSite) GetPageBuilderData() map[string]interface{} {
	builders := make(map[string]interface{})

	elementorPages := 0
	gutenbergBlocks := 0
	vcPages := 0

	for _, item := range site.Channel.Items {
		// Elementor
		if item.HasMeta("_elementor_edit_mode") || item.HasMeta("_elementor_data") {
			elementorPages++
		}

		// Gutenberg blocks (look for block content)
		if containsGutenbergBlocks(item.Content) {
			gutenbergBlocks++
		}

		// Visual Composer
		if item.HasMeta("_wpb_vc_js_status") {
			vcPages++
		}
	}

	if elementorPages > 0 {
		builders["elementor"] = map[string]int{"pages": elementorPages}
	}
	if gutenbergBlocks > 0 {
		builders["gutenberg"] = map[string]int{"pages": gutenbergBlocks}
	}
	if vcPages > 0 {
		builders["visual_composer"] = map[string]int{"pages": vcPages}
	}

	return builders
}

// Helper function to check for inline styles
func hasInlineStyles(content string) bool {
	return strings.Contains(content, "style=") || strings.Contains(content, "<style")
}

// Helper function to extract inline styles (basic implementation)
func extractInlineStyles(content string) string {
	// This is a simplified extraction - in practice, you'd want more robust parsing
	styles := ""

	// Look for <style> tags
	styleStart := strings.Index(content, "<style")
	for styleStart != -1 {
		styleEndTag := strings.Index(content[styleStart:], "</style>")
		if styleEndTag != -1 {
			styleContent := content[styleStart : styleStart+styleEndTag+8]
			styles += styleContent + "\n"
		}
		styleStart = strings.Index(content[styleStart+1:], "<style")
		if styleStart != -1 {
			styleStart += styleStart + 1
		}
	}

	return styles
}

// Helper function to check for Gutenberg blocks
func containsGutenbergBlocks(content string) bool {
	return strings.Contains(content, "<!-- wp:") || strings.Contains(content, "wp-block-")
}

// SanitizeWordPressContent removes Gutenberg block comments and cleans HTML
func SanitizeWordPressContent(content string) string {
	// Remove Gutenberg block comments
	content = removeGutenbergComments(content)

	// Clean up extra whitespace and line breaks
	content = cleanWhitespace(content)

	// Remove empty paragraphs
	content = removeEmptyParagraphs(content)

	return content
}

// CleanHTML sanitizes WordPress HTML content while preserving HTML structure
// This returns valid, cleaned HTML that can be used for display or further processing
func CleanHTML(content string) string {
	// First, use the base sanitize function to remove Gutenberg blocks
	content = SanitizeWordPressContent(content)

	// Fix common HTML issues in WordPress content

	// Fix unclosed or improperly nested tags (simple cases)
	content = fixUnclosedTags(content)

	// Remove empty or unnecessary attributes
	content = cleanAttributes(content)

	// Convert deprecated tags to modern HTML5 equivalents
	content = updateDeprecatedTags(content)

	// Fix malformed lists
	content = fixListStructure(content)

	return content
}

// Helper function to fix unclosed or improperly nested tags
func fixUnclosedTags(content string) string {
	// This is a simplified approach - for a complete solution,
	// a proper HTML parser would be needed

	// Fix common unclosed tags
	unclosedTags := []string{"p", "div", "span", "li", "ol", "ul", "strong", "em", "a"}
	for _, tag := range unclosedTags {
		openCount := strings.Count(content, "<"+tag)
		closeCount := strings.Count(content, "</"+tag)

		// If there are more opening than closing tags, add missing closing tags
		if openCount > closeCount {
			for i := 0; i < (openCount - closeCount); i++ {
				content = content + "</" + tag + ">"
			}
		}
	}

	return content
}

// Helper function to clean up HTML attributes
func cleanAttributes(content string) string {
	// Remove empty class attributes
	content = regexp.MustCompile(`class=["']\s*["']`).ReplaceAllString(content, "")

	// Remove WordPress-specific classes with more thorough pattern
	content = regexp.MustCompile(`class=["']([^"']*?)(wp-[a-zA-Z0-9_-]+|align[a-z]+|block-[a-z0-9-]+)([^"']*?)["']`).
		ReplaceAllStringFunc(content, func(match string) string {
			// Extract the class attribute
			re := regexp.MustCompile(`class=["']([^"']*)["']`)
			classes := re.FindStringSubmatch(match)
			if len(classes) > 1 {
				// Remove WordPress-specific classes
				classList := strings.Fields(classes[1])
				var cleanClasses []string
				for _, class := range classList {
					if !strings.HasPrefix(class, "wp-") &&
						!strings.HasPrefix(class, "block-") &&
						!regexp.MustCompile(`^align(left|right|center|none|wide|full)$`).MatchString(class) {
						cleanClasses = append(cleanClasses, class)
					}
				}

				// If no classes remain, return empty string to remove the attribute
				if len(cleanClasses) == 0 {
					return ""
				}

				// Otherwise, return the cleaned class attribute
				return fmt.Sprintf(`class="%s"`, strings.Join(cleanClasses, " "))
			}
			return match
		})

	// Clean up empty class attributes after removing WordPress classes
	content = regexp.MustCompile(`class=["']\s*["']`).ReplaceAllString(content, "")

	// Clean up style attributes with empty content
	content = regexp.MustCompile(`style=["']\s*["']`).ReplaceAllString(content, "")

	// Remove data attributes - more thorough pattern for WordPress-specific ones
	content = regexp.MustCompile(`\s+data-(wp-[a-zA-Z0-9_-]+=["'][^"']*["']|align=["'][^"']*["']|block=["'][^"']*["'])`).ReplaceAllString(content, "")

	// Remove other WordPress-specific data attributes
	content = regexp.MustCompile(`\s+data-[a-zA-Z0-9_-]+=["'][^"']*["']`).ReplaceAllString(content, "")

	// Remove empty attributes
	content = cleanEmptyAttributes(content)

	return content
}

// Helper function to clean empty attributes
func cleanEmptyAttributes(content string) string {
	// Clean up attributes that are empty
	attrNames := []string{"id", "class", "style", "title"}
	for _, attr := range attrNames {
		// Remove the attribute if it's empty
		pattern := fmt.Sprintf(`%s=["']\s*["']`, attr)
		re := regexp.MustCompile(pattern)
		content = re.ReplaceAllString(content, "")
	}

	// Clean up multiple spaces between attributes
	content = regexp.MustCompile(`\s{2,}`).ReplaceAllString(content, " ")

	// Clean up spaces around tag brackets
	content = regexp.MustCompile(`\s+>`).ReplaceAllString(content, ">")

	return content
}

// Helper function to update deprecated HTML tags
func updateDeprecatedTags(content string) string {
	// Convert <center> to styled div
	content = regexp.MustCompile(`<center([^>]*)>(.*?)</center>`).
		ReplaceAllString(content, `<div$1 style="text-align: center;">$2</div>`)

	// Convert <font> to span
	content = regexp.MustCompile(`<font([^>]*)>(.*?)</font>`).
		ReplaceAllString(content, `<span$1>$2</span>`)

	// Convert <b> to <strong>
	content = regexp.MustCompile(`<b([^>]*)>(.*?)</b>`).
		ReplaceAllString(content, `<strong$1>$2</strong>`)

	// Convert <i> to <em>
	content = regexp.MustCompile(`<i([^>]*)>(.*?)</i>`).
		ReplaceAllString(content, `<em$1>$2</em>`)

	return content
}

// Helper function to fix list structure
func fixListStructure(content string) string {
	// Fix common list issues:

	// 1. Ensure list items are inside lists
	content = regexp.MustCompile(`(?s)<li([^>]*)>(.*?)</li>`).
		ReplaceAllStringFunc(content, func(match string) string {
			// Check if this list item is not already inside a list
			before := regexp.MustCompile(`<[ou]l[^>]*>`).MatchString(match)
			after := strings.Contains(content, match+"</ul>") || strings.Contains(content, match+"</ol>")

			if !before && !after {
				re := regexp.MustCompile(`(?s)<li([^>]*)>(.*?)</li>`)
				matches := re.FindStringSubmatch(match)
				if len(matches) > 2 {
					return "<ul>" + match + "</ul>"
				}
			}
			return match
		})

	// 2. Fix nested lists
	content = regexp.MustCompile(`(?s)(<[ou]l[^>]*>)(.*?)(<[ou]l[^>]*>)(.*?)(</[ou]l>)(.*?)(</[ou]l>)`).
		ReplaceAllString(content, `$1$2<li>$3$4$5</li>$6$7`)

	return content
}

// Helper function to remove empty paragraphs
func removeEmptyParagraphs(content string) string {
	// Remove empty <p></p> tags
	re1 := regexp.MustCompile(`<p>\s*</p>`)
	content = re1.ReplaceAllString(content, "")

	// Remove paragraphs with only whitespace or <br> tags
	re2 := regexp.MustCompile(`<p>\s*(<br\s*/?>\s*)*</p>`)
	content = re2.ReplaceAllString(content, "")

	return content
}

// Helper function to remove Gutenberg block comments
func removeGutenbergComments(content string) string {
	// Remove opening block comments like <!-- wp:paragraph -->
	re1 := regexp.MustCompile(`<!-- wp:[^>]*-->\s*`)
	content = re1.ReplaceAllString(content, "")

	// Remove closing block comments like <!-- /wp:paragraph -->
	re2 := regexp.MustCompile(`\s*<!-- /wp:[^>]*-->`)
	content = re2.ReplaceAllString(content, "")

	return content
}

// Helper function to clean up whitespace
func cleanWhitespace(content string) string {
	// Replace multiple consecutive newlines with double newlines
	re1 := regexp.MustCompile(`\n{3,}`)
	content = re1.ReplaceAllString(content, "\n\n")

	// Remove trailing spaces at end of lines
	re2 := regexp.MustCompile(`[ \t]+\n`)
	content = re2.ReplaceAllString(content, "\n")

	// Clean up spaces around HTML tags
	re3 := regexp.MustCompile(`>\s+<`)
	content = re3.ReplaceAllString(content, "><")

	return strings.TrimSpace(content)
}

// ConvertToPlainText removes all HTML tags and returns plain text
func ConvertToPlainText(content string) string {
	// First check if we're dealing with nested lists
	if strings.Contains(content, "<ul") && strings.Contains(content, "<li") &&
		regexp.MustCompile(`<li[^>]*>.*?<ul[^>]*>`).MatchString(content) {
		return NestedListsToPlainText(content)
	}

	// If not nested lists, proceed with normal conversion
	// First sanitize WordPress content
	content = SanitizeWordPressContent(content)

	// Convert ordered lists to numbered lists
	content = convertOrderedListsToPlainText(content)

	// Convert unordered lists to bullet points
	content = convertUnorderedListsToPlainText(content)

	// Replace <br> tags with newlines
	re1 := regexp.MustCompile(`<br\s*/?>\s*`)
	content = re1.ReplaceAllString(content, "\n")

	// Replace </p> tags with double newlines
	re2 := regexp.MustCompile(`</p>\s*`)
	content = re2.ReplaceAllString(content, "\n\n")

	// Remove all remaining HTML tags
	re4 := regexp.MustCompile(`<[^>]*>`)
	content = re4.ReplaceAllString(content, "")

	// Decode HTML entities
	content = html.UnescapeString(content)

	// Clean up whitespace
	content = cleanWhitespace(content)

	return content
}

// ConvertToMarkdown converts WordPress content to Markdown format
func ConvertToMarkdown(content string) string {
	// First check if we're dealing with nested lists
	if strings.Contains(content, "<ul") && strings.Contains(content, "<li") &&
		regexp.MustCompile(`<li[^>]*>.*?<ul[^>]*>`).MatchString(content) {
		return NestedListsToMarkdown(content)
	}

	// If not nested lists, proceed with normal conversion
	// First sanitize WordPress content
	content = SanitizeWordPressContent(content)

	// Convert lists first (before other processing that might interfere)
	content = convertOrderedLists(content)
	content = convertUnorderedLists(content)

	// Convert images (handle various attribute orders and scenarios)
	content = regexp.MustCompile(`<img[^>]*src=["']([^"']*)["'][^>]*alt=["']([^"']*)["'][^>]*/?>`).ReplaceAllString(content, "![$2]($1)")
	content = regexp.MustCompile(`<img[^>]*alt=["']([^"']*)["'][^>]*src=["']([^"']*)["'][^>]*/?>`).ReplaceAllString(content, "![$1]($2)")
	content = regexp.MustCompile(`<img[^>]*src=["']([^"']*)["'][^>]*/?>`).ReplaceAllString(content, "![]($1)")
	// Convert headings (ensure proper spacing)
	content = regexp.MustCompile(`(?s)<h1[^>]*>(.*?)</h1>`).ReplaceAllString(content, "\n# $1\n\n")
	content = regexp.MustCompile(`(?s)<h2[^>]*>(.*?)</h2>`).ReplaceAllString(content, "\n## $1\n\n")
	content = regexp.MustCompile(`(?s)<h3[^>]*>(.*?)</h3>`).ReplaceAllString(content, "\n### $1\n\n")
	content = regexp.MustCompile(`(?s)<h4[^>]*>(.*?)</h4>`).ReplaceAllString(content, "\n#### $1\n\n")
	content = regexp.MustCompile(`(?s)<h5[^>]*>(.*?)</h5>`).ReplaceAllString(content, "\n##### $1\n\n")
	content = regexp.MustCompile(`(?s)<h6[^>]*>(.*?)</h6>`).ReplaceAllString(content, "\n###### $1\n\n")
	// Convert blockquotes (handle multi-line content better)
	content = regexp.MustCompile(`(?s)<blockquote[^>]*>(.*?)</blockquote>`).ReplaceAllStringFunc(content, func(match string) string {
		re := regexp.MustCompile(`(?s)<blockquote[^>]*>(.*?)</blockquote>`)
		matches := re.FindStringSubmatch(match)
		if len(matches) > 1 {
			// Clean inner content and add > to each line
			inner := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(matches[1], "")
			inner = html.UnescapeString(strings.TrimSpace(inner))
			lines := strings.Split(inner, "\n")
			result := "\n"
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					result += "> " + line + "\n"
				}
			}
			return result + "\n"
		}
		return match
	})
	// Convert code blocks (handle language attributes)
	content = regexp.MustCompile(`(?s)<pre[^>]*><code[^>]*class=["']language-([^"']*)["'][^>]*>(.*?)</code></pre>`).ReplaceAllString(content, "\n```$1\n$2\n```\n\n")
	content = regexp.MustCompile(`(?s)<pre[^>]*><code[^>]*>(.*?)</code></pre>`).ReplaceAllString(content, "\n```\n$1\n```\n\n")
	content = regexp.MustCompile(`(?s)<pre[^>]*>(.*?)</pre>`).ReplaceAllString(content, "\n```\n$1\n```\n\n")

	// Convert inline code
	content = regexp.MustCompile(`(?s)<code[^>]*>(.*?)</code>`).ReplaceAllString(content, "`$1`")

	// Convert text formatting (handle nested tags better)
	content = regexp.MustCompile(`(?s)<strong[^>]*>(.*?)</strong>`).ReplaceAllString(content, "**$1**")
	content = regexp.MustCompile(`(?s)<b[^>]*>(.*?)</b>`).ReplaceAllString(content, "**$1**")
	content = regexp.MustCompile(`(?s)<em[^>]*>(.*?)</em>`).ReplaceAllString(content, "*$1*")
	content = regexp.MustCompile(`(?s)<i[^>]*>(.*?)</i>`).ReplaceAllString(content, "*$1*")
	content = regexp.MustCompile(`(?s)<u[^>]*>(.*?)</u>`).ReplaceAllString(content, "<u>$1</u>") // Keep underline as HTML
	content = regexp.MustCompile(`(?s)<del[^>]*>(.*?)</del>`).ReplaceAllString(content, "~~$1~~")
	content = regexp.MustCompile(`(?s)<s[^>]*>(.*?)</s>`).ReplaceAllString(content, "~~$1~~")

	// Convert links (handle various scenarios)
	content = regexp.MustCompile(`(?s)<a[^>]*href=["']([^"']*)["'][^>]*>(.*?)</a>`).ReplaceAllString(content, "[$2]($1)")

	// Convert tables (basic table conversion)
	content = convertTables(content)

	// Convert horizontal rules
	content = regexp.MustCompile(`<hr\s*/?>`).ReplaceAllString(content, "\n---\n\n")
	// Convert line breaks (use two spaces + newline for Markdown line breaks)
	content = regexp.MustCompile(`<br\s*/?>`).ReplaceAllString(content, "  \n")

	// Convert paragraphs (use (?s) flag to match across newlines)
	content = regexp.MustCompile(`(?s)<p[^>]*>(.*?)</p>`).ReplaceAllString(content, "$1\n\n")

	// Convert divs to paragraphs (basic conversion)
	content = regexp.MustCompile(`(?s)<div[^>]*>(.*?)</div>`).ReplaceAllString(content, "$1\n\n")

	// Remove any remaining HTML tags
	content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")

	// Decode HTML entities
	content = html.UnescapeString(content)

	// Clean up excessive whitespace
	content = cleanWhitespace(content)

	// Final cleanup - ensure proper spacing around headings and blocks
	content = regexp.MustCompile(`\n{3,}`).ReplaceAllString(content, "\n\n")

	// Clean up any leading/trailing whitespace around block elements
	content = regexp.MustCompile(`\n\s+\n`).ReplaceAllString(content, "\n\n")

	return strings.TrimSpace(content)
}

// Helper function to convert ordered lists to markdown
func convertOrderedLists(content string) string {
	// This is a simplified conversion - for complex nested lists, you'd need more sophisticated parsing
	re := regexp.MustCompile(`<ol[^>]*>(.*?)</ol>`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		itemRe := regexp.MustCompile(`<li[^>]*>(.*?)</li>`)
		items := itemRe.FindAllStringSubmatch(match, -1)

		result := "\n" // Start with a newline for proper Markdown spacing
		for i, item := range items {
			if len(item) > 1 {
				// Clean the item content
				itemContent := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(item[1], "")
				itemContent = html.UnescapeString(itemContent)
				itemContent = strings.TrimSpace(itemContent)
				result += fmt.Sprintf("%d. %s\n", i+1, itemContent)
			}
		}
		result += "\n" // Add newline after the list
		return result
	})

	return content
}

// Helper function to convert unordered lists to markdown
func convertUnorderedLists(content string) string {
	// First, mark nested lists for special handling
	// This will enable us to preserve nested list structure
	content = regexp.MustCompile(`<li[^>]*>(.*?)<ul[^>]*>(.*?)</ul>(.*?)</li>`).
		ReplaceAllString(content, `<li data-has-child="true">$1<ul>$2</ul>$3</li>`)
	// Now handle outer lists
	re := regexp.MustCompile(`<ul[^>]*>(.*?)</ul>`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		itemRe := regexp.MustCompile(`<li[^>]*>(.*?)</li>`)
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

// Helper function to convert ordered lists to numbered plain text
func convertOrderedListsToPlainText(content string) string {
	// First, mark nested lists for special handling
	content = regexp.MustCompile(`<li[^>]*>(.*?)<ol[^>]*>(.*?)</ol>(.*?)</li>`).
		ReplaceAllString(content, `<li data-has-child="true">$1<ol>$2</ol>$3</li>`)

	re := regexp.MustCompile(`<ol[^>]*>(.*?)</ol>`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		itemRe := regexp.MustCompile(`<li[^>]*>(.*?)</li>`)
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

// Helper function to convert unordered lists to bullet point plain text
func convertUnorderedListsToPlainText(content string) string {
	// First, mark nested lists for special handling
	content = regexp.MustCompile(`<li[^>]*>(.*?)<ul[^>]*>(.*?)</ul>(.*?)</li>`).
		ReplaceAllString(content, `<li data-has-child="true">$1<ul>$2</ul>$3</li>`)

	re := regexp.MustCompile(`<ul[^>]*>(.*?)</ul>`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract list items
		itemRe := regexp.MustCompile(`<li[^>]*>(.*?)</li>`)
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

// Helper function to convert tables to markdown
func convertTables(content string) string {
	// Basic table conversion - this is simplified and may need enhancement for complex tables
	re := regexp.MustCompile(`<table[^>]*>(.*?)</table>`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract table content
		tableContent := regexp.MustCompile(`<table[^>]*>(.*?)</table>`).FindStringSubmatch(match)
		if len(tableContent) < 2 {
			return match
		}

		inner := tableContent[1]
		result := "\n"

		// Extract table rows
		rowRe := regexp.MustCompile(`<tr[^>]*>(.*?)</tr>`)
		rows := rowRe.FindAllStringSubmatch(inner, -1)

		isFirstRow := true
		for _, row := range rows {
			if len(row) < 2 {
				continue
			}

			// Extract cells (th or td)
			cellRe := regexp.MustCompile(`<t[hd][^>]*>(.*?)</t[hd]>`)
			cells := cellRe.FindAllStringSubmatch(row[1], -1)

			if len(cells) == 0 {
				continue
			}

			// Build row
			result += "|"
			for _, cell := range cells {
				if len(cell) > 1 {
					// Clean cell content
					cellContent := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(cell[1], "")
					cellContent = html.UnescapeString(strings.TrimSpace(cellContent))
					// Escape pipe characters in cell content
					cellContent = strings.ReplaceAll(cellContent, "|", "\\|")
					result += " " + cellContent + " |"
				}
			}
			result += "\n"

			// Add header separator after first row
			if isFirstRow {
				result += "|"
				for range cells {
					result += " --- |"
				}
				result += "\n"
				isFirstRow = false
			}
		}

		return result + "\n"
	})

	return content
}

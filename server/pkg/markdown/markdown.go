package markdown

import (
	"regexp"
	"strings"
)

// TOCItem 目录项
type TOCItem struct {
	Level    int
	Text     string
	Anchor   string
	Children []*TOCItem
}

// Result Markdown 解析结果
type Result struct {
	HTML      string
	TOC       []*TOCItem
	Excerpt   string
	WordCount int
}

// Parse 解析 Markdown
func Parse(content string) *Result {
	result := &Result{
		TOC: ExtractTOC(content),
	}

	// 提取纯文本摘要（前 200 字符）
	result.Excerpt = extractPlainText(content, 200)

	// 统计字数（中文字符 + 英文单词）
	result.WordCount = countWords(content)

	// HTML 渲染（简化实现，实际使用 markdown 库）
	result.HTML = renderToHTML(content)

	return result
}

// ExtractTOC 从 Markdown 提取目录
func ExtractTOC(content string) []*TOCItem {
	var items []*TOCItem
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line[0] != '#' {
			continue
		}

		// 计算标题级别
		level := 0
		for i := 0; i < len(line) && line[i] == '#'; i++ {
			level++
		}

		if level > 6 {
			continue
		}

		// 提取标题文本
		text := strings.TrimSpace(line[level:])
		anchor := generateAnchor(text)

		item := &TOCItem{
			Level:  level,
			Text:   text,
			Anchor: anchor,
		}
		items = append(items, item)
	}

	// 构建层级结构
	return buildTOCTree(items)
}

// buildTOCTree 构建目录树
func buildTOCTree(items []*TOCItem) []*TOCItem {
	if len(items) == 0 {
		return nil
	}

	var root []*TOCItem
	var stack []*TOCItem

	for _, item := range items {
		// 找到父节点
		for len(stack) > 0 && stack[len(stack)-1].Level >= item.Level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			// 顶层节点
			root = append(root, item)
		} else {
			// 添加到父节点的 Children
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, item)
		}

		stack = append(stack, item)
	}

	return root
}

// generateAnchor 生成锚点ID
func generateAnchor(text string) string {
	// 转换为小写
	anchor := strings.ToLower(text)

	// 替换空格为连字符
	anchor = strings.ReplaceAll(anchor, " ", "-")

	// 移除非字母数字连字符字符
	var result strings.Builder
	for _, r := range anchor {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// extractPlainText 提取纯文本
func extractPlainText(content string, maxLen int) string {
	// 移除代码块
	content = removeCodeBlocks(content)

	// 移除 Markdown 标记
	content = strings.ReplaceAll(content, "#", "")
	content = strings.ReplaceAll(content, "*", "")
	content = strings.ReplaceAll(content, "_", "")
	content = strings.ReplaceAll(content, "`", "")
	content = strings.ReplaceAll(content, "[", "")
	content = strings.ReplaceAll(content, "]", "")
	content = strings.ReplaceAll(content, "(", "")
	content = strings.ReplaceAll(content, ")", "")

	// 合并空白
	content = strings.Join(strings.Fields(content), " ")

	if len(content) > maxLen {
		return content[:maxLen] + "..."
	}
	return content
}

// removeCodeBlocks 移除代码块
func removeCodeBlocks(content string) string {
	var result strings.Builder
	inCodeBlock := false
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if !inCodeBlock {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.String()
}

// countWords 统计字数
func countWords(content string) int {
	count := 0
	
	// 统计英文单词
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
	count += len(wordRegex.FindAllString(content, -1))
	
	// 统计中文字符
	for _, r := range content {
		if r >= '\u4e00' && r <= '\u9fff' {
			count++
		}
	}
	
	return count
}

// renderToHTML 渲染为 HTML（简化实现）
func renderToHTML(content string) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")
	inCodeBlock := false
	inList := false
	listType := "" // "ul" or "ol"

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 代码块
		if strings.HasPrefix(trimmed, "```") {
			if inCodeBlock {
				result.WriteString("</code></pre>\n")
				inCodeBlock = false
			} else {
				result.WriteString("<pre><code>")
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			result.WriteString(escapeHTML(line))
			result.WriteString("\n")
			continue
		}

		// 空行
		if trimmed == "" {
			if inList {
				result.WriteString("</" + listType + ">\n")
				inList = false
				listType = ""
			}
			result.WriteString("<p></p>\n")
			continue
		}

		// 标题
		if strings.HasPrefix(trimmed, "#") {
			level := 0
			for level < len(trimmed) && trimmed[level] == '#' {
				level++
			}
			if level <= 6 {
				text := strings.TrimSpace(trimmed[level:])
				anchor := generateAnchor(text)
				result.WriteString("<h")
				result.WriteString(string(rune('0' + level)))
				result.WriteString(` id="`)
				result.WriteString(anchor)
				result.WriteString(`">`)
				result.WriteString(escapeHTML(text))
				result.WriteString("</h")
				result.WriteString(string(rune('0' + level)))
				result.WriteString(">\n")
				continue
			}
		}

		// 列表
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			if !inList || listType != "ul" {
				if inList {
					result.WriteString("</" + listType + ">\n")
				}
				result.WriteString("<ul>\n")
				inList = true
				listType = "ul"
			}
			text := strings.TrimPrefix(trimmed, "- ")
			text = strings.TrimPrefix(text, "* ")
			result.WriteString("<li>")
			result.WriteString(renderInline(escapeHTML(text)))
			result.WriteString("</li>\n")
			continue
		}

		// 有序列表
		orderedListRegex := regexp.MustCompile(`^\d+\.\s`)
		if orderedListRegex.MatchString(trimmed) {
			if !inList || listType != "ol" {
				if inList {
					result.WriteString("</" + listType + ">\n")
				}
				result.WriteString("<ol>\n")
				inList = true
				listType = "ol"
			}
			text := orderedListRegex.ReplaceAllString(trimmed, "")
			result.WriteString("<li>")
			result.WriteString(renderInline(escapeHTML(text)))
			result.WriteString("</li>\n")
			continue
		}

		// 普通段落
		if inList {
			result.WriteString("</" + listType + ">\n")
			inList = false
			listType = ""
		}
		result.WriteString("<p>")
		result.WriteString(renderInline(escapeHTML(trimmed)))
		result.WriteString("</p>\n")

		// 段落间添加空行（除了最后一行）
		if i < len(lines)-1 && strings.TrimSpace(lines[i+1]) == "" {
			result.WriteString("\n")
		}
	}

	// 关闭未关闭的列表
	if inList {
		result.WriteString("</" + listType + ">\n")
	}

	return result.String()
}

// escapeHTML HTML 转义
func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	return text
}

// renderInline 行内元素渲染（粗体、斜体、代码、链接）
func renderInline(text string) string {
	// 粗体 **text**
	boldRegex := regexp.MustCompile(`\*\*(.+?)\*\*`)
	text = boldRegex.ReplaceAllString(text, "<strong>$1</strong>")

	// 斜体 *text* 或 _text_
	italicRegex := regexp.MustCompile(`\*([^*]+)\*`)
	text = italicRegex.ReplaceAllString(text, "<em>$1</em>")

	// 行内代码 `code`
	codeRegex := regexp.MustCompile("`([^`]+)`")
	text = codeRegex.ReplaceAllString(text, "<code>$1</code>")

	// 链接 [text](url)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	text = linkRegex.ReplaceAllString(text, `<a href="$2">$1</a>`)

	return text
}

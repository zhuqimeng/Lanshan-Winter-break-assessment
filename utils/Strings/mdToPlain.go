package Strings

import (
	"regexp"
	"strings"
)

func MdToPlainText(md string) string {
	// 1. 移除代码块 (```code```)
	md = regexp.MustCompile("(?s)```.*?```").ReplaceAllString(md, " ")

	// 2. 移除行内代码 (`code`)
	md = regexp.MustCompile("`[^`]*`").ReplaceAllString(md, " ")

	// 3. 移除图片 ![alt](url)
	md = regexp.MustCompile(`!\[.*?]\(.*?\)`).ReplaceAllString(md, " ")

	// 4. 移除链接但保留文本 [text](url)
	md = regexp.MustCompile(`\[(.*?)]\(.*?\)`).ReplaceAllString(md, "$1")

	// 5. 移除标题标记 (# 标题)
	md = regexp.MustCompile(`(?m)^#{1,6}\s+`).ReplaceAllString(md, "")

	// 6. 移除粗体/斜体标记 (**text** 或 *text*)
	md = regexp.MustCompile(`\*\*([^*]+)\*\*`).ReplaceAllString(md, "$1")
	md = regexp.MustCompile(`\*([^*]+)\*`).ReplaceAllString(md, "$1")
	md = regexp.MustCompile(`__([^_]+)__`).ReplaceAllString(md, "$1")
	md = regexp.MustCompile(`_([^_]+)_`).ReplaceAllString(md, "$1")

	// 7. 移除列表标记 (- 或 * 或 1.)
	md = regexp.MustCompile(`(?m)^\s*[-*+]\s+`).ReplaceAllString(md, "")
	md = regexp.MustCompile(`(?m)^\s*\d+\.\s+`).ReplaceAllString(md, "")

	// 8. 移除引用标记 (>)
	md = regexp.MustCompile(`(?m)^>\s+`).ReplaceAllString(md, "")

	// 9. 移除分隔线 (--- 或 ***)
	md = regexp.MustCompile(`(?m)^\s*[-*_]{3,}\s*$`).ReplaceAllString(md, "")

	// 10. 移除HTML标签（如果有）
	md = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(md, " ")

	// 11. 将多个连续空格/换符合并为一个空格
	md = regexp.MustCompile(`\s+`).ReplaceAllString(md, " ")

	// 12. 移除首尾空格
	return strings.TrimSpace(md)
}

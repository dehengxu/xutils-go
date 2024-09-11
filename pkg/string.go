package pkg

import (
	"fmt"
	"regexp"
)

func XSprintf(format string, args ...interface{}) string {
	// 使用正则表达式匹配 %s、%v 等格式化占位符
	re := regexp.MustCompile(`%[vdsTtbcxXqgGeEfpPs]`)
	matches := re.FindAllString(format, -1)

	// 裁剪 args 的长度，使之与占位符数量匹配
	n1 := len(matches)
	// fmt.Printf("matches: %v\n", matches)
	if n1 < len(args) {
		args = args[:n1]
	}
	// fmt.Printf("new args: %v\n", args)

	// 使用 fmt.Sprintf 格式化字符串
	return fmt.Sprintf(format, args...)
}

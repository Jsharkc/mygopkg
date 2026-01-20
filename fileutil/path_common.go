package fileutil

import (
	"path/filepath"
	"strings"
)

// HasCommonPath 判断两个路径是否有公共部分（不包含根路径 /）
// 这是最准确和可靠的实现方法
func HasCommonPath(path1, path2 string) bool {
	// 标准化路径
	cleanPath1 := filepath.Clean(path1)
	cleanPath2 := filepath.Clean(path2)

	// 确保路径以 / 开头（处理相对路径）
	if !strings.HasPrefix(cleanPath1, "/") {
		cleanPath1 = "/" + cleanPath1
	}
	if !strings.HasPrefix(cleanPath2, "/") {
		cleanPath2 = "/" + cleanPath2
	}

	// 如果路径完全相同，认为有公共部分
	if cleanPath1 == cleanPath2 {
		return cleanPath1 != "/"
	}

	// 从 path1 开始，逐级向上查找
	current := cleanPath1
	for current != "/" && current != "." {
		if strings.HasPrefix(cleanPath2, current+"/") || cleanPath2 == current {
			return current != "/"
		}
		current = filepath.Dir(current)
	}

	// 从 path2 开始，逐级向上查找
	current = cleanPath2
	for current != "/" && current != "." {
		if strings.HasPrefix(cleanPath1, current+"/") || cleanPath1 == current {
			return current != "/"
		}
		current = filepath.Dir(current)
	}

	return false
}

// GetCommonPath 获取两个路径的公共部分
func GetCommonPath(path1, path2 string) string {
	// 标准化路径
	cleanPath1 := filepath.Clean(path1)
	cleanPath2 := filepath.Clean(path2)

	// 确保路径以 / 开头
	if !strings.HasPrefix(cleanPath1, "/") {
		cleanPath1 = "/" + cleanPath1
	}
	if !strings.HasPrefix(cleanPath2, "/") {
		cleanPath2 = "/" + cleanPath2
	}

	// 如果路径完全相同，返回路径本身
	if cleanPath1 == cleanPath2 {
		if cleanPath1 == "/" {
			return ""
		}
		return cleanPath1
	}

	// 从 path1 开始，逐级向上查找
	current := cleanPath1
	for current != "/" && current != "." {
		if strings.HasPrefix(cleanPath2, current+"/") || cleanPath2 == current {
			if current == "/" {
				return ""
			}
			return current
		}
		current = filepath.Dir(current)
	}

	// 从 path2 开始，逐级向上查找
	current = cleanPath2
	for current != "/" && current != "." {
		if strings.HasPrefix(cleanPath1, current+"/") || cleanPath1 == current {
			if current == "/" {
				return ""
			}
			return current
		}
		current = filepath.Dir(current)
	}

	return ""
}

// HasCommonPathSimple 简化版本：只判断是否有公共部分
func HasCommonPathSimple(path1, path2 string) bool {
	// 标准化路径
	cleanPath1 := filepath.Clean(path1)
	cleanPath2 := filepath.Clean(path2)

	// 确保路径以 / 开头
	if !strings.HasPrefix(cleanPath1, "/") {
		cleanPath1 = "/" + cleanPath1
	}
	if !strings.HasPrefix(cleanPath2, "/") {
		cleanPath2 = "/" + cleanPath2
	}

	// 如果路径完全相同，认为有公共部分
	if cleanPath1 == cleanPath2 {
		return cleanPath1 != "/"
	}

	// 检查一个是否是另一个的前缀
	if strings.HasPrefix(cleanPath1, cleanPath2+"/") || strings.HasPrefix(cleanPath2, cleanPath1+"/") {
		return true
	}

	// 找到公共前缀
	commonPrefix := findCommonPrefix(cleanPath1, cleanPath2)
	return commonPrefix != "" && commonPrefix != "/"
}

// findCommonPrefix 辅助函数：找到两个字符串的公共前缀
func findCommonPrefix(s1, s2 string) string {
	minLen := len(s1)
	if len(s2) < minLen {
		minLen = len(s2)
	}

	for i := 0; i < minLen; i++ {
		if s1[i] != s2[i] {
			// 找到最后一个 / 的位置
			lastSlash := strings.LastIndex(s1[:i], "/")
			if lastSlash == -1 {
				return ""
			}
			return s1[:lastSlash+1]
		}
	}

	// 如果一个是另一个的前缀
	if len(s1) == len(s2) {
		return s1
	}

	// 找到较短字符串的最后一个 / 的位置
	shorter := s1
	if len(s2) < len(s1) {
		shorter = s2
	}
	lastSlash := strings.LastIndex(shorter, "/")
	if lastSlash == -1 {
		return ""
	}
	return shorter[:lastSlash+1]
}

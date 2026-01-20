package fileutil_test

import (
	"fmt"
	"testing"

	"github.com/Jsharkc/mygopkg/fileutil"
)

func TestHasCommonPath(t *testing.T) {
	// 测试用例
	testCases := []struct {
		path1, path2 string
		expected     bool
	}{
		{"/a/b/c", "/a/b/d", true},     // 有公共部分 /a/b/
		{"/a/b/c", "/a/b/c/d", true},   // 有公共部分 /a/b/c/
		{"/a/b/c", "/x/y/z", false},    // 无公共部分
		{"/a/b", "/a/b", true},         // 完全相同
		{"/a", "/a/b", true},           // 一个是另一个的前缀
		{"/", "/a/b", false},           // 根路径不算
		{"/a/b/c", "/a/b/c", true},     // 完全相同
		{"/a/b/c", "/a/b/c/d/e", true}, // 有公共部分
		{"/a/b/c", "/a/b/x/y", true},   // 有公共部分 /a/b/
		{"/a/b/c", "/x/y/z", false},    // 无公共部分
		{"a/b/c", "a/b/d", true},       // 相对路径
		{"a/b/c", "x/y/z", false},      // 相对路径无公共部分
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_vs_%s", tc.path1, tc.path2), func(t *testing.T) {
			result := fileutil.HasCommonPath(tc.path1, tc.path2)
			if result != tc.expected {
				t.Errorf("HasCommonPath(%q, %q) = %v, want %v", tc.path1, tc.path2, result, tc.expected)
			}
		})
	}
}

func TestGetCommonPath(t *testing.T) {
	// 测试用例
	testCases := []struct {
		path1, path2 string
		expected     string
	}{
		{"/a/b/c", "/a/b/d", "/a/b"},       // 有公共部分 /a/b/
		{"/a/b/c", "/a/b/c/d", "/a/b/c"},   // 有公共部分 /a/b/c/
		{"/a/b/c", "/x/y/z", ""},           // 无公共部分
		{"/a/b", "/a/b", "/a/b"},           // 完全相同
		{"/a", "/a/b", "/a"},               // 一个是另一个的前缀
		{"/", "/a/b", ""},                  // 根路径不算
		{"/a/b/c", "/a/b/c", "/a/b/c"},     // 完全相同
		{"/a/b/c", "/a/b/c/d/e", "/a/b/c"}, // 有公共部分
		{"/a/b/c", "/a/b/x/y", "/a/b"},     // 有公共部分 /a/b/
		{"/a/b/c", "/x/y/z", ""},           // 无公共部分
		{"a/b/c", "a/b/d", "/a/b"},         // 相对路径
		{"a/b/c", "x/y/z", ""},             // 相对路径无公共部分
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_vs_%s", tc.path1, tc.path2), func(t *testing.T) {
			result := fileutil.GetCommonPath(tc.path1, tc.path2)
			if result != tc.expected {
				t.Errorf("GetCommonPath(%q, %q) = %q, want %q", tc.path1, tc.path2, result, tc.expected)
			}
		})
	}
}

func TestHasCommonPathSimple(t *testing.T) {
	// 测试用例
	testCases := []struct {
		path1, path2 string
		expected     bool
	}{
		{"/a/b/c", "/a/b/d", true},     // 有公共部分 /a/b/
		{"/a/b/c", "/a/b/c/d", true},   // 有公共部分 /a/b/c/
		{"/a/b/c", "/x/y/z", false},    // 无公共部分
		{"/a/b", "/a/b", true},         // 完全相同
		{"/a", "/a/b", true},           // 一个是另一个的前缀
		{"/", "/a/b", false},           // 根路径不算
		{"/a/b/c", "/a/b/c", true},     // 完全相同
		{"/a/b/c", "/a/b/c/d/e", true}, // 有公共部分
		{"/a/b/c", "/a/b/x/y", true},   // 有公共部分 /a/b/
		{"/a/b/c", "/x/y/z", false},    // 无公共部分
		{"a/b/c", "a/b/d", true},       // 相对路径
		{"a/b/c", "x/y/z", false},      // 相对路径无公共部分
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_vs_%s", tc.path1, tc.path2), func(t *testing.T) {
			result := fileutil.HasCommonPathSimple(tc.path1, tc.path2)
			if result != tc.expected {
				t.Errorf("HasCommonPathSimple(%q, %q) = %v, want %v", tc.path1, tc.path2, result, tc.expected)
			}
		})
	}
}

func TestExtraCases(t *testing.T) {
	// 额外测试用例
	extraTests := []struct {
		path1, path2   string
		description    string
		expectedHas    bool
		expectedCommon string
	}{
		{"/usr/local/bin", "/usr/local/lib", "同级目录", true, "/usr/local"},
		{"/home/user", "/home/user/documents", "父子目录", true, "/home/user"},
		{"/var/log", "/var/cache", "同级目录", true, "/var"},
		{"/etc", "/etc/passwd", "父子目录", true, "/etc"},
		{"/tmp", "/tmp/file.txt", "父子目录", true, "/tmp"},
		{"/opt", "/usr", "不同根目录", false, ""},
	}

	for _, tc := range extraTests {
		t.Run(tc.description, func(t *testing.T) {
			hasCommon := fileutil.HasCommonPath(tc.path1, tc.path2)
			common := fileutil.GetCommonPath(tc.path1, tc.path2)

			if hasCommon != tc.expectedHas {
				t.Errorf("HasCommonPath(%q, %q) = %v, want %v", tc.path1, tc.path2, hasCommon, tc.expectedHas)
			}

			if common != tc.expectedCommon {
				t.Errorf("GetCommonPath(%q, %q) = %q, want %q", tc.path1, tc.path2, common, tc.expectedCommon)
			}
		})
	}
}

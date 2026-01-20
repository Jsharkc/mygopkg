package fileutil

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// Exist 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func GetFilenameWithoutEx(filePath string) string {
	fullName := path.Base(filePath)
    dotIndex := strings.LastIndex(fullName, ".")
    if dotIndex > 0 {
        return fullName[:dotIndex]
    }
    return fullName
}

// RemoveLocalFiles 删除本地文件
func RemoveLocalFiles(fileList ...string) error {
	if len(fileList) == 0 {
		return nil
	}

	for _, filePath := range fileList {
		if !Exist(filePath) {
			continue
		}

		err := os.Remove(filePath)
		if err != nil {
			return fmt.Errorf("filePath: %s, err = %v", filePath, err)
		}
	}

	return nil
}

// DetectContentTypeByReader 从 r 中读取数据并检测其中的内容类型。
// 注意：r 中的数据被读取后，外部再次读取，可能会返回 io.EOF
// 可能需要支持 io.Seek 操作
func DetectContentTypeByReader(r io.Reader) ([]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return DetectContentType(data)
}

// DetectContentType 检测字节流 data 达标的类型
// 返回的类型后缀包含[.]。如 .png
func DetectContentType(data []byte) ([]string, error) {
	ext := mimetype.Detect(data).Extension()
	if ext != "" {
		return []string{ext}, nil
	}

	typ := http.DetectContentType(data)
	return mime.ExtensionsByType(typ)
}

// IsAllowType 判断 data 代表的内容（来自文件或数据流）类型是否在 acceptTypes 内
// acceptTypes 的格式类似：
//
//	map[string]strcut{}{
//			"png": {},
//			...
//	}
func IsAllowType(data []byte, acceptTypes map[string]struct{}) (bool, error) {
	exts, err := DetectContentType(data)
	if err != nil {
		return false, err
	}

	for _, ext := range exts {
		ext = strings.TrimPrefix(ext, ".")
		if _, ok := acceptTypes[ext]; ok {
			return true, nil
		}
	}

	return false, nil
}

// CheckAndGetType 检查 data 代表的内容（来自文件或数据流）类型是否在 acceptTypes 内
// 并返回类型扩展名
func CheckAndGetType(data []byte, acceptTypes map[string]struct{}) (string, error) {
	exts, err := DetectContentType(data)
	if err != nil {
		// 读取失败时，类型返回空
		return "", err
	}

	if len(acceptTypes) == 0 && len(exts) > 0 {
		// 如果没有传可用类型，默认返回解析得到的第一个文件类型
		return exts[0], nil
	}

	for _, ext := range exts {
		ext = strings.TrimPrefix(ext, ".")
		if _, ok := acceptTypes[ext]; ok {
			// 返回可用类型中命中的类型
			return ext, nil
		}
	}
	// 如果不在可用类型里面，返回文件或数据流的类型，并返回报错
	return fmt.Sprintf("%v", exts), errors.New("file type not allow")
}

// GetExt 获取 uri（可以是文件路径或 URL） 的后缀，不包含 "."
// 如果没有，返回 defExts[0] 的值
func GetExt(uri string, defExts ...string) (ext string) {
	defer func() {
		if ext == "" && len(defExts) > 0 {
			ext = defExts[0]
		}
	}()

	fpath := uri
	if strings.HasPrefix(uri, "http") {
		u, err := url.Parse(uri)
		if err != nil {
			log.Printf("Parse url:%s error:%v\n", uri, err)
			return
		}
		fpath = u.Path
	} else {
		data, err := os.ReadFile(fpath)
		if err != nil {
			log.Printf("ReadFile %s error:%v\n", fpath, err)
			return
		}
		exts, _ := DetectContentType(data)
		if len(exts) > 0 {
			return strings.TrimPrefix(exts[0], ".")
		}
	}

	// url 或上面没有检测出具体类型，则用后缀
	ext = strings.TrimPrefix(filepath.Ext(fpath), ".")

	return
}

func InferRootDir(exePath string) string {
	// 这里的技巧可以学习下：先定义一个函数变量，然后在赋值函数时，里面使用了这个函数变量
	var infer func(d string) string
	infer = func(d string) string {
		if Exist(d + "/go.mod") {
			return d
		}

		if d == "/" {
			panic("请确保在项目根目录或子目录下运行程序，当前在：" + d)
		}

		return infer(filepath.Dir(d))
	}

	return infer(exePath)
}

func CheckDirBeforeCreateFile(name string) error {
	var (
		err      error
		filePath = name
	)
	if !filepath.IsAbs(filePath) {
		filePath, err = filepath.Abs(filePath)
		if err != nil {
			return fmt.Errorf("get abs file path failed: %w", err)
		}
	}
	fileDir := filepath.Dir(filePath)
	if !Exist(fileDir) {
		err := os.MkdirAll(fileDir, 0755)
		if err != nil {
			return fmt.Errorf("create dir faield: %w", err)
		}
	}
	return nil
}

func Create(name string) (*os.File, error) {
	err := CheckDirBeforeCreateFile(name)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	if (flag & os.O_CREATE) == os.O_CREATE {
		err := CheckDirBeforeCreateFile(name)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(name, flag, perm)
}

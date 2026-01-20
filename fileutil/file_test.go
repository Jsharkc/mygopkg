package fileutil_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Jsharkc/mygopkg/fileutil"
)

func TestExist(t *testing.T) {
	filePath := `pkg/fileutil/file_test.go`
	exist := fileutil.Exist(filePath)
	if !exist {
		t.Errorf("filePath (%s) not exist", filePath)
		return
	}

	fmt.Printf("exist=%#v/n", exist)
}

func TestGetFilenameWithoutEx(t *testing.T) {
	filePath := `pkg/fileutil/file_test.go`
	filename := fileutil.GetFilenameWithoutEx(filePath)
	if filename != "file_test" {
		t.Errorf("expect: file_test, got:%s", filename)
	}
}

func TestRemoveLocalFiles(t *testing.T) {
	fileList := []string{
		`pkg/fileutil/tmp/temp1`,
	}

	err := fileutil.RemoveLocalFiles(fileList...)
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range fileList {
		if exist := fileutil.Exist(item); exist {
			t.Errorf("file_path=%s not be deleted", item)
			return
		}
	}
}

func TestGetExt(t *testing.T) {
	type args struct {
		uri    string
		defExt string
	}
	tests := []struct {
		name    string
		args    args
		wantExt string
	}{
		{"a", args{"/Users/polarisxu/integration.zsh", "zsh"}, "txt"},
		{"b", args{"abc/name", "mp3"}, "mp3"},
		{"c", args{"https://megaview.com/abc/name.txt", "mp3"}, "txt"},
		{"d", args{"https://megaview.com/abc/name", "mp3"}, "mp3"},
		{"e", args{"https://megaview.com/abc/name.txt?duration=3", "mp3"}, "txt"},
		{"f", args{"https://megaview.com/abc/name?duration=3", "mp3"}, "mp3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotExt := fileutil.GetExt(tt.args.uri, tt.args.defExt); gotExt != tt.wantExt {
				t.Errorf("GetExt() = %v, want %v", gotExt, tt.wantExt)
			}
		})
	}
}

func TestCheckAndGetType(t *testing.T) {
	type args struct {
		filename    string
		acceptTypes map[string]struct{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"a", args{"/Users/polarisxu/2.mp3", map[string]struct{}{"mp3": struct{}{}, "wav": struct{}{}}}, "mp3", false},
		{"b", args{"/Users/polarisxu/Downloads/05065c65-e42e-4f0c-8067-0f1271c86718.mp3", map[string]struct{}{"mp3": struct{}{}, "wav": struct{}{}}}, "mp3", false},
		{"c", args{"/Users/polarisxu/Downloads/8009244-20230424102053-17665077209-6126--record-medias_11-1682302853.21739.wav", map[string]struct{}{"mp3": struct{}{}, "wav": struct{}{}}}, "wav", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.args.filename)
			if err != nil {
				t.Errorf("CheckAndGetType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := fileutil.CheckAndGetType(data, tt.args.acceptTypes)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAndGetType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			if got != tt.want {
				t.Errorf("CheckAndGetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	filePath := "test_dir1/test_dir2/test_dir3/test_file"
	file, err := fileutil.Create(filePath)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	defer file.Close()
	fmt.Println(file)
}

func TestOpenFile(t *testing.T) {
	filePath := "test_dir1/test_dir2/test_dir3/test_file"
	file, err := fileutil.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	defer file.Close()
	fmt.Println(file)
}

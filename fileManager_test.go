package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

// Test: CreateDirectory, if the directory is created
func TestCreateDirectory(t *testing.T) {
	t.Run("check_if_path_empty", func(t *testing.T) {
		if got := IsStringEmpty(Path); got {
			t.Errorf("TestCreateDirectory = %v, should have a path set", Path)
		}
	})
	t.Run("should_create_a_directory", func(t *testing.T) {
		CreateDirectory()
		p, _ := os.Stat(Path)
		fmt.Println(p.Name(), Path)
		if p, err := os.Stat(Path); !p.IsDir() || err != nil {
			t.Errorf("TestCreateDirectory didn't create the directory, failed")
		}
	})
}

// Test: WriteToFile, should write something to file
func TestWriteToFile(t *testing.T) {
	f := fmt.Sprintf("%s/%s_file_manager_write_test_case.out", Path, programName)
	m := "Hello! World"
	t.Run("should_create_and_write_to_file", func(t *testing.T) {
		if err := WriteToFile(f, m); err != nil {
			t.Errorf("TestWriteToFile should create and write but got err: %v", err)
		}
	})
	t.Run("check_if_the_content_exists", func(t *testing.T) {
		c, err := ReadFile(f);
		if len(c) <= 0 {
			t.Errorf("TestWriteToFile %d content found in the file: %v", len(c), f)
		} else {
			if err != nil || c[0] != m {
				t.Errorf("TestWriteToFile either got err: %v, or has invalid content: %v", err, c[0])
			}
		}
	})
}

// Test: ListFile, should produce the list of file from the directory
func TestListFile(t *testing.T) {
	f := fmt.Sprintf("%s/%s_file_manager_list_test_case.out", Path, programName)
	err := WriteToFile(f, "")
	l, err := ListFile(Path, "*list_test_case.out")
	t.Run("should_give_us_valid_list_of_file_from_directory", func(t *testing.T) {
		if err != nil || len(l) <= 0 {
			t.Errorf("TestListFile either got err: %v, or return invalid count %v want > 0", err, len(l))
		}
	})
}

// Test: ReadFile, should provide the content from the given file
func TestReadFile(t *testing.T) {
	f := fmt.Sprintf("%s/%s_file_manager_read_test_case.out", Path, programName)
	m := `
	  First Line Hello! World
      Second Line Hello !World
	`
	_ = WriteToFile(f, m)
	c, err := ReadFile(f)
	fmt.Println()
	t.Run("should_give_us_valid_content_from_file", func(t *testing.T) {
		if len(c) <= 0 {
			t.Errorf("TestReadFile %d content found in the file: %v", len(c), f)
		} else {
			if err != nil || RemoveSpecialCharacters(strings.Join(c, ",")) != RemoveSpecialCharacters(m) {
				t.Errorf("TestListFile either got err: %v, or return invalid count: %v", err, c)
			}
		}
	})
}

// Test: CurrentDir(), check if the current directory return is valid
func TestCurrentDir(t *testing.T) {
	t.Run("should_return_a_valid_current_directory", func(t *testing.T) {
		_, filename, _, _ := runtime.Caller(0)
		cd := CurrentDir()
		if IsStringEmpty(cd) {
			t.Errorf("TestCurrentDir = %v, is empty, should return some value", cd)
		}
		if !strings.HasPrefix(filename, cd) {
			t.Errorf("TestCurrentDir = %v, should start with %v", cd, filename)
		}
	})
}
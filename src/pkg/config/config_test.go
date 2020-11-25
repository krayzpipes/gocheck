package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestGetConfigRootDir_NoEnvVar(t *testing.T) {
	path := "/fake/path"
	err := os.Setenv("GOCHECK_CONFIG_DIR", path)
	defer os.Unsetenv("GOCHECK_CONFIG_DIR")

	rootDir, err := GetConfigRootDir()

	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if rootDir != path {
		t.Fatalf("expected path %s, got %s", path, rootDir)
	}
}

func TestGetConfigRootDir_EmptyEnvVar(t *testing.T) {
	err := os.Setenv("GOCHECK_CONFIG_DIR", "")
	defer os.Unsetenv("GOCHECK_CONFIG_DIR")
	cwd, _ := os.Getwd()

	rootDir, err := GetConfigRootDir()

	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if rootDir != cwd {
		t.Fatalf("expected path %s, got %s", cwd, rootDir)
	}
}

func TestGetConfigRootDir_UnsetEnvVar(t *testing.T) {
	cwd, _ := os.Getwd()

	rootDir, err := GetConfigRootDir()

	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if rootDir != cwd {
		t.Fatalf("expected path %s, got %s", cwd, rootDir)
	}
}

func TestWalkConfigDirs(t *testing.T) {
	root, hclFiles, cleanUp := testGenerateConfigDirectory()
	defer cleanUp()

	matches, err := WalkConfigDirs(root, ".hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	sortSlice(hclFiles)
	sortSlice(matches)

	if !compareStringSlices(hclFiles, matches) {
		hclString := strings.Join(hclFiles, ",")
		matchesString := strings.Join(matches, ",")
		t.Fatalf("string slices do not match: %s != %s", hclString, matchesString)
	}
}

func testGenerateConfigDirectory() (string, []string, func()) {
	tempDir := os.TempDir()
	rootDir := filepath.Join(tempDir, "gocheck-test-root")
	_ = os.Mkdir(rootDir, 0777)
	childDir := filepath.Join(rootDir, "gocheck-test-child")
	_ = os.Mkdir(childDir, 0777)

	tempFileRoot := filepath.Join(rootDir, "test1.hcl")
	_, _ = os.Create(tempFileRoot)

	tempFileChild := filepath.Join(childDir, "testChild.hcl")
	_, _ = os.Create(tempFileChild)

	tempFileIgnored := filepath.Join(childDir, "ignore.txt")
	_, _ = os.Create(tempFileIgnored)

	hclFiles := []string{tempFileRoot, tempFileChild}

	return rootDir, hclFiles, func() {
		os.RemoveAll(rootDir)
	}

}

func compareStringSlices(sliceOne []string, sliceTwo []string) bool {
	if len(sliceOne) != len(sliceTwo) {
		return false
	}
	for i, item := range sliceOne {
		if item != sliceTwo[i] {
			return false
		}
	}
	return true
}

func sortSlice(slice []string) {
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
}

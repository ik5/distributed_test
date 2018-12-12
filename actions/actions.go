package actions

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ik5/gostrutils"
)

func isNumeric(name string) bool {
	re := regexp.MustCompile("^[0-9]+$")
	return re.MatchString(name)
}

func extractPid(name string) int {
	re := regexp.MustCompile("^/proc/([0-9]+)/")
	pid := re.FindStringSubmatch(name)
	if len(pid) != 2 || pid[1] == "" {
		return 0
	}

	return int(gostrutils.StrToInt64(pid[1], 0))
}

func commContent(path string) string {
	buff, _ := ioutil.ReadFile(path)
	result := strings.TrimSuffix(string(buff), "\n")
	return result
}

// ProcExists check if a given procName is currently running
func ProcExists(procName string) (bool, error) {
	pid := 0
	err := filepath.Walk("/proc/", func(path string, info os.FileInfo, err error) error {
		fname := filepath.Base(path)
		if fname == "comm" {
			content := commContent(path)
			if strings.Compare(content, procName) == 0 {
				pid = extractPid(path)
				return nil
			}
		}
		return nil
	})

	return pid > 0, err
}

// IsFileContains go over a given path and validate if substr exists
func IsFileContains(path, substr string) (bool, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	content := string(buff)
	return strings.Contains(content, substr), nil
}

// IsFileExists validate if a given file exists on the system
func IsFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

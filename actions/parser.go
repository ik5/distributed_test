package actions

import (
	"bufio"
	"bytes"
	"regexp"

	"github.com/ik5/distributed_test/types"
	"github.com/ik5/gostrutils"
)

func extractNumeric(line string) int {
	re := regexp.MustCompile("([0-9]+)$")
	found := re.FindStringSubmatch(line)
	if len(found) != 2 {
		return 0
	}

	return int(gostrutils.StrToInt64(found[1], 0))
}

// Parse translate a buffer to Structed
func Parse(buff []byte) types.RawAction {
	var result types.RawAction
	breader := bytes.NewReader(buff)
	r := bufio.NewReader(breader)
	readBody := false
	body := []byte{}
	for {
		var line []byte
		if !readBody {
			line, _, _ = r.ReadLine()
		}

		if string(line) == "" && !readBody {
			readBody = true
			continue
		}

		if readBody {
			r.Read(body)
			break
		}

		if bytes.HasPrefix(line, []byte("t: ")) {
			t := extractNumeric(string(line))
			result.Type = types.ActionType(t)
			continue
		}

		if bytes.HasPrefix(line, []byte("l: ")) {
			length := extractNumeric(string(line))
			result.Length = length
			body = make([]byte, length)
			continue
		}

	}

	result.Body = body
	return result
}

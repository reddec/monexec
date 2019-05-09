package pool

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Environment variables file format:
// pair: KEY=VALUE
// comment: line started with #
// empty lines or invalid (without = symbol) ignored
// there is no way to escape new line symbol
func ParseEnvironmentStream(stream io.Reader) map[string]string {
	ans := map[string]string{}
	reader := bufio.NewScanner(stream)
	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			// broken line
			continue
		}
		ans[kv[0]] = kv[1]
	}
	return ans
}

func ParseEnvironmentFile(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ParseEnvironmentStream(file), nil
}

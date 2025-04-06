package parsers

import (
	"bufio"
	"fmt"
	"strings"
)

func ParseIni(input string) (map[string]string, error) {
	resultMap := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		resultMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading string: %w", err)
	}

	return resultMap, nil
}

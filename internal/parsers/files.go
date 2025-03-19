package parsers

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func IsContainsInFile(path string, substring string) (bool, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open xl2tp.config file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, substring) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func IsFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func AppendToFile(file *os.File, data string) error {
	_, err := file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func AppendToFileByPath(path string, data string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := AppendToFile(file, data); err != nil {
		return err
	}

	return nil
}

func CopyFile(src, destination string) (err error) {
	data, err := os.ReadFile(src)
	if err != nil {
		log.Fatalf("failed to read src: %v", err)
	}

	if err := os.WriteFile(destination, data, 0644); err != nil {
		log.Fatalf("failed to write to target: %v", err)
	}

	return nil
}

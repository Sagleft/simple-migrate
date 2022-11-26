package simplemigrate

import (
	"fmt"
	"io/ioutil"
)

func readFile(filepath string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return fileBytes, err
}

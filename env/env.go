package env

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getKeyValue(property string) (string, string) {
	prop := strings.TrimSpace(property)
	keyValue := strings.SplitN(prop, "=", 2)
	return keyValue[0], keyValue[1]
}

func loadFile(rootPath, fileName string) (map[string]string, error) {
	filePath := fmt.Sprintf("%senv/%s.env", rootPath, fileName)
	file, err := os.Open(filePath)

	if err != nil {
		return nil, fmt.Errorf("error opening file %s.env: %v", fileName, err)
	}

	defer file.Close()

	properties := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		key, value := getKeyValue(scanner.Text())
		properties[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scaning file %s.env: %v", fileName, err)
	}

	return properties, nil
}

type Env struct {
	local map[string]string
}

func NewEnv(rootPath string, fileName string) *Env {
	props, err := loadFile(rootPath, fileName)

	if err != nil {
		log.Println(err)
	}

	return &Env{local: props}
}

func (env *Env) GetProp(key string) string {

	value, exists := os.LookupEnv(key)

	if exists {
		return value
	}

	return env.local[key]
}

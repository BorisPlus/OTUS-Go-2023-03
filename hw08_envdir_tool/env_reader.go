package main

import (
	"bufio"
	"os"
	"path"
	"strings"
)

// EnvValue - переменная окружения:
//   - Value - значение переменной окружения.
//   - NeedRemove - требование удаления переменной окружения.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// CreateEnvValue - конструктор переменной окружения, при этом:
//   - Правая табуляция удаляется.
//   - Символы нуль-байт заменяются на символ переноса.
func CreateEnvValue(value string, needRemove bool) EnvValue {
	envValue := EnvValue{}
	envValue.NeedRemove = needRemove
	if envValue.NeedRemove {
		return envValue
	}
	envValue.Value = strings.ReplaceAll(strings.TrimRight(value, " \t"), "\x00", "\n")
	return envValue
}

// Environment - окружение переменных.
type Environment map[string]EnvValue

// Apply - функция применения окружения переменных.
// Положительное значени индикатора NeedRemove указывает на необходимость удаления конкретной переменной из окружения.
func Apply(environment Environment) {
	for envName, envValue := range environment {
		_, exists := os.LookupEnv(envName)
		if exists {
			os.Unsetenv(envName)
		}
		if envValue.NeedRemove {
			continue
		}
		os.Setenv(envName, envValue.Value)
	}
}

// Clear - функция принудительной очистки окружения переменных.
func Clear(environment Environment) {
	for envName := range environment {
		_, exists := os.LookupEnv(envName)
		if exists {
			os.Unsetenv(envName)
		}
	}
}

// TODO: Test cover need

// UnsetOnlyNeeded - функция принудительной очистки окружения от требующих удаления переменных.
func UnsetOnlyNeeded(environment Environment) {
	for envName, envValue := range environment {
		_, exists := os.LookupEnv(envName)
		if exists && envValue.NeedRemove {
			os.Unsetenv(envName)
		}
	}
}

// TODO: Test cover need

// UpdateOrInsert - функция обновления или добавления переменных окружения, нетребующих удаления.
func UpdateOrInsert(environment Environment) {
	for envName, envValue := range environment {
		_, exists := os.LookupEnv(envName)
		if exists {
			os.Unsetenv(envName)
		}
		os.Setenv(envName, envValue.Value)
	}
}

// ReadEnvDir - функция обхода директории для заполнения окружения переменных по содержанию файлов.
// Файлы, содержащие в именовании знак равенства, игнорируются.
func ReadEnvDir(envDir string) (Environment, error) {
	environment := make(map[string]EnvValue)
	entries, err := os.ReadDir(envDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.Contains(e.Name(), "=") {
			continue
		}

		envValue, err := ReadEnvFile(path.Join(envDir, e.Name()))
		if err != nil {
			return nil, err
		}
		environment[e.Name()] = envValue
	}
	return environment, nil
}

// ReadEnvFile - функция парсинга первой строки файла для установки значения переменной окружения.
// При нулевом объеме файла переменная помечается как требуемая к удалению.
func ReadEnvFile(file string) (EnvValue, error) {
	envValue := EnvValue{}
	envFile, err := os.Open(file)
	if err != nil {
		return envValue, err
	}
	defer envFile.Close()
	fileInfo, err := envFile.Stat()
	if err != nil {
		return envValue, err
	}
	if fileInfo.Size() == 0 {
		return envValue, nil
	}
	scanner := bufio.NewScanner(envFile)
	value := ""
	for scanner.Scan() {
		value = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return envValue, err
	}
	envValue.Value = value
	envValue.NeedRemove = false
	return envValue, nil
}

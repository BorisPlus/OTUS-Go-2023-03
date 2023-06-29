package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/icrowley/fake"
)

// generateTestEnvironment - функция генерации уникальных тестовых наборов данных для каждого теста.
func generateTestEnvironment() Environment {
	return Environment{
		"GOLANG_HW08_TEST_EMAIL":     EnvValue{fake.EmailAddress(), false},
		"GOLANG_HW08_TEST_USERAGENT": EnvValue{fake.UserAgent(), false},
		"GOLANG_HW08_TEST_USERNAME":  EnvValue{fake.UserName(), true},
	}
}

func TestSet(t *testing.T) {
	name := "SOMEVAR"
	environment := make(map[string]EnvValue)
	environment[name] = CreateEnvValue("123", false)
	for envName, envValue := range environment {
		_, exists := os.LookupEnv(envName)
		if exists {
			os.Unsetenv(envName)
		}
		if envValue.NeedRemove {
			continue
		}
		os.Setenv(envName, envValue.Value)

		value, exists := os.LookupEnv(name)
		if !exists {
			t.Errorf("EnvVar %q was not set for %q.\n", name, value)
		}
	}
}

func TestGetOfNotExists(t *testing.T) {
	name := "HW08_SOMEVAR"
	value, exists := os.LookupEnv("HW08_SOMEVAR")
	if exists {
		t.Errorf("EnvVar %q must be do not set. but return %s.\n", name, value)
	}
}

// TestReadDir - тестирование ReadEnvDir.
func TestReadDir(t *testing.T) {
	ethalons := make(map[string]EnvValue)
	ethalons["BAR"] = CreateEnvValue("bar", false)
	ethalons["EMPTY"] = CreateEnvValue("", false)
	ethalons["FOO"] = CreateEnvValue("   foo\nwith new line", false)
	ethalons["HELLO"] = CreateEnvValue(`"hello"`, false)
	ethalons["UNSET"] = CreateEnvValue("", true)
	//
	environment, err := ReadEnvDir("./testdata/env")
	if err != nil {
		t.Errorf(err.Error())
	}
	for envName, envValue := range ethalons {
		if ethalons[envName] != environment[envName] {
			t.Errorf(
				"EnvVar %q with parsed-data %+v not equal ethalon-data %+v\n",
				envName, envValue, ethalons[envName])
		} else {
			fmt.Printf("OK. EnvVar %q: %+v \n", envName, envValue)
		}
	}
}

// TestReadDir - тестирование Clear (функции принудительной очистки окружения переменных).
func TestClearEnvironment(t *testing.T) {
	environment := generateTestEnvironment()

	Apply(environment)

	Clear(environment)

	for envName := range environment {
		_, exists := os.LookupEnv(envName)
		fmt.Printf("Check EnvVar %q\n", envName)
		if exists {
			t.Errorf("EnvVar %q was not unset", envName)
		}
	}
}

// TestApplyEnvironment - тестирование Apply (функции активации очистки окружения переменных).
func TestApplyEnvironment(t *testing.T) {
	environment := generateTestEnvironment()

	defer Clear(environment)

	Apply(environment)

	for envName, envValue := range environment {
		value, exists := os.LookupEnv(envName)
		fmt.Printf("Check EnvVar %q\n", envName)
		fmt.Printf("\tValue %v\n", envValue.Value)
		fmt.Printf("\tNeedRemove %v\n", envValue.NeedRemove)
		if !exists && !envValue.NeedRemove {
			t.Errorf("EnvVar %q was not setup", envName)
		}
		if exists && envValue.NeedRemove {
			t.Errorf("EnvVar %q must be removed", envName)
		}
		if !envValue.NeedRemove && value != envValue.Value {
			t.Errorf("EnvVar %q get value %q, but expected %q", envName, value, envValue.Value)
		}
		if exists {
			fmt.Printf("%s=%+v\n", envName, value)
		}
	}
}

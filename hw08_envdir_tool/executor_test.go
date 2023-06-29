package main

import (
	"fmt"
	"testing"
)

var cleanEnvironment = Environment{}

const (
	bashEnvVarCheckTemplate string = `
if [ ${GOLANG_HW08_TEST_EMAIL} == %s ]; then 
	exit 5; 
else 
	exit 17; 
fi;
`
	okExitCode   int = 5
	failExitCode int = 17
)

// TestRunCmd_ExitCode - тестирование факта перехвата функцией RunCmd кода завершения целевого процесса ОС.
func TestRunCmd_001_ExitCode(t *testing.T) {
	fmt.Println("Test RunCmd to catch exit code of OS process.")
	testCases := []struct {
		bashCommand      string
		expectedExitCode int
	}{
		{expectedExitCode: 0, bashCommand: "echo $GOLANG_HW08_TEST_EMAIL"},
		{expectedExitCode: 0, bashCommand: "echo echo $GOLANG_HW08_TEST_EMAIL"},
		{expectedExitCode: 0, bashCommand: "echo ${GOLANG_HW08_TEST_EMAIL}"},
		{expectedExitCode: 1, bashCommand: "exit 1"},
		{expectedExitCode: 0, bashCommand: "exit 0"},
		{expectedExitCode: 2, bashCommand: "exit 2"},
		{expectedExitCode: 5, bashCommand: "exit 5"},
	}
	for _, testCase := range testCases {
		exitCode := RunCmd(
			[]string{"bash", "-c", testCase.bashCommand}, generateTestEnvironment(),
		)
		if exitCode != testCase.expectedExitCode {
			t.Errorf(
				"Bash command %q return exit code %d, but expected %d",
				testCase.bashCommand,
				exitCode,
				testCase.expectedExitCode,
			)
		}
		fmt.Printf(
			"Bash command %q return expected exit code %d\n",
			testCase.bashCommand,
			testCase.expectedExitCode,
		)
	}
}

// TestRunCmd_InternalEnvironmentApply - тест видимости переменной окружения.
// Проверка с позиции системных процессов ОС видимости
// значения устанавливаемой переменной окружения
// в условии инициализации окружения переменных непосредственно внутри функции RunCmd.
func TestRunCmd_002_InternalEnvironmentApply(t *testing.T) {
	fmt.Println("Test RunCmd to use environment var into bash code with catch expected OK return code.")
	environment := generateTestEnvironment()
	email := environment["GOLANG_HW08_TEST_EMAIL"].Value
	bashEnvVarCheck := fmt.Sprintf(
		bashEnvVarCheckTemplate,
		email,
	)
	fmt.Println("RunCmd with internal environment apply")
	fmt.Printf("EnvVar was set to %q\n", email)
	exitCode := RunCmd(
		[]string{"bash", "-c", bashEnvVarCheck},
		environment,
	)
	if exitCode != okExitCode {
		t.Errorf("Bash command\n%s\nreturn unexpected exit code %d, expected %d\n", bashEnvVarCheck, exitCode, okExitCode)
	} else {
		fmt.Printf("Bash command\n%s\nreturn expected exit code %d\n", bashEnvVarCheck, okExitCode)
	}
}

// TestRunCmd_InternalEnvironmentApply_FalsePositive - тест видимости переменной окружения.
// Проверка с позиции системных процессов ОС видимости
// значения устанавливаемой переменной окружения
// в условии инициализации окружения переменных непосредственно внутри функции RunCmd
// и отсуствия ложноположитльного срабатывания (с заведомо ложным значением).
func TestRunCmd_003_InternalEnvironmentApply_FalsePositive(t *testing.T) {
	fmt.Println("Test RunCmd to use environment var into bash code with catch expected FAIL return code.")
	failEnvValValue := "i_am_apriory@failed"
	fmt.Printf("Fail EnvVar value is %q\n", failEnvValValue)
	bashEnvVarCheck := fmt.Sprintf(bashEnvVarCheckTemplate, failEnvValValue)
	environment := generateTestEnvironment()
	realEnvValValue := environment["GOLANG_HW08_TEST_EMAIL"].Value
	exitCode := RunCmd(
		[]string{"bash", "-c", bashEnvVarCheck},
		generateTestEnvironment(),
	)
	fmt.Println("RunCmd with internal environment apply")
	fmt.Printf("EnvVar value was set to %q\n", realEnvValValue)
	if exitCode != failExitCode {
		t.Errorf("Bash command\n%s\nreturn unexpected fail exit code %d for not valid environment var value, expected %d\n",
			bashEnvVarCheck,
			exitCode,
			failExitCode,
		)
	} else {
		fmt.Printf("Bash command\n%s\nreturn expected fail exit code %d for not valid environment var value\n",
			bashEnvVarCheck,
			failExitCode)
	}
}

// TestRunCmd_ExternalEnvironmentApply - тест видимости переменной окружения.
// Проверка с позиции системных процессов ОС видимости значения устанавливаемой переменной
// в условии инициализации окружения переменных вне функции RunCmd.
func TestRunCmd_004_ExternalEnvironmentApply(t *testing.T) {
	fmt.Println("Apply Environment out of at RunCmd-function executing")
	environment := generateTestEnvironment()
	Apply(environment)
	email := environment["GOLANG_HW08_TEST_EMAIL"].Value
	bashEnvVarCheck := fmt.Sprintf(bashEnvVarCheckTemplate, email)

	fmt.Println("Call RunCmd-function with empty environment")
	exitCode := RunCmd(
		[]string{"bash", "-c", bashEnvVarCheck},
		cleanEnvironment,
	)
	if exitCode != okExitCode {
		t.Errorf("Bash command\n%s\nreturn unexpected exit code %d, expected %d\n", bashEnvVarCheck, exitCode, okExitCode)
	} else {
		fmt.Printf("Bash command\n%s\nreturn expected exit code %d\n", bashEnvVarCheck, okExitCode)
	}
}

# Домашнее задание №8 «Утилита envdir»

[Задание](./README.md).

> **Для формирования данного отчета запустить**
>
> ```shell
> $ cd ../report_templator/
> $ go test templator.go hw08_envdir_tool.go
> ```

## Документация

> Сначала я реализовал более широкий вариант исполнения - с возможностью влияния функции исполнения на внешние относительно ее запуска переменные. В последующем разработал RunCmdVariant2, задействующую установление контекста посредством `command.Env`. В демонстрации оставил исходный вариант вызова.

```shell
cd ./hw08_envdir_tool/
go doc -all ./ > hw08_go_doc_-all.txt
```

```text


FUNCTIONS

func Apply(environment Environment)
    Apply - функция применения окружения переменных. Положительное значени
    индикатора NeedRemove указывает на необходимость удаления конкретной
    переменной из окружения.

func Clear(environment Environment)
    Clear - функция принудительной очистки окружения переменных.

func RunCmd(cmd []string, environment Environment) (returnCode int)
    RunCmd - функция запуска процесса ОС с применением к нему окружения
    переменных.

func RunCmdVariant2(cmd []string, environment Environment) (returnCode int)
    RunCmdVariant2 - функция запуска процесса ОС с применением к нему окружения
    переменных. В данной реализации задействуется установление контекста
    посредством `command.Env`.

func UnsetOnlyNeeded(environment Environment)
    UnsetOnlyNeeded - функция принудительной очистки окружения от требующих
    удаления переменных.

func UpdateOrInsert(environment Environment)
    UpdateOrInsert - функция обновления или добавления переменных окружения,
    нетребующих удаления.


TYPES

type EnvValue struct {
    Value      string
    NeedRemove bool
}
    EnvValue - переменная окружения:
      - Value - значение переменной окружения.
      - NeedRemove - требование удаления переменной окружения.

func CreateEnvValue(value string, needRemove bool) EnvValue
    CreateEnvValue - конструктор переменной окружения, при этом:
      - Правая табуляция удаляется.
      - Символы нуль-байт заменяются на символ переноса.

func ReadEnvFile(file string) (EnvValue, error)
    ReadEnvFile - функция парсинга первой строки файла для установки значения
    переменной окружения. При нулевом объеме файла переменная помечается как
    требуемая к удалению.

type Environment map[string]EnvValue
    Environment - окружение переменных.

func ReadEnvDir(envDir string) (Environment, error)
    ReadEnvDir - функция обхода директории для заполнения окружения переменных
    по содержанию файлов. Файлы, содержащие в именовании знак равенства,
    игнорируются.


```

## Тестирование

```shell
go test -v | sed 's/=== RUN/\n\n===RUN/g' > hw08_go_test_-v.txt
```

```text


===RUN   TestSet
--- PASS: TestSet (0.00s)


===RUN   TestGetOfNotExists
--- PASS: TestGetOfNotExists (0.00s)


===RUN   TestReadDir
OK. EnvVar "UNSET": {Value: NeedRemove:true} 
OK. EnvVar "BAR": {Value:bar NeedRemove:false} 
OK. EnvVar "EMPTY": {Value: NeedRemove:false} 
OK. EnvVar "FOO": {Value:   foo
with new line NeedRemove:false} 
OK. EnvVar "HELLO": {Value:"hello" NeedRemove:false} 
--- PASS: TestReadDir (0.00s)


===RUN   TestClearEnvironment
Check EnvVar "GOLANG_HW08_TEST_EMAIL"
Check EnvVar "GOLANG_HW08_TEST_USERAGENT"
Check EnvVar "GOLANG_HW08_TEST_USERNAME"
--- PASS: TestClearEnvironment (0.00s)


===RUN   TestApplyEnvironment
Check EnvVar "GOLANG_HW08_TEST_EMAIL"
    Value MargaretVasquez@Youfeed.mil
    NeedRemove false
GOLANG_HW08_TEST_EMAIL=MargaretVasquez@Youfeed.mil
Check EnvVar "GOLANG_HW08_TEST_USERAGENT"
    Value Mozilla/5.0 (Linux; U; Android 3.0; en-us; Xoom Build/HRI39) AppleWebKit/525.10  (KHTML, like Gecko) Version/3.0.4 Mobile Safari/523.12.2
    NeedRemove false
GOLANG_HW08_TEST_USERAGENT=Mozilla/5.0 (Linux; U; Android 3.0; en-us; Xoom Build/HRI39) AppleWebKit/525.10  (KHTML, like Gecko) Version/3.0.4 Mobile Safari/523.12.2
Check EnvVar "GOLANG_HW08_TEST_USERNAME"
    Value bHansen
    NeedRemove true
--- PASS: TestApplyEnvironment (0.00s)


===RUN   TestRunCmd_001_ExitCode
Test RunCmd to catch exit code of OS process.
qui_suscipit@Dabfeed.biz
Bash command "echo $GOLANG_HW08_TEST_EMAIL" return expected exit code 0
echo lRichardson@Linkbuzz.edu
Bash command "echo echo $GOLANG_HW08_TEST_EMAIL" return expected exit code 0
omnis@Layo.mil
Bash command "echo ${GOLANG_HW08_TEST_EMAIL}" return expected exit code 0
Bash command "exit 1" return expected exit code 1
Bash command "exit 0" return expected exit code 0
Bash command "exit 2" return expected exit code 2
Bash command "exit 5" return expected exit code 5
--- PASS: TestRunCmd_001_ExitCode (0.01s)


===RUN   TestRunCmd_002_InternalEnvironmentApply
Test RunCmd to use environment var into bash code with catch expected OK return code.
RunCmd with internal environment apply
EnvVar was set to "dolores_sed_vel@Einti.gov"
Bash command

if [ ${GOLANG_HW08_TEST_EMAIL} == dolores_sed_vel@Einti.gov ]; then 
    exit 5; 
else 
    exit 17; 
fi;

return expected exit code 5
--- PASS: TestRunCmd_002_InternalEnvironmentApply (0.00s)


===RUN   TestRunCmd_003_InternalEnvironmentApply_FalsePositive
Test RunCmd to use environment var into bash code with catch expected FAIL return code.
Fail EnvVar value is "i_am_apriory@failed"
RunCmd with internal environment apply
EnvVar value was set to "AliceSnyder@Devbug.mil"
Bash command

if [ ${GOLANG_HW08_TEST_EMAIL} == i_am_apriory@failed ]; then 
    exit 5; 
else 
    exit 17; 
fi;

return expected fail exit code 17 for not valid environment var value
--- PASS: TestRunCmd_003_InternalEnvironmentApply_FalsePositive (0.00s)


===RUN   TestRunCmd_004_ExternalEnvironmentApply
Apply Environment out of at RunCmd-function executing
Call RunCmd-function with empty environment
Bash command

if [ ${GOLANG_HW08_TEST_EMAIL} == dBarnes@Mynte.net ]; then 
    exit 5; 
else 
    exit 17; 
fi;

return expected exit code 5
--- PASS: TestRunCmd_004_ExternalEnvironmentApply (0.00s)
PASS
ok      github.com/BorisPlus/hw08_envdir_tool    0.027s

```

## Вывод

* реализован перехват возвращаемого значения запускаемой GOLANG-ом командой операционной системы
* проверено на факт доступа bash-скриптов к определяемым GOLANG-ом значениям переменных окружения

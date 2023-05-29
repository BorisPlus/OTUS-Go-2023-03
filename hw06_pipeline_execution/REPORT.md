# Домашнее задание №6 «Пайплайн»

Описание [задания](./README.md).

> **Для формирования данного отчета запустить**
>
> ```shell
> $ cd ../report_templator/
> $ go test templator.go hw06_pipeline_execution_test.go
> ```

## Реализация

Основная идея разработанного варианта решения в том, что согласно сигнатуре/интерфейсу функции пайплайна `func ExecutePipeline(in In, done In, stages ...Stage) Out` предварительно реализуются каналы на каждого Стейджа: входящий и исходящий. При этом исходящий канал для текущего Стейджа является входящим для Стейджа, следующего по списку из `...Stage` за текщим.

```go
func ExecutePipeline(in In, done In, stages ...Stage) Out {
    stageInput := in
    fmt.Println("Configute STAGING: in", stageInput)

    stageOutput := make(Bi)
    for stageID, stage := range stages {
        go func(stageId int, in In, done In, stage Stage, out Bi) {
            executePipepoint(stageId, in, done, stage, out)
        }(stageID, stageInput, done, stage, stageOutput)

        fmt.Println("Configute STAGING: stage", stageID)
        fmt.Println("Configute STAGING: stage", stageID, "with done", fmt.Sprintf("%p", done))
        fmt.Println("Configute STAGING: stage", stageID, "with in", fmt.Sprintf("%p", stageInput))
        fmt.Println("Configute STAGING: stage", stageID, "with out", fmt.Sprintf("%p", stageOutput))

        stageInput = stageOutput
        stageOutput = make(Bi)
    }

    fmt.Println("Configute STAGING: out", stageInput)
    fmt.Println()
    return stageInput
}
```

<details>
<summary>см. подробный код "pipeline.go"</summary>

```go
package hw06pipelineexecution

import (
    "fmt"
)

type (
    In  = <-chan interface{}
    Out = In
    Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
    stageInput := in
    fmt.Println("Configute STAGING: in", stageInput)

    stageOutput := make(Bi)
    for stageID, stage := range stages {
        go func(stageId int, in In, done In, stage Stage, out Bi) {
            executePipepoint(stageId, in, done, stage, out)
        }(stageID, stageInput, done, stage, stageOutput)

        fmt.Println("Configute STAGING: stage", stageID)
        fmt.Println("Configute STAGING: stage", stageID, "with done", fmt.Sprintf("%p", done))
        fmt.Println("Configute STAGING: stage", stageID, "with in", fmt.Sprintf("%p", stageInput))
        fmt.Println("Configute STAGING: stage", stageID, "with out", fmt.Sprintf("%p", stageOutput))

        stageInput = stageOutput
        stageOutput = make(Bi)
    }

    fmt.Println("Configute STAGING: out", stageInput)
    fmt.Println()
    return stageInput
}

func executePipepoint(stageID int, in In, done In, stage Stage, out Bi) {
    processor := func(stageId int, stage Stage, in In, out Bi) Out {
        staged := stage(in)
        terminated := make(Bi)
        go func() {
            defer func() {
                fmt.Println("stage", stageId, "processor", "end")

                fmt.Println("stage", stageId, "processor", "try to close(out)")
                close(out)
                fmt.Println("stage", stageId, "processor", "close(out)")

                fmt.Println("stage", stageId, "processor", "try setup terminated")
                close(terminated)
                fmt.Println("stage", stageId, "processor", "was setup terminated")
            }()
            for {
                select {
                case value, ok := <-staged:
                    fmt.Println("stage", stageId, "processor", "get from input", "value", value, "ok", ok)
                    if !ok {
                        fmt.Println("stage", stageId, "processor", "!ok - return")
                        return
                    }
                    fmt.Println("stage", stageId, "processor", "get from input", "value", value, "ok", ok, "try put to out")
                    out <- value
                    fmt.Println("stage", stageId, "processor", "get from input", "value", value, "ok", ok, "was put to out")
                case <-done:
                    fmt.Println("stage", stageId, "processor", "done - return")
                    return
                }
            }
        }()
        return terminated
    }
    terminated := processor(stageID, stage, in, out)
    fmt.Println("stage", stageID, "try", "terminated")
    <-terminated
    fmt.Println("stage", stageID, "was", "terminated")
}

```

</details>

в вербоуз-тестовом запуске будет продемонстирована построенная "цепочка" из каналов для каждого Этапа

```text
...
Configute STAGING: in 0xc00008e120 ====================╗ Это один
Configute STAGING: stage 0                             ║ и тот же канал.
Configute STAGING: stage 0 with done 0xc00008e180      ║
Configute STAGING: stage 0 with in 0xc00008e120    ====╝ --------╮ Тут срабатывает Стейдж 0 - как
Configute STAGING: stage 0 with out 0xc00008e1e0   ====╗ <-------╯ перекладчик из канала в канал.
Configute STAGING: stage 1                             ║ Это один  
Configute STAGING: stage 1 with done 0xc00008e180      ║ и тот же канал.
Configute STAGING: stage 1 with in 0xc00008e1e0    ====╝ --------╮ Тут срабатывает Стейдж 1 - как
Configute STAGING: stage 1 with out 0xc00008e240   ====╗ <-------╯ перекладчик из канала в канал.
Configute STAGING: stage 2                             ║ Это один 
Configute STAGING: stage 2 with done 0xc00008e180      ║ и тот же канал.
Configute STAGING: stage 2 with in 0xc00008e240    ====╝ --------╮ Тут срабатывает Стейдж 2 - как
Configute STAGING: stage 2 with out 0xc00008e2a0   ====╗ <-------╯ перекладчик из канала в канал.
Configute STAGING: stage 3                             ║ Это один  
Configute STAGING: stage 3 with done 0xc00008e180      ║ и тот же канал.
Configute STAGING: stage 3 with in 0xc00008e2a0    ====╝ --------╮ Тут срабатывает Стейдж 3 - как
Configute STAGING: stage 3 with out 0xc00008e300   ====╗ <-------╯ перекладчик из канала в канал.
Configute STAGING: out 0xc00008e300 ===================╝ Это один и тот же канал.
...
```

В целях неблокирования работы Стейджей их функционал обернут в горутины.

Внутри горутин:

* происходит обработка поступающей информации из входящего канала данного Стейджа;
* перенаправление результата в выходной канал данного Стейджа;
* посредством конструкиции `value, ok := <-staged` учитывается факт получения сигнала, что входящий канал Стейджа стал закрыт и было получено значение по умолчанию;
* учитывается ветвление логики с принудительным завершением Стейджа при срабатывании `done`-канала (он для всех Стейджей один и тот же).

Решение об отсутствии в необходимости в горутине принимается на основе доступности специального пустого канала `<-terminated` в результате его закрытия, происходящего совместного с закрытием входящего канала соотвествующего Стейджа.

Если взглянуть на журнал работы в вербоуз-теста, то видно, что:

* Стейджи запускаются изначально хаотично;
  
>```text
> ...
>stage 3 try terminated
>stage 1 try terminated
>stage 2 try terminated
>stage 0 try terminated
> ...
>```

* Стейджи работают параллельно;

>```text
> ...
>stage 0 processor get from input value 1 ok true                   Получили Стейджем № 0 как есть 1 - "Dummy"
>stage 0 processor get from input value 1 ok true try put to out    
>stage 0 processor get from input value 1 ok true was put to out    Переложили в выходной канал Стейджа № 0
>stage 0 processor get from input value 2 ok true                   Получили Стейджем № 0 как есть 2 - "Dummy"
>stage 0 processor get from input value 2 ok true try put to out
>stage 0 processor get from input value 2 ok true was put to out    Переложили в выходной канал Стейджа № 0
>stage 1 processor get from input value 2 ok true                   Преобразовали 1 в 2 Стейджем № 1 - "Multiplier (* 2)"
>stage 1 processor get from input value 2 ok true try put to out
>stage 1 processor get from input value 2 ok true was put to out    Переложили в выходной канал Стейджа № 1
>stage 1 processor get from input value 4 ok true                   Преобразовали 2 в 4 Стейджем № 1 - "Multiplier (* 2)"
>stage 1 processor get from input value 4 ok true try put to out
> ... и так далее
>```

Запустим тест без вербоуз-режима (инече мой лог будет большой)

```bash
go test -race -count=100 ./pipeline.go ./pipeline_test.go > N100TimesTesting.txt
```

Результат 100 запусков:

```text
ok      command-line-arguments    90.950s

```

### Дополнительное тестирование

Пусть Стейджи представлены двумя Слиперами по 2 и 8 секунд

```go
stages := []Stage{
    g("Sleep (2 sec)", func(v interface{}) interface{} { time.Sleep(2 * time.Second); return v }),
    g("Sleep (8_sec)", func(v interface{}) interface{} { time.Sleep(8 * time.Second); return v }),
}
```

Соответственно при запуске для данных

```go
data := []int{1, 2, 3}
```

Результат должен быть меньше 30 секунд

```go
start := time.Now()
for s := range ExecutePipeline(in, nil, stages...) {
    _ = s
}
elapsed := time.Since(start).Seconds()

require.Less(t,
    int64(elapsed),      // Засеченное время
    int64(10*len(data))) // 10*3 = 30 сек
```

Практический результат теста:

```bash
go test -run TestPipelineConcurencyTime ./pipeline.go ./pipeline_test.go
```

подтверждает теорию:

```text
ok      command-line-arguments    26.020s

```

## Вывод

Как продемонстрировано, присутствует конкурентный доступ горутин к каналам, неблокирующим параллельное выполнение Стейджей.

Если удалить комментарии, то объем кода станет соотвествовать требуемому в ~55 строк (авторская реализация не известна).

## На возможную доработку

На мой взгляд не хватает разве что ограничений на число одновременно работающих Стейджей, как делалось в предыдущей главе в задаче лимитирования числа одновременных `worker`.

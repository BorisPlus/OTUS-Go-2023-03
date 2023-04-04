# Вопрос относительно объема бинарного файла

## Посыл

Мне казалось, что объем самодостаточного варианта `./main.go`, который содержит весь код в себе и не зависит от внешних библиотек, должен быть меньше, чем вариант `./main.go`, имеющий зависимость от внешнего пакета.
Тем бодее ведь во внешнем пакете еще есть дополнительный код теста `reverse_test.go`.

Однако это не так.

## Исследование

Измерим объемы созаваемых альтернативных решений.

### Внешняя зависимость

Так, вариант "A" `./main.go` - имеющий зависимость от внешней библиотеки `stringutil`:

```go
package main

import (
    "fmt"

    "golang.org/x/example/stringutil"
)

func main() {
    row := "Hello, OTUS!"
    reversedRow := stringutil.Reverse(row)
    fmt.Println(reversedRow)
}
```

после компиляции:

```shell
go build main.go
```

занимает объем:

```bash
$ ls -la ./main
-rwxr-xr-x 1 b b 1820656 апр  1 00:19 ./main
```

1820656 байт.

### Самодостаточный

Тогда как вариант "B" `./question/main.go` - НЕ имеющий зависимость от внешних библиотек:

```go
package main

import "fmt"

func Reverse(s string) string {
    r := []rune(s)
    for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}

func main() {
    row := "Hello, OTUS!"
    reversedRow := Reverse(row)
    fmt.Println(reversedRow)
}
```

и содержащий абсолютно тот же самый алгоритм, что и внешний репозиторий:

```shell
go mod vendor -v
```

> ```text
> # golang.org/x/example v0.0.0-20220412213650-2e68773dfca0
> ## explicit; go 1.15
> golang.org/x/example/stringutil
> ```

```shell
cat ./vendor/golang.org/x/example/stringutil/reverse.go
```

> ```go
> // Package stringutil contains utility functions for working with strings.
> package stringutil
> 
> // Reverse returns its argument string reversed rune-wise left to right.
> func Reverse(s string) string {
>         r := []rune(s)
>         for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
>                 r[i], r[j] = r[j], r[i]
>         }
>         return string(r)
> }
> ```

```shell
cd
cat ./go/pkg/mod/golang.org/x/example@v0.0.0-20220412213650-2e68773dfca0/stringutil/reverse.go
```

> ```go
> // Package stringutil contains utility functions for working with strings.
> package stringutil
> 
> // Reverse returns its argument string reversed rune-wise left to right.
> func Reverse(s string) string {
>         r := []rune(s)
>         for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
>                 r[i], r[j] = r[j], r[i]
>         }
>         return string(r)
> }
> ```

после компиляции:

```shell
go build ./question/main.go 
```

занимает объем:

```bash
$ ls -la ./main
-rwxr-xr-x 1 b b 1821168 апр  1 00:29 ./main
```

1821168 байт, что немного больше 1820656 байт.

Теперь обратимся к варианту "C" - создадим свой публичный репозиторий и задействуем выше приведенный код

```go
package stringutil

func Reverse(s string) string {
    r := []rune(s)
    for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}
```

из него

```go
package main

import (
    "fmt"

    "github.com/account/stringutil"
)

func main() {
    row := "Hello, OTUS!"
    reversedRow := stringutil.Reverse(row)
    fmt.Println(reversedRow)
}
```

В итоге при компиляции

```bash
$ go build main.go
$ ls -la ./main
-rwxr-xr-x 1 b b 1820672 апр  4 00:07 ./main
```

вес будет 1820672 байт. Смешно, но это уже почти неустановленной причины "идеал", но между 1820656 и 1821168 байт.

А если код изменить так (вариант "D")

```go
package main

import (
    "fmt"

    "github.com/account/stringutil"
)

func main() {
    fmt.Println(stringutil.Reverse("Hello, OTUS!"))
}
```

то объем опять увеличится

```bash
$ go build main.go
$ ls -la ./main
-rwxr-xr-x 1 b b 1820952 апр  4 00:18 ./main
```

и станет 1820952 байт.

### Вывод

Объем одного и того же дистрибутива различается, хотя исходные коды скрипта и используемого метода абсолютно одинаковы

```text
A=1820656 байт (штатная внешняя зависимость "golang.org/x/example/stringutil") < 
  < C=1820672 байт (собственноручная внешняя зависимость "github.com/account/stringutil") <
      < D=1820952 байт (собственноручная внешняя зависимость "github.com/account/stringutil" с компактным исходным кодом) < 
          < B=1821168 байт (без зависимостей от каких-либо сторонних репозиториев, весь код зависимости внутри исходника) 
```

Мне кажется это весьма интересным, так как значит, что внешняя зависимость не только не приводит к росту объема итогового бинарного результата, но почему-то даже его уменьшает.

К однозначности в объеме приводит следующее:

```bash
$ go build -ldflags "-s -w" ./main.go 
$ ls -la ./main
-rwxr-xr-x 1 b b 1216512 апр  4 13:48 ./main
```

В итоге размер во всех ЧЕТЫРЕХ вариантах дистрибутива 1216512 байт.

```bash
go build -ldflags "-s -w" -o ./question/a/main ./question/a/main.go 
go build -ldflags "-s -w" -o ./question/b/main ./question/b/main.go 
go build -ldflags "-s -w" -o ./question/c/main ./question/c/main.go 
go build -ldflags "-s -w" -o ./question/d/main ./question/d/main.go 

b@b:~/vscode/OTUS-Go-2023-03/OTUS-Go-2023-03/hw01_hello_otus$ ls -la ./question/a/main*
    -rwxr-xr-x 1 b b 1216512 апр  4 14:06 ./question/a/main
    -rw-r--r-- 1 b b     174 апр  4 13:56 ./question/a/main.go

b@b:~/vscode/OTUS-Go-2023-03/OTUS-Go-2023-03/hw01_hello_otus$ ls -la ./question/b/main*
    -rwxr-xr-x 1 b b 1216512 апр  4 14:06 ./question/b/main
    -rw-r--r-- 1 b b     278 апр  3 23:52 ./question/b/main.go

b@b:~/vscode/OTUS-Go-2023-03/OTUS-Go-2023-03/hw01_hello_otus$ ls -la ./question/c/main*
    -rwxr-xr-x 1 b b 1216512 апр  4 14:06 ./question/c/main
    -rw-r--r-- 1 b b     174 апр  4 13:56 ./question/c/main.go

b@b:~/vscode/OTUS-Go-2023-03/OTUS-Go-2023-03/hw01_hello_otus$ ls -la ./question/d/main*
    -rwxr-xr-x 1 b b 1216512 апр  4 14:06 ./question/d/main
    -rw-r--r-- 1 b b     134 апр  4 13:56 ./question/d/main.go
```

Но поведение ранее, согласитесь, интересное.

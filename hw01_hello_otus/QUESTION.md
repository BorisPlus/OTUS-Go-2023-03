# Вопрос относительно объема бинарного файла

## Посыл

Мне казалось, что объем самодостаточного варианта `./main.go`, который содержит весь код в себе и не зависит от внешних библиотек, должен быть меньше, чем вариант `./main.go`, имеющий зависимость от внешнего пакета.
Тем бодее ведь во внешнем пакете еще есть дополнительный код теста `reverse_test.go`.

Однако это не так.

## Исследование

Измерим объемы созаваемых альтернативных решений.

### Внешняя зависимость

Так, вариант "А" `./main.go` - имеющий зависимость от внешней библиотеки `stringutil`:

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

Тогда как вариант "Б" `./question/main.go` - НЕ имеющий зависимость от внешних библиотек:

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
cat ./go/pkg/mod/golang.org/x/example@v0.0.0-20220412213650-2e68773dfca0/stringutil/reverse
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

Теперь обратимся к варианту "B" - создадим свой публичный репозиторий и задействуем выше приведенный код

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

А если код изменить так (вариант "Г")

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

```text
А=1820656 байт (штатная внешняя зависимость "golang.org/x/example/stringutil") < 
  < В=1820672 байт (собственноручная внешняя зависимость "github.com/account/stringutil") <
      < Г=1820952 байт (собственноручная внешняя зависимость "github.com/account/stringutil" с компактным исходным кодом) < 
          < Б=1821168 байт (без зависимостей от каких-либо сторонних репозиториев, весь код зависимости внутри исходника) 
```

Мне кажется это весьма интересным, так как значит, что внешняя зависимость не только не приводит к росту объема итогового бинарного результата, но почему-то даже его уменьшает.

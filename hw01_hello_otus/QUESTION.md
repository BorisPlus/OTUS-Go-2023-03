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

> ```
> # golang.org/x/example v0.0.0-20220412213650-2e68773dfca0
> ## explicit; go 1.15
> golang.org/x/example/stringutil
> ```

```shell
cat ./vendor/golang.org/x/example/stringutil/reverse.go
```

> 
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
> 


```shell
cd
cat ./go/pkg/mod/golang.org/x/example@v0.0.0-20220412213650-2e68773dfca0/stringutil/reverse
```

> 
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
> 

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

### Вывод

Мне кажется это весьма интересным, так как значит, что внешняя зависимость не только не приводит к росту объема итогового бинарного результата, но почему-то даже его уменьшает.

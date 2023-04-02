#

## Домашнее задание №1 «Hello, OTUS!»

Необходимо написать программу, печатающую в стандартный вывод перевернутую фразу

```text
Hello, OTUS!
```

Для переворота строки следует воспользоваться возможностями
[golang.org/x/example/stringutil](https://github.com/golang/example/tree/master/stringutil).

Кроме этого необходимо исправить **go.mod** так, чтобы для данного модуля работала
команда `go get`, а полученный **go.sum** закоммитить.

### Критерии оценки

- Пайплайн зелёный - 4 балла
- Используется `stringutil` - 4 балла
- Понятность и чистота кода - до 2 баллов

#### Зачёт от 7 баллов

### Подсказки

- `Reverse`

### Решение

Сначала я добавил в код

```go
import "github.com/golang/example/blob/master/stringutil"
```

но подсказка при исполнении дала

```bash
github.com/fixme_my_friend/hw01_hello_otus imports
        github.com/golang/example/blob/master/stringutil: github.com/golang/example@v0.0.0-20220412213650-2e68773dfca0: parsing go.mod:
        module declares its path as: golang.org/x/example
                but was required as: github.com/golang/example
```

Добавил как "golang.org/x/example/stringutil" в main.go и "require golang.org/x/example v0.0.0-20220412213650-2e68773dfca0" в go.mod

```bash
$ go mod tidy
$ gofmt -w main.go 
$ golangci-lint run .
$ go test -v -count=1 -race -timeout=1m .
?       github.com/fixme_my_friend/hw01_hello_otus      [no test files]
$ ./test.sh

+ expected='!SUTO ,olleH'
++ go run main.go
++ sed 's/> *//;s/ *$//'
+ result='!SUTO ,olleH'
+ '[' '!SUTO ,olleH' = '!SUTO ,olleH' ']'
+ echo PASS
PASS

$ go run ./main.go 
!SUTO ,olleH
```

<details>
<summary> Content of `go.sum`-file</summary>

```text
github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
golang.org/x/example v0.0.0-20220412213650-2e68773dfca0 h1:ygD+9PaH9IfzZUF131IxmiXGkxzuN/pphDjzh2LY8N8=
golang.org/x/example v0.0.0-20220412213650-2e68773dfca0/go.mod h1:+yakPl5KR9J+ysfUNADYwEU5qeqjUO473wDktD4xMYw=
golang.org/x/mod v0.3.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20210112183307-1e6ecd4bf1b0/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=

```

</details>

Таким образом задача выполнена: решение с `stringutil` приведено, `go get` - работает.

### Вопрос относительно объема бинарного файла

Мне казалось, что объем самодостаточного варианта `./main.go`, который содержит весь код в себе и не зависит от внешних библиотек, должен быть меньше, чем вариант `./main.go`, имеющий зависимость от внешнего пакета.
Тем бодее ведь во внешнем пакете еще есть дополнительный код теста `reverse_test.go`. 

Однако это не так.

Измерим объемы созаваемых альтернативных решений.

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

Мне кажется это весьма интересным, так как значит, что внешняя зависимость не только не приводит к росту объема итогового бинарного результата, но почему-то даже его уменьшает.

# Домашнее задание №7 «Утилита для копирования файлов»

Описание [задания](./README.md).

## Реализация

Реализованы варианты копирования файла:

* последовательного побайтного копирования файла [`Сopy()`](./copy.go)
* сегментированного с задействованием горутин [`СopySegmented()`](./copy_segmented.go) (файл сегментами выбирается и сегментами записывается, при этом это визуально выглядит как формирование файла по торрент-протоколу)

> В исходном коде `main.go` содержится обе реализации, но крайняя - для сегментированного, исходная простая - закомментирована.

```bash
go build -o copier.goc
```

>
> Для игнорирования бинарных Go-файлов на уровне репозитория ввел расширение - `*.goc`

### Реализация `Сopy()`

```bash
./copier.goc --help
        Usage "copier.goc":
        -from
          file to read from
        -to
          file to write to
        -limit
          byte-limit to copy (default: 4096)
        -offset
          byte-offset of input file (default: 0)
        -perc
          indicate percent processing (default: false)
        -v
          verbose log output (default: false)

        Example:
          ./copier.goc -from="/testdata/input.txt" -to="./testdata/output.txt" -limit=1000 -offset=0 -perc=true -v=false
```

```bash
./copier.goc -from="./testdata/input.txt" -to="./testdata/output.txt" -limit=10000 -perc=true
```

и ниже появится обновляющаяся автоматически строка с увеличивающимся процентом исполнения

```bash
Copy ...22.47%
```

и далее

```bash
Copy ...100.00%
```

### Реализация `СopySegmented()`

> **ПОДВОДНЫЙ КАМЕНЬ**: важно знать, что `bufio.NewReader` - в любом случае перемещает курсор при чтении на 4086 байт, даже если не считан буфером до конца. Объявление такого "ридера" и его буфера рядом внутри цикла чтения не позволит так просто последовательно получить все порции файла при сегментировании, меньшем 4096 байт. Моим решением стало использование `bufio.NewReaderSize`.

```bash
./copier.goc --help
        Usage "copier.goc":
        -from
          file to read from
        -to
          file to write to
        -limit
          byte-limit to copy (default: 4096)
        -offset
         byte-offset of input file (default: 0)
        -segment
          byte-segmentation count (default: 1)
        -writers
          parallel writers count(default: 1)
        -perc
          indicate percent processing (default: false)
        -v
          verbose log output (default: false)

        Example:
          ./copier.goc -from="./testdata/input.txt" \
                       -to="./testdata/output.txt" \
                       -limit=1000 \
                       -offset=0 \
                       -perc=true \
                       -v=false

```

```bash
./copier.goc -from="./testdata/input.txt" \
             -to="./testdata/output.txt" \
             -limit=1000 \
             -offset=0 \
             -perc=true \
             -v=false
```

и ниже появится обновляющаяся автоматически строка с увеличивающимся процентом исполнения

```bash
СopySegmented ...84.23%
```

и в конце

```bash
СopySegmented ...100.00%
```

Gif-анимация:

![alt="percentager"](./REPORT.files/percentager.gif)

Если включить `verbose`-режим, например, с такими флагами

> `perc=true` в режиме `v=true` игнорируется

```bash
./copier.goc -from="./testdata/input.txt" \
               -to="./testdata/output.txt" \
               -limit=1000 \
               -offset=0 \
               -segment=200 \
               -writers=2 \
               -perc=true \  
               -v=true
```

то будет продемонстрирован весь лог работы (чтение блоков и их запись в параллельном наборе горутин):

<details><summary>см. лог:</summary>

```text
2023/06/24 23:42:27 СopySegmented
2023/06/24 23:42:27 from = ./testdata/input.txt
2023/06/24 23:42:27 offset = 0
2023/06/24 23:42:27 limit = 1000
2023/06/24 23:42:27 fileInfo.Size() = 6617
2023/06/24 23:42:27 repairLimit = 1000
2023/06/24 23:42:27 
2023/06/24 23:42:27 READER (No. 1): initial segment size = 200
2023/06/24 23:42:27 STEP READ n = 1
2023/06/24 23:42:27 STEP READ n = 1 - repairLimit = 1000
2023/06/24 23:42:27 STEP READ n = 1 - inisegmentSize = 200
2023/06/24 23:42:27 STEP READ n = 1 - prevSegmentsSizesSum = 0
2023/06/24 23:42:27 READER (No. 1): segment size = 200
2023/06/24 23:42:27 READER (No. 1): expected read = 200
2023/06/24 23:42:27 READER (No. 1): result read = 200
2023/06/24 23:42:27 WRITER (No. 0): offset = 0
2023/06/24 23:42:27 WRITER (No. 0): data
00000000  47 6f 0a 44 6f 63 75 6d  65 6e 74 73 0a 50 61 63  |Go.Documents.Pac|
00000010  6b 61 67 65 73 0a 54 68  65 20 50 72 6f 6a 65 63  |kages.The Projec|
00000020  74 0a 48 65 6c 70 0a 42  6c 6f 67 0a 50 6c 61 79  |t.Help.Blog.Play|
00000030  0a 53 65 61 72 63 68 0a  0a 47 65 74 74 69 6e 67  |.Search..Getting|
00000040  20 53 74 61 72 74 65 64  0a 49 6e 73 74 61 6c 6c  | Started.Install|
00000050  20 74 68 65 20 47 6f 20  74 6f 6f 6c 73 0a 54 65  | the Go tools.Te|
00000060  73 74 20 79 6f 75 72 20  69 6e 73 74 61 6c 6c 61  |st your installa|
00000070  74 69 6f 6e 0a 49 6e 73  74 61 6c 6c 69 6e 67 20  |tion.Installing |
00000080  65 78 74 72 61 20 47 6f  20 76 65 72 73 69 6f 6e  |extra Go version|
00000090  73 0a 55 6e 69 6e 73 74  61 6c 6c 69 6e 67 20 47  |s.Uninstalling G|
000000a0  6f 0a 47 65 74 74 69 6e  67 20 68 65 6c 70 0a 44  |o.Getting help.D|
000000b0  6f 77 6e 6c 6f 61 64 20  74 68 65 20 47 6f 20 64  |ownload the Go d|
000000c0  69 73 74 72 69 62 75 74                           |istribut|

2023/06/24 23:42:27 WRITER (No. 0): must be len() = 200
2023/06/24 23:42:27 WRITER (No. 0): wrotet len()=200
2023/06/24 23:42:27 READER (No. 1): put in channel-segment with len() 200, at offset 0 with data
00000000  47 6f 0a 44 6f 63 75 6d  65 6e 74 73 0a 50 61 63  |Go.Documents.Pac|
00000010  6b 61 67 65 73 0a 54 68  65 20 50 72 6f 6a 65 63  |kages.The Projec|
00000020  74 0a 48 65 6c 70 0a 42  6c 6f 67 0a 50 6c 61 79  |t.Help.Blog.Play|
00000030  0a 53 65 61 72 63 68 0a  0a 47 65 74 74 69 6e 67  |.Search..Getting|
00000040  20 53 74 61 72 74 65 64  0a 49 6e 73 74 61 6c 6c  | Started.Install|
00000050  20 74 68 65 20 47 6f 20  74 6f 6f 6c 73 0a 54 65  | the Go tools.Te|
00000060  73 74 20 79 6f 75 72 20  69 6e 73 74 61 6c 6c 61  |st your installa|
00000070  74 69 6f 6e 0a 49 6e 73  74 61 6c 6c 69 6e 67 20  |tion.Installing |
00000080  65 78 74 72 61 20 47 6f  20 76 65 72 73 69 6f 6e  |extra Go version|
00000090  73 0a 55 6e 69 6e 73 74  61 6c 6c 69 6e 67 20 47  |s.Uninstalling G|
000000a0  6f 0a 47 65 74 74 69 6e  67 20 68 65 6c 70 0a 44  |o.Getting help.D|
000000b0  6f 77 6e 6c 6f 61 64 20  74 68 65 20 47 6f 20 64  |ownload the Go d|
000000c0  69 73 74 72 69 62 75 74                           |istribut|
.
2023/06/24 23:42:27 
2023/06/24 23:42:27 READER (No. 2): initial segment size = 200
2023/06/24 23:42:27 STEP READ n = 2
2023/06/24 23:42:27 STEP READ n = 2 - repairLimit = 1000
2023/06/24 23:42:27 STEP READ n = 2 - inisegmentSize = 200
2023/06/24 23:42:27 STEP READ n = 2 - prevSegmentsSizesSum = 200
2023/06/24 23:42:27 READER (No. 2): segment size = 200
2023/06/24 23:42:27 READER (No. 2): expected read = 200
2023/06/24 23:42:27 READER (No. 2): result read = 200
2023/06/24 23:42:27 READER (No. 2): put in channel-segment with len() 200, at offset 200 with data
00000000  69 6f 6e 0a 44 6f 77 6e  6c 6f 61 64 20 47 6f 0a  |ion.Download Go.|
00000010  43 6c 69 63 6b 20 68 65  72 65 20 74 6f 20 76 69  |Click here to vi|
00000020  73 69 74 20 74 68 65 20  64 6f 77 6e 6c 6f 61 64  |sit the download|
00000030  73 20 70 61 67 65 0a 4f  66 66 69 63 69 61 6c 20  |s page.Official |
00000040  62 69 6e 61 72 79 20 64  69 73 74 72 69 62 75 74  |binary distribut|
00000050  69 6f 6e 73 20 61 72 65  20 61 76 61 69 6c 61 62  |ions are availab|
00000060  6c 65 20 66 6f 72 20 74  68 65 20 46 72 65 65 42  |le for the FreeB|
00000070  53 44 20 28 72 65 6c 65  61 73 65 20 31 30 2d 53  |SD (release 10-S|
00000080  54 41 42 4c 45 20 61 6e  64 20 61 62 6f 76 65 29  |TABLE and above)|
00000090  2c 20 4c 69 6e 75 78 2c  20 6d 61 63 4f 53 20 28  |, Linux, macOS (|
000000a0  31 30 2e 31 30 20 61 6e  64 20 61 62 6f 76 65 29  |10.10 and above)|
000000b0  2c 20 61 6e 64 20 57 69  6e 64 6f 77 73 20 6f 70  |, and Windows op|
000000c0  65 72 61 74 69 6e 67 20                           |erating |
.
2023/06/24 23:42:27 
2023/06/24 23:42:27 READER (No. 3): initial segment size = 200
2023/06/24 23:42:27 STEP READ n = 3
2023/06/24 23:42:27 STEP READ n = 3 - repairLimit = 1000
2023/06/24 23:42:27 STEP READ n = 3 - inisegmentSize = 200
2023/06/24 23:42:27 STEP READ n = 3 - prevSegmentsSizesSum = 400
2023/06/24 23:42:27 READER (No. 3): segment size = 200
2023/06/24 23:42:27 READER (No. 3): expected read = 200
2023/06/24 23:42:27 READER (No. 3): result read = 200
2023/06/24 23:42:27 READER (No. 3): put in channel-segment with len() 200, at offset 400 with data
00000000  73 79 73 74 65 6d 73 20  61 6e 64 20 74 68 65 20  |systems and the |
00000010  33 32 2d 62 69 74 20 28  33 38 36 29 20 61 6e 64  |32-bit (386) and|
00000020  20 36 34 2d 62 69 74 20  28 61 6d 64 36 34 29 20  | 64-bit (amd64) |
00000030  78 38 36 20 70 72 6f 63  65 73 73 6f 72 20 61 72  |x86 processor ar|
00000040  63 68 69 74 65 63 74 75  72 65 73 2e 0a 0a 49 66  |chitectures...If|
00000050  20 61 20 62 69 6e 61 72  79 20 64 69 73 74 72 69  | a binary distri|
00000060  62 75 74 69 6f 6e 20 69  73 20 6e 6f 74 20 61 76  |bution is not av|
00000070  61 69 6c 61 62 6c 65 20  66 6f 72 20 79 6f 75 72  |ailable for your|
00000080  20 63 6f 6d 62 69 6e 61  74 69 6f 6e 20 6f 66 20  | combination of |
00000090  6f 70 65 72 61 74 69 6e  67 20 73 79 73 74 65 6d  |operating system|
000000a0  20 61 6e 64 20 61 72 63  68 69 74 65 63 74 75 72  | and architectur|
000000b0  65 2c 20 74 72 79 20 69  6e 73 74 61 6c 6c 69 6e  |e, try installin|
000000c0  67 20 66 72 6f 6d 20 73                           |g from s|
.
2023/06/24 23:42:27 
2023/06/24 23:42:27 READER (No. 4): initial segment size = 200
2023/06/24 23:42:27 STEP READ n = 4
2023/06/24 23:42:27 STEP READ n = 4 - repairLimit = 1000
2023/06/24 23:42:27 STEP READ n = 4 - inisegmentSize = 200
2023/06/24 23:42:27 STEP READ n = 4 - prevSegmentsSizesSum = 600
2023/06/24 23:42:27 READER (No. 4): segment size = 200
2023/06/24 23:42:27 READER (No. 4): expected read = 200
2023/06/24 23:42:27 READER (No. 4): result read = 200
2023/06/24 23:42:27 WRITER (No. 1): offset = 400
2023/06/24 23:42:27 WRITER (No. 1): data
00000000  73 79 73 74 65 6d 73 20  61 6e 64 20 74 68 65 20  |systems and the |
00000010  33 32 2d 62 69 74 20 28  33 38 36 29 20 61 6e 64  |32-bit (386) and|
00000020  20 36 34 2d 62 69 74 20  28 61 6d 64 36 34 29 20  | 64-bit (amd64) |
00000030  78 38 36 20 70 72 6f 63  65 73 73 6f 72 20 61 72  |x86 processor ar|
00000040  63 68 69 74 65 63 74 75  72 65 73 2e 0a 0a 49 66  |chitectures...If|
00000050  20 61 20 62 69 6e 61 72  79 20 64 69 73 74 72 69  | a binary distri|
00000060  62 75 74 69 6f 6e 20 69  73 20 6e 6f 74 20 61 76  |bution is not av|
00000070  61 69 6c 61 62 6c 65 20  66 6f 72 20 79 6f 75 72  |ailable for your|
00000080  20 63 6f 6d 62 69 6e 61  74 69 6f 6e 20 6f 66 20  | combination of |
00000090  6f 70 65 72 61 74 69 6e  67 20 73 79 73 74 65 6d  |operating system|
000000a0  20 61 6e 64 20 61 72 63  68 69 74 65 63 74 75 72  | and architectur|
000000b0  65 2c 20 74 72 79 20 69  6e 73 74 61 6c 6c 69 6e  |e, try installin|
000000c0  67 20 66 72 6f 6d 20 73                           |g from s|

2023/06/24 23:42:27 WRITER (No. 1): must be len() = 200
2023/06/24 23:42:27 WRITER (No. 1): wrotet len()=200
2023/06/24 23:42:27 WRITER (No. 1): offset = 600
2023/06/24 23:42:27 WRITER (No. 1): data
00000000  6f 75 72 63 65 20 6f 72  20 69 6e 73 74 61 6c 6c  |ource or install|
00000010  69 6e 67 20 67 63 63 67  6f 20 69 6e 73 74 65 61  |ing gccgo instea|
00000020  64 20 6f 66 20 67 63 2e  0a 0a 53 79 73 74 65 6d  |d of gc...System|
00000030  20 72 65 71 75 69 72 65  6d 65 6e 74 73 0a 47 6f  | requirements.Go|
00000040  20 62 69 6e 61 72 79 20  64 69 73 74 72 69 62 75  | binary distribu|
00000050  74 69 6f 6e 73 20 61 72  65 20 61 76 61 69 6c 61  |tions are availa|
00000060  62 6c 65 20 66 6f 72 20  74 68 65 73 65 20 73 75  |ble for these su|
00000070  70 70 6f 72 74 65 64 20  6f 70 65 72 61 74 69 6e  |pported operatin|
00000080  67 20 73 79 73 74 65 6d  73 20 61 6e 64 20 61 72  |g systems and ar|
00000090  63 68 69 74 65 63 74 75  72 65 73 2e 20 50 6c 65  |chitectures. Ple|
000000a0  61 73 65 20 65 6e 73 75  72 65 20 79 6f 75 72 20  |ase ensure your |
000000b0  73 79 73 74 65 6d 20 6d  65 65 74 73 20 74 68 65  |system meets the|
000000c0  73 65 20 72 65 71 75 69                           |se requi|

2023/06/24 23:42:27 WRITER (No. 1): must be len() = 200
2023/06/24 23:42:27 WRITER (No. 0): offset = 200
2023/06/24 23:42:27 WRITER (No. 0): data
00000000  69 6f 6e 0a 44 6f 77 6e  6c 6f 61 64 20 47 6f 0a  |ion.Download Go.|
00000010  43 6c 69 63 6b 20 68 65  72 65 20 74 6f 20 76 69  |Click here to vi|
00000020  73 69 74 20 74 68 65 20  64 6f 77 6e 6c 6f 61 64  |sit the download|
00000030  73 20 70 61 67 65 0a 4f  66 66 69 63 69 61 6c 20  |s page.Official |
00000040  62 69 6e 61 72 79 20 64  69 73 74 72 69 62 75 74  |binary distribut|
00000050  69 6f 6e 73 20 61 72 65  20 61 76 61 69 6c 61 62  |ions are availab|
00000060  6c 65 20 66 6f 72 20 74  68 65 20 46 72 65 65 42  |le for the FreeB|
00000070  53 44 20 28 72 65 6c 65  61 73 65 20 31 30 2d 53  |SD (release 10-S|
00000080  54 41 42 4c 45 20 61 6e  64 20 61 62 6f 76 65 29  |TABLE and above)|
00000090  2c 20 4c 69 6e 75 78 2c  20 6d 61 63 4f 53 20 28  |, Linux, macOS (|
000000a0  31 30 2e 31 30 20 61 6e  64 20 61 62 6f 76 65 29  |10.10 and above)|
000000b0  2c 20 61 6e 64 20 57 69  6e 64 6f 77 73 20 6f 70  |, and Windows op|
000000c0  65 72 61 74 69 6e 67 20                           |erating |

2023/06/24 23:42:27 WRITER (No. 0): must be len() = 200
2023/06/24 23:42:27 WRITER (No. 1): wrotet len()=200
2023/06/24 23:42:27 READER (No. 4): put in channel-segment with len() 200, at offset 600 with data
00000000  6f 75 72 63 65 20 6f 72  20 69 6e 73 74 61 6c 6c  |ource or install|
00000010  69 6e 67 20 67 63 63 67  6f 20 69 6e 73 74 65 61  |ing gccgo instea|
00000020  64 20 6f 66 20 67 63 2e  0a 0a 53 79 73 74 65 6d  |d of gc...System|
00000030  20 72 65 71 75 69 72 65  6d 65 6e 74 73 0a 47 6f  | requirements.Go|
00000040  20 62 69 6e 61 72 79 20  64 69 73 74 72 69 62 75  | binary distribu|
00000050  74 69 6f 6e 73 20 61 72  65 20 61 76 61 69 6c 61  |tions are availa|
00000060  62 6c 65 20 66 6f 72 20  74 68 65 73 65 20 73 75  |ble for these su|
00000070  70 70 6f 72 74 65 64 20  6f 70 65 72 61 74 69 6e  |pported operatin|
00000080  67 20 73 79 73 74 65 6d  73 20 61 6e 64 20 61 72  |g systems and ar|
00000090  63 68 69 74 65 63 74 75  72 65 73 2e 20 50 6c 65  |chitectures. Ple|
000000a0  61 73 65 20 65 6e 73 75  72 65 20 79 6f 75 72 20  |ase ensure your |
000000b0  73 79 73 74 65 6d 20 6d  65 65 74 73 20 74 68 65  |system meets the|
000000c0  73 65 20 72 65 71 75 69                           |se requi|
.
2023/06/24 23:42:27 
2023/06/24 23:42:27 READER (No. 5): initial segment size = 200
2023/06/24 23:42:27 STEP READ n = 5
2023/06/24 23:42:27 STEP READ n = 5 - repairLimit = 1000
2023/06/24 23:42:27 STEP READ n = 5 - inisegmentSize = 200
2023/06/24 23:42:27 STEP READ n = 5 - prevSegmentsSizesSum = 800
2023/06/24 23:42:27 READER (No. 5): segment size = 200
2023/06/24 23:42:27 READER (No. 5): expected read = 200
2023/06/24 23:42:27 READER (No. 5): result read = 200
2023/06/24 23:42:27 WRITER (No. 0): wrotet len()=200
2023/06/24 23:42:27 READER (No. 5): put in channel-segment with len() 200, at offset 800 with data
00000000  72 65 6d 65 6e 74 73 20  62 65 66 6f 72 65 20 70  |rements before p|
00000010  72 6f 63 65 65 64 69 6e  67 2e 20 49 66 20 79 6f  |roceeding. If yo|
00000020  75 72 20 4f 53 20 6f 72  20 61 72 63 68 69 74 65  |ur OS or archite|
00000030  63 74 75 72 65 20 69 73  20 6e 6f 74 20 6f 6e 20  |cture is not on |
00000040  74 68 65 20 6c 69 73 74  2c 20 79 6f 75 20 6d 61  |the list, you ma|
00000050  79 20 62 65 20 61 62 6c  65 20 74 6f 20 69 6e 73  |y be able to ins|
00000060  74 61 6c 6c 20 66 72 6f  6d 20 73 6f 75 72 63 65  |tall from source|
00000070  20 6f 72 20 75 73 65 20  67 63 63 67 6f 20 69 6e  | or use gccgo in|
00000080  73 74 65 61 64 2e 0a 0a  4f 70 65 72 61 74 69 6e  |stead...Operatin|
00000090  67 20 73 79 73 74 65 6d  09 41 72 63 68 69 74 65  |g system.Archite|
000000a0  63 74 75 72 65 73 09 4e  6f 74 65 73 0a 46 72 65  |ctures.Notes.Fre|
000000b0  65 42 53 44 20 31 30 2e  33 20 6f 72 20 6c 61 74  |eBSD 10.3 or lat|
000000c0  65 72 09 61 6d 64 36 34                           |er.amd64|
.
2023/06/24 23:42:27 all writers Wait()
2023/06/24 23:42:27 WRITER (No. 1): offset = 800
2023/06/24 23:42:27 WRITER (No. 1): data
00000000  72 65 6d 65 6e 74 73 20  62 65 66 6f 72 65 20 70  |rements before p|
00000010  72 6f 63 65 65 64 69 6e  67 2e 20 49 66 20 79 6f  |roceeding. If yo|
00000020  75 72 20 4f 53 20 6f 72  20 61 72 63 68 69 74 65  |ur OS or archite|
00000030  63 74 75 72 65 20 69 73  20 6e 6f 74 20 6f 6e 20  |cture is not on |
00000040  74 68 65 20 6c 69 73 74  2c 20 79 6f 75 20 6d 61  |the list, you ma|
00000050  79 20 62 65 20 61 62 6c  65 20 74 6f 20 69 6e 73  |y be able to ins|
00000060  74 61 6c 6c 20 66 72 6f  6d 20 73 6f 75 72 63 65  |tall from source|
00000070  20 6f 72 20 75 73 65 20  67 63 63 67 6f 20 69 6e  | or use gccgo in|
00000080  73 74 65 61 64 2e 0a 0a  4f 70 65 72 61 74 69 6e  |stead...Operatin|
00000090  67 20 73 79 73 74 65 6d  09 41 72 63 68 69 74 65  |g system.Archite|
000000a0  63 74 75 72 65 73 09 4e  6f 74 65 73 0a 46 72 65  |ctures.Notes.Fre|
000000b0  65 42 53 44 20 31 30 2e  33 20 6f 72 20 6c 61 74  |eBSD 10.3 or lat|
000000c0  65 72 09 61 6d 64 36 34                           |er.amd64|

2023/06/24 23:42:27 WRITER (No. 1): must be len() = 200
2023/06/24 23:42:27 WRITER (No. 1): wrotet len()=200
2023/06/24 23:42:27 close(percents)
2023/06/24 23:42:27 percenter Wait()
```

</details>

### Тестирование

В рамках тестирования проводится проверка:

* соответствия результатов работы алгоритмов уже имеющимся в репозитории эталонам с задействованием в качестве функции эдентичности MD5-хеша от файла (`v := md5.New()`, ..., `v.Sum(nil)` пакета `crypto/md5`).
* времени работы алгоритмов с различными параметрами

#### Тестирование соответствия

```bash
go test -v ./
```

##### Соответствие результатов `Сopy()`

```text
=== RUN   TestCopy
OK. Результат соотвествует эталону: testdata/out_offset0_limit0_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset0_limit10_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset0_limit1000_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset0_limit10000_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset100_limit1000_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset6000_limit1000_test_copy.txt
--- PASS: TestCopy (0.15s)
```

##### Соответствие результатов `СopySegmented()`

```text
=== RUN   TestCopySegmented
OK. Результат соотвествует эталону: testdata/out_offset0_limit0_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset0_limit10_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset0_limit1000_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset0_limit10000_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset100_limit1000_test_copy.txt
OK. Результат соотвествует эталону: testdata/out_offset6000_limit1000_test_copy.txt
--- PASS: TestCopySegmented (0.18s)

=== RUN   TestCopySegmentedCustomParams
OK. Результат соотвествует эталону: testdata/out_offset0_limit10000_test_copy_segmented.0.txt
OK. Результат соотвествует эталону: testdata/out_offset0_limit10000_test_copy_segmented.1.txt
OK. Результат соотвествует эталону: testdata/out_offset0_limit10000_test_copy_segmented.2.txt
--- PASS: TestCopySegmentedCustomParams (0.16s)

=== RUN   TestCopySegmentedBigFile
Run [ID 1]: Побайтное копирование 50000 байт с 1 врайтером.
=== RUN   TestCopySegmentedBigFile/Побайтное_копирование_50000_байт_с_1_врайтером.
OK. Результат testdata/output.1.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
Run [ID 2]: Копирование 50000 байт с буфером 256-байт с 1 врайтером.
=== RUN   TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_256-байт_с_1_врайтером.
OK. Результат testdata/output.2.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
Run [ID 3]: Копирование 50000 байт с буфером 256-байт с 4 врайтерами.
=== RUN   TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_256-байт_с_4_врайтерами.
OK. Результат testdata/output.3.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
Run [ID 4]: Копирование 50000 байт с буфером 256-байт с 10 врайтерами.
=== RUN   TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_256-байт_с_10_врайтерами.
OK. Результат testdata/output.4.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
Run [ID 5]: Копирование 50000 байт с буфером 500-байт с 1 врайтерами.
=== RUN   TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_500-байт_с_1_врайтерами.
OK. Результат testdata/output.5.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
Run [ID 6]: Копирование 50000 байт с буфером 500-байт с 10 врайтером.
=== RUN   TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_500-байт_с_10_врайтером.
OK. Результат testdata/output.6.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
Run [ID 7]: Копирование 50000 байт с буфером 500-байт с 100 врайтерами.
=== RUN   TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_500-байт_с_100_врайтерами.
OK. Результат testdata/output.7.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
Run [ID 8]: Копирование 50000 байт с буфером 1000-байт с 5 врайтерами.
=== RUN   TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_1000-байт_с_5_врайтерами.
OK. Результат testdata/output.8.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
Run [ID 9]: Копирование 50000 байт с буфером 1000-байт с 50 врайтерами.
=== RUN   TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_1000-байт_с_50_врайтерами.
OK. Результат testdata/output.9.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
Run [ID 10]: Копирую 50000 байт в 1 врайтер :)
=== RUN   TestCopySegmentedBigFile/Копирую_50000_байт_в_1_врайтер_:)
OK. Результат testdata/output.10.txt соотвествует эталону: ./testdata/alice29.ethalon.txt
--- PASS: TestCopySegmentedBigFile (1.01s)
    --- PASS: TestCopySegmentedBigFile/Побайтное_копирование_50000_байт_с_1_врайтером. (0.91s)
    --- PASS: TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_256-байт_с_1_врайтером. (0.01s)
    --- PASS: TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_256-байт_с_4_врайтерами. (0.01s)
    --- PASS: TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_256-байт_с_10_врайтерами. (0.01s)
    --- PASS: TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_500-байт_с_1_врайтерами. (0.03s)
    --- PASS: TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_500-байт_с_10_врайтером. (0.01s)
    --- PASS: TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_500-байт_с_100_врайтерами. (0.01s)
    --- PASS: TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_1000-байт_с_5_врайтерами. (0.01s)
    --- PASS: TestCopySegmentedBigFile/Копирование_50000_байт_с_буфером_1000-байт_с_50_врайтерами. (0.01s)
    --- PASS: TestCopySegmentedBigFile/Копирую_50000_байт_в_1_врайтер_:) (0.01s)
PASS
```

##### Вывод о соответствии

Обе реализации верны.

#### Тестирование производительности

```bash
go test -bench=.
```

>
> Для повышения наглядности скорости работы реализованных подходов добавил в репозиторий большой файл (литературный текст из открытых источников).

##### Скорость работы `Сopy()`

```text
Run [ID 1]: Побайтное копирование 50000 байт [0:5000].
BenchmarkСopy/Побайтное_копирование_50000_байт_[0:5000].-4      1000000000  0.3671 ns/op
Run [ID 2]: Побайтное копирование 50000 байт [100:5100].
BenchmarkСopy/Побайтное_копирование_50000_байт_[100:5100].-4    1000000000  0.3867 ns/op
Run [ID 3]: Побайтное копирование 50000 байт [1000:6000].
BenchmarkСopy/Побайтное_копирование_50000_байт_[1000:6000].-4   1000000000  0.4789 ns/op
```

##### Скорость работы `СopySegmented()`

```text
Run [ID 1]: Побайтное копирование 50000 байт с 1 врайтером.
BenchmarkСopySegmented/Побайтное_копирование_50000_байт_с_1_врайтером.-4                1000000000  0.6298 ns/op
Run [ID 2]: Копирование 50000 байт с буфером 256-байт с 1 врайтером.
BenchmarkСopySegmented/Копирование_50000_байт_с_буфером_256-байт_с_1_врайтером.-4      1000000000  0.01251 ns/op
Run [ID 3]: Копирование 50000 байт с буфером 256-байт с 4 врайтерами.
BenchmarkСopySegmented/Копирование_50000_байт_с_буфером_256-байт_с_4_врайтерами.-4     1000000000  0.008512 ns/op
Run [ID 4]: Копирование 50000 байт с буфером 256-байт с 10 врайтерами.
BenchmarkСopySegmented/Копирование_50000_байт_с_буфером_256-байт_с_10_врайтерами.-4    1000000000  0.01186 ns/op
Run [ID 5]: Копирование 50000 байт с буфером 500-байт с 1 врайтерами.
BenchmarkСopySegmented/Копирование_50000_байт_с_буфером_500-байт_с_1_врайтерами.-4     1000000000  0.01098 ns/op
Run [ID 6]: Копирование 50000 байт с буфером 500-байт с 10 врайтером.
BenchmarkСopySegmented/Копирование_50000_байт_с_буфером_500-байт_с_10_врайтером.-4     1000000000  0.007875 ns/op
Run [ID 7]: Копирование 50000 байт с буфером 500-байт с 100 врайтерами.
BenchmarkСopySegmented/Копирование_50000_байт_с_буфером_500-байт_с_100_врайтерами.-4   1000000000  0.006355 ns/op
Run [ID 8]: Копирование 50000 байт с буфером 1000-байт с 5 врайтерами.
BenchmarkСopySegmented/Копирование_50000_байт_с_буфером_1000-байт_с_5_врайтерами.-4    1000000000  0.006963 ns/op
Run [ID 9]: Копирование 50000 байт с буфером 1000-байт с 50 врайтерами.
BenchmarkСopySegmented/Копирование_50000_байт_с_буфером_1000-байт_с_50_врайтерами.-4   1000000000  0.007719 ns/op
Run [ID 10]: Смотрите как быстро копирую 50000 байт в 1 врайтер :)
BenchmarkСopySegmented/Смотрите_как_быстро_копирую_50000_байт_в_1_врайтер_:)-4          1000000000  0.005542 ns/op
```

#### Выводы о производительности

* `СopySegmented()` (сегментированное чтение с параллельным процессом сегментированной записи) по скорости формирования результирующего файла производительнее, чем `Сopy()` (последовательное чтение и запись в одном потоке исполнения).
* Чем больше у `СopySegmented()` буфер копирования, тем копирование быстрее (строка отчета выше - "Смотрите как быстро копирую 50000 байт в 1 врайтер :)". Этот вариант быстрее всех остальных, хотя там индикация итоговая состоит из двух значений - 0% и сразу 100%.).
* В данной `СopySegmented()`-задаче узкое место - "запись в один файл", и увеличение числа одновременных горутин-врайтеров значительно не повышает скорость формирования итогового файла (так увеличение числа врайтеров с 10 на 100 существенно не сократило время работы - "Копирование 50000 байт с буфером 500-байт" с 10 врайтерами `0.007875 ns/op` или с 100 врайтерами `0.006355 ns/op` одного порядка).

## Общий вывод

Исследован механизм чтения и записи в файл, разработан модуль копирования файла в соответствии с условиями задачи.

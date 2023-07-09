# Домашние задания курса [OTUS «Разработчик Golang» 2023-03)](https://otus.ru/lessons/golang-professional/)

```text
"Go, Go, Johnny, Go, Go, Go!"

                  Marty McFly
       ("Back to the Future")
```

1) [«Hello, OTUS!»](./hw01_hello_otus)

   > [Исполнено](./hw01_hello_otus/README.md). В рамках решения задачи дополнительно исследовал объем результирующего бинарного файла, его зависимость от видов импорта и реализации алгоритма ([см. отчет](./hw01_hello_otus/QUESTION.md)).

2) [«Распаковка строки»](./hw02_unpack_string)

   > [Исполнено](./hw02_unpack_string/REPORT.md). Для решения задачи применил подход [порождающей грамматики](https://ru.wikipedia.org/wiki/%D0%9F%D0%BE%D1%80%D0%BE%D0%B6%D0%B4%D0%B0%D1%8E%D1%89%D0%B0%D1%8F_%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B0%D1%82%D0%B8%D0%BA%D0%B0), формализовав грамматические правила составления слов и предложений из них.

3) [«Частотный анализ»](./hw03_frequency_analysis)

   > [Исполнено](./hw03_frequency_analysis/README.md). Объем покрытия кода [тестами](./hw03_frequency_analysis/README.md#%D0%B4%D0%B5%D0%BC%D0%BE%D0%BD%D1%81%D1%82%D1%80%D0%B0%D1%86%D0%B8%D1%8F-%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%BE%D1%81%D0%BF%D0%BE%D1%81%D0%BE%D0%B1%D0%BD%D0%BE%D1%81%D1%82%D0%B8) составляет **97.2%**.

4) [«LRU-кэш»](./hw04_lru_cache)

   > [Исполнено](./hw04_lru_cache/REPORT.md). Для демонстрации O(1) реализовал метрики тестирования Benchmark-нагрузки и возможность сторонней обработки результатов Benchmark (см. [подробнее](./hw04_lru_cache/REPORT.md#benchmark-или-как-я-01-сложность-предъявлял)).

5) [«Параллельное исполнение»](./hw05_parallel_execution)

   > [Исполнено](./hw05_parallel_execution/REPORT.md). [Изначально](https://github.com/BorisPlus/OTUS-Go-2023-03/blob/master/hw05_parallel_execution/REPORT.md#%D0%B2%D1%81%D0%BF%D0%BE%D0%BC%D0%BE%D0%B3%D0%B0%D1%82%D0%B5%D0%BB%D1%8C%D0%BD%D0%BE%D0%B5) для демонстрации работы реализовал структуру учета статистики исполнения с инкапсулированием логики за счет `Mutex`. Она также принимает решение о критерии останова и вносит правки в количество воркеров, если задач априори меньше заданного числа горутин. Покрыта тестами. [Второй вариант](https://github.com/BorisPlus/OTUS-Go-2023-03/blob/master/hw05_parallel_execution/REPORT.md#%D1%81%D0%BE%D0%BA%D1%80%D0%B0%D1%89%D0%B5%D0%BD%D0%B8%D0%B5-%D0%BA%D0%BE%D0%B4%D0%B0) реализации появился в результате "борьба" между наглядностью и малым объемом кодовой базы и основан на `atomic`.

6) [«Пайплайн»](./hw06_pipeline_execution)

   > [Исполнено](./hw06_pipeline_execution/REPORT.md). Основаная идея в подготовке каналов Стейджей и их замыкание в одну "цепочку": выходной канал предыдущего Стейджа является входным каналом последующего Стейджа.

7) [«Утилита для копирования файлов»](./hw07_file_copying)

   > [Исполнено](./hw07_file_copying/REPORT.md). В рамках решения проведено исследование с сегментированным копированием файла в исходный (по типу торрент протокола).

8) [«Утилита envdir»](./hw08_envdir_tool)
9) [«Валидатор структур»](./hw09_struct_validator)

   > [Исполнено](./hw09_struct_validator/REPORT.md). Реализован параметризуемый валидатор структур с выводом стека из всех ошибок валидации.

10) [«Оптимизация программы»](./hw10_program_optimization)
11) [«Клиент TELNET»](./hw11_telnet_client)
12) [«Заготовка сервиса Календарь»](./hw12_13_14_15_calendar/docs/12_README.md)
13) [«API к Календарю»](./hw12_13_14_15_calendar/docs/13_README.md)
14) [«Кроликизация Календаря»](./hw12_13_14_15_calendar/docs/14_README.md)
15) [«Докеризация и интеграционное тестирование Календаря»](./hw12_13_14_15_calendar/docs/15_README.md)
16) [«Проект»](https://github.com/OtusGolang/final_project)

---
[Инструкция по сдаче ДЗ](https://github.com/OtusGolang/home_work/wiki#%D0%A1%D1%82%D1%83%D0%B4%D0%B5%D0%BD%D1%82%D0%B0%D0%BC).

---
Используемая версия [golangci-lint](https://golangci-lint.run/usage/install/#other-ci): __v1.50.1__

```shell
$ golangci-lint version
golangci-lint has version 1.50.1 built from 8926a95 on 2022-10-22T10:48:48Z
```

---
Авторы ДЗ:

* [Дмитрий Смаль](https://github.com/mialinx)
* [Антон Телышев](https://github.com/Antonboom)

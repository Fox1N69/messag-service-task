# Сервис приема и обработки сообщений

В данном проекте я реализовал сервис, которые принимает сообщения и отправляет их в kafka для дальнейшей обработки.
После обработки статус сообщений меняется в базе данных.

#### Документация к api  - {BASE_URL}/swagger/index.html#/

## Запуск сервиса
Перед запуском приложения поменяйте конфигурационные данные в файлe который находится по пути [./config/config.json](./config/config.json)

**Сборка docker контрейнера**

```console
make docker-build
```

**Запуск docker контрейнера**

```console
make docker-run
```

**Запуск миграции зависимостей**

```console
make dep
```

**Запуск проекта локально**

```console
make run-algosync
```

**Запуск тестов**

```console
make test
```

**Сборка проекта**

```console
make build-algosync
```

**Запуск с hot reload**

Переменуйте example.air.toml в air.tomal

```console
air
```

**Запус hey для тестирования rps и время отклика сервера**

```console
./hey.sh
```

## Используемые библиотеки

| Library    | Usage             |
| ---------- | ----------------- |
| fiber      | FastHTTP framework|
| sqlc       | SQL library       |
| postgres   | Database          |
| logrus     | Logger library    |
| viper      | Config library    |
| redis      | Cache storage     |


## Оптимизация
Для оптимизации сервиса было принято использовать следующие технологии:
Prefork - технологолия при которой создается несколько копий сервиса для многопроцессорной обработки.
SQLC - библиотека для генерации кода на go из sqlc запросов. Как показало тестирование, при большом количестве запросов к бд, sqlc обходит по скорости аналоги.

## Защиты

### Защита от DoS и DDoS
В конфигурации config.json есть возможность включить "ip_block_config" поставить параметер enable на true.
Это моя реализация защиты от dos и ddos атак. Когда запрос с одного ip адресса будет привышать параметер "max_requests" ip будет добавлятся в кэш redis и блокироватся, блокировака продлится 10 минут.
Даная реализация необходима на случай если хакер будет атаковать сервис напрямую, обойдя балансер нагрузки.

### Защита от SQL иньекций
Что бы защитить свой сервис от sql иньекций, я использую заготовленый sql запросы.
И не использую денамические sql запросы.

## Benchmark
Summary:
  Total:        5.6101 secs
  Slowest:      0.0404 secs
  Fastest:      0.0003 secs
  Average:      0.0056 secs
  Requests/sec: 17824.9432

  Total data:   4900000 bytes
  Size/request: 49 bytes

Response time histogram:
  0.000 [1]     |
  0.004 [12796] |■■■■■■
  0.008 [82458] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.012 [3725]  |■■
  0.016 [628]   |
  0.020 [226]   |
  0.024 [69]    |
  0.028 [29]    |
  0.032 [14]    |
  0.036 [25]    |
  0.040 [29]    |


Latency distribution:
  10% in 0.0043 secs
  25% in 0.0046 secs
  50% in 0.0054 secs
  75% in 0.0060 secs
  90% in 0.0067 secs
  95% in 0.0082 secs
  99% in 0.0126 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0003 secs, 0.0404 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0046 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0008 secs
  resp wait:    0.0056 secs, 0.0003 secs, 0.0330 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0005 secs

Status code distribution:
  [201] 100000 responses
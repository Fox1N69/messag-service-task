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
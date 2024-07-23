# Сервис приема и обработки сообщений


В данном проекте я реализовал сервис, которые принимает сообщения и отправляет их в kafka для дальнейшей обработки.
После обработки статус сообщений меняется в базе данных.



> [!WARNING]  
>
> Перед проведением нагрузочного тестирования убедитесь, что параметр "enable" в [config файле](./config/config.json) установлен в положение "false". 
>
> В противном случае, ваш IP-адрес будет заблокирован и предоставление доступа к сервису будет приостановлено на 10 минут в соответствии с механизмом защиты от DDoS-атак.



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

## API ручки
CreateMessage()

Request
```json
  
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



## Защитa

### Защита от DoS и DDoS

В конфигурационном файле [config.json](./config/config.json) имеется возможность включить блокировку IP-адресов ("ip_block_config") установив параметр "enable" в значение true. 
Это моя реализация защиты от атак DoS и DDoS. Если запросы с одного IP-адреса превышают параметр "max_requests", IP-адрес будет добавлен в кэш Redis и заблокирован на 10 минут. Эта реализация необходима для предотвращения напрямую направленных атак на сервис, обходящих балансировщик нагрузки.

### Защита от SQL иньекций

Для защиты сервиса от SQL инъекций, я использую предварительно подготовленные SQL запросы и избегаю динамических SQL запросов.



# Результаты тестирования нагрузки

Тестирование проводилось локально на **mac mini m1**
На сервер было отправлено 100k post запросов по 100 за раз.

**Скрипт теста** 
 ```console
hey -n 100000 -c 100 -m POST -d '{"message":"Lorem ipsum dolor sit amet, consectetur adipiscing elit", "status_id": 1}' -H "Content-Type: application/json" http://localhost:4000/api/message/

```

### Основные Показатели

- **Общее время теста**: 5.2101 секунд
- **Самый медленный запрос**: 0.0404 секунд
- **Самый быстрый запрос**: 0.0001 секунд
- **Среднее время запроса**: 0.0056 секунд
- **Запросов в секунду**:  19824.9432

### Размеры данных

- **Всего данных**: 4,700,000 байт
- **Размер запроса**: 47 байт

### Распределение задержек

- **10% запросов**: до 0.0043 секунд
- **25% запросов**: до 0.0046 секунд
- **50% запросов**: до 0.0052 секунд (медиана)
- **75% запросов**: до 0.0060 секунд
- **90% запросов**: до 0.0067 секунд
- **95% запросов**: до 0.0074 секунд
- **99% запросов**: до 0.0103 секунд

### Распределение кодов состояния

- **201**: 100,000 ответов
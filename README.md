# Сервис синхронизации пользовательских алгоритмов

В данном проекте я реализовал микросервисный api для синхронизации пользовательских алгоритмов.
У сервиса есть обработчик который раз в 5 минунт смотрит статусы алгоритмов и если алгоритм включен, то создается соответствующий pod, если алгоритм выключен, то pod удаляется

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
| gin        |  Web framework    |
| database/sql | SQL library       |
| postgres   | Database          |
| logrus     | Logger library    |
| viper      | Config library    |

## Комментарий
В проекте подразумевалось использованние Redis для кэширования частых запросов, но из-за нехватки времяни, было принято решение не использовать его.

Он прописан в config файле так как был подключен, но при этом он не инициализирован, и по тому нет необходимости устанавливать его для запуска сервиса!!!

# DNS Manager (dnsmgr)
## Описание

Клиент-серверное приложение на базе **gRPC** для удаленного управления конфигурацией DNS-серверов (файл `/etc/resolv.conf`).

## Возможности
1. На сервере:
- Добавить DNS-сервер
- Удалить DNS-сервер
- Получить список всех DNS-серверов

2. CLI клиент позволяет выполнять функции на удалённом компьютере. Используйте `-h` или `-help` для справки.

## Установка
1. Склонируйте репозиторий:
```bash
git clone https://github.com/IliaSotnikov2005/dnsmgr.git
cd dnsmgr
```
2. Установите зависимости:
```bash
go mod download
```

## Сборка бинарных файлов
```bash
mkdir bin
go build -o bin/dnsmgr ./client/cmd/main.go
go build -o bin/dnsmgr-server ./server/cmd/main.go
```

## Конфигурация сервера
Сервер использует файл конфигурации в формате YAML. Для работы нужны права на чтение/запись файла /etc/resolv.conf (или иного пути, указанного в конфигурации).

**Пример файла конфигурации**
```yaml
storage_path: "/etc/resolv.conf"  # Путь к файлу DNS-серверов
grpc:
  port: "50051"                   # Порт, который будет слушать сервер
  timeout_seconds: 5              # Таймаут обработки запроса
```

## Запуск
**1. Сервер**
```bash
go run server/main.go --config="path/to/config.yaml"
```

**2. Клиент**

Клиент предоставляет интерфейс командной строки для взаимодействия с сервером.
```bash
# Добавить DNS
go run client/main.go -add 8.8.8.8

# Удалить DNS
go run client/main.go -rm 1.1.1.1

# Список всех DNS
go run client/main.go -ls

# Справка по всем флагам
go run client/main.go -h
```

## Тестирование
```bash
go test -v ./...
```

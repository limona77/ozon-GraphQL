# Тестовое Задание

1. Склонируй репозиторий

```bash 
   git clone https://github.com/limona77/ozon-GraphQL
```
2. настрой env файл и установи все зависимости

```bash 
   go mod tidy
```

3. Запусти докер композ

```bash 
   docker-compose up
```

4. Создай папку storage в корне проект и запусти миграции
```bash 
   go run ./cmd/migrator --storage-path=./storage/postgres.db --migrations-path=./migrations --action=down
```


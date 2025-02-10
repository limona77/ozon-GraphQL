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

4. Запусти миграции
```bash 
   go run ./cmd/migrator  --migrations-path=./migrations --action=up
```


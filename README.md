# lab4KR

## Используемые инструменты

- Go.
- GoLand.
- PostgreSQL.
- pgx/v5 — драйвер для работы с PostgreSQL.
- golang-migrate — инструмент управления миграциями.
- Git и GitHub.

- ## Реализация

- `internal/model` — модели User, DefectReason, Defect.
- `internal/storage/postgres` — подключение к БД и DefectRepository с методами List, Get, Create, Update, Delete.
- `migrations` — создание и удаление таблиц и SQL-функций.
- `cmd/app` — проверка подключения и получения списка брака.

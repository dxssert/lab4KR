# lab4KR

REST API для ведения журнала производственного брака. 

Контроллер и схемы запросов и ответов находятся в `internal/httpapi`. Проверка данных выполняется в `internal/service`. Репозиторий `internal/storage/postgres` работает с PostgreSQL через pgx/v5.

## Запуск

Для запуска требуются Go, PostgreSQL и `golang-migrate` с поддержкой PostgreSQL. Создайте базу `lab4kr` и выполните команды в PowerShell из корня проекта. Замените `YOUR_PASSWORD` паролем пользователя PostgreSQL. Если миграция уже применена, команду `migrate` можно пропустить.

```powershell
$env:DATABASE_URL = 'postgres://postgres:YOUR_PASSWORD@127.0.0.1:5432/lab4kr?sslmode=disable'
migrate -path ./migrations -database $env:DATABASE_URL up
go run ./cmd/app
```

Сервер доступен по адресу `http://127.0.0.1:9095`. Адрес задаётся переменной `HTTP_ADDR`.

## Операции API

| Метод и путь | Операция | Код ответа |
| --- | --- | --- |
| `GET /api/v1/defects` | Список записей | `200 OK` |
| `GET /api/v1/defects/{id}` | Получение записи | `200 OK` |
| `POST /api/v1/defects` | Создание записи | `201 Created` |
| `PUT /api/v1/defects/{id}` | Обновление записи | `200 OK` |
| `DELETE /api/v1/defects/{id}` | Удаление записи | `204 No Content` |

## Проверка в Postman

Выполните `scripts/demo.sql` в базе `lab4kr`. Подставьте идентификаторы из результата запроса в поля `production_output_id`, `reason_id` и `created_by`.

Отправьте `POST http://127.0.0.1:9095/api/v1/defects` с телом в формате JSON (**Body → raw → JSON**).

```json
{
  "production_output_id": 1,
  "reason_id": 2,
  "defect_grams": 500,
  "comment": "Подгоревшая корка",
  "created_by": 1
}
```

После ответа `201 Created` подставьте возвращённый `id` в запрос `GET http://127.0.0.1:9095/api/v1/defects/{id}`. Сервер вернёт `200 OK` и сохранённую запись.

Тело запроса `PUT` содержит те же поля, кроме `created_by`. При обновлении автор записи сохраняется.

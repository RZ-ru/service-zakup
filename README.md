# service-zakup

Проект реализует систему задач с авторизацией, проверкой прав доступа, кэшированием permissions, асинхронными уведомлениями через RabbitMQ и синхронным межсервисным взаимодействием через gRPC.

## Стек

* Go
* Gin
* PostgreSQL
* Redis
* RabbitMQ
* gRPC / Protobuf
* Docker Compose
* JWT
* Transactional Outbox Pattern

## Сервисы

```text
auth-service
  Выдаёт JWT-токены.

user-service
  Создаёт пользователей и отдаёт данные пользователя по id.

task-service
  Основной сервис задач.
  Создаёт, читает, изменяет и удаляет задачи.
  Пишет события в outbox_events.
  Публикует события в RabbitMQ через outbox worker.
  Проверяет права доступа через permission-service по gRPC.

permission-service
  Создаёт и проверяет права доступа к задачам.
  Использует PostgreSQL как источник истины.
  Использует Redis как cache-aside для проверки permissions.
  Предоставляет gRPC API для task-service.

notification-service
  Слушает события задач из RabbitMQ.
  Получает email владельца задачи через user-service.
  Отправляет уведомление на email пользователя.
```

## Архитектура

```text
                  ┌──────────────┐
                  │ auth-service │
                  │   JWT token  │
                  └──────┬───────┘
                         │
                         ▼
┌──────────────┐   ┌──────────────┐        gRPC         ┌────────────────────┐
│    client    │──▶│ task-service │────────────────────▶│ permission-service │
└──────────────┘   └──────┬───────┘                     └─────────┬──────────┘
                          │                                       │
                          │                                       ▼
                          │                              ┌─────────────────┐
                          │                              │ Redis + Postgres │
                          │                              └─────────────────┘
                          │
                          ▼
                  ┌──────────────┐
                  │ task_db      │
                  │ tasks        │
                  │ outbox_events│
                  └──────┬───────┘
                         │
                         ▼
                  ┌──────────────┐
                  │ outbox worker│
                  └──────┬───────┘
                         │
                         ▼
                  ┌──────────────┐
                  │ RabbitMQ     │
                  │ task.events  │
                  └──────┬───────┘
                         │
                         ▼
              ┌──────────────────────┐
              │ notification-service │
              └──────────┬───────────┘
                         │ HTTP
                         ▼
                  ┌──────────────┐
                  │ user-service │
                  └──────────────┘
```

## Основной сценарий

### Создание задачи

```text
POST /tasks
↓
task-service валидирует JWT
↓
task-service создаёт задачу и outbox event в одной транзакции
↓
task-service вызывает permission-service по gRPC для создания owner permission
↓
outbox worker публикует task.created в RabbitMQ
↓
notification-service получает событие
↓
notification-service получает email пользователя из user-service
↓
notification-service логирует fake email
```

### Изменение задачи

```text
PATCH /tasks/:id
↓
task-service проверяет доступ через permission-service по gRPC
↓
task-service обновляет задачу и пишет task.updated в outbox_events
↓
outbox worker публикует событие в RabbitMQ
↓
notification-service отправляет fake email владельцу задачи
```

### Удаление задачи

```text
DELETE /tasks/:id
↓
task-service проверяет доступ через permission-service по gRPC
↓
task-service удаляет задачу и пишет task.deleted в outbox_events
↓
outbox worker публикует событие в RabbitMQ
↓
notification-service отправляет fake email владельцу задачи
```

## Авторизация

`auth-service` выдаёт JWT с полями:

```json
{
  "user_id": "uuid",
  "role": "user|admin",
  "exp": 1234567890
}
```

JWT передаётся во внешних HTTP-запросах:

```text
Authorization: Bearer <token>
```

При gRPC-вызовах `task-service → permission-service` токен передаётся через metadata:

```text
authorization: Bearer <token>
```

`permission-service` сам валидирует JWT и достаёт `user_id`. `user_id` не передаётся между сервисами отдельным параметром.

## Роли и permissions

В JWT используется системная роль пользователя:

```text
user
admin
```

В `permission-service` хранится роль пользователя относительно задачи:

```text
owner
```

Правила доступа:

```text
admin
  имеет доступ ко всем задачам
  обрабатывается в task-service

owner
  имеет доступ только к своей задаче
  проверяется через permission-service
```

## Redis

Redis используется в `permission-service` для кэширования проверки доступа.

Ключ:

```text
perm:<user_id>:<task_id>
```

Значение:

```text
true | false
```

Источник истины — PostgreSQL. Redis используется только как кэш.

## RabbitMQ и Outbox

Для надёжной доставки событий используется Transactional Outbox Pattern.

При изменении задачи `task-service` в одной транзакции пишет:

```text
1. основную запись в tasks
2. событие в outbox_events
```

Outbox worker читает события со статусом `pending`, публикует их в RabbitMQ и после успешной публикации помечает как `processed`.

Статусы outbox:

```text
pending    событие создано, но ещё не опубликовано
processed  событие успешно опубликовано
failed     событие не удалось опубликовать после нескольких попыток
```

Exchange:

```text
task.events
```

Тип exchange:

```text
topic
```

События:

```text
task.created
task.updated
task.deleted
```

Очередь notification-service:

```text
notification.task.events
```

## gRPC

`task-service` использует gRPC для синхронного взаимодействия с `permission-service`.

Proto-контракт:

```proto
syntax = "proto3";

package permission;

service PermissionService {
  rpc CreatePermission(CreatePermissionRequest) returns (CreatePermissionResponse);
  rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
}

message CreatePermissionRequest {
  string task_id = 1;
}

message CreatePermissionResponse {
  bool created = 1;
}

message CheckPermissionRequest {
  string task_id = 1;
}

message CheckPermissionResponse {
  bool allowed = 1;
}
```

`user_id` не передаётся в gRPC request. Он достаётся из JWT на стороне `permission-service`.

## Запуск

```bash
docker compose up --build -d
```

Проверить контейнеры:

```bash
docker compose ps
```

Логи всех сервисов:

```bash
docker compose logs -f
```

Логи конкретного сервиса:

```bash
docker compose logs -f task-service
```

Пересобрать только один сервис:

```bash
docker compose up --build -d task-service
```

## Основные URL

```text
auth-service          http://localhost:8083
user-service          http://localhost:8081
task-service          http://localhost:8080
permission-service    http://localhost:8082
RabbitMQ UI           http://localhost:15672
```

RabbitMQ UI:

```text
login:    guest
password: guest
```

## API

### Создать пользователя

Перед получением JWT сначала нужно создать пользователя через `user-service`. `auth-service` выдаёт токен по уже существующему `user_id`.

```http
POST /users
```

```json
{
  "email": "tom@mail.com",
  "name": "Tom"
}
```

В ответе нужно сохранить поле `id` — оно используется как `user_id` при получении JWT.

Пример:

```json
{
  "id": "13837266-38c2-40b4-8374-024f15680f9f",
  "email": "tom@mail.com",
  "name": "Tom"
}
```

### Получить пользователя

```http
GET /users/:id
```

### Получить JWT

После создания пользователя нужно передать его `id` в `auth-service`.

```http
POST /login
```

```json
{
  "user_id": "<id из ответа POST /users>"
}
```

Ответ:

```json
{
  "token": "<jwt>"
}
```

### Создать задачу

```http
POST /tasks
Authorization: Bearer <token>
```

```json
{
  "title": "buy milk",
  "description": "go to shop"
}
```

### Получить задачу

```http
GET /tasks/:id
Authorization: Bearer <token>
```

### Изменить задачу

```http
PATCH /tasks/:id
Authorization: Bearer <token>
```

```json
{
  "title": "Updated task title",
  "description": "Updated task description",
  "status": "done"
}
```

Можно передавать только изменяемые поля.

### Удалить задачу

```http
DELETE /tasks/:id
Authorization: Bearer <token>
```

## Проверка outbox

Зайти в PostgreSQL:

```bash
docker exec -it postgres-db psql -U postgres -d task_db
```

Посмотреть последние события:

```sql
SELECT
  event_type,
  routing_key,
  status,
  attempts,
  payload,
  created_at,
  processed_at
FROM outbox_events
ORDER BY created_at DESC
LIMIT 10;
```

Ожидаемый результат:

```text
task.created | task.created | processed
task.updated | task.updated | processed
task.deleted | task.deleted | processed
```

## Проверка notification-service

```bash
docker compose logs -f notification-service
```

Примеры логов:

```text
EMAIL TO=tom@mail.com SUBJECT=Создана новая задача BODY=Создана задача: buy milk
EMAIL TO=tom@mail.com SUBJECT=Задача изменена BODY=Изменена задача: Updated task title
EMAIL TO=tom@mail.com SUBJECT=Задача удалена BODY=Удалена задача: Updated task title
```

## Проверка RabbitMQ

Открыть UI:

```text
http://localhost:15672
```

Проверить:

```text
Exchanges → task.events
Queues → notification.task.events
```

## Генерация gRPC-кода

Для `permission-service`:

```powershell
cd "D:\Go projects\service-zakup\services\permission-service"

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/permission.proto
```

Для `task-service`:

```powershell
cd "D:\Go projects\service-zakup\services\task-service"

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/permission.proto
```

## Что уже реализовано

* Микросервисная структура
* JWT-аутентификация
* Проверка доступа через permission-service
* Redis cache-aside для permissions
* gRPC между task-service и permission-service
* PostgreSQL для сервисов
* RabbitMQ для асинхронных событий
* Transactional Outbox Pattern
* Notification-service consumer
* Fake email sender
* Docker Compose

## Возможные следующие улучшения

* Graceful shutdown для HTTP/gRPC/worker/consumer
* Request ID / correlation ID в outbox events и RabbitMQ headers
* Notification history в отдельной БД
* Идемпотентность notification-service через event_id
* Dead Letter Queue для RabbitMQ
* Mailpit/SMTP вместо fake email sender
* Approval-service для согласования задач
* Вынос proto-контрактов в общий shared-модуль
* Healthcheck endpoints для сервисов

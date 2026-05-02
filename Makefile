# DB URLs
USER_DB=postgres://postgres:postgres@localhost:5432/user_db?sslmode=disable
TASK_DB=postgres://postgres:postgres@localhost:5432/task_db?sslmode=disable

# USER SERVICE

migrate-user-up:
	migrate -path services/user-service/migrations -database "$(USER_DB)" up

migrate-user-down:
	migrate -path services/user-service/migrations -database "$(USER_DB)" down

# TASK SERVICE

migrate-task-up:
	migrate -path services/task-service/migrations -database "$(TASK_DB)" up

migrate-task-down:
	migrate -path services/task-service/migrations -database "$(TASK_DB)" down
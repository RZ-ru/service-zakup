CREATE TABLE permissions (
    user_id TEXT NOT NULL,
    task_id TEXT NOT NULL,
    role TEXT NOT NULL,
    PRIMARY KEY (user_id, task_id)
);
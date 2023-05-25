DROP TABLE IF EXISTS tasks_labels, tasks, labels, users;

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE labels(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE tasks(
    id BIGSERIAL PRIMARY KEY,
    opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
    closed BIGINT DEFAULT 0,
    author_id BIGINT REFERENCES users(id) DEFAULT 0,
    assigned_id BIGINT REFERENCES users(id) DEFAULT 0,
    title TEXT NOT NULL,
    content TEXT
);

CREATE TABLE tasks_labels(
    task_id BIGINT REFERENCES tasks(id),
    label_id BIGINT REFERENCES labels(id)
);
CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users (
    user_id            SERIAL                  PRIMARY KEY,
    user_version       BIGINT        NOT NULL  DEFAULT 1,
    user_full_name     VARCHAR(100)  NOT NULL  CHECK (char_length(user_full_name) BETWEEN 3 AND 100),
    user_phone_number  VARCHAR(15)             CHECK (
        user_phone_number ~ '^\+[0-9]+$'
        AND
        char_length(user_phone_number) BETWEEN 10 AND 15
    )
);

CREATE TABLE todoapp.tasks (
    task_id           SERIAL                    PRIMARY KEY,
    task_version      BIGINT         NOT NULL   DEFAULT 1,
    task_title        VARCHAR(100)   NOT NULL   CHECK (char_length(task_title) BETWEEN 1 AND 100),
    task_description  VARCHAR(1000)             CHECK (char_length(task_description) BETWEEN 1 AND 1000), 
    task_completed    BOOLEAN        NOT NULL,
    task_created_at   TIMESTAMPTZ    NOT NULL,
    task_completed_at TIMESTAMPTZ,

    CHECK (
        (task_completed=FALSE AND task_completed_at IS NULL)
        OR
        (task_completed=TRUE AND task_completed_at IS NOT NULL AND task_completed_at >= task_created_at)
    ),

    author_user_id    INTEGER        NOT NULL  REFERENCES todoapp.users(user_id)
);
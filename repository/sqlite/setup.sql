-- users TABLE --
CREATE TABLE IF NOT EXISTS "user" (
    "id" INTEGER NOT NULL UNIQUE,
    "username" TEXT NOT NULL UNIQUE,
    "email" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    PRIMARY KEY("id" AUTOINCREMENT)
);
-- post TABLE --
CREATE TABLE IF NOT EXISTS "post" (
    "id" INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER NOT NULL,
    "author" TEXT NOT NULL,
    "title" TEXT NOT NULL UNIQUE,
    "content" TEXT NOT NULL,
    "category" TEXT NOT NULL,
    "date" TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id)
);
-- comment TABLE --
CREATE TABLE IF NOT EXISTS "comment" (
    "id" INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER NOT NULL,
    "post_id" INTEGER NOT NULL,
    "author" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "date" TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (post_id) REFERENCES post(id)
)
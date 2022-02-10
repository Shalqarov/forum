CREATE TABLE IF NOT EXISTS "User" (
    "ID" INTEGER NOT NULL UNIQUE,
    "Login" TEXT NOT NULL UNIQUE,
    "Email" TEXT NOT NULL UNIQUE,
    PRIMARY KEY("ID" AUTOINCREMENT)
);
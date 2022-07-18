CREATE TABLE IF NOT EXISTS "user" (
    "user_id" BIGSERIAL PRIMARY KEY,
    "email" VARCHAR(320) NOT NULL UNIQUE,
    "username" VARCHAR(20) NOT NULL,
    "password" TEXT NOT NULL,
    "avatar" TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS "post" (
    "post_id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT NOT NULL REFERENCES "user"("user_id"),
    "title" VARCHAR(255) NOT NULL,
    "content" VARCHAR(2000) NOT NULL,
    "category" TEXT NOT NULL,
    "date" DATE NOT NULL,
    "image" TEXT
);
CREATE TABLE IF NOT EXISTS "comment" (
    "comment_id" BIGSERIAL PRIMARY KEY,
    "user_id" BIGINT NOT NULL REFERENCES "user"("user_id"),
    "post_id" BIGINT NOT NULL REFERENCES "post"("post_id"),
    "content" VARCHAR(255) NOT NULL,
    "created_at" DATE NOT NULL
);
CREATE TABLE IF NOT EXISTS "post_vote" (
    "post_vote_id" BIGSERIAL NOT NULL PRIMARY KEY,
    "user_id" BIGINT NOT NULL REFERENCES "user"("user_id"),
    "post_id" BIGINT NOT NULL REFERENCES "post"("post_id"),
    "vote" INT NOT NULL
);
CREATE TABLE IF NOT EXISTS "comment_vote" (
    "comment_vote_id" BIGSERIAL NOT NULL PRIMARY KEY,
    "user_id" BIGINT NOT NULL REFERENCES "user"("user_id"),
    "comment_id" BIGINT NOT NULL REFERENCES "comment"("comment_id"),
    "vote" INT NOT NULL
);
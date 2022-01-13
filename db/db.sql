CREATE extension IF NOT EXISTS CITEXT;

CREATE TABLE "user" (
    id serial not null UNIQUE,
    nickname TEXT UNIQUE,
    fullname TEXT not null,
    about TEXT,
    email TEXT not null UNIQUE
);

CREATE INDEX ON "user" (nickname, email);

CREATE TABLE "forum" (
    id serial not null UNIQUE,
    title varchar(255) not NULL ,
    user_id int REFERENCES "user" (id) on DELETE CASCADE not NULL,
    slug TEXT not NULL UNIQUE,
    posts int not null DEFAULT 0,
    threads int not null DEFAULT 0
);

CREATE INDEX ON "forum" (slug);

CREATE TABLE "thread" (
    id serial not NULL UNIQUE,
    title TEXT not NULL,
    user_id int REFERENCES "user" (id) on DELETE CASCADE not NULL,
    forum_id int REFERENCES "forum" (id) on DELETE CASCADE not NULL,
    message TEXT not NULL,
    votes int DEFAULT 0,
    slug CITEXT UNIQUE,
    created TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX ON "thread" (slug);

CREATE Table "post" (
    id serial not Null UNIQUE,
    parent integer,
    user_id int REFERENCES "user" (id) on DELETE CASCADE not null,
    message TEXT not NULL,
    thread_id int REFERENCES "thread" (id) on DELETE CASCADE not null,
    edited boolean DEFAULT false,
    created TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX ON post(id);

CREATE TABLE "vote" (
    id serial not NULL UNIQUE,
    thread_id int REFERENCES "thread" (id) on DELETE CASCADE not NULL,
    user_id int REFERENCES "user" (id) on DELETE CASCADE NOT NULL,
    voice SMALLINT 
);

CREATE OR REPLACE FUNCTION  update_vote() RETURNS TRIGGER AS $update_vote$
BEGIN
    IF OLD.voice = NEW.voice
    THEN
        RETURN NULL;
    END IF;
    UPDATE "thread"
    SET
        votes = votes + NEW.voice * 2
    WHERE id = NEW.thread_id;
    RETURN NULL;
END;
$update_vote$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_vote ON "vote";
CREATE TRIGGER update_vote AFTER UPDATE ON "vote" FOR EACH ROW EXECUTE PROCEDURE update_vote();

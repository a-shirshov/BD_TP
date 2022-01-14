CREATE extension IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS "forum_user";
DROP TABLE IF EXISTS "vote";
DROP TABLE IF EXISTS "post";
DROP TABLE IF EXISTS "thread";
DROP TABLE IF EXISTS "forum";
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
    id serial not null UNIQUE,
    nickname CITEXT UNIQUE not null,
    fullname TEXT not null,
    about TEXT,
    email CITEXT not null UNIQUE
);

CREATE TABLE "forum" (
    id serial not null UNIQUE,
    title varchar(255) not NULL ,
    user_id int REFERENCES "user" (id) on DELETE CASCADE not NULL,
    slug CITEXT not NULL UNIQUE,
    posts int not null DEFAULT 0,
    threads int not null DEFAULT 0
);


CREATE TABLE "thread" (
    id serial not NULL UNIQUE,
    title TEXT not NULL,
    user_id int REFERENCES "user" (id) on DELETE CASCADE not NULL,
    forum_id int REFERENCES "forum" (id) on DELETE CASCADE not NULL,
    message TEXT not NULL,
    votes int DEFAULT 0,
    slug CITEXT,
    created TIMESTAMP WITH TIME ZONE DEFAULT now()
);



CREATE Table "post" (
    id serial not Null UNIQUE,
    parent integer,
    user_id int REFERENCES "user" (id) on DELETE CASCADE not null,
    message TEXT not NULL,
    thread_id int REFERENCES "thread" (id) on DELETE CASCADE not null,
    forum citext references "forum"(slug) on delete cascade not null,
    edited boolean DEFAULT false,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    path int[]
);



CREATE TABLE "vote" (
    id serial not NULL UNIQUE,
    thread_id int REFERENCES "thread" (id) on DELETE CASCADE not NULL,
    user_id int REFERENCES "user" (id) on DELETE CASCADE NOT NULL,
    voice SMALLINT,
    UNIQUE(thread_id,user_id)
);

CREATE TABLE "forum_user"(
    id          serial primary key,
    forum       citext references "forum"(slug) on delete cascade not null,
    "user"      int references "user" (id) on delete cascade not null,
    UNIQUE (forum, "user")
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

CREATE or REPLACE FUNCTION update_thread_votes() returns TRIGGER as $update_thread_votes$
begin
    update thread set votes = (votes + new.voice) where id = new.thread_id;
    return null;
end;
$update_thread_votes$ language plpgsql;

DROP TRIGGER IF EXISTS update_thread_vote ON "vote";
create TRIGGER update_thread_vote after insert on vote for each row execute procedure update_thread_votes();

create or replace function create_post_with_path() returns trigger as $create_post_with_path$
declare
    parent_path int[];
begin
    update forum set posts = posts + 1 where slug = new.forum;
    insert into forum_user (forum, "user") values (new.forum, new.user_id)
    on conflict do nothing;
    parent_path = (select path from post where id = new.parent limit 1);
    new.path = parent_path || NEW.id;
    return NULL;
end;
$create_post_with_path$ language plpgsql;

DROP TRIGGER IF EXISTS create_post_with_path ON "post";
create trigger create_post_with_path before insert on post for each row execute procedure create_post_with_path();

CREATE INDEX ON "user" (nickname, email);
CREATE INDEX ON "forum" (slug);
CREATE INDEX ON "thread" (slug);
CREATE INDEX ON "post"(id);

SELECT * FROM pg_indexes WHERE tablename = 'user';
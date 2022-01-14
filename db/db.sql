CREATE extension IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS "forum_user";
DROP TABLE IF EXISTS "vote";
DROP TABLE IF EXISTS "post";
DROP TABLE IF EXISTS "thread";
DROP TABLE IF EXISTS "forum";
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
    id serial not null UNIQUE,
    nickname CITEXT collate "C" UNIQUE not null,
    fullname TEXT not null,
    about TEXT,
    email CITEXT collate "C" not null UNIQUE
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
    edited boolean DEFAULT false,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum_id int REFERENCES "forum" (id) on Delete Cascade not null,
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
    id          serial not null unique,
    forum_id    int references "forum" (id) on delete cascade not null,
    user_id     int references "user" (id) on delete cascade not null,
    UNIQUE (forum_id, user_id)
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
    update "forum" set posts = posts + 1 where id = new.forum_id;
    insert into "forum_user" (forum_id, user_id) values (new.forum_id, new.user_id)
    on conflict do nothing;
    parent_path = (select path from post where id = new.parent limit 1);
    new.path = parent_path || NEW.id;
    return new;
end;
$create_post_with_path$ language plpgsql;

create trigger create_post_with_path
    before insert on post
    for each row execute procedure create_post_with_path();

create or replace function increment_forum_threads() returns trigger as $increment_forum_threads$
begin
    update "forum" set threads = threads + 1 where forum.id = new.forum_id;
    insert into "forum_user" (forum_id, user_id) values (new.forum_id, new.user_id)
    on conflict do nothing;
    return null;
end;
$increment_forum_threads$ language plpgsql;

create trigger increment_forum_threads
    after insert on thread
    for each row execute procedure increment_forum_threads();

drop index if exists index_user_on_nickname;
create index if not exists index_user_on_nickname on "user" using hash(nickname);
drop index if exists index_user_on_email;
create index if not exists index_user_on_email on "user" using hash(email);

drop index if exists index_forum_slug;
create index if not exists index_forum_slug on forum using hash(slug);

drop index if exists index_thread_on_created;
create index if not exists index_thread_on_created on thread(created);
drop index if exists index_thread_on_slug;
create index if not exists index_thread_on_slug on thread using hash(slug);
drop index if exists index_thread_on_forum_and_created;
create index if not exists index_thread_on_forum_and_created on thread(forum_id, created);
drop index if exists index_thread_on_forum;
create index if not exists index_thread_on_forum on thread using hash(forum_id);
drop index if exists thread_on_id_and_forum;
create index if not exists thread_on_id_and_forum ON thread(id, forum_id);

drop index if exists index_post_on_id;
create unique index if not exists index_post_on_id ON post(id);

drop index if exists index_post_on_thread;
create index if not exists index_post_on_thread on post using hash (thread_id);

drop index if exists index_post_on_thread_and_path_and_id;
create index if not exists index_thread_on_forum_and_created on "thread"(forum_id, created);

drop index if exists index_post_on_thread_and_path_and_id;
create index if not exists index_post_on_thread_and_path_and_id on "post"(thread_id, path, id);

drop index if exists index_post_on_parent;
create index if not exists index_post_on_parent on "post"(parent, id);

drop index if exists index_post_on_parent_path_and_path;
create index if not exists index_post_on_parent_path_and_path on "post"((path[1]), path);

drop index if exists index_post_on_parent_path_and_path_and_id;
create index if not exists  index_post_on_parent_path_and_path_and_id ON post ((path[1]), path, id);

drop index if exists index_post_on_parent_and_thread;
create index if not exists index_post_on_parent_and_thread ON post(parent, thread_id);

VACUUM;
VACUUM ANALYSE;
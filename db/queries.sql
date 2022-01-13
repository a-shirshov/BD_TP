insert into "user" (nickname, fullname, about, email) values ('pirate', 'Jack Sparrow', 'About', 'pirate@mail.ru');
insert into "user" (nickname, fullname, about, email) values ('admin', 'Artyom Shirshov', 'Proger', 'admin@mail.ru');

select * from "user";

insert into "forum" (title,slug,user_id) values ('Threasures','pirate stories',1);
insert into "forum" (title,slug,user_id) values ('Code','Go',2);

select * from "forum";

insert into "thread" (title,message,created,slug,user_id,forum_id) values ('Coins','Golden Coins are here','2022-01-09','pirate1',1,1);
insert into "thread" (title,message,created,slug,user_id,forum_id) values ('Chests','Chests are somethere','2022-01-05','pirate2',1,1);
insert into "thread" (title,message,created,slug,user_id,forum_id) values ('Big Money','Progers love pirate money','2022-01-11','pirate3',2,1);
insert into "thread" (title,message,created,slug,user_id,forum_id) values ('BD_TP','So Tired','2022-01-09','pirate4',2,2);

select * from "thread";

--Передай f.slug и всё. Получишь Треды. 
select t.id, t.title,u.nickname as author,f.slug as forum,t.message,t.votes,f.slug,t.created from "thread" as t 
join "forum" as f on t.forum_id = f.id  
join "user" as u on u.id = t.user_id where f.slug = 'pirate stories' 
ORDER BY t.title DESC
LIMIT 2;

--Posts by id create
insert into "post" (parent,message,user_id,thread_id,created) 
    select 0,'first_post',u.id,t.id,'now' from "user" as u, "thread" as t 
    where u.nickname = 'pirate' AND t.id = 1;

insert into "post" (parent,message,user_id,thread_id,created) 
    select 0,'second_post',u.id,t.id,'now' from "user" as u, "thread" as t 
    where u.nickname = 'pirate' AND t.id = 2;

insert into "post" (parent,message,user_id,thread_id,created) 
    select 0,'third_post',u.id,t.id,'now' from "user" as u, "thread" as t 
    where u.nickname = 'admin' AND t.id = 1;

insert into "post" (parent,message,user_id,thread_id,created) 
    select 0,'fourth_post',u.id,t.id,'now' from "user" as u, "thread" as t 
    where u.nickname = 'pirate' AND t.id = 1 returning id;

select * from "post";

select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,p.thread_id as thread,p.created from "post" as p
    join "thread" as t on t.id = p.thread_id
    join "forum" as f on f.id = t.forum_id
    join "user" as u on u.id = p.user_id
    where p.id = 1;

--Select posts id
select p.id, p.edited, f.slug as forum, p.created from "post" as p
    join "thread" as t on p.thread_id = t.id 
    join "forum" as f on f.id = t.forum_id
    where t.id = 2;

--Posts by id create
insert into "post" (parent,message,user_id,thread_id,created) 
    select 0,'first_post',u.id,t.id,'now' from "user" as u, "thread" as t 
    where u.nickname = 'pirate' AND t.id = 1;

--Thread details by id
select t.title,u.nickname as author,f.slug as forum,t.message,t.votes,t.slug,t.created from "thread" as t
    join "user" as u on u.id = t.user_id
    join "forum" as f on f.id = t.forum_id
    where t.id = 2;

--Posts details update by id
update "thread" set title = 'ewrw',message = 'updated message' 
    where id = 1;

--Vote insert
insert into "vote" (voice,user_id,thread_id)
    select -1, u.id, t.id from "user" as u, "thread" as t
    where u.nickname = 'pirate' and t.id = 1;

insert into "post" (parent,message,user_id,thread_id,created) 
    select 0,'message',u.id,t.id,'now' from "user" as u, "thread" as t 
    where u.nickname = 'pirate' AND t.slug = 'PirateBattle/thread13baa5ca-bfe6-46e3-b910-81871d3e4fbb/5dee5d37-8e10-4be1-a997-c7af889ea3cc' returning id;

select p.id, p.edited, f.slug as forum, p.created,t.id as thread from "post" as p
    join "thread" as t on p.thread_id = t.id 
    join "forum" as f on f.id = t.forum_id
    where t.id = 1 and p.id = 1;



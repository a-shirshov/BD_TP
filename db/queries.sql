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
join "user" as u on u.id = t.user_id where f.slug = 'i-ojlRKE2u8Us' 
ORDER BY t.title DESC
LIMIT 100;

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

select * from "vote"
    join "user" as u on "vote".user_id = u.id
    join "thread" as t on t.id = "vote".thread_id
    where u.nickname = 'silentio.j0vCx4kA1trxpD' and t.id = 75;

update "vote"
    set voice = -1 
    from (select u.nickname,t.id from "vote" as table_vote
        join "user" as u on u.id = table_vote.user_id
        join "thread" as t on t.id = table_vote.thread_id
    ) as subquery
    where subquery.nickname = 'silentio.j0vCx4kA1trxpD';

select u.nickname,t.id from "vote" as table_vote
        join "user" as u on u.id = table_vote.user_id
        join "thread" as t on t.id = table_vote.thread_id
        where u.nickname = 'silentio.j0vCx4kA1trxpD'

update "vote" 
        join "user" as u on u.id = "vote" .user_id
        join "thread" as t on t.id = "vote" .thread_id
    set voice = 1
    where u.nickname = 'silentio.j0vCx4kA1trxpD'

update "vote"
    set voice = 1 
    where vote.t_id = ;


select p.id,p.parent,u.nickname as author,p.message,p.edited,f.slug as forum,t.id as thread,p.created
	from "post" as p
	join "thread" as t on p.thread_id = t.id
	join "forum" as f on f.id = t.forum_id
	join "user" as u on u.id = p.user_id
	where p.thread_id = $1 and p.id = $2
	order by p.id desc limit $3;


insert into "post" (parent,message,user_id,forum,thread_id,created) 
    select 0,'qwe',u.id,'O-z-M0A2-98y82',t.id,'2021-01-21T12:22:38.680+03:00' from "user" as u
    join "thread" as t on t.user_id = u.id
    join "forum" as f on f.slug = t.forum
    where u.nickname = 'bonam.NHQ6kGl1h87sJU' AND t.id = 68  returning id;


select u.id,'O-z-M0A2-98y82',t.id from "user" as u, "thread" as t 
    where u.nickname = 'bonam.NHQ6kGl1h87sJU' AND t.id = 68 and t.forum = 'O-z-M0A2-98y82';

select * from "forum" as f  
join "thread" as t on f.id = t.forum_id
join "user" as u on t.user_id = u.id
where f.slug ='O-z-M0A2-98y82' and t.id = 68 and u.nickname = 'bonam.NHQ6kGl1h87sJU'
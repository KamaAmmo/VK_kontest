
select * from people;

select * from films;

select * from films_actors;

delete from people where id = 13;

select p.name, f.title from people AS p LEFT JOIN films_actors fa ON p.id = fa.actor_id
LEFT JOIN films f ON f.id = fa.film_id
GROUP BY p.id, f.title
ORDER BY p.id;


UPDATE people SET (name, gender, birth_date) = ('Lool', 'Male', '2001-09-28')
where id = 19;



select * from people;

select * from films ORDER BY rating DESC;

select * from films_actors;

delete from people where id = 13;

select p.name, f.title from people AS p LEFT JOIN films_actors fa ON p.id = fa.actor_id
LEFT JOIN films f ON f.id = fa.film_id
GROUP BY p.id, f.title
ORDER BY p.id;


UPDATE people SET (name, gender, birth_date) = ('Lool', 'Male', '2001-09-28')
where id = 19;

SELECT DISTINCT f.title FROM films f JOIN films_actors fa ON fa.film_id = f.id
JOIN people p ON p.id = fa.actor_id WHERE  p.name LIKE '%Стрип%';
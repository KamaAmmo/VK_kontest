INSERT INTO people (name, gender, birth_date)
VALUES 
('Leonardo DiCaprio', 'Male', '1974-11-11'),
('Meryl Streep', 'Female', '1949-07-08'),
('Robert Downey Jr.', 'Male', '1965-04-04'),
('Annette Bening', 'Female', '1949-06-24'),
('John Travolta', 'Male', '1954-02-18'),
('Matthew McFadden', 'Male', '1974-06-19'),
('Julia Roberts', 'Female', '1967-10-28'),
('Brad Pitt', 'Male', '1963-12-18'),
('Angela Lansbury', 'Female', '1985-03-14'),
('Kate Blanchett', 'Female', '1969-08-14'),
('Bruce Willis', 'Male', '1955-09-27'),
('Anne Hathaway', 'Female', '1981-05-15');

INSERT INTO films (title, description, release_date, rating)
VALUES 
('Titanic', 'A love story', '1997-12-19', 7.8),
('The Aviator', 'The biography of Charlie Chaplin', '2004-12-10', 7.6),
('The Social Network', 'The story of the creation of Facebook', '2010-10-01', 7.7),
('Taxi Driver', 'A story of love', '1998-06-19', 7.6),
('Interstellar', 'A journey into space', '2014-11-07', 8.6),
('Pirates of the Caribbean', 'A pirate romance', '2000-12-15', 7.2),
('The Devil Wears Prada', 'The harsh world of fashion', '2005-09-21', 7.6),
('Die Hard', 'An unusual day to die', '1990-11-23', 7.5);

INSERT INTO films_actors (film_id, actor_id)
VALUES 
(1, 1), -- Леонардо ДиКаприо в "Титанике"
(2, 2), -- Мэрил Стрип в "Авиаторе"
(2, 10), -- Кейт Бланшетт в "Интерстелларе"
(2, 7), -- Джулия Робертс в "Авиаторе"
(3, 3), -- Роберт Дауни-младший в "Социальной сети"
(3, 8), -- Брэд Питт в "Социальной сети"
(4, 4), -- Аннетт Бенинг в "Таксисте"
(4, 9), -- Ангула Лэнсберри в "Таксисте"
(5, 5), -- Джон Траволта в "Интерстелларе"
(6, 6), -- Мэттью Макфэдиен в "Острове проклятых"
(7, 12), --Энн Хэтуэй в Дьявол носит прада 
(7, 2), -- Мерил Стрип в дьявол носит прада
(8, 11); --Брюс Уильяс в КО




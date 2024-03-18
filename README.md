# Инструкция
## Подготовка
Создаем и поднимаем контейнеры:
`docker build -t app:1 . && docker compose up`   
Далее необходимо создать таблицы в контейнере database и (по желанию заполнить ее) - для этого выполняем:  
`psql -h localhost -p 5431 -U vk_admin --d vk_db`
Пароль: `vk_pass`
После чего выполняем команды указанные в файле **prepare.sql**
## Работа с приложением 
Необходимо перейти по адресу http://localhost:5000/swagger  
Сначала необходимо зарегестрироваться, выбрать логин, пароль и роль. Потом залогиниться (вводим логин и пароль) - получим jwt токен.  
Далее вводим полученный токен для авторизации и тестим приложение.
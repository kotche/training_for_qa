## Проект для тренировок QA: поиск багов, составление тест-кейсов и т.д.

## Команды docker-compose
>docker-compose up --build -d (запустить контейнер)     
>docker-compose down (оставить контейнер)    
>docker-compose down -v (оставить контейнер и удалить volume)    

## Проверка доступа к БД
>docker exec -it postgres_qa psql -U youruser -d yourdb  (вызвать напрямую pg из контейнера)      
>psql -h localhost -p 5433 -U youruser -d yourdb 

Тестировать по пути http://localhost:8080/   

### Версия 1   
register.html - форма регистрации    

### Версия 2   
login.html - форма входа в систему   
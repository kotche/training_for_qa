## Проект для тренировок QA: поиск багов, составление тест-кейсов и т.д.

## Команды docker-compose
>docker-compose up --build -d (запустить контейнер)     
>docker-compose down (оставить контейнер)    
>docker-compose down -v (оставить контейнер и удалить volume)    

## Проверка доступа к БД
>docker exec -it postgres_qa psql -U youruser -d yourdb  (вызвать напрямую pg из контейнера)      
>psql -h localhost -p 5433 -U youruser -d yourdb 

### Версия 1
При заходе на http://localhost:8080/ показывает форму регистрации
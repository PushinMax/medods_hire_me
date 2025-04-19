# medods_hire_me


## **Запуск**
```bash
docker build -t medods-server .

docker run -p 8080:8080 medods-server
```

## **Подключиться к БД**
```bash
docker exec -it postgres_db psql -U postgres
```
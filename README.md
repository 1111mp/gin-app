# gin-app

A template for rapidly creating RESTful web services with Gin.

# Prerequisites

## PostgreSQL

- Install PostgreSQL use docker with the following command:
  ```bash
  docker run -d --name pg-server -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=db_name -p 5432:5432 -v pgdata:/var/lib/postgresql postgres
  ```

## Redis

- Install Redis use docker with the following command:
  ```bash
  docker run -d --name redis -p 6379:6379 -v redis-data:/data redis redis-server --requirepass yourpassword --appendonly yes
  ```

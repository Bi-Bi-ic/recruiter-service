[![CircleCI](https://circleci.com/gh/victorsteven/Go-JWT-Postgres-Mysql-Restful-API.svg?style=svg)](https://circleci.com/gh/victorsteven/Go-JWT-Postgres-Mysql-Restful-API)

# Go-JWT-Postgres-Mysql-Restful-API
This is a application build with golang, jwt, gorm, postgresql, mysql

## How To Use
- Running in local:
    - In a root Wordir (have main.go)
    ```go
    go run main.go
    ```

- Build with Docker-Compose:
    - 1. Have Docker, Docker-Compose installed:
        - [Links to Docker Website](https://www.docker.com/)
    - 2. If you have PostgreSQL installed on your local machine:
        - Stop the PostgreSQL service:
            - Linux
                ```
                systemctl stop postgresql
                ```
            - Windows
                - First, you need to find the PostgreSQL database directory, it can be something like "C:\Program Files\PostgreSQL\10.4\data" . Then open Command Prompt and execute this command:
                ```cmd
                pg_ctl -D "C:\Program Files\PostgreSQL\10.4\data" stop
                ```

            - MacOS
                - If you installed PostgreSQL via Homebrew:
                ```terminal
                pg_ctl -D /usr/local/var/postgres stop
                ```
    - 3. Buid it Or Update, Run:
            - ```
              sudo bash ./deploy.sh
              ```
Note: At Default:
- we use port 8080 for api
- we use port 5432 for database

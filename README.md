## подключение бд в докер
     docker run postgres:17-bookworm
    
## задаем пароль
    docker run -e POSTGRES_PASSWORD=111 postgres:17-bookworm

## задаем порт
    docker run -e POSTGRES_PASSWORD=111 -p 5432:5432 -d postgres:17-bookworm

## создаем volume для postgres чтобы сохранялись наши таблицы создаем папку out в нашей директории
    docker run -e POSTGRES_PASSWORD=111 -p 5432:5432 -v ./out/pgdata:/var/lib/postgresql -d postgres:17-bookworm

## выключение локальных портов бд check
    sudo -u postgres /Library/PostgreSQL/17/bin/pg_ctl -D /Library/PostgreSQL/17/data status/stop 

## удаляем все контейнеры docker 
    docker container prune






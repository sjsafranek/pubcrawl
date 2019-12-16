psql -c "CREATE USER crawluser WITH PASSWORD 'dev'"
psql -c "CREATE DATABASE crawldb"
psql -c "GRANT ALL PRIVILEGES ON DATABASE crawldb to crawluser"
psql -c "ALTER USER crawluser WITH SUPERUSER"

PGPASSWORD=dev psql -d crawldb -U crawluser -f db_setup.sql

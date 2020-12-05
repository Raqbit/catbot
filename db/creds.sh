user=catbot
password=catbot
host=localhost
port=5432
database=catbot
options="sslmode=disable"

export SOURCE="postgres://${user}:${password}@${host}:${port}/${database}?${options}"

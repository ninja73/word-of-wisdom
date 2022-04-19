# word of wisdom
Для реализации защиты от DDOS атак используется алгоритм hashcash.
На данный ограничения начинают действовать при достижении определенного порога, для всех запросов,
в дальнейшем лучше предусмотреть применение ограничения в разрезе пользователя.

Основной стек: **Golang, Redis, Docker**

Запуск приложения: **docker-compose up --build**

Конфигурация серверной части:  **server-config.toml**

```
store-file = "./wow-db.txt" - файл со словорем word-of-wisdom

[server]
address = ":8080" - адрес на котором стартует tcp сервер
bit-strength = 20 - сложность алгоритма hashcash
secret-key = "GLx%y@~z5mR6V3p6" - серетный ключ для подписи сигнатуры challenge
timeout = "2s" - таймаут tcp содениения 
expiration = "4m" - время жизни challenge
rate-limit = 10 - порог максимальной нагрузки, после которой мы будем отдавать challenge

[cache-redis]
host = "redis"
port = 6379
password = ""
db = 0
pool-size = 100
```

Конфигурация клиенткой части:  **client-config.toml**

```
server-address = "server:8080" - адрес tcp сервера
timeout = "2s" - таймаут tcp содениения 
clients = 10 - колличество параллельных клиентов
```
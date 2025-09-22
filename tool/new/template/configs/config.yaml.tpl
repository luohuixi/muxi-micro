server:
  grpc:
    addr: 0.0.0.0:50051
    timeout: 5s
data:
  mysql:
    dsn: root:root@tcp(localhost:3306)/circle2?parseTime=true&charset=utf8mb4&loc=Local
  redis:
    addr: localhost:6379
    password: 123456
    num: 0
    read: 0.2s
    write: 0.2s
log:
    dir: ./logs

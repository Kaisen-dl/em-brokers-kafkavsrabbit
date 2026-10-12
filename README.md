# Notifications Demo

## Запуск инфраструктуры

```bash
docker compose up -d
```

```bash
docker compose ps
```

## Kafka

```bash
go run ./kafka/consumer A
```

```bash
go run ./kafka/consumer B
```

```bash
go run ./kafka/producer
```

## RabbitMQ

```bash
go run ./rabbit/email_consumer E1
```

```bash
go run ./rabbit/email_consumer E2
```

```bash
go run ./rabbit/push_consumer P1
```

```bash
go run ./rabbit/push_consumer P2
```

```bash
go run ./rabbit/publisher
```

## RabbitMQ UI

```
http://localhost:15672
guest / guest
```

## Остановка

```bash
docker compose down
```

```bash
docker compose down -v
```
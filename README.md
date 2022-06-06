# RabbitMQ Churn Benckmark

### Benchmark to compare maximum churn capacity between two RabbitMQ versions

##### Start the two RabbitMQ instances
`docker-compose build && docker-compose up -d`

##### Benchmark RabbitMQ 3.10.5:
`docker-compose up bench-new`

##### Benchmark RabbitMQ 3.7.28:
`docker-compose up bench-old`

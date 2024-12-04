# airbnb-cli

This is a project called [airbnb-cli](https://docs.google.com/document/d/1JqpxeSFDkspIiW16zqwTCLaIYEByVUSy7CmT-yCJRBw/edit?tab=t.0#heading=h.fg54jvgi09lw)



## project outline
There are a produder and a consumer involved in this project.
Besides, we also need a message queue (Kafka) and database (Mysql or something).

## prerequisites

### Message Queue Server

Since we are using producer-consumer model, then we need a queue like structure, kafka is a good choice for it. 

I am running a kafka server in container mode. To simplify the deployment, I am running it in single node mode. 

#### Deploy Kafka
```
docker run -d --name kafka \
  -p 9092:9092 \
  -e KAFKA_CFG_SASL_ENABLED_MECHANISMS=PLAIN \
  -e KAFKA_CFG_SASL_MECHANISM_INTER_BROKER_PROTOCOL=PLAIN \
  -e KAFKA_CFG_AUTHORIZER_CLASS_NAME=kafka.security.authorizer.AclAuthorizer \
  -e KAFKA_CFG_ALLOW_EVERYONE_IF_NO_ACL_FOUND=false \
  -e KAFKA_CFG_SUPER_USERS=User:kafka_admin \
  -e KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,SASL_PLAINTEXT:SASL_PLAINTEXT \
  -e KAFKA_CFG_LISTENERS=SASL_PLAINTEXT://:9092 \
  -e KAFKA_CFG_ADVERTISED_LISTENERS=SASL_PLAINTEXT://localhost:9092 \
  -e KAFKA_CFG_SASL_JAAS_CONFIG="org.apache.kafka.common.security.plain.PlainLoginModule required \
  username=\"ctw_kafka_admin\" password=\"ctw_admin_pass\";" \
  bitnami/kafka:latest
```

#### Create Topic in Kafka
Set partitions to 3 for dev mode.
```
docker exec -it kafka kafka-topics.sh --create \
  --topic airbnb-topic \
  --bootstrap-server localhost:9092 \
  --partitions 3 --replication-factor 1
```


## producer
The producer will mainly doing one jobs, it will traverse the li

./airbnb-cli --command=producer --data=tasks.json --kafka=localhost:9092 --topic=airbnb-topic

## consumer
The consumer will get links from kafka topic, and made request;
It will extract all required information and store them into DBs.

./airbnb-cli --command=consumer --workers=10 --kafka=localhost:9092 --topic=airbnb-topic

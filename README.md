# Docker Compose Setup for AMQP Producer, Consumer, RabbitMQ, and MailHog

This Docker Compose configuration sets up a development environment for an AMQP producer, an AMQP consumer, RabbitMQ, and MailHog.
## Services
### MailHog
- **Image:** mailhog/mailhog
- **Ports:**
  - 8025: Web UI
  - 1025: SMTP access

### Producer
- **Image:** ghcr.io/kilianp07/amqp_producer:v0.1
- **Restart:** always
- **Depends On:**
  - rabbitmq (condition: service_healthy)
- **Environment Variables:**
  - RABBITMQ_HOST: rabbitmq
  - SMTP_HOST: mailhog
- **Command:** amqpp
- **Volumes:**
  - ./data:/tmp
- **Ports:**
  - 8080: API access

### Consumer
- **Image:** ghcr.io/kilianp07/amqp_consumer:v0.1
- **Restart:** always
- **Depends On:**
  - rabbitmq (condition: service_healthy)
- **Environment Variables:**
  - RABBIT_HOST: rabbitmq
  - SMTP_HOST: mailhog
- **Command:** amqpc

## Usage

1. Ensure you have Docker and Docker Compose installed on your machine.
2. Clone the repository and navigate to the directory containing the Docker Compose file.
3. Run `docker-compose up -d` to start the services in detached mode.
4. Wait for the producer and consumer to start up and for RabbitMQ and MailHog to become healthy. You'll see the following message in the logs when the producer and consumer are ready to use:
```
amqp_producer-producer-1  | POST   /api/register             --> main.registerHandler (3 handlers)
```
5. Access MailHog's web UI at [http://localhost:8025](http://localhost:8025) to view emails sent by the producer.

## API Usage

To use the API, register a user using the following `curl` command:

```bash
curl -X POST -H "Content-Type: application/json" -d '{"username":"votre_utilisateur", "password":"votre_mot_de_passe", "mail":"mail@mail.com"}' http://localhost:8080/api/register
```

You can then check the MailHog web UI at [http://localhost:8025](http://localhost:8025) to view the received email.
package main

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var smtpHost string

func main() {
	log.SetOutput(os.Stdout)
	rabbitHost := os.Getenv("RABBIT_HOST")
	if rabbitHost == "" {
		log.Fatalf("RABBIT_HOST environment variable is not set")
		panic(errors.New("RABBIT_HOST environment variable is not set"))
	}

	smtpHost = os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		log.Fatalf("SMTP_HOST environment variable is not set")
		panic(errors.New("SMTP_HOST environment variable is not set"))
	}

	amqpConn, err := amqp.Dial("amqp://guest:guest@" + rabbitHost + ":5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer amqpChannel.Close()

	// Déclarez une file RabbitMQ à écouter
	queue, err := amqpChannel.QueueDeclare(
		"users", // Remplacez par le nom de votre file RabbitMQ
		false,   // Durable
		false,   // Auto-delete
		false,   // Non-exclusif
		false,   // No-wait
		nil,     // Arguments supplémentaires
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Configurez un consommateur pour la file spécifiée
	msgs, err := amqpChannel.Consume(
		queue.Name, // Nom de la file
		"",         // Nom du consommateur (laissez vide pour un nom généré automatiquement)
		true,       // Auto-acknowledge (le message sera marqué comme confirmé automatiquement)
		false,      // Exclusive
		false,      // No-local
		false,      // No-wait
		nil,        // Arguments supplémentaires
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Boucle infinie pour écouter les messages
	for msg := range msgs {
		// Traitez le message ici
		email := string(msg.Body)
		sendVerificationEmail(email)
		fmt.Println("Message reçu:", email)
	}
}

func sendVerificationEmail(toEmail string) {
	from := "your_email@example.com"
	password := "your_email_password"

	msg := "From: " + from + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: Vérification d'e-mail\n\n" +
		"Bonjour,\n\n" +
		"Merci de vous être enregistré. Veuillez cliquer sur le lien suivant pour vérifier votre adresse e-mail.\n" +
		"http://localhost:8025/ (MailHog Web UI)\n\n" +
		"Cordialement,\n" +
		"Votre Application"

	err := smtp.SendMail(smtpHost+":1025", smtp.CRAMMD5Auth(from, password), from, []string{toEmail}, []byte(msg))
	if err != nil {
		fmt.Println("Failed to send verification email:", err)
	}
}

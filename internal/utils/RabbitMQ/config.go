package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	menunotificationmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/menu_notif_model"
	resetnotifmodel "github.com/KusakinDev/Catering-Notif-Service/internal/models/reset_notif_model"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	ResetConsumer <-chan amqp.Delivery
	MenuConsumer  <-chan amqp.Delivery
}

func (rmq *RabbitMQ) InitConnection() error {
	var err error

	err = godotenv.Load("dev.env")
	if err != nil {
		log.Fatalf("Error load .env: %v", err)
	}
	host := os.Getenv("RMQ_HOST")
	user := os.Getenv("RMQ_USER")
	password := os.Getenv("RMQ_PASSWORD")
	port := os.Getenv("RMQ_PORT")
	dsn := "amqp://" + user + ":" + password + "@" + host + ":" + port + "/"
	rmq.connection, err = amqp.Dial(dsn)
	if err != nil {
		logrus.Errorln("RabbitMQ.connecction: ", err)
		return err
	}
	logrus.Infoln("RabbitMQ.chanel SUCCESS CONNECT")
	return nil
}

func (rmq *RabbitMQ) InitChannel() error {
	var err error
	rmq.channel, err = rmq.connection.Channel()
	if err != nil {
		logrus.Errorln("RabbitMQ.chanel: ", err)
		return err
	}
	logrus.Infoln("RabbitMQ.chanel SUCCESS OPEN")
	return nil
}

func (rmq *RabbitMQ) DeclareQueue(queueName string) error {
	_, err := rmq.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (rmq *RabbitMQ) InitConsumer(queueName string, consumerType string) error {
	if err := rmq.DeclareQueue(queueName); err != nil {
		logrus.Errorln("Error declaring queue: ", err)
		return err
	}

	var err error
	consumer, err := rmq.channel.Consume(
		queueName, // queue name
		"",        // consumer tag
		true,      // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		logrus.Errorln("RabbitMQ.Queue: ", err)
		return err
	}

	switch consumerType {
	case "reset":
		rmq.ResetConsumer = consumer
	case "menu":
		rmq.MenuConsumer = consumer
	default:
		return fmt.Errorf("unknown consumer type: %s", consumerType)
	}

	return nil
}

func (rmq *RabbitMQ) Publish(body []byte, queueName string) error {
	err := rmq.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		logrus.Errorln("RabbitMQ.Publish: ", err)
	}
	logrus.Infoln("RabbitMQ.Publish: Success publish")
	return err
}

func (rmq *RabbitMQ) ConsumeReset() {
	for d := range rmq.ResetConsumer {
		var notif resetnotifmodel.ResetNotification
		err := json.Unmarshal(d.Body, &notif)
		if err != nil {
			logrus.Errorln("Error decoding notification: ", err)
			continue
		}
		notif.Send()
		logrus.Infoln("RabbitMQ.ConsumeReset: Success send email")
	}
}

func (rmq *RabbitMQ) ConsumeMenu() {
	for d := range rmq.MenuConsumer {
		var notif menunotificationmodel.MenuNotification
		err := json.Unmarshal(d.Body, &notif)
		if err != nil {
			logrus.Errorln("Error decoding notification: ", err)
			continue
		}
		notif.Send()
	}
}

func (rmq *RabbitMQ) Close() {
	if rmq.channel != nil {
		rmq.channel.Close()
	}
	if rmq.connection != nil {
		rmq.connection.Close()
	}
}

package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
)

const (
	networkTCP = "tcp"
	partition  = 0
	topicName  = "order"
)

type config struct {
	User     string `env:"KAFKA_USER" envDefault:"pds"`
	Password string `env:"KAFKA_PASS,required"`
	Host     string `env:"KAFKA_HOST,required"`
	Port     string `env:"KAFKA_PORT" envDefault:"9092"`
}

type PxOrder struct {
	OrderId     int    `bson:"orderid,omitempty"`
	Email       string `bson:"email,omitempty"`
	Main        string `bson:"main,omitempty"`
	Side1       string `bson:"side1,omitempty"`
	Side2       string `bson:"side2,omitempty"`
	Drink       string `bson:"drink,omitempty"`
	Restaurant  string `bson:"restaurant,omitempty"`
	Date        string `bson:"date,omitempty"`
	Street1     string `bson:"street1,omitempty"`
	Street2     string `bson:"street2,omitempty"`
	City        string `bson:"city,omitempty"`
	State       string `bson:"state,omitempty"`
	Zip         string `bson:"zip,omitempty"`
	OrderStatus string `bson:"orderstatus,omitempty"`
}

func KafkaInit() {
	// Read config.
	fmt.Println("Checking Kafka Broker")
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("ERROR: failed to parse config: %v\n", err)
	}

	brokerURL := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	fmt.Println("Broker URL is : " + brokerURL)
	fmt.Println("Host is : " + cfg.Host)
	fmt.Println("Port is : " + cfg.Port)
	dialer := &kafka.Dialer{
		SASLMechanism: plainMechanism(cfg.User, cfg.Password),
		Timeout:       10 * time.Second,
		DualStack:     true,
	}

	// Create topic.
	err := createTopic(dialer, brokerURL, topicName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Wait few seconds to sync the topic and to avoid the "Unknown Topic Or Partition" error.
	time.Sleep(5 * time.Second)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
}

func connectToController(dialer *kafka.Dialer, url string) (*kafka.Conn, *kafka.Conn, error) {
	ctx := context.Background()

	// Connecting to broker url.
	//log.Printf("Connecting to %s\n", url)
	conn, err := dialer.DialContext(ctx, networkTCP, url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a connection: %s", err)
	}

	// Connecting to controller.
	controller, err := conn.Controller()
	if err != nil {
		return conn, nil, fmt.Errorf("failed to get the current controller: %s", err)
	}
	controllerUrl := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	//log.Printf("Connecting to controller %s\n", controllerUrl)
	controllerConn, err := dialer.DialContext(ctx, networkTCP, controllerUrl)
	if err != nil {
		return conn, nil, fmt.Errorf("failed to open a connection to the controller: %s", err)
	}

	return conn, controllerConn, err
}

func createTopic(dialer *kafka.Dialer, brokerURL, topicName string) error {
	// Connect to controller.
	conn, controllerConn, err := connectToController(dialer, brokerURL)
	if conn != nil {
		defer conn.Close()
	}
	if controllerConn != nil {
		defer controllerConn.Close()
	}
	if err != nil {
		return err
	}

	// Create topic.
	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topicName,
		NumPartitions:     1,
		ReplicationFactor: 3,
	})
	if err != nil {
		return fmt.Errorf("failed to create the %s topic: %s", topicName, err)
	}
	fmt.Println("Kafka Topic " + topicName + " is ready!")

	return nil
}

func deleteTopic(dialer *kafka.Dialer, brokerURL, topicName string) error {
	// Connect to controller.
	conn, controllerConn, err := connectToController(dialer, brokerURL)
	if conn != nil {
		defer conn.Close()
	}
	if controllerConn != nil {
		defer controllerConn.Close()
	}
	if err != nil {
		return err
	}

	// Delete topic.
	err = controllerConn.DeleteTopics(topicName)
	if err != nil {
		return fmt.Errorf("failed to delete the %s topic: %s", topicName, err)
	}

	return nil
}

func writeMessages(dialer *kafka.Dialer, url string, topic string, msg PxOrder) error {

	conn, controllerConn, err := connectToController(dialer, url)
	if conn != nil {
		defer conn.Close()
	}
	if controllerConn != nil {
		defer controllerConn.Close()
	}
	if err != nil {
		return err
	}

	// Find leader node for topic.
	ctx := context.Background()
	leader, err := dialer.LookupLeader(ctx, networkTCP, url, topic, partition)
	//leader := kafkaHost + ":" + kafkaPort
	if err != nil {
		return fmt.Errorf("failed to find leader for topic %s: %v", topic, err)
	}
	leaderURL := net.JoinHostPort(leader.Host, strconv.Itoa(leader.Port))

	//log.Printf("write messages (%s)\n", leader)

	w := newWriter(leaderURL, topic, dialer)

	fmt.Println("topic is : " + w.Topic)
	defer w.Close()

	//convert msg to a json object and store it in payload
	//payload, err := json.Marshal(msg)
	//if err != nil {
	//	fmt.Println("Failed to Marshal json")
	//}

	key := fmt.Sprintf("key%d", time.Now().Nanosecond())
	message := kafka.Message{
		Key:   []byte(key),
		Value: []byte("this is message" + key), //TEST DATA THIS SHOULD BE THE MSG JSON WHEN DONE
	}

	err = w.WriteMessages(ctx, message)
	if err != nil {
		return fmt.Errorf("write failed: %v", err)
	}

	return err
}

func readMessages(dialer *kafka.Dialer, url, topic string, count int) error {
	// Find leader node for topic.
	ctx := context.Background()
	leader, err := dialer.LookupLeader(ctx, networkTCP, url, topic, partition)
	if err != nil {
		return fmt.Errorf("failed to find leader for topic %s: %v", topic, err)
	}
	leaderURL := net.JoinHostPort(leader.Host, strconv.Itoa(leader.Port))

	log.Printf("read messages (%s)\n", leaderURL)
	r := newReader(leaderURL, topic, partition, dialer)
	defer r.Close()
	start := time.Now()
	for i := 0; i < count; i++ {
		_, err := r.ReadMessage(context.Background())
		if err != nil {
			return fmt.Errorf("read #%d failed: %v", i, err)
		}
	}
	stop := time.Now()
	log.Printf("%d reads done in %v\n", count, stop.Sub(start))
	return nil
}

func plainMechanism(user, password string) sasl.Mechanism {
	return plain.Mechanism{
		Username: user,
		Password: password,
	}
}

func newWriter(url string, topic string, dialer *kafka.Dialer) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(url),
		Topic:    "order",
		Balancer: &kafka.Hash{},
		Transport: &kafka.Transport{
			SASL: dialer.SASLMechanism,
		},
		BatchTimeout: 20 * time.Millisecond,
	}
}

func newReader(url string, topic string, partition int, dialer *kafka.Dialer) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{url},
		Topic:     topic,
		Dialer:    dialer,
		Partition: partition,
	})
}

func SubmitOrder(orderNum int, orderDate string, email string, restaurant string, main string, side1 string, side2 string, drink string, street1 string, street2 string, city string, state string, zipcode string) {

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("ERROR: failed to parse config: %v\n", err)
	}
	fmt.Println("cfg.host is : " + cfg.Host)

	msg := PxOrder{
		Email:       email,
		OrderId:     orderNum,
		Restaurant:  restaurant,
		Main:        main,
		Side1:       side1,
		Side2:       side2,
		Drink:       drink,
		Date:        orderDate,
		Street1:     street1,
		Street2:     street2,
		City:        city,
		State:       state,
		Zip:         zipcode,
		OrderStatus: "Pending",
	}

	dialer := &kafka.Dialer{
		SASLMechanism: plainMechanism(cfg.User, cfg.Password),
		Timeout:       10 * time.Second,
		DualStack:     true,
	}

	// create a new writer that writes to the topic "order" on localhost:9092
	brokerURL := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{brokerURL},
		Topic:        "order",
		Balancer:     &kafka.LeastBytes{},
		Dialer:       dialer,
		RequiredAcks: 1,
	})

	// marshal the payload into json
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to marshal payload: %v", err)
	}

	// write the message to the topic
	err = w.WriteMessages(context.Background(),
		kafka.Message{
			Value: payload,
		},
	)
	if err != nil {
		log.Fatalf("Failed to write message: %v", err)
	}

	w.Close()

}

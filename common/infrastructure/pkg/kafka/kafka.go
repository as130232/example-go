package pkgKafka

import (
	"git-new.okkia.site/crk/decimal-cricket-common/global"
	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/consts/env"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"

	"git-new.okkia.site/crk/decimal-cricket-common/infrastructure/utils"
)

func NewDecimalCricketKafkaWriter(topicName string) *kafka.Writer {
	servers := utils.GetKafkaAddress()
	return &kafka.Writer{
		Addr:                   kafka.TCP(servers...),
		Topic:                  topicName,
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
		BatchSize:              1000,
		BatchBytes:             10485760,
		BatchTimeout:           10 * time.Millisecond,
		ReadTimeout:            30 * time.Second,
		WriteTimeout:           30 * time.Second,
		MaxAttempts:            30,
		RequiredAcks:           kafka.RequireOne,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				actionId := uuid.NewString()
				log.Printf("kafka writer actionId:%s, error:%+v, message count:%d", actionId, err, len(messages))
				for message := range messages {
					log.Printf("kafka writer actionId:%s, error:%+v, messages:%+v", actionId, err, message)
				}
			}
		},
	}
}

func NewDecimalCricketKafkaWriterList(topicName string) []*kafka.Writer {
	var result []*kafka.Writer
	var kafkaArr []string

	if global.ServerConfig.AppEnv == env.Dev || global.ServerConfig.AppEnv == env.Local {
		kafkaArr = append(kafkaArr, global.ServerConfig.DecimalCricketKafkaEnv.Dev)
	} else if global.ServerConfig.AppEnv == env.Prod {
		kafkaArr = append(kafkaArr, global.ServerConfig.DecimalCricketKafkaEnv.Cqa)
		kafkaArr = append(kafkaArr, global.ServerConfig.DecimalCricketKafkaEnv.Uat)
		kafkaArr = append(kafkaArr, global.ServerConfig.DecimalCricketKafkaEnv.Prod)
	}

	for _, v := range kafkaArr {
		if len(v) == 0 {
			continue
		}
		servers := utils.GetKafkaAddressByArg(v)
		kafkaWriter := &kafka.Writer{
			Addr:                   kafka.TCP(servers...),
			Topic:                  topicName,
			AllowAutoTopicCreation: true,
			Balancer:               &kafka.LeastBytes{},
			BatchTimeout:           10 * time.Millisecond,
			RequiredAcks:           kafka.RequireOne,
		}
		result = append(result, kafkaWriter)
	}
	return result
}

func NewDecimalKafkaReader(groupId string, topicName string) *kafka.Reader {
	servers := utils.GetKafkaAddress()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:                servers,
		GroupID:                groupId,
		Topic:                  topicName,
		MaxBytes:               10e6, // 10MB
		QueueCapacity:          1000,
		SessionTimeout:         120 * time.Second,
		RebalanceTimeout:       60 * time.Second,
		WatchPartitionChanges:  true,
		PartitionWatchInterval: 5 * time.Minute,
	})

	return reader
}

func CreateTopic(topicName string, numPartitions int, replicationFactor int) {
	servers := utils.GetKafkaAddress()
	connection, err := kafka.Dial("tcp", servers[0])
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	controller, err := connection.Controller()
	if err != nil {
		panic(err)
	}
	var controllerConnection *kafka.Conn
	controllerConnection, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err)
	}
	defer controllerConnection.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topicName,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
		},
	}
	log.Printf("%+v", topicConfigs)
	err = controllerConnection.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err)
	}
}

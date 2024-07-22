package infra

import (
	"context"
	"errors"
	"messaggio/infra/kafka"
	"messaggio/internal/repository"
	"messaggio/pkg/util/logger"
	"messaggio/storage/postgres"
	"sync"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Infra interface {
	Config() *viper.Viper
	GetLogger() logger.Logger
	SetMode() string
	Port() string
	RedisClient() *redis.Client
	PSQLClient() (*postgres.PSQLClient, error)
	KafkaProducer() *kafka.KafkaProducer
	KafkaConsumer() *kafka.KafkaConsumer
}

// infra is an implementation of the Infra interface.
// It provides access to configuration, logging, and various infrastructure clients.
type infra struct {
	configFile string
	psqlClient *postgres.PSQLClient
}

// New creates a new instance of infra with the specified configuration file path.
// It initializes the PostgreSQL client and handles any initialization errors.
//
// configFile - The path to the configuration file.
//
// Returns:
// Infra - A new instance of infra.
func New(configFile string) Infra {
	i := &infra{configFile: configFile}
	var err error
	i.psqlClient, err = i.PSQLClient() // Инициализация psqlClient
	if err != nil {
		logrus.Fatalf("[infra][New] %v", err)
	}
	return i
}

var (
	vprOnce sync.Once
	vpr     *viper.Viper
)

// Config returns the Viper configuration instance.
// It initializes and reads the configuration from the specified configFile path.
//
// Returns:
// *viper.Viper - The Viper configuration instance.
func (i *infra) Config() *viper.Viper {
	vprOnce.Do(func() {
		viper.SetConfigFile(i.configFile)
		if err := viper.ReadInConfig(); err != nil {
			logrus.Fatalf("[infra][Config][viper.ReadInConfig] %v", err)
		}

		vpr = viper.GetViper()
	})

	return vpr
}

// GetLogger returns the application-wide logger instance.
//
// Returns:
// logger.Logger - The application-wide logger instance.
func (i *infra) GetLogger() logger.Logger {
	log := logger.GetLogger()
	return log
}

var (
	modeOnce    sync.Once
	mode        string
	development = "dev"
	production  = "release"
)

// SetMode sets the application mode based on the environment configuration.
// It retrieves the mode from the environment settings and configures the Gin framework accordingly.
//
// Returns:
// string - The current mode of the Gin framework.
func (i *infra) SetMode() string {
	modeOnce.Do(func() {
		env := i.Config().Sub("environment").GetString("mode")
		if env == development {
			mode = gin.DebugMode
		} else if env == production {
			mode = gin.ReleaseMode
		} else {
			logrus.Fatalf("[infa][SetMode] %v", errors.New("environment not setup"))
		}

		gin.SetMode(mode)
	})

	return mode
}

var (
	portOnce sync.Once
	port     string
)

// Port retrieves the server port from the configuration.
// It initializes the port once and returns it prefixed with ':'.
//
// Returns:
// string - The server port, prefixed with ':'.
func (i *infra) Port() string {
	portOnce.Do(func() {
		port = i.Config().Sub("server").GetString("port")
	})

	return ":" + port
}

var (
	rdbOnce sync.Once
	rdb     *redis.Client
)

// RedisClient returns a Redis client instance configured based on the environment settings.
// It initializes the Redis client once using the configured address, password, and database.
// It also verifies the connection by pinging the Redis server.
//
// Returns:
// *redis.Client - The Redis client instance.
func (i *infra) RedisClient() *redis.Client {
	rdbOnce.Do(func() {
		config := i.Config().Sub("redis")
		addr := config.GetString("addr")
		password := config.GetString("password")
		db := config.GetInt("db")

		rdb = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		if _, err := rdb.Ping(context.Background()).Result(); err != nil {
			logrus.Fatalf("[infra][RedisClient][rdb.Ping] %v", err)
		}

		logrus.Println("Connected to Redis")
	})

	return rdb
}

// PSQLClient returns a PostgreSQL client instance initialized with the configuration settings.
// It creates a new PostgreSQL client and establishes a connection using provided credentials.
//
// Returns:
// *postgres.PSQLClient - The PostgreSQL client instance.
// error - An error if there was an issue creating the PostgreSQL client or connecting to the database.
func (i *infra) PSQLClient() (*postgres.PSQLClient, error) {
	if i.psqlClient == nil {
		config := i.Config().Sub("database")
		user := config.GetString("user")
		pass := config.GetString("pass")
		host := config.GetString("host")
		port := config.GetString("port")
		name := config.GetString("name")

		psqlClient := postgres.NewPSQLClient()
		err := psqlClient.Connect(user, pass, host, port, name)
		if err != nil {
			return nil, err
		}
		i.psqlClient = psqlClient
	}
	return i.psqlClient, nil
}

var (
	kafkaProducerOnce sync.Once
	kafkaProducer     *kafka.KafkaProducer

	kafkaConsumerOnce sync.Once
	kafkaConsumer     *kafka.KafkaConsumer
)

// KafkaProducer returns a Kafka producer instance configured based on the environment settings.
// It initializes the Kafka producer once using the configured bootstrap servers and settings.
//
// Returns:
// *kafka.KafkaProducer - The Kafka producer instance.
func (i *infra) KafkaProducer() *kafka.KafkaProducer {
	kafkaProducerOnce.Do(func() {
		brokers := i.Config().Sub("kafka").GetStringSlice("bootstrap_servers")

		config := sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Retry.Max = 5
		config.Producer.Return.Successes = true

		producer, err := sarama.NewSyncProducer(brokers, config)
		if err != nil {
			logrus.Fatalf("[infra][KafkaProducer] %v", err)
		}

		kafkaProducer = kafka.NewKafkaProducer(producer)
	})

	return kafkaProducer
}

// KafkaConsumer returns a Kafka consumer instance configured based on the environment settings.
// It initializes the Kafka consumer once using the configured bootstrap servers, group ID, and topic.
// It also initializes a message repository and message handler for processing messages.
//
// Returns:
// *kafka.KafkaConsumer - The Kafka consumer instance.
// error - An error if there was an issue creating the Kafka consumer or initializing the message repository.
func (i *infra) KafkaConsumer() *kafka.KafkaConsumer {
	kafkaConsumerOnce.Do(func() {
		if i.psqlClient == nil {
			_, err := i.PSQLClient()
			if err != nil {
				logrus.Fatalf("[infra][KafkaConsumer] %v", err)
			}
		}

		config := i.Config().Sub("kafka")
		brokers := config.GetStringSlice("bootstrap_servers")
		groupID := config.GetString("group_id")
		topic := config.GetString("topic")

		messageRepo := repository.NewMessageRepository(i.psqlClient.Queries)

		handler := kafka.NewKafkaMessageHandler(messageRepo)

		var err error
		kafkaConsumer, err = kafka.NewKafkaConsumer(brokers, groupID, topic, handler)
		if err != nil {
			logrus.Fatalf("[infra][KafkaConsumer] %v", err)
		}
	})

	return kafkaConsumer
}

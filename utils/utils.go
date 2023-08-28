package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EnvironmentConfig struct {
	ServerPort    string
	DBURL         string
	CacheUrl      string
	CachePassword string
	CacheUsername string
	DBName        string
	DBUsername    string
	DBPassword    string
	RABBITMQ_HOST string
	REDIS_HOST    string
	l             *log.Logger
}

func LoadEnv(l *log.Logger) *EnvironmentConfig {
	if err := godotenv.Load(); err != nil {
		l.Fatalln("Error loading env file")
	}

	return &EnvironmentConfig{
		ServerPort:    os.Getenv("SERVER_PORT"),
		DBURL:         os.Getenv("DB_URL"),
		CacheUrl:      os.Getenv("CACHE_URL"),
		CachePassword: os.Getenv("CACHE_PASSWORD"),
		CacheUsername: os.Getenv("CACHE_USERNAME"),
		DBName:        os.Getenv("DB_NAME"),
		DBUsername:    os.Getenv("DB_USERNAME"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		REDIS_HOST:    os.Getenv("REDIS_HOST"),
		RABBITMQ_HOST: os.Getenv("RABBITMQ_HOST"),
		l:             l,
	}
}
func (c *EnvironmentConfig) InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.REDIS_HOST,
		Password: "",
		DB:       0,
	})
}

func (c *EnvironmentConfig) InitQueue() *amqp.Connection {

	conn, err := amqp.Dial(c.RABBITMQ_HOST)
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

func (env *EnvironmentConfig) ConnectToDB() *mongo.Database {
	env.l.Println("Starting connection to db")

	client, err := mongo.NewClient(options.Client().ApplyURI(env.DBURL))

	if err != nil {
		env.l.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		env.l.Fatalln(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		env.l.Fatalln(err)
	}

	env.l.Println("Connected to db")

	return client.Database(env.DBName)

}

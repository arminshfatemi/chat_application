package main

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Content   string             `bson:"content"`
	SenderID  primitive.ObjectID `bson:"sender_id"`
	RoomID    primitive.ObjectID `bson:"room_id"`
	Timestamp time.Time          `bson:"timestamp"`
}

const (
	TypeGetOldMessages     = "get:old_messages"
	TypeArchiveOldMessages = "archive:old_messages"
)

func initMongoClient() (*mongo.Collection, *mongo.Collection) {
	mongoURI := os.Getenv("MONGO_URI")
	opts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	messagesCollection := client.Database("chat_app").Collection("messages")
	articleCollection := client.Database("chat_app").Collection("archived_messages")
	return messagesCollection, articleCollection
}

func handleGetOldMessagesTask(ctx context.Context, t *asynq.Task) error {
	msgCollection, _ := initMongoClient()

	coll := msgCollection
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	filter := bson.M{"timestamp": bson.M{"$lt": oneHourAgo}}

	cur, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return err
	}

	defer cur.Close(context.Background())

	var messages []Message
	if err := cur.All(context.Background(), &messages); err != nil {
		log.Println("oldsss", err)
		return err
	}

	for _, msg := range messages {
		log.Printf("Found old message: %v", msg)
	}

	return nil
}

func handleArchiveOldMessagesTask(ctx context.Context, t *asynq.Task) error {
	coll, ar := initMongoClient()
	archiveColl := ar
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	filter := bson.M{"timestamp": bson.M{"$lt": oneHourAgo}}

	cur, err := coll.Find(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return err
	}

	defer cur.Close(context.Background())

	var messages []Message
	if err := cur.All(context.Background(), &messages); err != nil {
		return err
	}

	log.Println(messages)

	for _, msg := range messages {
		log.Println(msg)
		_, err := archiveColl.InsertOne(context.Background(), msg)
		if err != nil {
			log.Println("archive", err)
			return err
		}

		//_, err = coll.DeleteOne(context.Background(), bson.M{"_id": msg.ID})
		//if err != nil {
		//	return err
		//}
		log.Printf("Archived message: %v", msg)
	}

	return nil
}

func main() {
	redisURI := os.Getenv("REDIS_URI")
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisURI},
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeGetOldMessages, handleGetOldMessagesTask)
	mux.HandleFunc(TypeArchiveOldMessages, handleArchiveOldMessagesTask)

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run asynq server: %v", err)
		}
	}()

	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisURI},
		&asynq.SchedulerOpts{},
	)

	_, err := scheduler.Register("* * * * *", asynq.NewTask(TypeGetOldMessages, nil))
	if err != nil {
		log.Fatalf("could not register get old messages task: %v", err)
	}

	_, err = scheduler.Register("* * * * *", asynq.NewTask(TypeArchiveOldMessages, nil))
	if err != nil {
		log.Fatalf("could not register archive old messages task: %v", err)
	}

	if err := scheduler.Run(); err != nil {
		log.Fatalf("could not run scheduler: %v", err)
	}

	select {}

}

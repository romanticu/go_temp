package funcs

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongodbStream() {
	ctx := context.TODO()

	// clientOptions := options.Client().ApplyURI("mongodb://root:root@192.168.88.39:27017/highTech_Source?authSource=admin&authMechanism=PLAIN")

	credential := options.Credential{
		Username: "mroot",
		Password: "pacman",
	}
	clientOptions := options.Client().ApplyURI("mongodb://192.168.88.122:27017").SetAuth(credential)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	collection := client.Database("saic").Collection("socials")

	stream, err := collection.Watch(ctx, mongo.Pipeline{})

	if err != nil {
		log.Fatal(err)
	}
	log.Print("waiting for changes")
	var changeDoc map[string]interface{}
	for stream.Next(ctx) {
		if e := stream.Decode(&changeDoc); e != nil {
			log.Printf("error decoding: %s", e)
		}
		log.Printf("change: %+v", changeDoc)
	}
}

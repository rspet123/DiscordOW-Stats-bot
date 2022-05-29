package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func db() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connectString := fmt.Sprintf("mongodb+srv://dbuser:%s@ow-stats.b2vm3.mongodb.net", os.Getenv("DB_PASSWORD"))
	fmt.Println(connectString)
	fmt.Println(os.Getenv("DB_PASSWORD"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	ping := client.Ping(ctx, nil)
	fmt.Println(ping)
	return client
}

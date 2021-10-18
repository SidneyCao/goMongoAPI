package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	host       = flag.String("ho", "", "mongo host")
	user       = flag.String("u", "", "mongo user")
	password   = flag.String("p", "", "mongo password")
	db         = flag.String("d", "", "mongo db")
	collection = flag.String("c", "", "mongo collection")
	file       = flag.String("f", "", "read file")
)

//创建wait group
var waitGroup sync.WaitGroup

func Process(client *mongo.Client, collection *mongo.Collection, line string) {
	defer waitGroup.Done()
	date := strings.Split(strings.Split(line, string(uint64(1)))[2], " ")[0]
	sid := strings.Split(line, string(uint64(1)))[4]
	uid := strings.Split(line, string(uint64(1)))[7]
	actType := strings.Split(line, string(uint64(1)))[11]
	subType := strings.Split(line, string(uint64(1)))[12]
	fmt.Printf("%s,%s,%s,%s,%s\n", date, sid, uid, actType, subType)
	filter := bson.D{{"_id", date + "_" + sid + "_" + uid}}
	var result bson.D
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Printf("failed to search: %v", err)
	}
	if err == errors.New("mongo: no documents in result") {
		fmt.Println("匹配")
	}
	fmt.Println(result)
}

func main() {
	//获取参数
	flag.Parse()

	//wait group中始终有n+1个counter
	waitGroup.Add(1)

	//新建mongo client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017", *user, *password, *host)))
	if err != nil {
		log.Panicf("failed to connect mongo: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panicf("failed to disconnect mongo: %v", err)
		}
	}()
	//获取collection
	collection := client.Database(*db).Collection(*collection)

	//读取差异文件
	f, err := os.Open(*file)
	if err != nil {
		log.Panicf("failed to open diff file: %v", err)
	}
	defer f.Close()
	//按行读取文件
	br := bufio.NewReader(f)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		waitGroup.Add(1)
		go Process(client, collection, string(line))
	}

	waitGroup.Done()
	waitGroup.Wait()
}

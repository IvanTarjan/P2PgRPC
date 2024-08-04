package main

import (
	"context"
	"time"

	"log"
	// "time"

	"github.com/IvanTarjan/P2PgRPC/node"
	"github.com/IvanTarjan/P2PgRPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(time.Second)
	makeNode(":4000", []string{":3000"})
	time.Sleep(time.Second*3)
	makeNode(":5000", []string{":4000"})
	select{}
	// node := node.NewNode()
	
	// go func ()  {
	// 	for {
	// 		time.Sleep(2 * time.Second)
	// 		makeTransaction()
	// 	}
	// }()

	// log.Fatal(node.Start(":3000"))
}

func makeNode(listenString string, bootstrapNodes []string) *node.Node{
	n := node.NewNode()
	go n.Start(listenString, bootstrapNodes)
	return n
}

func makeTransaction(){
	client, err := grpc.NewClient("127.0.0.1:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	c := proto.NewNodeClient(client)
	// tx := &proto.Transaction{
	// 	Version: 1,
	// }

	version := &proto.Version{
		Version: "0.1",
		Height: 1,
	}

	_, err =c.Handshake(context.TODO(), version)
	if err != nil {
		log.Fatalln(err)
	}

}
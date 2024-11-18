package main

import (
	"bufio"
	"fmt"
	"log"
	"mandatory-handin-5/pb"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func client() {
	conn, err := grpc.Dial("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := pb.NewActionhouseClient(conn)

	if err != nil {
		log.Fatalf("Connection failure")
	}
	fmt.Printf("Welcome to ze bidding hall.")
	defer conn.Close()

	CommandlineInput(client)
}

//func bid(*pb.Request)(pb.ack) {}

func status(ctx context) {

	fmt.Printf()
}

func CommandlineInput(client pb.ActionhouseClient) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		if scanner.Scan() {
			scanText := scanner.Text()
			text := strings.Fields(scanText) // Use Fields to handle spaces correctly
			if len(text) > 1 && text[0] == "bid" {
				bidVal, err := strconv.Atoi(text[1])
				if err != nil || bidVal <= 0 {
					log.Printf("Cant bid 0 or below: %v", err)
					continue
				}
				res, err := client.Bid(&pb.Request{
					Bid:    int32(bidVal),
					UserId: 0,
				})
				if err != nil {
					log.Printf("Error Bidding: %v", err)
				}
			}
			if len(text) == 1 && text[0] == "result" {
				// in here we fetch the result.
				// result, err := client.Result(&pb.Result{})
				//				if err != nil {
				//	log.Printf("Error Requesting Result: %v", err)
				//}
			} else {
				log.Printf("Unknown command, Use 'bid <int>' to place a bid. or 'result' to view current hights bid/final result ")
			}
		}
	}
}

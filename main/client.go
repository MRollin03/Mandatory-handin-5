package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"mandatory-handin-5/pb"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
var Nnodes = 3; //Number of server nodes
var ID = 0;

func main(){
	client()
}

func client() {
	ConnectedClients := make([]pb.ActionhouseClient, Nnodes);

	for i := 0; i < Nnodes; i++ {

		conn, err := grpc.Dial("localhost:" + strconv.Itoa(5001 + i), grpc.WithTransportCredentials(insecure.NewCredentials()))
		client := pb.NewActionhouseClient(conn)

		ConnectedClients[i] = client;

		if err != nil {
			//log.Fatalf("Connection failure")
		}
		defer conn.Close()
	}
	
	fmt.Println("Wilkommen to ze bidding hall.")
	

	CommandlineInput(ConnectedClients)

	fmt.Printf("RUN fOR dA HILLS!")
}

func setbid(clients []pb.ActionhouseClient, bidVal int32) {
	var errors int32
	var successes int32

	done := make(chan struct{}, Nnodes)

	for i := 0; i < Nnodes; i++ {
		go func(client pb.ActionhouseClient, nodeID int) {
			defer func() { done <- struct{}{} }() 

			res, err := client.Bid(context.Background(), &pb.Request{
				Bid:    bidVal,
				UserId: int32(ID), // User ID is assigned by the server
			})
			if err != nil {
				log.Printf("Error from node %d: %v\n", nodeID, err)
				errors++
				return
			}
			successes++
			ID = int(res.UserId)
		}(clients[i], i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < Nnodes; i++ {
		<-done
	}

	// Determine the final outcome
	if successes > 0 {
		fmt.Println("Bid placed!")
	} else {
		fmt.Println("Nodes failed.")
	}
}


func getresult(clients []pb.ActionhouseClient) {
	var errors int32
	var successes int32
	var result string
	var mu sync.Mutex // Mutex to ensure safe access to shared variables

	done := make(chan struct{}, Nnodes)

	for i := 0; i < Nnodes; i++ {
		go func(client pb.ActionhouseClient, nodeID int) {
			defer func() { done <- struct{}{} }() // Signal completion

			res, err := client.Result(context.Background(), &pb.Empty{})
			mu.Lock()
			defer mu.Unlock() // Ensure thread-safe access to shared variables

			if err != nil {
				log.Printf("Error fetching results from node %d: %v", nodeID, err)
				errors++
				return
			}

			successes++
			if res.UserId == -1 {
				result = fmt.Sprintf("An auction is ongoing & the current highest bid is: %d\n", res.Bid)
			} else {
				if res.Bid == 0 {
					result = "No winners, no bids were placed.\n"
				} else {
					result = fmt.Sprintf("An auction has ended, the winning bid was %d from user %d\n", res.Bid, res.UserId)
				}
			}
		}(clients[i], i)
	}

	// Wait for all goroutines to finish
	for i := 0; i < Nnodes; i++ {
		<-done
	}

	// Determine the final outcome
	if successes > 0 {
		fmt.Print(result)
	} else {
		fmt.Println("Nodes failed.")
	}
}


func CommandlineInput(clients []pb.ActionhouseClient) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter command: ")
	for {
		if scanner.Scan() {
			scanText := scanner.Text()
			text := strings.Fields(scanText) // Use Fields to handle spaces correctly
			if len(text) > 1 && text[0] == "bid" {
				bidVal, err := strconv.Atoi(text[1])
				if err != nil || bidVal <= 0 {
					log.Printf("Cant bid 0 or below: %v\n", err)
					continue
				}
				setbid(clients, int32(bidVal))
				if err != nil {
					log.Print("Error Bidding: %v\n", err)
				}
			}else if len(text) == 1 && text[0] == "result" {
				getresult(clients);
			} else {
				log.Println("Unknown command, Use 'bid <int>' to place a bid. or 'result' to view current hights bid/final result ")
			}
		}
		//panic("error ignored")
		time.Sleep(time.Millisecond * 10);
		fmt.Print("Enter command: ")
	}
}

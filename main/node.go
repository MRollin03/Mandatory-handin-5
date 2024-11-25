package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"mandatory-handin-5/pb"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
)
var (
	NodePort          = flag.Int("port", 5001, "The server port, we append it to 'localhost:' later in the program")
	NodeCrash         = flag.Bool("crash", false, "Will the application simulate a crash when called?")
	CID               = int32(1) // Client ID = Bidder ID, calling it BID would be cool, but confusing.
) 
// this means that when we launch the app, we can add ' -port=5002 -crash=True'
// though crash should not be set frequently on launch!
type Node struct {
	pb.UnimplementedActionhouseServer

	nodeID          int
	address         string
	timer           int
	bid             int32
	highestbidder   int32
	ongoingAuction  bool
	mu              sync.Mutex // Mutex for thread-safe operations
}

func main(){
	flag.Parse()

	SetupNode(1, *NodePort)
	
}

func SetupNode(nodeid int, addr int){
	log.Println("trying to setup " + strconv.Itoa(addr)+ "..." )
	
	n := &Node{
		nodeID:         0,
		address:        fmt.Sprintf("localhost:%d", addr),
		bid:            0,
		highestbidder:  0,
		ongoingAuction: false,
	}

	lis, err := net.Listen("tcp", n.address)
	if err != nil{
		log.Fatalf("Server was unable to start/setup")
	}


	// Create a gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterActionhouseServer(grpcServer, n)

	log.Println("Server is running on: " + n.address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


func (n *Node) Bid(ctx context.Context, req *pb.Request) (*pb.Ack, error){
	if (req.UserId == int32(0)){
		req.UserId = CID
		CID++
		// checks if it is a first time user, and gives 'em an ID if that's the case.

	}
	if(*NodeCrash){	
		log.Println("Sending error >:|")
		os.Exit(418)
		return &pb.Ack{},errors.New("Failed to place bid")
	}
	
	if !n.ongoingAuction{
		n.startAuction(25)
		return &pb.Ack{
			Msg: "No ongoing auctions, starting one now",
			UserId: req.UserId, 
		}, nil
	}

	if n.isBidBigger(req){
		n.bid = req.Bid;
		n.highestbidder = req.UserId;
		n.resetAuctionTimer();
		return &pb.Ack{ Msg: "Bid Successful!", UserId: req.UserId,}, nil 
	}

	return &pb.Ack{
		Msg: "Your bid is lower than current bid" , UserId: req.UserId,}, nil
}
/*
func (n *Node) Status() (*pb.Current){
	return *pb.Current{
		bid: n.bid,
	};
}*/

func (n *Node) startAuction(duration int) {
	n.ongoingAuction = true
	n.timer = duration
	n.bid = 0
	n.highestbidder = 0

	go n.auctionTimer()
}

func (n *Node) auctionTimer() {
	for n.timer > 0 {
		time.Sleep(1 * time.Second)
		n.timer--

		if(n.timer == 1){log.Println("1 Second left")}
		if(n.timer == 2){log.Println("2 Seconds left")}
		if(n.timer == 3){log.Println("3 Seconds left")}
		if(n.timer == 4){log.Println("4 Seconds left")}
		if(n.timer == 5){log.Println("5 Seconds left")}
		if(n.timer == 10){log.Println("10 Seconds left")}
		if(n.timer == 20){log.Println("20 Seconds left")}
		if(n.timer == 50){log.Println("50 Seconds left")}

	}
	n.endAuction()
}

// End the auction
func (n *Node) endAuction() {
	n.ongoingAuction = false
	log.Println(fmt.Sprintf("Auction ended.\n Winning bid: %d by user %d\n", n.bid, n.highestbidder))
}

// Reset auction timer
func (n *Node) resetAuctionTimer() {
	if (n.timer <25){n.timer = 25} // Reset to 25 seconds or any desired duration
}

func (n *Node) Result(ctx context.Context, req *pb.Empty) ( *pb.Outcome, error) {
	
	if(*NodeCrash){	
		log.Println("Sending error >:|")
		os.Exit(418)
		return &pb.Outcome{},errors.New("Failed to fetch result ")
	}
	if (n.ongoingAuction){
		return &pb.Outcome{
			Bid: n.bid,
			UserId: -1,
		}, nil
	}
	return &pb.Outcome{
		Bid: n.bid,
		UserId: n.highestbidder,
	}, nil
}

func (n *Node) isBidBigger(req *pb.Request) (bool){

	value := req.Bid;
	if value > n.bid{
		return true;
	}
	return false;
}
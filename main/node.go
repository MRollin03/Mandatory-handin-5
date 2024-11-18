package main

import (
	"flag"
	"fmt"
	"mandatory-handin-5/pb"
	"os"
	"time"
)
var (
	NodePort          = flag.Int("port", 5001, "The server port, we append it to 'localhost:' later in the program")
	NodeCrash         = flag.Bool("crash", false, "Will the application simulate a crash when called?")
) 
// this means that when we launch the app, we can add ' -port= 5002 -crash= True'
// though crash should not be set frequently on launch!
type Node struct{
	nodeID int
	Address string
	timer int
	bid int32
	highestbidder string
	ongoingAuction bool
	clients pb.ActionhouseClient
}

func main(){
	flag.Parse()
	if *NodeCrash == true{
		fmt.Println("crash")
		os.Exit(418)
		}// this should be moved further in, so the program crashed when receive a call instead
	}


func (n *Node) Bid(req *pb.Request) (ack *pb.acknowledge){
	if !n.ongoingAuction{
		return "No ongoing auctions"
	}

	if (pb.Request.UserId != ){

	}
	
	if n.isBidBigger(req){
		n.bid = req.bid;
		n.highestbidder = req.clientname;
		n.auctionTimer();
		return "Bid Successful" 
	}

	return "Your bid is lower than current bid" 
}

func (n *Node) Status() (*pb.Current){
	return *pb.Current{
		bid: n.bid,
	};
}

func (n *Node) startAuction(){
	go n.auctionTimer();
	defer n.Result();
	n.ongoingAuction = false;
}

func (n *Node) auctionTimer(){		//TODO: we might need an "updateAuctionTimer" func instead of calling this mulitple times.
	n.ongoingAuction = true;
	n.timer = 25;		//TODO: review time
	time.Sleep(time.Duration(n.timer))
}


func (n *Node) Result() (){

	return *pb.Outcome{
		bid: n.bid,
		UserId: n.highestbidder,
	}
	
}

func (n *Node) isBidBigger(req *pb.Request) (bool){

	value := req.Bid;
	if value > n.bid{
		return true;
	}
	return false;
}
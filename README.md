# Mandatory-handin-5

# Project Setup and Running Instructions

## 1. Setup

To get started with this project, follow these steps:

1. **Clone the Repository**

   ```bash
   git clone https://github.com/MRollin03/Mandatory-handin-5
   cd Mandatory-handin-5
   ```

2. import dependencies

- `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
- follow next steps

## 2. ''How to run the nodes''

To run the program, execute the following commands:

Navigate to the Main Directory

- cd from root
- `cd main`

- start 3 nodes you using:
- `go run node.go` -port=[insert port number here] e.g. `go node.go -port=5002`

default port is 5001, meaning you can omit the flag that node.
Our clients expect 3 servers with the address 5001, 5002, 5003.

## 3. ''How to run the client''

before starting a client, please ensure that all servers are up and running.

- run `go run client.go`
- available commands are:
  - `bid [bid amount]`
  - `result`
    The first called bid command will start the auction. This bid will not be placed
    as the auction was not ongoing then it was placed.
    This is works as a "start" command. So if the nodes tell you no "No ongoing auctions, starting one now"
    your bid command was in fact a start command, and no bid has been placed.

You can run as many clients as you want.

package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/joho/godotenv"
)

var logger *log.Logger
var servers string
var server_list []string

func main() {
	logger = log.Default()
	if err := godotenv.Load(); err != nil {
		logger.Fatalf("Error loading .env: %v\n", err)
	}

	address := os.Getenv("MEDIVPN_LISTENING_ADDRESS") // TODO: `tailscale ip` for dynamic assignment
	servers = os.Getenv("MEDIVPN_SERVERS")
	server_list = strings.Split(servers, ",")

	listener, err := net.Listen("tcp4", address)
	if err != nil {
		logger.Fatalf("Error starting server: %v\n", err)
	}
	defer listener.Close()
	logger.Printf("Server listening on %v\n", address)

	for {
		connection, err := listener.Accept()
		if err != nil {
			os.Exit(1)
		}
		go handleConnection(connection)
	}
}

func handleConnection(client_connection net.Conn) {
	defer client_connection.Close()

	client_address := client_connection.RemoteAddr()
	logger.Printf("Connection established with %v\n", client_address)
	defer logger.Printf("Connection closed with %v\n", client_address)

	scanner := bufio.NewScanner(client_connection)

	for scanner.Scan() {
		message := scanner.Text()
		logger.Printf("%v: received '%v'\n", client_address, message)

		args := strings.Split(message, " ")
		command := args[0]

		switch command {
		case "server":
			serverHandler(client_connection, client_address, args[1:])
		default:
			client_connection.Write([]byte("invalid\n"))
		}
	}
}

func serverHandler(connection net.Conn, address net.Addr, args []string) {
	logger.Printf("%v >> serverHandler()", address)

	if len(args) != 1 {
		logger.Printf("%v: invalid or missing args (%v); sending server list to client\n", address, args)
		connection.Write([]byte(servers + "\n"))
		return
	}

	server := args[0]

	if !slices.Contains(server_list, server) {
		logger.Printf("%v: server not found; sending server list to client\n", address)
		connection.Write([]byte(servers + "\n"))
		return
	}

	for _, s := range server_list {
		exec.Command("/usr/bin/wg-quick", "down", s).Run()
	}
	logger.Printf("dropped all interfaces\n")

	err := exec.Command("/usr/bin/wg-quick", "up", server).Run()
	if err != nil {
		logger.Printf("err: %v\n", err)
	}
}

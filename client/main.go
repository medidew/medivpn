package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/joho/godotenv"
)

const usage string = `
medivpn - Personal vpn configuration tool, because I'm sick of typing out the longer command.

Usage:

	medivpn <command>

Commands:

	on		sets tailscale exit node according to config
	off		unsets tailscale exit node
	server	lists the possible exit node servers
`

var tailscale_binary_path string

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%v\n", usage)
		os.Exit(0)
	}

	if err := godotenv.Load(); err != nil {
		err := godotenv.Load("/etc/medivpn/client.env")

		if err != nil {
			fmt.Printf("Error loading .env: %v\n", err)
			os.Exit(1)
		}
	}
	server_ip := os.Getenv("MEDIVPN_SERVER_ADDRESS")
	server_port := os.Getenv("MEDIVPN_SERVER_PORT")
	tailscale_binary_path = os.Getenv("TAILSCALE_BINARY_PATH")

	command := os.Args[1]

	var err error
	switch command {
	case "on":
		err = setExitNode(server_ip)
	case "off":
		err = setExitNode("")
	case "server":
		err = serverHandler(server_ip + ":" + server_port)
	default:
		fmt.Printf("%v\n", usage)
	}

	if err != nil {
		fmt.Printf("Error setting exit node: %v\n", err)
		os.Exit(1)
	}
}

func setExitNode(ipaddr string) error {
	exit_node_arg := fmt.Sprintf("--exit-node=%v", ipaddr)
	cmd := exec.Command(tailscale_binary_path, "set", exit_node_arg)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

const server_usage string = `
Usage:

	medivpn server <server>

Servers:

`

func serverHandler(address string) error {
	connection, err := net.Dial("tcp4", address)
	if err != nil {
		return err
	}
	defer connection.Close()

	connection.Write([]byte("server\n"))
	scanner := bufio.NewScanner(connection)
	scanner.Scan()
	response := scanner.Text()
	server_list := strings.Split(response, ",")

	server_string := server_usage
	for _, server := range server_list {
		server_string += "    " + server + "\n"
	}

	if len(os.Args) != 3 {
		fmt.Printf("%v\n", server_string)
		return nil
	}

	server := os.Args[2]

	if !slices.Contains(server_list, server) {
		fmt.Printf("%v\n", server_string)
		return nil
	}

	if err := changeServer(connection, server); err != nil {
		return err
	}

	return nil
}

func changeServer(connection net.Conn, server string) error {
	connection.Write([]byte("server " + server + "\n"))

	//scanner := bufio.NewScanner(connection)
	//scanner.Scan()
	//response := scanner.Text()

	return nil
}

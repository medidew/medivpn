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

const usage string = `medivpn - Personal vpn configuration tool, because I'm sick of typing out the longer command.

Usage:

	medivpn <command>

Commands:

	on		sets tailscale exit node according to config
	off		unsets tailscale exit node
	server	lists the possible exit node servers
	lan		controls the Tailscale exit node lan access property
`

// TODO: would be neat if I could actually query the current vpn server/status through this CLI

var tailscale_binary_path string

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%v", usage)
		os.Exit(0)
	}

	if err := godotenv.Load(); err != nil {
		err := godotenv.Load("/etc/medivpn/client.env") // TODO: replace env with .conf because convention

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
	case "lan":
		err = lanAccess()
	default:
		fmt.Printf("%v", usage)
	}

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func setExitNode(ipaddr string) error {
	exit_node_arg := fmt.Sprintf("--exit-node=%v", ipaddr)
	cmd := exec.Command(tailscale_binary_path, "set", exit_node_arg)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error setting exit node: %v", err)
	}
	return nil
}

const server_usage string = `Usage:

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
	server_list := strings.Split(response, ",") // TODO: this check doesn't really need to exist when the server already handles it

	server_string := server_usage
	for _, server := range server_list {
		server_string += "    " + server + "\n"
	}

	if len(os.Args) != 3 {
		fmt.Printf("%v", server_string)
		return nil
	}

	server := os.Args[2]

	if !slices.Contains(server_list, server) {
		fmt.Printf("%v", server_string)
		return nil
	}

	if err := changeServer(connection, server); err != nil {
		return fmt.Errorf("error switching server: %v", err)
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

const lan_usage string = `Usage:

	medivpn lan <command>

Commands:

	on		allows exit node LAN access through Tailscale
	off		forbids exit node LAN access through Tailscale
`

func lanAccess() error {
	if len(os.Args) != 3 {
		fmt.Printf("%v", lan_usage)
		return nil
	}

	arg := os.Args[2]

	var err error
	switch arg {
	case "on":
		cmd := exec.Command(tailscale_binary_path, "set", "--exit-node-allow-lan-access")
		err = cmd.Run()
	case "off":
		cmd := exec.Command(tailscale_binary_path, "set", "--exit-node-allow-lan-access=false")
		err = cmd.Run()
	default:
		fmt.Printf("%v", lan_usage)
	}

	return err
}

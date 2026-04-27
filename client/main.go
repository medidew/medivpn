package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

const usage string = `
medivpn - Personal vpn configuration tool, because I'm sick of typing out the longer command.

Usage:

	medivpn <on|off>

Commands:

	on		sets tailscale exit node according to config
	off		unsets tailscale exit node
`

var tailscale_binary_path string

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("%v\n", usage)
		os.Exit(0)
	}

	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env: %v\n", err)
	}
	exit_node_ip := os.Getenv("MEDIVPN_EXIT_NODE")
	tailscale_binary_path = os.Getenv("TAILSCALE_BINARY_PATH")

	command := os.Args[1]

	var err error
	switch command {
	case "on":
		err = setExitNode(exit_node_ip)
	case "off":
		err = setExitNode("")
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

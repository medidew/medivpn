package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"slices"
	"strings"
	"time"

	"medidew.name/medivpn-client/internal/config"
	"tailscale.com/client/local"
	"tailscale.com/ipn"
	"tailscale.com/tailcfg"
)

const usage string = `medivpn - Personal vpn configuration tool, because I'm sick of typing out the longer command.

Usage:

	medivpn <command>

Commands:

	on		sets tailscale exit node according to config
	off		unsets tailscale exit node
	server		lists the possible exit node servers
	lan		controls the Tailscale exit node lan access property
`

// TODO: would be neat if I could actually query the current vpn server/status through this CLI

var tailscale_client *local.Client
var cfg *config.Config
var ctx context.Context

func main() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Printf("%v", usage)
		os.Exit(0)
	}

	var stop context.CancelFunc
	ctx, stop = signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tailscale_client = &local.Client{}

	command := os.Args[1]
	switch command {
	case "on":
		mp := &ipn.MaskedPrefs{
			ExitNodeIDSet: true,
			Prefs: ipn.Prefs{
				ExitNodeID: tailcfg.StableNodeID(cfg.Server_id),
			},
		}
		_, err = tailscale_client.EditPrefs(ctx, mp)
	case "off":
		err = tailscale_client.SetUseExitNode(ctx, false)
	case "server":
		err = serverHandler(cfg.Server_id + ":" + cfg.Server_port)
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

	if _, err = connection.Write([]byte("server\n")); err != nil {
		return fmt.Errorf("failed to write to server: %v", err)
	}
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
	if _, err := connection.Write([]byte("server " + server + "\n")); err != nil {
		return fmt.Errorf("failed to write to server: %v", err)
	}

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

	mp := &ipn.MaskedPrefs{}

	arg := os.Args[2]
	var err error
	switch arg {
	case "on":
		mp.ExitNodeAllowLANAccess = true
		_, err = tailscale_client.EditPrefs(ctx, mp)
	case "off":
		mp.ExitNodeAllowLANAccess = false
		_, err = tailscale_client.EditPrefs(ctx, mp)
	default:
		fmt.Printf("%v", lan_usage)
	}

	return err
}

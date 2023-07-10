package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/statechannels/go-nitro/internal/chain"
	"github.com/statechannels/go-nitro/internal/utils"
	"github.com/statechannels/go-nitro/types"
)

type participant string

const (
	alice participant = "alice"
	bob   participant = "bob"
	irene participant = "irene"
	ivan  participant = "ivan"
)

type color string

const (
	black   color = "[30m"
	red     color = "[31m"
	green   color = "[32m"
	yellow  color = "[33m"
	blue    color = "[34m"
	magenta color = "[35m"
	cyan    color = "[36m"
	white   color = "[37m"
	gray    color = "[90m"
)

var (
	participants     = []participant{alice, bob, irene, ivan}
	participantColor = map[participant]color{alice: blue, irene: green, ivan: cyan, bob: yellow}
)

const (
	FUNDED_TEST_PK = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	// TODO: This token must be generated by running `lotus auth api-info --perm admin
	CHAIN_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.iKUOdDak2Gwhyh_neCukKS-nsJjy41XdUTTM7vZRKTU"
	// This is the RPC endpoint for the local filecoin devnet
	FILECOIN_DEVNET_URL = "ws://127.0.0.1:1234/rpc/v1"
)

func main() {
	running := []*exec.Cmd{}

	naAddress, vpaAddress, caAddress, err := chain.DeployContracts(context.Background(), FILECOIN_DEVNET_URL, CHAIN_TOKEN, FUNDED_TEST_PK)
	if err != nil {
		utils.StopCommands(running...)
		panic(err)
	}

	for _, p := range participants {
		client, err := setupRPCServer(p, participantColor[p], naAddress, vpaAddress, caAddress)
		if err != nil {
			utils.StopCommands(running...)
			panic(err)
		}
		running = append(running, client)
	}

	waitForKillSignal()

	utils.StopCommands(running...)
}

// waitForKillSignal blocks until we receive a kill or interrupt signal
func waitForKillSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	fmt.Printf("Received signal %s, exiting..\n", sig)
}

// setupRPCServer starts up an RPC server for the given participant
func setupRPCServer(p participant, c color, na, vpa, ca types.Address) (*exec.Cmd, error) {
	args := []string{"run", ".", "-naaddress", na.String()}
	args = append(args, "-vpaaddress", vpa.String())
	args = append(args, "-caaddress", ca.String())
	args = append(args, "-config", fmt.Sprintf("./scripts/test-configs/%s.toml", p))

	cmd := exec.Command("go", args...)
	cmd.Stdout = newColorWriter(c, os.Stdout)
	cmd.Stderr = newColorWriter(c, os.Stderr)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// colorWriter is a writer that writes to the underlying writer with the given color
type colorWriter struct {
	writer io.Writer
	color  color
}

func (cw colorWriter) Write(p []byte) (n int, err error) {
	_, err = cw.writer.Write([]byte("\033" + string(cw.color) + string(p) + "\033[0m"))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// newColorWriter creates a writer that colors the output with the given color
func newColorWriter(c color, w io.Writer) colorWriter {
	return colorWriter{
		writer: w,
		color:  c,
	}
}

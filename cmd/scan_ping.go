package cmd

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"

	"github.com/ayanrajpoot10/bugscanx-go/pkg/queuescanner"
)

var pingCmd = &cobra.Command{
	Use:     "ping",
	Short:   "Scan hosts using TCP ping.",
	Example: "  bugscanx-go ping -f hosts.txt\n  bugscanx-go ping -f hosts.txt --port 443 --timeout 5",
	Run:     pingRun,
}

var (
	pingFlagFilename string
	pingFlagTimeout  int
	pingFlagOutput   string
	pingFlagPort     int
)

func init() {
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().StringVarP(&pingFlagFilename, "filename", "f", "", "domain list filename")
	pingCmd.Flags().IntVar(&pingFlagTimeout, "timeout", 2, "timeout in seconds")
	pingCmd.Flags().StringVarP(&pingFlagOutput, "output", "o", "", "output result")
	pingCmd.Flags().IntVar(&pingFlagPort, "port", 80, "port to use")
	
	pingCmd.MarkFlagRequired("filename")
}

func pingHost(ctx *queuescanner.Ctx, params *queuescanner.QueueScannerScanParams) {
	host := params.Data.(string)
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", pingFlagPort)), time.Duration(pingFlagTimeout)*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	remoteAddr := conn.RemoteAddr()
	ip, _, err := net.SplitHostPort(remoteAddr.String())
	if err != nil {
		ip = remoteAddr.String()
	}

	formatted := fmt.Sprintf("%-16s %-20s", ip, host)
	ctx.ScanSuccess(formatted)
	ctx.Log(formatted)
}

func pingRun(cmd *cobra.Command, args []string) {
	hosts, err := ReadLinesFromFile(pingFlagFilename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	fmt.Printf("%-16s %-20s\n", "IP Address", "Host")
	fmt.Printf("%-16s %-20s\n", "----------", "----")

	scanner := queuescanner.NewQueueScanner(globalFlagThreads, pingHost)
	for _, host := range hosts {
		scanner.Add(&queuescanner.QueueScannerScanParams{Name: host, Data: host})
	}

	scanner.SetOutputFile(pingFlagOutput)
	scanner.Start()
}

package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	host    string
	startP  int
	endP    int
	timeout time.Duration
	workers int
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "portscanner",
		Short:   "Escaneia portas TCP em um host",
		Version: "1.0.1",
	}

	scanCmd := &cobra.Command{
		Use:   "scan",
		Short: "Executa o scan de portas",
		RunE:  runScan,
	}

	scanCmd.Flags().StringVarP(&host, "host", "H", "scanme.nmap.org", "Host a escanear (IP ou domínio)")
	scanCmd.Flags().IntVarP(&startP, "start", "s", 1, "Porta inicial")
	scanCmd.Flags().IntVarP(&endP, "end", "e", 1024, "Porta final")
	scanCmd.Flags().DurationVarP(&timeout, "timeout", "t", time.Second, "Timeout para cada tentativa")
	scanCmd.Flags().IntVarP(&workers, "workers", "w", 100, "Número de goroutines concorrentes")

	rootCmd.AddCommand(scanCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func runScan(cmd *cobra.Command, args []string) error {
	ports := make(chan int, workers)
	results := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range ports {
				addr := net.JoinHostPort(host, fmt.Sprintf("%d", p))
				conn, err := net.DialTimeout("tcp", addr, timeout)
				if err == nil {
					conn.Close()
					results <- p
				}
			}
		}()
	}

	go func() {
		for p := startP; p <= endP; p++ {
			ports <- p
		}
		close(ports)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	startTime := time.Now()
	var openPorts []int
	for p := range results {
		openPorts = append(openPorts, p)
	}

	sort.Ints(openPorts)
	duration := time.Since(startTime)

	if len(openPorts) == 0 {
		fmt.Printf("Nenhuma porta aberta encontrada em %s\n", host)
	} else {
		fmt.Printf("Portas abertas em %s:\n", host)
		for _, p := range openPorts {
			fmt.Printf("  - %d\n", p)
		}
	}

	fmt.Printf("Scan concluído em %s\n", duration)
	return nil
}

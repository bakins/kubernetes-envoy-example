package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bakins/kubernetes-envoy-example/frontend"
	"github.com/spf13/cobra"
)

var address = ":8080"

var rootCmd = &cobra.Command{
	Use:   "frontend",
	Short: "simple HTTP frontend",
	Run:   runServer,
}

func main() {
	f := rootCmd.Flags()
	f.StringVarP(&address, "address", "a", address, "listening address")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {

	s, err := frontend.New(frontend.SetAddress(address))

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.Stop()
	}()

	if err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

}

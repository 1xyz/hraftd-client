package cmd

import (
	"fmt"
	"github.com/1xyz/hraftd-client/client"
	"github.com/1xyz/hraftd-client/config"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

const version = "1.0.0"

func NewCmdRoot(cfg *config.Config) *cobra.Command {
	var serverURLs = ""

	var rootCmd = &cobra.Command{
		Version: version,
		Use:     "client for hraftd",
		Short:   "Work seamlessly with the aardy service from the command line",
	}

	var putCmd = &cobra.Command{
		Use:   "put",
		Short: "Put the specified key and value",
		Args:  MinimumArgs(2, ""),
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverURLs == "" {
				serverURLs = cfg.URL
			}
			urls := parseServerURLs(serverURLs)
			log.Printf("Using url %v from %v\n", urls[0], serverURLs)
			cli := client.NewHttpClient(urls[0], cfg.Timeout)
			if err := cli.Put(args[0], args[1]); err != nil {
				log.Printf("Put key = %v error = %v\n", args[0], err)
				return err
			}
			return nil
		},
	}

	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get the specified key and value",
		Args:  MinimumArgs(1, ""),
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverURLs == "" {
				serverURLs = cfg.URL
			}
			urls := parseServerURLs(serverURLs)
			log.Printf("Using url %v from %v\n", urls[0], serverURLs)
			cli := client.NewHttpClient(urls[0], cfg.Timeout)
			value, err := cli.Get(args[0])
			if err != nil {
				log.Printf("Put key = %v error = %v\n", args[0], err)
				return err
			}
			log.Printf("value = %v\n", value)
			return nil
		},
	}

	var loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Run a simple load test",
		Args:  MinimumArgs(1, ""),
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverURLs == "" {
				serverURLs = cfg.URL
			}
			leaderURL, err := getLeaderURL(serverURLs, cfg.Timeout)
			if err != nil {
				return err
			}

			duration, err := time.ParseDuration(args[0])
			if err != nil {
				return err
			}

			cli := client.NewHttpClient(leaderURL, cfg.Timeout)
			return client.RunLoadTest(cli, duration)
		},
	}

	var serverInfoCmd = &cobra.Command{
		Use:   "server-info",
		Short: "Query a server info",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if serverURLs == "" {
				serverURLs = cfg.URL
			}
			urls := parseServerURLs(serverURLs)

			for _, url := range urls {
				cli := client.NewHttpClient(url, cfg.Timeout)
				info, err := cli.GetInfo()
				if err != nil {
					return err
				}

				fmt.Printf("URL %v -- Leader %v State %v\n", url, info.Leader, info.State)
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&serverURLs, "serverURLs", "s",
		"http://localhost:11001", "Server URL of the master")
	rootCmd.AddCommand(putCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(serverInfoCmd)
	return rootCmd
}

func Execute(cfg *config.Config) error {
	rootCmd := NewCmdRoot(cfg)
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}

func MinimumArgs(n int, msg string) cobra.PositionalArgs {
	if msg == "" {
		return cobra.MinimumNArgs(1)
	}

	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			return fmt.Errorf("number of args %d is less than expected(%d)", len(args), n)
		}
		return nil
	}
}

func parseServerURLs(serverURLs string) []string  {
	tokens := strings.Split(serverURLs, ",")
	result := make([]string, 0)
	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if len(tok) == 0 {
			continue
		}
		result = append(result, tok)
	}
	return result
}

func getLeaderURL(serverURLs string, timeout time.Duration) (string, error)  {
	urls := parseServerURLs(serverURLs)
	for _, url := range urls {
		cli := client.NewHttpClient(url, timeout)
		info, err := cli.GetInfo()
		if err != nil {
			log.Printf("warn URL %v error %v", url, err)
		}

		fmt.Printf("URL %v -- Leader %v State %v\n", url, info.Leader, info.State)
		if info.State == "Leader" {
			return url, nil
		}
	}

	return "", fmt.Errorf("leader URL not found")
}

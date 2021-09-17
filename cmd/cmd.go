package cmd

import (
	"fmt"
	"github.com/1xyz/hraftd-client/client"
	"github.com/1xyz/hraftd-client/config"
	"github.com/spf13/cobra"
	"log"
	"time"
)

const version = "1.0.0"

func NewCmdRoot(cfg *config.Config) *cobra.Command {
	var serverURL = ""

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
			if serverURL == "" {
				serverURL = cfg.URL
			}
			cli := client.NewHttpClient(serverURL, cfg.Timeout)
			if err := cli.Put(args[0], args[1]); err != nil {
				log.Printf("Put key = %v error = %w", args[0], err)
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
			if serverURL == "" {
				serverURL = cfg.URL
			}
			cli := client.NewHttpClient(serverURL, cfg.Timeout)
			value, err := cli.Get(args[0])
			if err != nil {
				log.Printf("Put key = %v error = %w", args[0], err)
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
			if serverURL == "" {
				serverURL = cfg.URL
			}

			duration, err := time.ParseDuration(args[0])
			if err != nil {
				return err
			}

			cli := client.NewHttpClient(serverURL, cfg.Timeout)
			return client.RunLoadTest(cli, duration)
		},
	}

	rootCmd.PersistentFlags().StringVarP(&serverURL, "serverURL", "s",
		"http://localhost:11001", "Server URL of the master")
	rootCmd.AddCommand(putCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(loadCmd)
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

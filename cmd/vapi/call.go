package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chriscow/vapi-go"
	"github.com/spf13/cobra"
)

func newCallCmd() *cobra.Command {
	var outFile string

	cmd := &cobra.Command{
		Use:   "call [call-id]",
		Short: "Get call details",
		Long:  "Retrieve details for a specific call by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			callID := args[0]
			call, err := vapi.GetCall(context.Background(), callID)
			if err != nil {
				return fmt.Errorf("failed to get call: %w", err)
			}

			// Marshal the call with indentation
			data, err := json.MarshalIndent(call, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal call data: %w", err)
			}

			// If outFile is specified, write to file
			if outFile != "" {
				if err := os.WriteFile(outFile, data, 0644); err != nil {
					return fmt.Errorf("failed to write to file: %w", err)
				}
				fmt.Printf("Call data written to %s\n", outFile)
				return nil
			}

			// Otherwise print to stdout
			fmt.Println(string(data))
			return nil
		},
	}

	cmd.Flags().StringVarP(&outFile, "out", "o", "", "Output file path to save the call data")

	return cmd
}

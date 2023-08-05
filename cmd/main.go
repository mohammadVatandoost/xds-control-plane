package main

import (
	"fmt"
	"os"
)

const serviceName = "xds_control_plane"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

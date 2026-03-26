package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"cldap/session"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear the saved session",
	RunE:  runLogout,
}

func runLogout(_ *cobra.Command, _ []string) error {
	if err := session.Clear(); err != nil {
		return err
	}
	fmt.Println("Logged out.")
	return nil
}

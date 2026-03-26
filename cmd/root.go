package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "cldap",
	Short: "LDAP CLI — configure, login, and query your directory",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(groupsCmd)
	rootCmd.AddCommand(logoutCmd)
}

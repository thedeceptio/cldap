package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thedeceptio/cldap/config"
	ldapclient "github.com/thedeceptio/cldap/ldap"
	"github.com/thedeceptio/cldap/session"
)

var emailFlag string

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List security groups for the logged-in user or a given email",
	RunE:  runGroups,
}

func init() {
	groupsCmd.Flags().StringVar(&emailFlag, "email", "", "Email address to look up groups for")
}

func runGroups(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("not configured — run 'cldap configure' first")
	}

	sess, err := session.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Not logged in — run 'cldap login' first")
		os.Exit(1)
	}

	client, err := ldapclient.Connect(cfg)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer client.Close()

	if err := client.Bind(sess.BindDN, sess.Password); err != nil {
		fmt.Fprintln(os.Stderr, "Session expired — run 'cldap login' again")
		os.Exit(1)
	}

	var groups []string
	if emailFlag != "" {
		groups, err = client.GroupsForEmail(emailFlag)
	} else {
		groups, err = client.GroupsForUsername(sess.Username)
	}
	if err != nil {
		return err
	}

	if len(groups) == 0 {
		fmt.Println("No groups found.")
		return nil
	}

	for _, g := range groups {
		fmt.Println(ldapclient.ExtractCN(g))
	}
	return nil
}

package cmd

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/thedeceptio/cldap/config"
	ldapclient "github.com/thedeceptio/cldap/ldap"
	"github.com/thedeceptio/cldap/session"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the LDAP server",
	RunE:  runLogin,
}

func runLogin(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("not configured — run 'cldap configure' first")
	}

	fmt.Print("Username: ")
	var username string
	fmt.Scanln(&username)
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return fmt.Errorf("read password: %w", err)
	}
	password := string(passwordBytes)

	// For AD UPN-style logins (user@domain or userPrincipalName attribute),
	// bind directly with the username — it is already the bind identity.
	// Otherwise construct a DN like: sAMAccountName=user,dc=...
	var bindDN string
	if cfg.UsernameAttr == "userPrincipalName" || strings.Contains(username, "@") {
		bindDN = username
	} else {
		bindDN = fmt.Sprintf("%s=%s,%s", cfg.UsernameAttr, username, cfg.UserSearchBase)
	}

	client, err := ldapclient.Connect(cfg)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer client.Close()

	if err := client.Bind(bindDN, password); err != nil {
		fmt.Fprintln(os.Stderr, "Login failed: invalid credentials")
		os.Exit(1)
	}

	if err := session.Save(&session.Session{
		BindDN:   bindDN,
		Password: password,
		Username: username,
	}); err != nil {
		return fmt.Errorf("save session: %w", err)
	}

	fmt.Printf("Logged in as %s\n", username)
	return nil
}

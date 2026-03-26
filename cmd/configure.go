package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"cldap/config"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Set LDAP endpoint configuration",
	RunE:  runConfigure,
}

func runConfigure(_ *cobra.Command, _ []string) error {
	scanner := bufio.NewScanner(os.Stdin)

	// Load existing config as defaults (ignore error if none exists yet)
	existing, _ := config.Load()
	if existing == nil {
		existing = &config.Config{
			Port:         389,
			UsernameAttr: "uid",
		}
	}

	c := &config.Config{}
	c.Host = prompt(scanner, "LDAP Host", existing.Host)

	portStr := prompt(scanner, "Port", strconv.Itoa(existing.Port))
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid port %q", portStr)
	}
	c.Port = port

	c.BaseDN = prompt(scanner, "Base DN (e.g. dc=example,dc=com)", existing.BaseDN)
	c.UserSearchBase = prompt(scanner, "User Search Base DN (e.g. ou=users,dc=example,dc=com)", existing.UserSearchBase)
	c.UsernameAttr = prompt(scanner, "Username attribute (uid / sAMAccountName)", existing.UsernameAttr)
	c.UseTLS = promptBool(scanner, "Use TLS (LDAPS)", existing.UseTLS)
	if !c.UseTLS {
		c.UseStartTLS = promptBool(scanner, "Use StartTLS", existing.UseStartTLS)
	}

	if err := config.Save(c); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	fmt.Println("Configuration saved.")
	return nil
}

func prompt(scanner *bufio.Scanner, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	scanner.Scan()
	val := strings.TrimSpace(scanner.Text())
	if val == "" {
		return defaultVal
	}
	return val
}

func promptBool(scanner *bufio.Scanner, label string, defaultVal bool) bool {
	d := "n"
	if defaultVal {
		d = "y"
	}
	val := prompt(scanner, label+" (y/n)", d)
	return strings.ToLower(strings.TrimSpace(val)) == "y"
}

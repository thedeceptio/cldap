package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"

	"github.com/thedeceptio/cldap/config"
)

type Client struct {
	conn *ldap.Conn
	cfg  *config.Config
}

func Connect(cfg *config.Config) (*Client, error) {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	var (
		conn *ldap.Conn
		err  error
	)

	if cfg.UseTLS {
		conn, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	} else {
		conn, err = ldap.Dial("tcp", address)
	}
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", address, err)
	}

	if !cfg.UseTLS && cfg.UseStartTLS {
		if err := conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil { //nolint:gosec
			conn.Close()
			return nil, fmt.Errorf("startTLS: %w", err)
		}
	}

	return &Client{conn: conn, cfg: cfg}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Bind(bindDN, password string) error {
	return c.conn.Bind(bindDN, password)
}

// GroupsForUsername returns the memberOf values for the given username.
func (c *Client) GroupsForUsername(username string) ([]string, error) {
	filter := fmt.Sprintf("(%s=%s)", c.cfg.UsernameAttr, ldap.EscapeFilter(username))
	return c.searchGroups(filter)
}

// GroupsForEmail returns the memberOf values for the user with the given email.
func (c *Client) GroupsForEmail(email string) ([]string, error) {
	filter := fmt.Sprintf("(mail=%s)", ldap.EscapeFilter(email))
	return c.searchGroups(filter)
}

func (c *Client) searchGroups(filter string) ([]string, error) {
	req := ldap.NewSearchRequest(
		c.cfg.UserSearchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"memberOf"},
		nil,
	)
	sr, err := c.conn.Search(req)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return sr.Entries[0].GetAttributeValues("memberOf"), nil
}

// ExtractCN returns the CN value from a full DN string.
// e.g. "CN=mygroup,OU=groups,DC=example,DC=com" → "mygroup"
func ExtractCN(dn string) string {
	for _, part := range strings.Split(dn, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToUpper(part), "CN=") {
			return part[3:]
		}
	}
	return dn
}

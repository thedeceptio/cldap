# cldap

A CLI tool for querying LDAP / Active Directory. Configure once, then run queries without repeating connection flags every time.

## Requirements

- Go 1.18 or newer

## Install

```sh
go install github.com/thedeceptio/cldap@latest
```

Or clone and build locally:

```sh
git clone https://github.com/thedeceptio/cldap.git
cd cldap
go install .
```

The binary is placed in `$GOPATH/bin` (typically `~/go/bin`). Make sure that directory is in your `$PATH`.

## Usage

### 1. Configure

```sh
cldap configure
```

Interactive prompts — press Enter to accept the default in brackets:

```
LDAP Host: ldap.example.com
Port [389]:
Base DN (e.g. dc=example,dc=com): dc=example,dc=com
User Search Base DN: dc=example,dc=com
Username attribute (uid / sAMAccountName / userPrincipalName): userPrincipalName
Use TLS (LDAPS) (y/n) [n]:
Use StartTLS (y/n) [n]:
```

Config is saved to `~/.cldap/config.json` (mode 0600).

### 2. Login

```sh
cldap login
```

Prompts for username and password. Session is saved to `~/.cldap/session.json` (mode 0600).

For Active Directory, enter your UPN (`user@domain.com`) or `sAMAccountName` depending on how you configured the username attribute.

### 3. List groups

```sh
# groups for yourself
cldap groups

# groups for another user by email
cldap groups --email colleague@example.com
```

### 4. Logout

```sh
cldap logout
```

Removes the saved session.

## Notes

- Credentials are stored locally in `~/.cldap/` with restrictive permissions (0600). Keep that directory secure.
- TLS certificate verification is skipped by default (matching the behaviour of most internal AD tools). Do not use over untrusted networks without a proper cert.

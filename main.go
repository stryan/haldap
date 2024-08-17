package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/charmbracelet/log"
	ldap "github.com/go-ldap/ldap/v3"
	"github.com/kelseyhightower/envconfig"
)

type User struct {
	Name      string
	Group     string
	LocalOnly bool
}

type Config struct {
	URL      string `envconfig:"HALDAP_URL"`
	BaseDN   string `envconfig:"HALDAP_BASEDN"`
	UserAttr string `envconfig:"HALDAP_USERATTR"`
}

func (u *User) String() string {
	return fmt.Sprintf("name = %v\ngroup = %v\nlocal_only = %v", u.Name, u.Group, u.LocalOnly)
}

func main() {
	logger := log.New(os.Stderr)
	var c Config
	err := envconfig.Process("haldap", &c)
	if err != nil {
		logger.Fatal(err)
	}
	if c.URL == "" {
		// try getting flags
		urlPtr := flag.String("url", "", "LDAP server URL with scheme")
		basePtr := flag.String("basedn", "", "LDAP server base DN")
		attrPtr := flag.String("userattr", "", "LDAP user attribute")
		flag.Parse()
		if urlPtr == nil || basePtr == nil || attrPtr == nil {
			logger.Fatal("all config options failed")
		}
		c.URL = *urlPtr
		c.BaseDN = *basePtr
		c.UserAttr = *attrPtr
	}
	uname, ok := os.LookupEnv("username")
	if !ok {
		logger.Warn("couldn't get username")
		return
	}
	pwd, ok := os.LookupEnv("password")
	if !ok {
		logger.Warn("couldn't get password")
		return
	}
	userdn := fmt.Sprintf("%v=%v,%v", c.UserAttr, uname, c.BaseDN)

	l, err := ldap.DialURL(c.URL)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() { _ = l.Close() }()
	err = l.Bind(userdn, pwd)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("bind successful")
	u := &User{}
	result, err := l.Search(ldap.NewSearchRequest(
		userdn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", c.UserAttr, uname),
		[]string{"displayName", "memberOf"},
		nil,
	))
	if err != nil {
		logger.Fatal(err)
	}
	if len(result.Entries) > 1 {
		logger.Fatal("expected one result got many")
	}
	entry := result.Entries[0]
	groups := entry.GetAttributeValues("memberOf")
	re := regexp.MustCompile(`cn=([0-9-a-zA-Z]+),`)
	for _, group := range groups {
		g := re.FindStringSubmatch(group)[1]
		if g == "hausers" {
			u.Group = "system-users"
		} else if g == "haadmins" {
			u.Group = "system-admin"
		}
	}
	if u.Group == "" {
		logger.Info("user not enabled")
		return
	}
	u.Name = entry.GetAttributeValue("displayName")
	u.LocalOnly = false
	fmt.Printf("%v\n", u)
}

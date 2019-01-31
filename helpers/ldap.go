package helpers

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/ldap.v3"
)

type LDAP struct {
	L	*ldap.Conn
	DN  string
}

func (l *LDAP) SearchForSlackUID(slackUID string) (*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(l.DN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, fmt.Sprintf("(slackuid=%s)", slackUID), []string{"uid", "cn"},nil)
	sr, err := l.L.Search(searchRequest)

	if err != nil {
		return nil, err
	}

	for _, entry := range sr.Entries {
		spew.Dump(entry)
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("uid"))
	}

	return sr.Entries[0], nil
}
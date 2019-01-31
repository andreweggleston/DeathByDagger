package helpers

import (
	"fmt"
	"gopkg.in/ldap.v3"
)

type LDAP struct {
	L	*ldap.Conn
	DN  string
}

func (l *LDAP) SearchForSlackUID(slackUID string) (*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(l.DN, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, fmt.Sprintf("(&(slackuid=%S))", slackUID), []string{"dn", "cn"},nil)
	sr, err := l.L.Search(searchRequest)

	if err != nil {
		return nil, err
	}

	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}

	return sr.Entries[0], nil
}
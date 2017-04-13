package puppet

import (
	"bytes"
	"errors"
	"fmt"
)

// KEYWORD represents the keywords supported in Puppetfile
type KEYWORD string

const (
	KW_GIT            KEYWORD = ":git"
	KW_LATEST         KEYWORD = ":latest"
	KW_REF            KEYWORD = ":ref"
	KW_TAG            KEYWORD = ":tag"
	KW_COMMIT         KEYWORD = ":commit"
	KW_BRANCH         KEYWORD = ":branch"
	KW_DEFAULT_BRANCH KEYWORD = ":default_branch"
)

// isValid checks that the KEYWORD is supported by this api
func (p KEYWORD) isValid() bool {
	return p == KW_GIT ||
		p == KW_LATEST ||
		p == KW_REF ||
		p == KW_TAG ||
		p == KW_COMMIT ||
		p == KW_BRANCH ||
		p == KW_DEFAULT_BRANCH
}

// Property represents a single property in a module other than the version.
type property struct {
	key   KEYWORD
	value string
}

// String prints the Puppetfile representation of modules property
func (p *property) String() string {
	return fmt.Sprintf("\t%s => %s\n", string(p.key), formatValue(p.value))
}

type Module struct {
	Name       string
	Version    string
	properties []*property
}

// New Module creates a new puppet.Module
func NewModule() *Module {
	return &Module{}
}

// SetProperty will add, update or remove a property from the module. If the value is empty,
// it would be removed. If a key does not exist, it will be added. If a key exists, it will
// be updated. Note that this function will maintain order and will set the order of the properties
// in the order they are added to the Module.
func (m *Module) SetProperty(keyword KEYWORD, value string) error {
	if !keyword.isValid() {
		return errors.New(fmt.Sprintf("Invalid or Unsupported keyword provided in puppetfile: %s", keyword))

	}
	found := false
	for i, prop := range m.properties {
		if prop.key == keyword {
			found = true
			if value == "" {
				m.properties = append(m.properties[:i], m.properties[i+1:]...)
				break
			}
			prop.value = value
			break
		}
	}
	if !found {
		m.properties = append(m.properties, &property{key: keyword, value: value})
	}

	return nil
}

// String prints the module in the format that can be used in the Puppetfile.
// If no properties or version exits, the String property will return an empty string.
func (m *Module) String() string {
	if len(m.properties) == 0 && m.Version == "" {
		return ""
	}
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("mod '%s',", m.Name))

	// Prints the version only if specified
	if m.Version != "" {
		b.WriteString(" ")
		b.WriteString(formatValue(m.Version))

	} else {
		// Prints Properties if no version
		b.WriteString("\n")
		for _, v := range m.properties {
			b.WriteString(v.String())

		}
	}
	return b.String()
}

// formatValue formats a value of one of the Modules property values
func formatValue(s string) string {
	if KEYWORD(s).isValid() {
		return s
	}
	return fmt.Sprintf("'%s'", s)
}

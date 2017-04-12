package puppet

import (
	"bytes"
	"errors"
	"fmt"
)

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

func (p KEYWORD) isValid() bool {
	return p == KW_GIT ||
		p == KW_LATEST ||
		p == KW_REF ||
		p == KW_TAG ||
		p == KW_COMMIT ||
		p == KW_BRANCH ||
		p == KW_DEFAULT_BRANCH
}

type Module struct {
	Name       string
	Version    string
	properties map[KEYWORD]string
}

func NewModule() *Module {
	return &Module{
		properties: make(map[KEYWORD]string),
	}
}

func (m *Module) SetProperty(keyword KEYWORD, value string) error {
	if !keyword.isValid() {
		return errors.New(fmt.Sprintf("Invalid or Unsupported keyword provided in puppetfile: %s", keyword))

	}
	m.properties[keyword] = value
	return nil
}

func (m *Module) String() string {
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("mod '%s',", m.Name))

	// Prints the version only if specified
	if m.Version != "" {
		b.WriteString(" ")
		b.WriteString(m.formatValue(m.Version))

	} else { // Prints Properties if no version
		b.WriteString("\n")
		for k, v := range m.properties {
			b.WriteString(m.formatProperty(k, v))

		}
	}
	return b.String()
}

func (m *Module) formatValue(s string) string {
	if KEYWORD(s).isValid() {
		return s
	}
	return fmt.Sprintf("'%s'", s)
}

func (m *Module) formatProperty(k KEYWORD, s string) string {
	return fmt.Sprintf("\t%s => %s\n", string(k), m.formatValue(s))
}

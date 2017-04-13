package puppet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// PuppetParse is used to Write and Read a Puppetfile
type PuppetParse struct{}

// WriteModules takes a list of Puppetfile module definitions and writes
// them to the given io.Writer
func (p PuppetParse) WriteModules(w io.Writer, modules []*Module) error {
	writer := bufio.NewWriter(w)
	for _, mod := range modules {
		if _, err := writer.WriteString(mod.String()); err != nil {
			return err
		}
	}
	return writer.Flush()
}

// ReadModules reads all of the modules as they are defined in puppet from
// an io.Reader
func (p PuppetParse) ReadModules(r io.Reader) ([]*Module, error) {
	reader := bufio.NewReader(r)
	var modules []*Module
	var currentModule *Module
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		// start New Module
		if strings.Contains(line, "mod '") {
			parts := strings.Split(line, ",")
			currentModule = NewModule()
			modules = append(modules, currentModule)
			// Set Module Name
			currentModule.Name, err = p.parseQuotedText(parts[0])
			if err != nil {
				return nil, err
			}
			// Versioned Module
			if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
				// Get Version
				currentModule.Version, err = p.parseValueText(parts[1])
				if err != nil {
					return nil, err
				}
				// Bad Format
			} else if len(parts) > 2 {
				return nil, errors.New(fmt.Sprintf("Invalid TestPuppetfile, too many parts for module declaration, Could not parse module from line: %s", line))
			}
			// Add to existing
		} else {
			parts := strings.Split(line, "=>")
			if len(parts) != 2 {
				return nil, errors.New(fmt.Sprintf("Invalid TestPuppetfile, Git parameters must contain =>: line: %s", line))
			}
			label, err := p.parseAndValidateKeywordText(parts[0])
			if err != nil {
				return nil, err
			}
			value, err := p.parseValueText(parts[1])
			if err != nil {
				return nil, err
			}
			currentModule.SetProperty(KEYWORD(label), value)

		}
	}
	return modules, nil
}

// parseValueText is used to attempt to extract text from a puppet file variable value. It
// first tries to pull quoted text from the string and if the text does not exist, it will
// attempt to determine if a keyword was used as a value.
func (p PuppetParse) parseValueText(s string) (string, error) {
	result, err := p.parseQuotedText(s)
	if err != nil {
		result, err = p.parseAndValidateKeywordText(s)
	}
	return result, err

}

// parseQuotedText attempts to parse the value of a puppet file property that
// is a string
func (p PuppetParse) parseQuotedText(s string) (string, error) {
	start := strings.Index(s, "'")
	end := strings.LastIndex(s, "'")
	if start == -1 || end == -1 || start == end-1 {
		return "", errors.New(fmt.Sprintf("Invalid TestPuppetfile, Could not parse text from line: %s", s))
	}
	return s[start+1 : end], nil

}

// parseAndValidateKeywordText attempts to parse a keyword from the puppet file
// and returns an error if the keyword is not supported.
func (p PuppetParse) parseAndValidateKeywordText(s string) (string, error) {
	keyword := ""
	start := strings.Index(s, ":")
	end := strings.LastIndex(s, " =>")
	if start == -1 {
		return "", errors.New(fmt.Sprintf("Invalid TestPuppetfile, Could not parse colon text from line: %s", s))
	} else if end == -1 {
		keyword = strings.TrimSpace(s[start:])
	} else {
		keyword = strings.TrimSpace(s[start:end])
	}
	if !KEYWORD(keyword).isValid() {
		return "", errors.New(fmt.Sprintf("Invalid or Unsupported keyword provided in puppetfile: %s", keyword))
	}
	return keyword, nil
}

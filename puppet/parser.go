package puppet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type PuppetParse struct{}

// ReadModules reads all of the modules as they are definied in puppet
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
				currentModule.Version, err = p.extractText(parts[1])
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
			value, err := p.extractText(parts[1])
			if err != nil {
				return nil, err
			}
			currentModule.properties[KEYWORD(label)] = value

		}
	}
	return modules, nil
}

func (p PuppetParse) extractText(s string) (string, error) {
	result, err := p.parseQuotedText(s)
	if err != nil {
		result, err = p.parseAndValidateKeywordText(s)
	}
	return result, err

}

func (p PuppetParse) parseQuotedText(s string) (string, error) {
	start := strings.Index(s, "'")
	end := strings.LastIndex(s, "'")
	if start == -1 || end == -1 || start == end-1 {
		return "", errors.New(fmt.Sprintf("Invalid TestPuppetfile, Could not parse text from line: %s", s))
	}
	return s[start+1 : end], nil

}

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

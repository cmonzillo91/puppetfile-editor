package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"strings"

	"github.com/cmonzillo91/puppetfile-editor/puppet"
	"github.com/pkg/errors"
)

func main() {
	log.SetOutput(os.Stdout)
	conf := &Config{}
	conf.LoadFromFlags()
	if err := conf.ValidateConfig(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Open file
	log.Println("Reading Puppetfile")
	file, err := os.Open(conf.FileName)
	if err != nil {
		log.Fatal("Could not read TestPuppetfile from disk: %s", err)
		os.Exit(1)
	}

	// Parse Puppetfile
	log.Println("Parsing Puppetfile")
	parser := puppet.PuppetParse{}
	modules, err := parser.ReadModules(file)
	if err != nil {
		file.Close()
		log.Fatal("Could not parse Puppetfile: %s", err)
		os.Exit(1)
	}
	file.Close()

	// Update module data
	log.Println("Modifying Puppetfile")
	UpdateModuleProperty(modules, conf.Module, conf.Key, conf.Value)

	// Create New Files
	log.Println("Writing New Modules")
	file, err = os.Create(conf.FileName)
	if err != nil {
		file.Close()
		log.Fatal("Could not open new Puppetfile: %s", err)
		os.Exit(1)
	}
	if err := parser.WriteModules(file, modules); err != nil {
		file.Close()
		log.Fatal("Could write new modules to the Puppetfile: %s", err)
		os.Exit(1)
	}
	file.Close()

}

// UpdateModuleProperty, updates a single property in a single module
func UpdateModuleProperty(modules []*puppet.Module, module, key, value string) {
	for _, mod := range modules {
		if strings.ToLower(mod.Name) == module {
			if strings.ToLower(key) == "version" {
				mod.Version = value
			} else if err := mod.SetProperty(puppet.KEYWORD(key), value); err != nil {
				log.Fatalf("Could not set property: %s", err)
			}
		}
	}
}

// Config contains the command line properties passed into the program
type Config struct {
	FileName string
	Module   string
	Key      string
	Value    string
}

// LoadFromFlags loads the configuration from commandline flags
func (c *Config) LoadFromFlags() {
	fileName := flag.String("puppetfile", "", "Original PuppetFile")
	module := flag.String("module", "", "The module whos properties to update")
	key := flag.String("key", "", "The key of the property to change")
	value := flag.String("value", "", "Value of the property that will be set")
	flag.Parse()
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
		flag.PrintDefaults()
	}
	if flag.Lookup("help") != nil {
		flag.Usage()
		os.Exit(1)
	}

	c.FileName = *fileName
	c.Module = *module
	c.Key = *key
	c.Value = *value
}

// ValidateConfig validates that the config is properly set
func (c *Config) ValidateConfig() error {
	if c.FileName == "" {
		return errors.New("Puppetfile not provided")
	}
	if c.Module == "" {
		return errors.New("Module not provided")
	}
	if c.Key == "" {
		return errors.New("Key not provided")
	}
	return nil
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cmonzillo91/puppetfile-editor/puppet"
)

func main() {

	fileName := flag.String("puppetfile", "", "Original PuppetFile")
	//overwrite := flag.Bool("overwrite", false, "Overwrite current file")
	//output := flag.String("o", "TestPuppetfile.new", "New file")

	if *fileName == "" {
		log.Fatalf("Puppetfile not provided")
		os.Exit(1)
	}
	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal("Could not read TestPuppetfile from disk: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	parser := puppet.PuppetParse{}
	modules, err := parser.ReadModules(file)
	if err != nil {
		log.Fatal("Could not parse TestPuppetfile: %s", err)
		os.Exit(1)
	}
	fmt.Println(modules[0])
}

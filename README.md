# puppetfile-editor
A simple utility to modify properties in a puppetfile for puppet. It supports update the properties
for individual modules in the puppet file. Currently, this tool only supports updating on property in
one module at a time.

# Usage
Below is an example of updating a `:ref` property of a Puppetfile module named `hieradata` with the value `1.0.1.0`.
```
puppetfile-editor --puppetfile TestPuppetfile --module hieradata --key :ref --value 1.0.1.0
```

Resulting Puppetfile
```
mod 'hieradata',
	:git => 'test@github.com.git'
	:ref => '1.0.1.0'
mod 'test_module',
	:git => 'test@github.com.git'
	:ref => '1.0.0.10'
```

# Help
For help, you can use the `--help` flags
```
# ./puppetfile-editor --help
Usage of ./puppetfile-editor:
  -key string
    	The key of the property to change
  -module string
    	The module whos properties to update
  -puppetfile string
    	Original PuppetFile
  -value string
    	Value of the property that will be set
```
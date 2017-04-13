# puppetfile-editor
A simple utility to modify properties in a puppetfile for puppet. It supports update the properties
for individual modules in the puppet file. Currently, this tool only supports updating on property in
one module at a time.

# Usage
Below is an example of updating a :ref property of a Puppetfile module named hieradata.
```
puppetfile-editor --puppetfile TestPuppetfile --module hieradata --key :ref --value 1.0.1.0
```
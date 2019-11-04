# FileCryptoGolang
This project contains a cli tool for encrypting/decrypting and bruteforcing files.

## Usage:
  -b string\
    	&ensp;File to bruteforce\
  -c string\
    	&ensp;Characters for bruteforcing\
  -d string\
    	&ensp;File to decrypt\
  -e string\
    	&ensp;File to encrypt\
  -l int\
    	&ensp;Max password length for brute force\
  -o string\
    	&ensp;File to save decrypted content into\
  -p string\
    	&ensp;Password to encrypt/decrypt file\

## Compile and Install:
  To compile: run ```go build main.go```
  (You can specify the "-o output_filename" flag to name the executable differently or simply rename it using ```mv```)
  
  Install the go way: ```go install main.go``` to install compile and install the binary to your GOBIN directory (usually this is &nbsp;&nbsp;$HOME/go/bin)
  
  Install the bash way: ```install <executable_name> /usr/local/bin``` (make sure /usr/local/bin is in your $PATH or simply change &nbsp;&nbsp;it to a directory that is in your path)

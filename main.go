package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(filename string, outputFilename string, passphrase string) {
	if plaintext, err := ioutil.ReadFile(filename); err != nil {
		panic("Cant read from file!")
	} else {
		f, _ := os.Create(outputFilename)
		defer f.Close()
		f.Write(encrypt(plaintext, passphrase))
	}
}

func decryptFile(filename string, outputFilename string, passphrase string) {
	if ciphertext, err := ioutil.ReadFile(filename); err != nil {
		log.Fatalf("Failed to read from file!\n")
	} else {
		if plaintext := decrypt(ciphertext, passphrase); string(plaintext) != "" {
			if f, err := os.Create(outputFilename); err != nil {
				log.Fatalf("Failed to create the output file!")
			} else {
				defer f.Close()
				if _, err := f.Write(plaintext); err != nil {
					log.Fatalf("Failed to write to file!\n")
				} else {
					fmt.Printf("Successfully decrypted %s into %s\n", filename, outputFilename)
					return
				}
			}
		} else {
			log.Fatalf("Failed to decrypt file!\n")
		}
	}
}

func GenerateCombinations(alphabet string, length int) <-chan string {
	c := make(chan string)
	go func(c chan string) {
		defer close(c)

		AddLetter(c, "", alphabet, length)
	}(c)
	return c
}

func AddLetter(c chan string, combo string, alphabet string, length int) {
	if length <= 0 {
		return
	}

	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		c <- newCombo
		AddLetter(c, newCombo, alphabet, length-1)
	}
}

func bruteForce(filename string, outputFilename string, characters string, passlen int) {
	if data, err := ioutil.ReadFile(filename); err != nil {
		log.Fatalf("Failed to read from file!\n")
	} else {
		for passGuess := range GenerateCombinations(characters, passlen) {
			fmt.Printf("\033[2K\rTrying: %s", passGuess)
			key := []byte(createHash(passGuess))
			if block, err := aes.NewCipher(key); err != nil {
				log.Fatalf("\nFailed to create a new AES Cipher!\n")
			} else {
				if gcm, err := cipher.NewGCM(block); err != nil {
					log.Fatalf("\nFailed to create a new GCM Block!\n")
				} else {
					nonceSize := gcm.NonceSize()
					nonce, ciphertext := data[:nonceSize], data[nonceSize:]
					if plaintext, err := gcm.Open(nil, nonce, ciphertext, nil); err != nil {
						continue
					} else {
						if f, err := os.Create(outputFilename); err != nil {
							log.Fatalf("\nSuccessfully decrypted the file %s, but failed to create the output to %s!\nThe decrypted content is:\n%s",filename, outputFilename, string(plaintext))
						} else {
							defer f.Close()
							if _, err := f.Write(plaintext); err != nil {
								log.Fatalf("\nSuccessfully decrypted the file %s, but failed to write the output to %s!\nThe decrypted content is:\n%s",filename, outputFilename, string(plaintext))
							} else {
								fmt.Printf("\nSuccessfully decrypted the file %s and write it's output to %s\nThe decrypted content is:\n%s\n", filename, outputFilename, string(plaintext))
								return
							}
						}
					}
				}
			}
		}
		log.Fatalf("\033[2K\r\033[0;31mNo matching password found using the provided characters: '%s' with maximum password length: %d\033[0m\n", characters, passlen)
	}

}

func checkMutuallyExclusiveArgsString(args []string) bool{
	for i := range args {
		for x := range args {
			if x != i {
				if args[i] != "" && args[x] != "" {
					fmt.Printf("Comparing %s and %s\n", args[i], args[x])
					return false
				}
			}
		}
	}
	return true
}

func requireFlags(flags []string) {
	required := flags
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			log.Fatalf("missing required -%s argument/flag\n", req)
		}
	}
}

func handleArgs() {
	encryptFlag := flag.String("e", "", "File to encrypt")
	decryptFlag := flag.String("d", "", "File to decrypt")
	bruteforceFlag := flag.String("b", "", "File to bruteforce")

	passwordFlag := flag.String("p", "", "Password to encrypt/decrypt file")
	characterFlag := flag.String("c", "", "Characters for bruteforcing")
	outputFlag := flag.String("o", "", "File to save decrypted content into")
	lengthFlag := flag.Int("l", 0, "Max password length for brute force")

	flag.Parse()

	if checkMutuallyExclusiveArgsString([]string{*encryptFlag, *decryptFlag, *bruteforceFlag}) {
		if *encryptFlag != "" {
			requireFlags([]string{"e", "p", "o"})
			encryptFile(*encryptFlag, *outputFlag, *passwordFlag)
		} else if *decryptFlag != "" {
			requireFlags([]string{"d", "o", "p"})
			decryptFile(*decryptFlag, *outputFlag, *passwordFlag)
		} else if *bruteforceFlag != "" {
			requireFlags([]string{"b", "c", "l", "o"})
			bruteForce(*bruteforceFlag, *outputFlag, *characterFlag, *lengthFlag)
		}
	} else {

	}
}

func main() {
	handleArgs()
}
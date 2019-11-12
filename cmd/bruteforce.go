package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/spf13/cobra"
	"gocrypter/hash"
	"gocrypter/log"
	"io/ioutil"
	"os"
)

var (
	bruteForceOutputFile        string
	bruteForceOutputDestination *os.File
	characters                  string
	passlen                     int
)

func init() {
	rootCmd.AddCommand(bruteForceCmd)
	bruteForceCmd.Flags().StringVarP(&bruteForceOutputFile, "output-file", "o", "", "the file to write the decrypted content to (default Stdout)")
	bruteForceCmd.Flags().StringVarP(&characters, "characters", "c", "", "the characters to include for brute forcing")
	bruteForceCmd.Flags().IntVarP(&passlen, "password-length", "l", 1, "the maximum password length (default is 1)")
}

var bruteForceCmd = &cobra.Command{
	Use:              "bruteforce",
	Aliases:          []string{"b", "brute"},
	SuggestFor:       nil,
	Short:            "bruteforce a file",
	Long:             "",
	Example:          "",
	ValidArgs:        nil,
	Args:             cobra.ExactArgs(1),
	ArgAliases:       nil,
	PersistentPreRun: nil,
	PreRun: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			log.Fatalf("%s: no such file or directory", args[0])
		} else if os.IsPermission(err) {
			log.Fatalf("%s: permission denied", args[0])
		} else if err != nil {
			log.Panic(err)
		}

		if bruteForceOutputFile != "" {
			file, err := os.Stat(bruteForceOutputFile)
			if os.IsNotExist(err) {
				bruteForceOutputDestination, err = os.Create(bruteForceOutputFile)
			} else if os.IsPermission(err) {
				log.Fatal(err)
			} else if os.IsExist(err) && !func() bool {
				bytes, err := ioutil.ReadFile(file.Name())
				if err != nil {
					log.Panic(err)
				}
				return bytes == nil
			}() {
				log.Fatalf("%s exists, but is not empty", file.Name())
			}
		} else {
			bruteForceOutputDestination = os.Stdout
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if data, err := ioutil.ReadFile(args[0]); err != nil {
			log.Fatalf("Failed to read from file!\n")
		} else {
			for passGuess := range generateCombinations(characters, passlen) {
				log.Printf("\033[2K\rTrying: %s", passGuess)
				key := []byte(hash.Create(passGuess))
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
							if _, err := bruteForceOutputDestination.Write(plaintext); err != nil {
								log.Printf("\nSuccessfully decrypted the file %s, but failed to write the output to %s!\nThe decrypted content is:\n%s", args[0], bruteForceOutputDestination.Name(), string(plaintext))
								log.Fatal(err)
							} else {
								log.Printf("\nSuccessfully decrypted the file %s and write it's output to %s\nThe decrypted content is:\n%s\n", args[0], bruteForceOutputDestination.Name(), string(plaintext))
								return
							}
						}
					}
				}
			}
			log.Fatalf("\033[2K\r\033[0;31mNo matching password found using the provided characters: '%s' with maximum password length: %d\033[0m\n", characters, passlen)
		}
	},
	PostRun:           nil,
	PersistentPostRun: nil,
}

func generateCombinations(alphabet string, length int) <-chan string {
	c := make(chan string)
	go func(c chan string) {
		defer close(c)

		addLetter(c, "", alphabet, length)
	}(c)
	return c
}

func addLetter(c chan string, combo string, alphabet string, length int) {
	if length <= 0 {
		return
	}

	var newCombo string
	for _, ch := range alphabet {
		newCombo = combo + string(ch)
		c <- newCombo
		addLetter(c, newCombo, alphabet, length-1)
	}
}

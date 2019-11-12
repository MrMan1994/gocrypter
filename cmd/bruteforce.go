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

var(
	bruteForceFile string
	bruteForceOutputFile string
	characters string
	passlen int
)

func init() {
	rootCmd.AddCommand(bruteForceCmd)
	rootCmd.Flags().StringVarP(&bruteForceOutputFile, "output-file", "o", "", "the file to write the decrypted content to (default Stdout)")
	rootCmd.Flags().StringVarP(&characters, "characters", "c", "", "the characters to include for brute forcing")
	rootCmd.Flags().IntVarP(&passlen, "password-length", "l", 1, "the maximum password length (default is 1)")
}

var bruteForceCmd = & cobra.Command{
	Use:                        "bruteforce",
	Aliases:                    []string{"b", "brute"},
	SuggestFor:                 nil,
	Short:                      "bruteforce a file",
	Long:                       "",
	Example:                    "",
	ValidArgs:                  nil,
	Args:                       cobra.ExactArgs(1),
	ArgAliases:                 nil,
	PersistentPreRun:           nil,
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		if data, err := ioutil.ReadFile(bruteForceFile); err != nil {
			log.Fatalf("Failed to read from file!\n")
		} else {
			for passGuess := range GenerateCombinations(characters, passlen) {
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
							if f, err := os.Create(bruteForceOutputFile); err != nil {
								log.Fatalf("\nSuccessfully decrypted the file %s, but failed to create the output to %s!\nThe decrypted content is:\n%s",bruteForceFile, bruteForceOutputFile, string(plaintext))
							} else {
								defer f.Close()
								if _, err := f.Write(plaintext); err != nil {
									log.Fatalf("\nSuccessfully decrypted the file %s, but failed to write the output to %s!\nThe decrypted content is:\n%s",bruteForceFile, bruteForceOutputFile, string(plaintext))
								} else {
									log.Printf("\nSuccessfully decrypted the file %s and write it's output to %s\nThe decrypted content is:\n%s\n", bruteForceFile, bruteForceOutputFile, string(plaintext))
									return
								}
							}
						}
					}
				}
			}
			log.Fatalf("\033[2K\r\033[0;31mNo matching password found using the provided characters: '%s' with maximum password length: %d\033[0m\n", characters, passlen)
		}
	},
	PostRun:                    nil,
	PersistentPostRun:          nil,
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
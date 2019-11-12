package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/spf13/cobra"
	"gocrypter/hash"
	"gocrypter/log"
	"io"
	"io/ioutil"
	"os"
)

var(
	encryptPassword string
	encryptOutputFile string
	encryptOutputDestination os.File
)

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVarP(&encryptPassword, "password", "p", "", "the password to use for encryption")
	encryptCmd.Flags().StringVarP(&encryptOutputFile, "output-file", "o", "", "the file to save the encrypted content into (default is: the input filename + .encrypted)")
	if err := encryptCmd.MarkFlagRequired("password"); err != nil {
		log.Panic(err)
	}
}

var encryptCmd = &cobra.Command{
	Use:                        "encrypt",
	Aliases:                    []string{"e"},
	SuggestFor:                 nil,
	Short:                      "encrypt a file",
	Long:                       "",
	Example:                    "",
	ValidArgs:                  nil,
	Args:                       cobra.ExactArgs(1),
	ArgAliases:                 nil,
	PersistentPreRun:           nil,
	PreRun: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			log.Fatalf("%s: no such file or directory", args[0])
		} else if os.IsPermission(err) {
			log.Fatalf("%s: permission denied", args[0])
		} else if err != nil {
			log.Panic(err)
		}

		if encryptOutputFile != "" {
			_, err := os.Stat(encryptOutputFile)
			if os.IsNotExist(err) {
				_, err = os.Create(encryptOutputFile)
			} else if os.IsPermission(err) {
				log.Fatal(err)
			} else if err != nil{
				log.Panic(err)
			}
		} else {
			file, err := os.Create(args[0]+".encrypted")
			if err != nil {
				log.Panic(err)
			}
			encryptOutputDestination = *file
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_, err := encryptOutputDestination.Write(func(data []byte, password string) []byte {
			block, _ := aes.NewCipher([]byte(hash.Create(password)))
			gcm, err := cipher.NewGCM(block)
			if err != nil {
				log.Panic(err)
			}
			nonce := make([]byte, gcm.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
				log.Panic(err)
			}
			ciphertext := gcm.Seal(nonce, nonce, data, nil)
			return ciphertext
		}(func() []byte {
			plaintext, err := ioutil.ReadFile(args[0])
			if err != nil {
				log.Panic(err)
			}
			return plaintext
		}(), encryptPassword))
		if err != nil {
			log.Panic(err)
		}
	},
	PostRun:                    nil,
	PersistentPostRun:          nil,
}
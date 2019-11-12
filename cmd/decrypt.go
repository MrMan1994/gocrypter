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
	decryptOutputFile        string
	decryptPassword          string
	decryptOutputDestination *os.File
)

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringVarP(&decryptOutputFile, "output-file", "o", "", "the output filename (default is Stdout)")
	decryptCmd.Flags().StringVarP(&decryptPassword, "password", "p", "", "the password to use for decryption")
}

var decryptCmd = &cobra.Command{
	Use:              "decrypt",
	Aliases:          []string{"d"},
	SuggestFor:       nil,
	Short:            "decrypt a file",
	Long:             "",
	Example:          "",
	ValidArgs:        nil,
	Args:             nil,
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

		if decryptOutputFile != "" {
			file, err := os.Stat(decryptOutputFile)
			if os.IsNotExist(err) {
				decryptOutputDestination, err = os.Create(decryptOutputFile)
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
			decryptOutputDestination = os.Stdout
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		encryptedText, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatalf("failed to read from file\n")
		}
		plaintext := func() []byte {
			key := []byte(hash.Create(decryptPassword))
			block, err := aes.NewCipher(key)
			if err != nil {
				log.Panic(err)
			}
			gcm, err := cipher.NewGCM(block)
			if err != nil {
				log.Panic(err)
			}
			nonceSize := gcm.NonceSize()
			nonce, ciphertext := encryptedText[:nonceSize], encryptedText[nonceSize:]
			plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
			if err != nil {
				log.Panic(err)
			}
			return plaintext
		}()
		if string(plaintext) != "" {
			if _, err := decryptOutputDestination.Write(plaintext); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalln("failed to decrypt file")
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		log.Printf("successfully decrypted %s into %s\n", args[0], decryptOutputDestination.Name())
	},
	PersistentPostRun: nil,
}

package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/codegangsta/cli"
	"github.com/sam701/awstools/config"
	"github.com/sam701/awstools/sess"
)

func kmsAction(c *cli.Context) error {
	txt := c.Args().First()
	if txt == "" {
		cli.ShowCommandHelp(c, "kms")
		return nil
	}
	quiet := c.Bool("quiet")
	cl := kms.New(sess.FromEnvVar())
	if c.Bool("decrypt") {
		bb, err := base64.StdEncoding.DecodeString(txt)
		if err != nil {
			log.Fatalln("ERROR", err)
		}

		out, err := cl.Decrypt(&kms.DecryptInput{
			CiphertextBlob: bb,
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		if !quiet {
			fmt.Print("Decrypted: ")
		}
		fmt.Println(string(out.Plaintext))
	} else if c.Bool("encrypt") {
		keyId := c.String("key-id")
		if keyId == "" {
			keyId = config.Current.DefaultKmsKey
		}
		if keyId == "" {
			log.Fatalln("No key-id provided")
		}
		out, err := cl.Encrypt(&kms.EncryptInput{
			KeyId:     aws.String(keyId),
			Plaintext: []byte(txt),
		})
		if err != nil {
			log.Fatalln("ERROR", err)
		}
		if !quiet {
			fmt.Print("Encrypted: ")
		}
		fmt.Println(base64.StdEncoding.EncodeToString(out.CiphertextBlob))
	} else {
		cli.ShowCommandHelp(c, "kms")
	}
	return nil
}

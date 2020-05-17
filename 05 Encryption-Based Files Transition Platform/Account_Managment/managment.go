// managment
package main

import (
	"os"
	"fmt"
	"io/ioutil" // file io
	
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateKey(who , folder string, size int) {
	newPrivkey, err := rsa.GenerateKey(rand.Reader, size)
	Check(err)
    
    // create pem file
    PemRrikey := []byte(RsaPriToPem(newPrivkey))
    
    	newPubkey := &newPrivkey.PublicKey
	PemPubkey := []byte(RsaPubToPem(newPubkey))
	
	// https : //golang.org/pkg/io/ioutil/#WriteFile
	err = ioutil.WriteFile("../"+folder+"/key/"+who+"pri_key.pem", PemRrikey, 0644)
	
	if err != nil {
		fmt.Printf("Error creating Key file!")
		os.Exit(0)
	} else {
		err = ioutil.WriteFile("../PublicKeys/"+who+"pub_key.pem", PemPubkey, 0644)
	}
}

func RsaPriToPem (rsaPrivkey *rsa.PrivateKey) string {
	// *rsa.PrivateKey to []byte：x509.MarshalPKCS1PrivateKey
	// https://stackoverflow.com/questions/13555085/save-and-load-crypto-rsa-privatekey-to-and-from-the-disk
    privkey_bytes := x509.MarshalPKCS1PrivateKey(rsaPrivkey)
    
    privkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PRIVATE KEY",
                    Bytes: privkey_bytes,
            },
    )
    return string(privkey_pem)
}

func RsaPubToPem (rsaPubkey *rsa.PublicKey) string {
	// *rsa.PublicKey to []byte：x509.MarshalPKIXPublicKey
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(rsaPubkey)
	Check(err)
	
    pubkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PUBLIC KEY",
                    Bytes: pubkey_bytes,
            },
    )
    
    return string(pubkey_pem)
}

func main() {
	// Creat keys
	// 角色、資料夾名稱、key長度(byte)
	CreateKey("sender", "Encrypt", 4096)
	CreateKey("receiver", "Decrypt", 4096)
	
	fmt.Println("Done!")
}

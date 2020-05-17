// file
package main

import (
	"os"
	"fmt"
	"io/ioutil" // file io
	"errors"
	"archive/zip"
	// example:https://www.itread01.com/content/1546726868.html
	
	"crypto/rand"
	"crypto/rsa"
	// "crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

var (
	ReceiverPubkey *rsa.PublicKey
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func PemToRsaPub (pubPEM string) (*rsa.PublicKey, error) {
    block, _ := pem.Decode([]byte(pubPEM))
    if block == nil {
            return nil, errors.New("failed to parse PEM block containing the key")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
            return nil, err
    }
	
	fmt.Println("pub=",pub)
	
    switch pub := pub.(type) {
    case *rsa.PublicKey:
            return pub, nil
    default:
            break // fall through
    }
    
    return nil, errors.New("Key type is not RSA")
}

// Read Files
func ReadFiles(folder string) ([][]byte, []string) {
	
	dir := folder+"/"
	files, err := ioutil.ReadDir(dir)
	fmt.Println(files,"\nlen=",len(files))
	
    Check(err)
    
    FileData := [][]byte{}
    FileName := []string{}
    for _, file := range files {
    		filecontent, err := ioutil.ReadFile(dir + file.Name())
        Check(err)
        
    		FileData = append(FileData, filecontent)
		FileName = append(FileName, "Crypt_"+ file.Name())
    }
    
    return FileData, FileName
}

// Encrypt Files and Compress
func encryptfile(folder string) {

	// Read files
	FileData, FileName := ReadFiles(folder)
	
	// encrypt and zip files
	fzip, _ := os.Create(folder+"(encrypt).zip")
	w := zip.NewWriter(fzip)
    defer w.Close()
    for i, filedata := range FileData {
        fw, _ := w.Create(FileName[i])
        
        // encrypt file
        // hash_data := sha256.Sum256(filedata)
        // https://ithelp.ithome.com.tw/articles/10188698
        encryptedmsg, err := rsa.EncryptPKCS1v15(rand.Reader, ReceiverPubkey, filedata[:])
		Check(err)
		
        _, err = fw.Write(encryptedmsg)
        Check(err)
        // fmt.Println(n)
    }
}

func main() {
	
	fmt.Println("[]byte vs []unit8")
	// uint8  the set of all unsigned  8-bit integers (0 to 255)
	// byte   alias for uint8
	a := []byte{1, 1}
	b := []uint8{2, 2}

	B := func(b []byte) []byte {
		return b
	}(b)

	fmt.Printf("\ntype of a:%T\ntype of b:%T\ntype of B:%T\n", a, b, B)
	
	// Read Receiver Public Key file
	receiver_pubPEM, err := ioutil.ReadFile("../PublicKeys/receiverpub_key.pem")
	Check(err)
	
	ReceiverPubkey, err = PemToRsaPub(string(receiver_pubPEM))
	Check(err)
	
	// encrypt file in "Files" folder
	encryptfile("Files")
	fmt.Println("Done!")
	
}

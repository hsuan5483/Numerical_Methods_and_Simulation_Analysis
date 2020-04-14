// file
package main

//參考https://www.wandouip.com/t5i259203/
//參考http://blog.studygolang.com/2013/01/go%E5%8A%A0%E5%AF%86%E8%A7%A3%E5%AF%86%E4%B9%8Brsa/
//各種形式加密解密https://blog.csdn.net/u013565368/article/details/53081195
//函數形式參考https://github.com/Vaultpls/Go-File-Encryption-AES/blob/master/enc.go
import (
	// "encoding/pem"
	"errors"
	"fmt"
	"io/ioutil" // file io
	"os"

	// "log"
	"archive/zip"
	// example:https://www.itread01.com/content/1546726868.html

	// "crypto"
	"crypto/rand"
	"crypto/rsa"

	// "crypto/sha256"

	"crypto/x509"

	"encoding/pem"
	"io"
	// "path/filepath"
)

var (
	privateKey, _ = ioutil.ReadFile("./key/receiverpri_key.pem")
	publicKey, _  = ioutil.ReadFile("./key/receiverpub_key.pem")
)

func main() {

	fmt.Println(privateKey)
	fmt.Printf("privateKey orginal TYPE : %T\n", privateKey)

	priv, err := ParseRsaPrivateKeyFromPemStr(string(privateKey))

	fmt.Println(priv)
	fmt.Printf("privateKey converted TYPE : %T\n", priv)

	encryptedFile := "./Files(encrypt).zip" //加密後檔案之壓縮檔
	path := "./DeCompressZip"               //解壓縮檔儲存位址
	depath := "./dcrypted"                  //解密檔儲存位址

	DeCompressZip(encryptedFile, path) //解壓縮

	os.Mkdir(depath, 0777) //一定要打這個才會創建新路徑

	list, err := ioutil.ReadDir(path) //ioutil.ReadDir讀資料夾、ioutil.ReadFile讀單一檔案 ; 用list叫出資料夾的基本性質
	if err != nil {                   //https://blog.csdn.net/yxys01/article/details/78136295
		fmt.Println("Read DeCompressZip")
		panic(err)
	}

	for _, files := range list {
		fmt.Println(files.Name())                                         //確認讀入檔名
		encrypted_file, err := ioutil.ReadFile(path + "/" + files.Name()) //用.Name()去抓出資料夾內的檔案的名稱並讀取
		if err != nil {
			fmt.Println("read dir error")
			return
		}
		fmt.Printf("%T\n", encrypted_file) //確認讀入內容

		TESTa, err := rsa.DecryptPKCS1v15(rand.Reader, priv, encrypted_file) //解密	//方法一
		// TESTa, err := RsaDecrypt(encrypted_file)							 		//方法二
		// fmt.Println("TESTa = ", TESTa) //確認TESTa內容

		if err != nil {
			fmt.Println("rsa.DecryptPKCS1v15 error")
			panic(err)
		}

		fw, _ := os.Create(depath + "/" + "de" + files.Name()) //開檔案位置
		n, err := fw.Write(TESTa)                              //將內容寫入檔案
		if err != nil {
			fmt.Println("fw.Write error")
			fmt.Println(err)
		}
		defer fw.Close()
		fmt.Println(n)

	}

}

//解壓縮https://blog.csdn.net/wangshubo1989/article/details/71743374
//解壓縮https://studygolang.com/articles/7471

func DeCompressZip(File, dir string) {

	os.Mkdir(dir, 0777) //建立一個目錄（資料夾）

	cf, err := zip.OpenReader(File) //讀取zip檔案
	if err != nil {
		fmt.Println("Func DeCompressZip rader")
		fmt.Println(err)
	}
	defer cf.Close()
	for _, file := range cf.File {
		rc, err := file.Open() //開檔
		if err != nil {
			fmt.Println("Func DeCompressZip open")
			fmt.Println(err)
		}

		f, err := os.Create(dir + "/" + file.Name) //開檔案位置
		if err != nil {
			fmt.Println("Func DeCompressZip os.Create")
			fmt.Println(err)
		}
		defer f.Close()          //關檔
		n, err := io.Copy(f, rc) //覆寫
		if err != nil {
			fmt.Println("Func DeCompressZip io.Copy")
			fmt.Println(err)
		}
		fmt.Println(n)
	}

}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		fmt.Println("Func ParseRsaPrivateKeyFromPemStr Decode")
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Func ParseRsaPrivateKeyFromPemStr x509.ParsePKCS1PrivateKey")
		return nil, err
	}

	return priv, nil
}

func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		fmt.Println("Func RsaDecrypt Decode")
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	fmt.Println(priv)
	if err != nil {
		fmt.Println("Func RsaDecrypt x509.ParsePKCS1PrivateKey")
		return nil, err

	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)

}

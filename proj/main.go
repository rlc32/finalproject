package main
import (
	"fmt"
	"os"
	"crypto/rsa"
	"crypto/rand"
	"log"
    "crypto/md5"
    "crypto/x509"
	"hash"
	"io/ioutil"
	"encoding/pem"
	)
func main(){
var repeat bool
repeat = true
var answer string
str1 := "Welcome to a second class password keeper!\nDo you want to add a password or see a password?\n"
str2 := "Please enter the word 'add' to add a new password or enter 'see' to see a password\n or enter 'exit' to exit:\n"
fmt.Printf(str1)

	for repeat == true {
		fmt.Scanf("%s", &answer)
		if answer == "add"{
			addpass()
		break
		}
		if answer == "see"{
			seepass()
		break
		} 
		if answer == "exit"{
		break		
		}else{
			fmt.Printf(str2)
		}	
		
	}

}


/* This function takes no arguments:
			designed to ask for a website or product to store username and password.
			stores username in a separate file with key.
			asks for password and sends it to be encrypted using the key.*/

func addpass(){
	var password, productName string
	var encPassword []byte
	fmt.Printf("please enter the name of the product that you want to store: \n")
	fmt.Scanf("%s", &productName)
	fmt.Printf("Please enter the password for %s: \n", productName)
	fmt.Scanf("%s", &password)
	encPassword = encryption(password, productName)
	f, err := os.Create("/Users/russclousing/214/text/ " + productName + ".txt")
	check(err)
	defer f.Close()
	
	
	n1, err := f.Write(encPassword)
	check(err)
	fmt.Printf("wrote %d bytes\n", n1)
	f.Sync()
	
	
}


/*This function takes no arguments:
			designed to ask for a website or product to return the username and password
			retrieves the username and key from file
			uses key to decrypt the password stored in separate file*/

func seepass(){
productName := ""
fmt.Printf("please enter the name of the product that you want to see: \n")
fmt.Scanf("%s", &productName)
password := decryption(productName)
fmt.Printf(password)
}




//this is the function that takes a encrypted password and returns a unencrypted password
func decryption(product string)(password string){
var private_key *rsa.PrivateKey
var pem_data, encrypted, decrypted, label []byte
var pem_file_path, encryptedFileName string
var err error
var block *pem.Block

pem_file_path = product + ".pub"
pem_data, err = ioutil.ReadFile(pem_file_path)
 
if err != nil {
        log.Fatalf("Error reading pem file: %s", err)
    }
    block, _ = pem.Decode(pem_data)
    
if block == nil || block.Type != "RSA PUBLIC KEY" {
        log.Fatal("No valid PEM data found")
    }
 private_key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
if err != nil{
log.Fatalf("Private key can't be decoded: %s", err)
    }

encryptedFileName = "/Users/russclousing/214/text/ " + product + ".txt"


encrypted, err = ioutil.ReadFile(encryptedFileName)
if err != nil{
log.Fatalf("error in reading encrypted data: %s", err) 
}

decrypted = decrypt_oaep(private_key, encrypted, label)
password = string(decrypted)

return
}



/* this checks for errors in file usage*/
func check(e error){
	if e != nil{
		panic(e)}
}


/* This function encrypts the password and returns the password as a encrypted string of numbers
			Input: this funciton takes a string
			Output: THis returns a string of numerals that is the encrypted string*/
			
func encryption(password string, product string)(encrypted []byte){ 
	var private_key *rsa.PrivateKey
	var public_key *rsa.PublicKey
	var err error
	var plain_text, label []byte
	var filename string
	plain_text = []byte(password)
	label = []byte(product)
	
	//Generate the Private Key
    if private_key, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {
        log.Fatal(err)
    }

    // Precompute some calculations -- Calculations that speed up private key operations in the future
    private_key.Precompute()

    //Validate Private Key -- Sanity checks on the key
    if err = private_key.Validate(); err != nil {
        log.Fatal(err)
    }

    //Public key address (of an RSA key)
    public_key = &private_key.PublicKey

    encrypted = encrypt_oaep(public_key, plain_text, label)

	filename = product + ".pub"

	pemData := pem.EncodeToMemory(&pem.Block{
    		Type:  "RSA PUBLIC KEY",
    		Bytes: x509.MarshalPKCS1PrivateKey(private_key),
	})
	
	ioutil.WriteFile(filename, pemData, 0644)

	return
}



//Code used from http://blog.giorgis.io/golang-rsa-encryption
//OAEP Encrypt
func encrypt_oaep(public_key *rsa.PublicKey, plain_text, label []byte) (encrypted []byte) {
    var err error
    var md5_hash hash.Hash

    md5_hash = md5.New()

    if encrypted, err = rsa.EncryptOAEP(md5_hash, rand.Reader, public_key, plain_text, label); err != nil {
        log.Fatal(err)
    }
    return
}
//Code used from http://blog.giorgis.io/golang-rsa-encryption
//OAEP Decryption
func decrypt_oaep(private_key *rsa.PrivateKey, encrypted, label []byte) (decrypted []byte) {
    var err error
    var md5_hash hash.Hash

    md5_hash = md5.New()
    if decrypted, err = rsa.DecryptOAEP(md5_hash, rand.Reader, private_key, encrypted, label); err != nil {
        log.Fatalf("failed in decrypt_oaep: %s", err)
    }
    return
}








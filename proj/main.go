package main
import (
	"fmt"
	"crypto/rsa"
	"crypto/rand"
	"log"
    "crypto/md5"
    "crypto/x509"
    "bytes"
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
	var password, productName, filename string
	var encPassword []byte
	fmt.Printf("please enter the name of the product that you want to store: \n")
	//gets the product name
	fmt.Scanf("%s", &productName)
	fmt.Printf("Please enter the password for %s: \n", productName)
	//gets the password that should be encrypted
	fmt.Scanf("%s", &password)
	encPassword = encryption(password, productName)
	filename = "/Users/username/214/text/" + productName + ".txt"
	ioutil.WriteFile(filename, encPassword, 0644)

	
}


/*This function takes no arguments:
			designed to ask for a website or product to return the password
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
//setting the file name of the private key and loading the private key from memory
pem_file_path = product + ".pub"
pem_data, err = ioutil.ReadFile(pem_file_path)
 //checks to make sure key is loaded right
if err != nil {
        log.Fatalf("Error reading pem file: %s", err)
    }
    block, _ = pem.Decode(pem_data)
    
if block == nil || block.Type != "RSA PRIVATE KEY" {
        log.Fatal("No valid PEM data found")
    }
 private_key, err = x509.ParsePKCS1PrivateKey(block.Bytes)

if err != nil{
log.Fatalf("Private key can't be decoded: %s", err)
    }

//location of the encrypted password
encryptedFileName = "/Users/username/214/text/" + product + ".txt"
// reads in the contents of the file
encrypted, err = ioutil.ReadFile(encryptedFileName)
// checks for errors in reading in the file
if err != nil{
log.Fatalf("error in reading encrypted data: %s", err) 
}

fmt.Printf("encrypted version is %x\n", encrypted)
encrypted = bytes.Trim(encrypted, "\r")
encrypted = bytes.Trim(encrypted, "\n")
encrypted = bytes.TrimSpace(encrypted)
//calls the decrypt_oaep function to decrypt the encrypted passwords
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
	//initiating variables
	var private_key *rsa.PrivateKey
	var public_key *rsa.PublicKey
	var err error
	var plain_text, label []byte
	var filename string
	//converting product and password to byte arrays from strings
	plain_text = []byte(password)
	label = []byte(product)
	
	//Generate the Private Key
    if private_key, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {
        log.Fatal(err)
    }

    // Precompute some calculations 
    private_key.Precompute()

    //Validate Private Key 
    if err = private_key.Validate(); err != nil {
        log.Fatal(err)
    }

    //Public key address 
    public_key = &private_key.PublicKey

    encrypted = encrypt_oaep(public_key, plain_text, label)
    

	filename = product + ".pub"

	pemData := pem.EncodeToMemory(&pem.Block{
    		Type:  "RSA PRIVATE KEY",
    		Bytes: x509.MarshalPKCS1PrivateKey(private_key),
	})
	ioutil.WriteFile(filename, pemData, 0644)
	fmt.Printf("OAEP Encrypted [%s] to \n[%x]\n", string(plain_text), encrypted)
	





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
    decrypted, err = rsa.DecryptOAEP(md5_hash, rand.Reader, private_key, encrypted, label)
    if err != nil {
        log.Fatalf("failed in decrypt_oaep: %s", err)
    }
    return
}








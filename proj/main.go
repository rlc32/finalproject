/************************************************************************************************************
 * This is a GoLang password keeper. It is a command line application.										*
 * The password keeper utilizes OEAP encryption to securily store the passwords 							*
 * OEAP is Opitmal Asymetric Encryptinon Padding. This is a form of RSA Encryption 							*
 * with padding. The password keeper the offers the ability to add, see, edit, list, and remove passwords.	*
 * also if the command help is entered it will give the user some helpful hints about the program.			*
 * Created By: Russell Clousing																				*
 * Created On: May 5, 2016																					*
 ************************************************************************************************************/







package main
import (
	"fmt"
	"crypto/rsa"
	"crypto/rand"
	"log"
    "crypto/md5"
    "crypto/x509"
  	"hash"
	"io/ioutil"
	"encoding/pem"
	"os"
	)

func main(){
var repeat bool
repeat = true
var answer string
str1 := "Welcome to password keeper!\nDo you want to add a password or see a password?\n"
str2 := "Please enter the word 'add' to add a new password or enter 'see' to see a password\n or enter 'exit' to exit the program. Enter 'help' for more options:\n"
str3 := "Would you like to perform another action or exit the program"
fmt.Printf(str1)

	for repeat == true {																// this if block steers the program to correct actions
		fmt.Scanf("%s", &answer)														// also provides feedback for incorrect inputs
		if answer == "add"{																// when incorrect input happens the user is given a warning that
			addpass()																	// gives them direction to use the help command
			fmt.Println(str3)
			continue
		}
		if answer == "edit"{
			edit()
			fmt.Println(str3)
			continue
		}
		if answer == "see"{
			seepass()
			fmt.Println(str3)
			continue
		} 
		if answer == "list"{
			listPasswords()
			fmt.Println(str3)
			continue
		}
		if answer == "help"{
			goToHelp()
			fmt.Println(str3)
			continue
		}
		if answer == "remove"{
			removePass()
			fmt.Println(str3)
			continue
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
			asks for password and sends it to be encrypted using the key.
			also will ask user if they want to overwrite an existing file if there is an existing file*/

func addpass(){
	var password, productName, filename string
	var encPassword []byte
	var edit bool
	var editvalue string
	edit = true
	fmt.Printf("please enter the name of the product that you want to store: \n")
	fmt.Scanf("%s", &productName)																//gets the product name
	files, err := ioutil.ReadDir("/Users/russclousing/214/text")								//reads in all files in directory
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {																//for loop to make sure that user is not overwriting
		fileNoExt := file.Name()[:len(file.Name())-4]											// an already declared password unintentianaly
		if productName == fileNoExt{
			fmt.Printf("There already exists a password for %s.\n", productName)				// warning issued
			fmt.Println("Would you like to overwrite the existing password?")
			fmt.Scanf("%s", &editvalue)															// gets decision
			if editvalue == "no" || editvalue == "n"{											// if answer is no then skip rest of function
				edit = false
			}
		}

	}
	if edit == true {																			// if answer yes then adds password anyways
	fmt.Printf("Please enter the password for %s: \n", productName)
	//gets the password that should be encrypted
	fmt.Scanf("%s", &password)																	//gets password to update 
	encPassword = encryption(password, productName)												//sends password to be encrypted
	filename = "/Users/russclousing/214/text/" + productName + ".txt"							//file name for saving the password after encryption
	ioutil.WriteFile(filename, encPassword, 0644)												// write to file
	}
}


/*This function takes no arguments:
			designed to ask for a website or product to return the password
			retrieves thekey from a PEM file
			uses key to decrypt the password stored in separate file*/

func seepass(){
	var productName string
	fmt.Printf("please enter the name of the product that you want to see: \n")
	fmt.Scanf("%s", &productName)																// gets product name
	password := decryption(productName)															// gets decrypted password
	fmt.Printf("The password for [%s] is [%s].\n", productName, password)						// prints password
}




/*this is the function that takes a product and returns the password unencrypted
							the function finds the correct private key and reads in the key
							the function then finds the file that holds teh encrypted version of the password
							the key and encrypted password are then sent to be decrypted in the OAEP_decryption function
							returns unencrypted password*/
func decryption(product string)(password string){
	var private_key *rsa.PrivateKey
	var pem_data, encrypted, decrypted, label []byte
	var pem_file_path, encryptedFileName string
	var err error
	var block *pem.Block

	label = []byte(product)																//sets
												
	pem_file_path = product + ".pub"						//setting the file name of the private key and loading the private key from memory
	pem_data, err = ioutil.ReadFile(pem_file_path)
							
		if err != nil {																	 //checks to make sure key is loaded right
	        log.Fatalf("Error reading pem file: %s", err)
    }
    block, _ = pem.Decode(pem_data)														//populates block with []byte from pem file
    
	if block == nil || block.Type != "RSA PRIVATE KEY" {
        log.Fatal("No valid PEM data found")
    }
    
	 private_key, err = x509.ParsePKCS1PrivateKey(block.Bytes) 							//creating RSA key from the bytes in block

	if err != nil{
	log.Fatalf("Private key can't be decoded: %s", err)
    }

	encryptedFileName = "/Users/russclousing/214/text/" + product + ".txt"				//location of the encrypted password
	encrypted, err = ioutil.ReadFile(encryptedFileName)   								// reads in the contents of the file
	if err != nil {																		// checks for errors in reading in the file
        log.Fatalf("Error reading in the password: %s", err)
    }

	decrypted = decrypt_oaep(private_key, encrypted, label)								//calls the decrypt_oaep function to decrypt the encrypted passwords
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
	
    if private_key, err = rsa.GenerateKey(rand.Reader, 1024); err != nil {				//Generate the Private Key
        log.Fatal(err)
    }


    private_key.Precompute()    														// Precompute some calculations. speeds up future runtime

    //Validate Private Key 
    if err = private_key.Validate(); err != nil {
        log.Fatal(err)
    }

    //Public key address 
    public_key = &private_key.PublicKey													//gets public key from inside the private key struct

    encrypted = encrypt_oaep(public_key, plain_text, label)								//sends to OAEP encryption
    
    //gets file name for key location
	filename = product + ".pub"

	pemData := pem.EncodeToMemory(&pem.Block{											//Sets up the information for 
    		Type:  "RSA PRIVATE KEY",
    		Bytes: x509.MarshalPKCS1PrivateKey(private_key),
	})
	ioutil.WriteFile(filename, pemData, 0644)											//stores private key in a PEM file
	fmt.Println("Password was succesfully encrypted and stored")			

	return
}



//Code used from http://blog.giorgis.io/golang-rsa-encryption
//OAEP Encrypt. Optimal Asymetric Encryptoin Padding
func encrypt_oaep(public_key *rsa.PublicKey, plain_text, label []byte) (encrypted []byte) {
    var err error
    var md5_hash hash.Hash

    md5_hash = md5.New()

    if encrypted, err = rsa.EncryptOAEP(md5_hash, rand.Reader, public_key, plain_text, label); err != nil {	//encrypts using OAEP encryptino
        log.Fatal(err)
    }
    return
}
//Code used from http://blog.giorgis.io/golang-rsa-encryption
//OAEP Decryption. Optimal Asymetric Encryptoin Padding
func decrypt_oaep(private_key *rsa.PrivateKey, encrypted, label []byte) (decrypted []byte) {
    var err error
    var md5_hash hash.Hash
    
    
    md5_hash = md5.New()																	
    decrypted, err = rsa.DecryptOAEP(md5_hash, rand.Reader, private_key, encrypted, label)	//decrypts the password using OEAP encryption
    if err != nil {
        log.Fatalf("failed in decrypt_oaep: %s", err)
    }
    return
}

/*This function takes no arguments;
						This function lists all of the files that are holding passwords
						Does not show files that start with a '.'
						removes file extensions */
func listPasswords(){
	var fileNoExt string
	fmt.Println("Here is a list of all the products that have passwords stored")	
	files, err := ioutil.ReadDir("/Users/russclousing/214/text")				// creates list of all files in the specified directory
	if err != nil {
		log.Fatal(err)
	}
for _, file := range files {													//iterates through list of files
	if file.Name()[:1] == "."{													// ignores files that start with '.'
		continue
	}
	fileNoExt = file.Name()[:len(file.Name())-4]								//removes file extension
	fmt.Println(fileNoExt)														// prints the file name without extension
	
	}
}
/*This function prints out the help documentations:
							list functions and rules for using the program */
func goToHelp(){
	file_contents, err := ioutil.ReadFile("helpDoc.txt") 						//reads in help doc file
	if err != nil {
		log.Fatalf("Failed to read file!")
	}
	fmt.Printf(string(file_contents))											//prints out file contents to console
}
/*This function allows the user to replace a password for a product.
							also will check if the password exists in the system.
							if no password exists that matches the product then a warning is issued.
							if a password exists that matches the product then a new password is asked for
							and it is encrypted using a new RSA key and stored */
func edit(){
	var password, productName, filename string
	var encPassword []byte
	var editvalue 	bool
	editvalue = false

	fmt.Printf("please enter the name of the product that you want to edit: \n")
	fmt.Scanf("%s", &productName)												//gets the product name
	files, err := ioutil.ReadDir("/Users/russclousing/214/text")				//reads in all files in the directory
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {												// iterates through the list of files
		fileNoExt := file.Name()[:len(file.Name())-4]							// removes file extension
		if productName == fileNoExt{											// checks to find if product is already in directory
			editvalue = true													
			break
		}
	}

	if editvalue == true{														//this if statment will allow editing if the product is found
		fmt.Printf("Please enter the new password for %s: \n", productName)		// gets new password
		fmt.Scanf("%s", &password)												//gets the password that should be encrypted
		encPassword = encryption(password, productName)							// encrypts the password
		filename = "/Users/russclousing/214/text/" + productName + ".txt"		
		ioutil.WriteFile(filename, encPassword, 0644)							//stores password
	}
	if editvalue == false{														//this if activates if product is not found
		fmt.Println("The product you tried does not exist. Please add the product before trying to edit")	// warning issued in this case
	}
}



func removePass(){
	var isfile bool
	var textFileName, keyFileName, confirmation, productName string
	isfile = false
	fmt.Printf("please enter the name of the product that you want to remove: \n")
	fmt.Scanf("%s", &productName)												//gets the product name
	files, err := ioutil.ReadDir("/Users/russclousing/214/text")				//reads in all files in the directory
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {												// iterates through the list of files
		fileNoExt := file.Name()[:len(file.Name())-4]							// removes file extension
		if productName == fileNoExt{											// checks to find if product is already in directory
			isfile = true													
			break
		}
	}
	if isfile == true{
		fmt.Printf("Are you sure that you want to remove this password?\nThis action cannot be undone.\n")
		fmt.Scanf("%s", &confirmation)
		if confirmation == "yes" || confirmation == "y"{
			textFileName = "/Users/russclousing/214/text/" + productName + ".txt"	//sets file name location
			keyFileName = productName + ".pub"		
			os.Remove(textFileName)													// removes the file containing the password 
			os.Remove(keyFileName)													// removes the file containing the RSA key associated with the password
			fmt.Println("Password succesfully removed")
		}else{fmt.Println("Remove canceled.")}										// confirmation that remove action canceled
	}else{fmt.Println("The password that you attempted to remove does not exist")}	// Warning that the file does not exist
}
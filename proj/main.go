package main
import (
	"fmt"
	/*"io/ioutil"*/
	"os"
	"time"
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
	password := ""
	productName := ""
	fmt.Printf("please enter the name of the product that you want to store: \n")
	fmt.Scanf("%s", &productName)
	/*err := ioutil.WriteFile("214/text/" + productName + ".txt", []byte(password), 0644)
	check(err)*/
	f, err := os.Create("/Users/russclousing/214/text/ " + productName + ".txt")
	check(err)
	defer f.Close()
	n1, err := f.Write([]byte(password))
	check(err)
	fmt.Printf("wrote %d bytes\n", n1)
	f.Sync()
	fmt.Printf("please work like i want it to\n")
	
}


/*This function takes no arguments:
			designed to ask for a website or product to return the username and password
			retrieves the username and key from file
			uses key to decrypt the password stored in separate file*/

func seepass(){

fmt.Printf("tell me this works as well\n")
}

/* this checks for errors in file usage*/
func check(e error){
	if e != nil{
		panic(e)}
}


/* This function encrypts the password and returns the password as a encrypted string of numbers
			Input: this funciton takes a string
			Output: THis returns a string of numerals that is the encrypted string*/
func encryption(password string){
	fmt.Printf(Now())


}






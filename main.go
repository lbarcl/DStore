package main

import (
	"fmt"
	"log"
)

func main() {
	//pattern, _ := regexp.Compile(`^(?:[a-zA-Z]:\\|\/)?(?:[^\/\\\0]+[\/\\])*(?:[^\/\\\0]+)$`)

	var input string
l:
	for {
		fmt.Scanln(&input)
		CallClear()
		switch input {
		case "help":
			// Prints commands
			fmt.Println("Welcome to the DClient unlimted storage system.")
			fmt.Println("Commands ----------------------------------------------------------------------------")
			fmt.Println("upload")
			fmt.Println("Uplooads the given file to the storage system.")
			fmt.Println("exit")
			fmt.Println("Closes the program.")
		case "upload":
			fmt.Print("Enter the file path: ")
			fmt.Scanln(&input)

			/*
				fmt.Println(pattern.MatchString(input))

				if !pattern.MatchString(input) {
					fmt.Println("The file path you entered is not valid!")
					continue
				}

			*/

			err := Upload(input, callBack)
			if err != nil {
				log.Fatalf("There was an error uploading the file: %v", err)
			}

		case "exit":
			// Breaks the loop
			break l
		}
	}
}

func callBack() {

}

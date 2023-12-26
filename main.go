package main

import (
	"fmt"
	"lbarcl/DStore/Records"
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
			fmt.Println("resume")
			fmt.Println("When the upload is intercepted or an error is thrown,\nIt flags the file as paused and you can resume it with File ID that you got with list command.")
			fmt.Println("list")
			fmt.Println("Lists the uploaded/uploading files.")
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
				} */

			err := Upload(input)
			if err != nil {
				fmt.Println("We have paused the upload process. You can resume it later")
				log.Fatalf("There was an error uploading the file: %v", err)
			}
		case "resume":
			fmt.Print("Enter the file ID to resume: ")
			fmt.Scanln(&input)

			file, err := Records.GetFile(input)
			if err != nil {
				log.Fatal(err)
			}

			if !file.Paused {
				fmt.Println("The file you are trying to resume is not paused!")
				continue
			}

			err = Resume(file)
			if err != nil {
				fmt.Println("We have paused the upload process. You can resume it later")
				log.Fatalf("There was an error uploading the file: %v", err)
			}

		case "list":
			files, err := Records.GetAllFiles()
			if err != nil {
				log.Fatal(err)
			}

			for i, f := range files {
				fmt.Printf("%d) %s -> %s\n", i, f.Name, f.ID)
			}

		case "exit":
			// Breaks the loop
			break l
		}
	}
}

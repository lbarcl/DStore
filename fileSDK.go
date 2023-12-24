package main

import (
	"lbarcl/DClient/Discord"
	"os"
)

const maxAttachmentSize = 20_000_000

type Callback func(map[string]Discord.Message, string, bool)

func Upload(filePath string, cb Callback) error {
	fileInfo, err := os.Lstat(filePath)

	if err != nil {
		return err
	}

	//fmt.Println("File found!")

	var parts int = 0
	d := fileInfo.Size() / maxAttachmentSize

	if d < 1 {
		parts = 1
	} else if d > 1 {
		if (d % maxAttachmentSize) != 0 {
			d += 1
		}

		parts = int(d)
	}

	//fmt.Println("File parts calculated!", parts)

	file, err := os.Open(filePath)

	//fmt.Println("File opend!")

	if err != nil {
		return err
	}

	defer file.Close()

	for i := 0; i < parts; i++ {
		var buffer []byte
		if parts == 1 {
			buffer = make([]byte, fileInfo.Size())
		} else if i == parts-1 {
			buffer = make([]byte, fileInfo.Size()-int64(maxAttachmentSize*i))
		} else {
			buffer = make([]byte, maxAttachmentSize)
		}

		file.ReadAt(buffer, int64(i*maxAttachmentSize))
		//fmt.Println("File read!", len(buffer))

		msg, err := Discord.SendMultipart(buffer)

		if i == parts-1 {
			cb(msg, "", true)
		} else {
			cb(msg, "", false)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

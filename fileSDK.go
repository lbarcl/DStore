package main

import (
	"io"
	"lbarcl/DStore/Discord"
	"lbarcl/DStore/Records"
	"os"

	"github.com/schollz/progressbar/v3"
)

const maxAttachmentSize = 10_000_000

func Upload(filePath string) error {
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

	fileRecord, err := Records.NewFile(fileInfo.Name(), fileInfo.Size(), filePath, parts)
	bar := progressbar.DefaultBytes(fileRecord.Size, "Uploading "+fileRecord.Name)

	if err != nil {
		return err
	}
	i := 0
	for i < parts {
		var buffer []byte
		if parts == 1 {
			buffer = make([]byte, fileInfo.Size())
		} else if i == parts-1 {
			buffer = make([]byte, fileInfo.Size()-int64(maxAttachmentSize*i))
		} else {
			buffer = make([]byte, maxAttachmentSize)
		}

		_, err := file.ReadAt(buffer, int64(i*maxAttachmentSize))
		//fmt.Println("Part readed ", len(buffer), " bytes. ", i, " Sequence.")
		if err != nil {
			if err == io.EOF {
				msg, err := Discord.SendMultipart(buffer)
				if err != nil {
					fileRecord.Pause()
					return err
				}
				fileRecord.AddPart(msg.Attachments[0].ID, msg.ID, i, Hash(buffer), int64(msg.Attachments[0].Size))
				bar.Add(len(buffer))
				break
			}
			return err
		}

		msg, err := Discord.SendMultipart(buffer)
		if err != nil {
			fileRecord.Pause()
			return err
		}
		fileRecord.AddPart(msg.Attachments[0].ID, msg.ID, i, Hash(buffer), int64(msg.Attachments[0].Size))
		bar.Add(len(buffer))
		i++
	}
	bar.Finish()
	return nil
}

func Resume(f *Records.File) error {
	file, err := os.Open(f.LocalPath)

	//fmt.Println("File opend!")

	if err != nil {
		return err
	}

	defer file.Close()

	bar := progressbar.DefaultBytes(f.Size, "Uploading "+f.Name)
	bar.Add(f.RecordedParts * maxAttachmentSize)

	if err != nil {
		return err
	}
	i := f.RecordedParts - 1
	for i < f.TotalParts {
		var buffer []byte
		if f.TotalParts == 1 {
			buffer = make([]byte, f.Size)
		} else if i == f.TotalParts-1 {
			buffer = make([]byte, f.Size-int64(maxAttachmentSize*i))
		} else {
			buffer = make([]byte, maxAttachmentSize)
		}

		_, err := file.ReadAt(buffer, int64(i*maxAttachmentSize))

		if err != nil {
			if err == io.EOF {
				msg, err := Discord.SendMultipart(buffer)
				if err != nil {
					f.Pause()
					return err
				}
				f.AddPart(msg.Attachments[0].ID, msg.ID, i, Hash(buffer), int64(msg.Attachments[0].Size))
				bar.Add(len(buffer))
				break
			}
			return err
		}

		msg, err := Discord.SendMultipart(buffer)
		if err != nil {
			f.Pause()
			return err
		}
		f.AddPart(msg.Attachments[0].ID, msg.ID, i, Hash(buffer), int64(msg.Attachments[0].Size))
		bar.Add(len(buffer))
		i++
	}
	bar.Finish()
	return nil
}

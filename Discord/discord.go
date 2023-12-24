package Discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

const channelID = "1037079784229978123"
const baseURL = "https://discord.com/api/v10"
const botToken = "Nzg4NDUyMDg4MTM2MjA0MzYx.GS0Rfj.rh5Ii_eY02Hbe5cR6OT0OMXpaue53mRFgl7iRs"

func SendMultipart(buffer []byte) (map[string]Message, error) {
	var body bytes.Buffer

	writer := multipart.NewWriter(&body)
	writer.SetBoundary("boundary")

	header := textproto.MIMEHeader{}
	header.Add("Content-Disposition", "form-data; name=\"payload_json\"")
	header.Set("Content-Type", "application/json")
	jsonWriter, err := writer.CreatePart(header)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"content": "hello world!",
		"attachments": []interface{}{map[string]interface{}{
			"id":          0,
			"description": "Part of a file",
			"filename":    "part.fbin",
		}},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	jsonWriter.Write(payloadBytes)

	bufferWriter, err := writer.CreateFormFile("files[0]", "part.fbin")
	if err != nil {
		return nil, err
	}

	bufferWriter.Write(buffer)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/channels/%s/messages", baseURL, channelID)
	request, err := http.NewRequest("POST", url, &body)
	//fmt.Println("Multipart created!", request)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("authorization", fmt.Sprintf("Bot %s", botToken))

	//fmt.Println(request.URL)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		//fmt.Println("Status OK")
		var result map[string]Message
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			return nil, err
		}

		return result, nil
	} else {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("error uploading file. status: %s\nresponse Body: %s", response.Status, body)
	}
}

package whatsapp

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"xanny-go-template/pkg/exceptions"
)

func Send(target string, message string) *exceptions.Exception {
	url := "https://api.fonnte.com/send"

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("target", target)
	writer.WriteField("message", message)
	writer.Close()

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return exceptions.NewException(http.StatusInternalServerError, fmt.Sprintf("error creating request: %v", err))
	}

	req.Header.Set("Authorization", os.Getenv("FONNTE_API_KEY"))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return exceptions.NewException(http.StatusInternalServerError, fmt.Sprintf("error sending request: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return exceptions.NewException(http.StatusInternalServerError, fmt.Sprintf("error reading response: %v", err))
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(body))

	return nil
}

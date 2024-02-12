package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Embed struct {
	URL string `json:"url"`
}

type CastBody struct {
	EmbedsDeprecated  []string `json:"embedsDeprecated"`
	Mentions          []int64  `json:"mentions"`
	ParentUrl         string   `json:"parentUrl"`
	Text              string   `json:"text"`
	MentionsPositions []int16  `json:"mentionsPositions"`
	Embeds            []Embed  `json:"embeds"`
}

type Data struct {
	Type        string   `json:"type"`
	FID         int64    `json:"fid"`
	Timestamp   int64    `json:"timestamp"`
	Network     string   `json:"network"`
	CastAddBody CastBody `json:"castAddBody"`
}

type Messages struct {
	Data            Data   `json:"data"`
	Hash            string `json:"hash"`
	HashScheme      string `json:"hashScheme"`
	SignatureScheme string `json:"signatureScheme"`
	Signer          string `json:"signer"`
}

type MemeResponse struct {
	Messages      []Messages `json:"messages"`
	NextPageToken string     `json:"nextPageToken"`
}

func GetRandomMeme() string {
	url := "https://hub.pinata.cloud/v1/castsByParent?url=chain://eip155:1/erc721:0xfd8427165df67df6d7fd689ae67c8ebf56d9ca61"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making the request: %v\n", err)
		return "Error"
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return "Error"
	}

	// Parse the JSON response into the ApiResponse struct
	var apiResponse MemeResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return "Error"
	}

	var matchingURLs []string

	for _, message := range apiResponse.Messages {
		// Iterate over the embeds in each message
		for _, embed := range message.Data.CastAddBody.Embeds {
			// Check if the URL ends with .png, .jpg, or .gif
			if strings.HasSuffix(embed.URL, ".png") || strings.HasSuffix(embed.URL, ".jpg") || strings.HasSuffix(embed.URL, ".gif") {
				fmt.Println("Found matching URL:", embed.URL)
				matchingURLs = append(matchingURLs, embed.URL)
			}
		}
	}

	// Check if there are any matching URLs
	if len(matchingURLs) == 0 {
		fmt.Println("No matching URLs found")
		return "Error"
	}

	// Generate a random index within the bounds of matchingURLs slice
	randomIndex := rand.Intn(len(matchingURLs))

	// Select the URL at the randomly generated index
	selectedURL := matchingURLs[randomIndex]

	// Use the data from the response
	fmt.Printf("Response Data: %+v\n", selectedURL)
	return selectedURL
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Parse templates
	templatesDir := "./templates"
	template := template.Must(template.ParseFiles(filepath.Join(templatesDir, "template.html")))

	// Handle requests on the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set content type to HTML
		w.Header().Set("Content-Type", "text/html")
		//	Get random meme
		randomMemeUrl := GetRandomMeme()

		if randomMemeUrl == "Error" {
			fmt.Println("Error getting meme")
			return
		}
		// Check the request method and serve the appropriate template
		switch r.Method {
		case "GET":
			if err := template.Execute(w, randomMemeUrl); err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
			}
		case "POST":
			if err := template.Execute(w, randomMemeUrl); err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
			}
		default:
			// Respond with a 405 Method Not Allowed if the method is not GET or POST
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "<h1>405 Method Not Allowed</h1>")
		}
	})

	// Start the HTTP server on port 8080 and handle errors
	fmt.Println("Server listening on port", port)
	httpError := http.ListenAndServe(":"+port, nil)
	if httpError != nil {
		fmt.Println("Error starting server: ", httpError)
		return
	}
}

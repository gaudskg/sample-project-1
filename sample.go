package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form input
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	containerNumber := r.FormValue("container_number")
	if containerNumber == "" {
		http.Error(w, "Container number is required", http.StatusBadRequest)
		return
	}

	// Make POST request to fetch data
	client := &http.Client{}
	data := strings.NewReader(`{"variables":{"{\"__wwtype\":\"f\",\"code\":\"variables['465276e0-a224-4a28-bf2c-098ddca8ca4b']\"}":"62bd38e5-6c23-4f5b-839e-57086cc492f1","{\"__wwtype\":\"f\",\"code\":\"variables['148f038e-244a-44c4-bcb9-750a64e74167-value']||variables['654fb0ae-1573-42f0-abce-648701780a13-value']||variables['3300cf30-27b8-4c61-9de6-b7d1bd59af48']\"}":"` + containerNumber + `","{\"__wwtype\":\"f\",\"code\":\"variables['1bb4ebcb-2021-4770-99f3-c7a00de242ca-value']||\\\"\\\"\"}":"","{\"__wwtype\":\"f\",\"code\":\"variables['a03e8a5d-fd4b-43b8-b206-1f95f441753c-value']||\\\"\\\"\"}":""}}`)
	req, err := http.NewRequest("POST", "https://tools.sinay.ai/ww/cms_data_sets/97105c3b-3a4b-41e4-a2a9-f102ad05c44f/fetch?limit=100&offset=0", data)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Set response content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Display response body in the HTML page
	tmpl := `<html>
<head>
    <title>Data Response</title>
</head>
<body>
    <h1>Container Details</h1>
    <pre>%s</pre>
</body>
</html>`
	fmt.Fprintf(w, tmpl, string(body))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler(w, r) // If method is POST, handle the form submission
			return
		}

		// For GET requests, display the form to input the container number
		tmpl := `<html>
<head>
    <title>Container Number Input</title>
</head>
<body>
    <h1>Enter Container Number</h1>
    <form method="post">
        <label for="container_number">Container Number:</label>
        <input type="text" id="container_number" name="container_number" required>
        <button type="submit">Submit</button>
    </form>
</body>
</html>`
		fmt.Fprint(w, tmpl)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port := "3000"
	}
	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

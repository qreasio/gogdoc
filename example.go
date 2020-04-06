package main

import (
	gDoc "gogdoc/pkg"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	credentialFile := "credentials.json" // This file is downloaded from Google API Console

	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes configuration, delete the previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tokenFile := "token.json"
	client := gDoc.GetClient(config, tokenFile)

	srv, err := docs.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Docs client: %v", err)
	}

	// Original Google doc url https://docs.google.com/document/d/1PrS4u6SBGVZj__GI3tYrqcOWBc2w0FQt7-Dxs_x1mp0/edit
	docId := "1PrS4u6SBGVZj__GI3tYrqcOWBc2w0FQt7-Dxs_x1mp0"
	doc, err := srv.Documents.Get(docId).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from document: %v", err)
	}
	log.Printf("The title of the doc is: %s\n", doc.Title)

	docJSON, err := doc.MarshalJSON()
	if err != nil {
		log.Fatalf("Unable to marshal JSON: %v", err)
	}

	f, err := os.Create("doc.json")
	if err != nil {
		log.Fatalf("Failed to create file")
	}
	defer f.Close()
	totalBytes, err := f.WriteString(string(docJSON))
	log.Printf("wrote %d bytes\n", totalBytes)

	gDoc := gDoc.Document{GDoc: *doc}
	if err != nil {
		log.Fatalf("Unable to marshal JSON: %v", err)
	}

	mdOutput := gDoc.Markdown()
	if err != nil {
		log.Fatalf("Unable to convert to Markdown: %v", err)
	}
	mdFile, err := os.Create("doc.md")
	if err != nil {
		log.Fatalf("Failed to create file markdown")
	}
	defer mdFile.Close()

	_, err = mdFile.WriteString(string(mdOutput))
	if err != nil {
		log.Fatalf("Failed to write file markdown")
	}

}

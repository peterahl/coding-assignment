package main

import (
	"log"
	"os"
	"sort"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	// "github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

var providerIndex *ProviderIndex

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
	)

	m := make(map[string]string)
	m["github"] = "Github"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	providerIndex = &ProviderIndex{Providers: keys, ProvidersMap: m}
	log.Println("Init for auth done", providerIndex)
}

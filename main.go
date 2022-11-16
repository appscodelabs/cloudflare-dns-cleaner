package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	domain := "appscode.cloud"
	flag.StringVar(&domain, "domain", domain, "Domain name")
	flag.Parse()

	// https://github.com/tamalsaha/cloudflare-dns-proxy
	cfProxy := "http://127.0.0.1:63707"
	// , cloudflare.BaseURL(cfProxy)
	fmt.Println(cfProxy)
	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"), cloudflare.BaseURL(cfProxy))
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch the zone ID
	id, err := api.ZoneIDByName(domain) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Zone ID:", id)
	os.Exit(1)

	records, err := api.DNSRecords(ctx, id, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range records {
		fmt.Println(r.Content)
		if err := api.DeleteDNSRecord(ctx, id, r.ID); err != nil {
			log.Fatal(err)
		}
	}
}

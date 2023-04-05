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
	fmt.Println("Domain to delete all DNS records from: ", domain)

	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
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
	
	zoneId := cloudflare.ZoneIdentifier(id)
	records, _, err := api.ListDNSRecords(ctx, zoneId, cloudflare.ListDNSRecordsParams{})
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range records {
		fmt.Printf("%s: %s", r.Name, r.Content)
		if err := api.DeleteDNSRecord(ctx, zoneId, r.ID); err != nil {
			log.Fatal(err)
		}
		fmt.Println(" DELETED")
	}
}

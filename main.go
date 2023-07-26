package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/publicsuffix"
)

func main() {
	fqdn := ""
	flag.StringVar(&fqdn, "fqdn", fqdn, "Domain name")
	flag.Parse()

	if fqdn == "" {
		panic("set --fqdn flag")
	}

	ListDNSRecords(fqdn, false)
}

func ListDNSRecords(fqdn string, del bool) error {
	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		return err
	}

	domain, err := publicsuffix.EffectiveTLDPlusOne(fqdn)
	if err != nil {
		return err
	}

	ctx := context.Background()

	id, err := api.ZoneIDByName(domain)
	if err != nil {
		return err
	}
	fmt.Println("Zone ID:", id)

	rc := cloudflare.ResourceContainer{
		Type:       cloudflare.ZoneType,
		Identifier: id,
	}
	records, _, err := api.ListDNSRecords(ctx, &rc, cloudflare.ListDNSRecordsParams{})
	if err != nil {
		return err
	}
	for _, r := range records {
		if r.Name == fqdn || strings.HasSuffix(r.Name, "."+fqdn) {
			fmt.Println(r.Type, r.Name)
			if del {
				if err := api.DeleteDNSRecord(ctx, &rc, r.ID); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

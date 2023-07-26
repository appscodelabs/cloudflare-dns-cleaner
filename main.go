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
	del := false
	flag.StringVar(&fqdn, "fqdn", fqdn, "Domain name")
	flag.BoolVar(&del, "delete", del, "If true, deletes DNS records")
	flag.Parse()

	if fqdn == "" {
		panic("set --fqdn flag")
	}

	if del {
		// ask for confirmation
		fmt.Println("Are you sure you want to delete DNS records? [Y/N]")
		if !askForConfirmation() {
			os.Exit(0)
		}
	}
	ListDNSRecords(fqdn, del)
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

// ref: https://gist.github.com/albrow/5882501

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		panic(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

// You might want to put the following two functions in a separate utility package.

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

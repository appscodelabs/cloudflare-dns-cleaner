# cloudflare-dns-cleaner

DELETE all dns records from a Cloudflare Zone

Usage:

Set your cloudflare token as an env variable
export CLOUDFLARE_API_TOKEN="..."

Run go with a -domain flag specifying your domain
go run . -domain your_domain.com
_

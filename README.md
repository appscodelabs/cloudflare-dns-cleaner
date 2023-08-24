# cloudflare-dns-cleaner

DELETE all dns records from a Cloudflare Zone

## Installation

```bash
go install github.com/appscodelabs/cloudflare-dns-cleaner@master
```

## Usage

```bash
export CLOUDFLARE_API_TOKEN=***

# list records
cloudflare-dns-cleaner --fqdn=example.com

# delete records
cloudflare-dns-cleaner --fqdn=example.com --delete
```

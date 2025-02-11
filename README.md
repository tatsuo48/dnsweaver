# dnsweaver

A tool for converting Cloudflare DNS records to Terraform configurations.

## Overview

This repository uses Go and Go templates to extract DNS records from Cloudflare and convert them into Terraform configuration files (`records.tf` and `import.tf`).

## Setup

1. Set your Cloudflare API token in the `CLOUDFLARE_API_TOKEN` environment variable. Ensure the token has DNS List permissions.
2. Set your Cloudflare Zone ID in the `CLOUDFLARE_ZONE_ID` environment variable.
3. Ensure Go is installed.
4. Clone the repository and build and run the project:

   ```bash
   go run main.go
   ```

## Files

- `main.go`: Main application logic.
- `go.tmpl`: Template for generating Terraform DNS record resources.
- `import.tmpl`: Template for generating Terraform import statements.
- `records.tf` and `import.tf`: Generated Terraform configuration files.

## License

MIT

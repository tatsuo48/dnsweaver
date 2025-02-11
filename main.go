package main

import (
	"context"
	"fmt"
	"os"
	"strings" // added import
	"text/template"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

type DNSRecord struct {
	ZoneID                string
	Name                  string
	Content               string
	DNSType               string
	Comment               string
	Proxied               bool
	TTL                   int
	ID                    string
	TerraformResourceName string
	Priority              int // added field for MX records
}

func main() {
	client := cloudflare.NewClient(option.WithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN")))
	// list dns records
	records, err := client.DNS.Records.List(context.TODO(), dns.RecordListParams{ZoneID: cloudflare.F(os.Getenv("CLOUDFLARE_ZONE_ID"))})
	if err != nil {
		panic(err)
	}
	// create file
	f, err := os.Create("terraform/records.tf")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f2, err := os.Create("terraform/import.tf")
	if err != nil {
		panic(err)
	}
	defer f2.Close()
	for _, record := range records.Result {
		// Escape double quotes in record.Content if present
		escapedContent := strings.Replace(record.Content, `"`, `\"`, -1)
		// Generate Terraform resource name, avoid starting with a digit
		terraformResourceName := fmt.Sprintf("%s_%s", strings.Replace(record.Name, ".", "_", -1), string(record.Type))
		if len(terraformResourceName) > 0 && terraformResourceName[0] >= '0' && terraformResourceName[0] <= '9' {
			terraformResourceName = "_" + terraformResourceName
		}
		// Set priority if type is MX
		priority := 0
		if string(record.Type) == "MX" {
			priority = int(record.Priority)
		}
		dnsrecord := DNSRecord{
			ZoneID:                os.Getenv("CLOUDFLARE_ZONE_ID"),
			Name:                  record.Name,
			Content:               escapedContent, // using escaped content
			DNSType:               string(record.Type),
			Comment:               record.JSON.Comment.Raw(),
			Proxied:               record.Proxied,
			TTL:                   int(record.TTL),
			ID:                    record.ID,
			TerraformResourceName: terraformResourceName,
			Priority:              priority,
		}
		t, err := template.New("records.tmpl").ParseFiles("records.tmpl")
		if err != nil {
			panic(err)
		}
		if err := t.Execute(f, dnsrecord); err != nil {
			panic(err)
		}
		fmt.Fprintln(f)

		t2, err := template.New("import.tmpl").ParseFiles("import.tmpl")
		if err != nil {
			panic(err)
		}
		if err := t2.Execute(f2, dnsrecord); err != nil {
			panic(err)
		}
		fmt.Fprintln(f2)
	}
}

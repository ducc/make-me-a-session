package main

import (
	"fmt"
	"github.com/respawner/peeringdb"
	"github.com/tcnksm/go-input"
	"os"
	"strconv"
	"strings"
)

func main() {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	asn, _ := ui.Ask("What ASN?", &input.Options{
		Default:  "51391",
		Required: true,
		Loop:     true,
	})

	asNum, err := strconv.ParseInt(asn, 10, 0)
	if err != nil {
		fmt.Println("Please enter the AS in the format '12345' as apposed to 'as12345'")
	}

	wantedIx, _ := ui.Ask("What IX?", &input.Options{
		Required: true,
		Loop:     true,
	})

	password, _ := ui.Ask("Any password? Enter if so", &input.Options{
		Required: false,
	})

	api := peeringdb.NewAPI()
	as := api.GetASN(int(asNum))
	fmt.Println("Making session for", as.Name)
	ixList := as.NetworkInternetExchangeLANSet

	for _, ixId := range ixList {
		ix, err := api.GetNetworkInternetExchangeLANByID(ixId)
		if err == nil {
			if strings.Contains(strings.ToLower(ix.Name), strings.ToLower(wantedIx)) {
				fmt.Println("Making session for", as.Name, "on", ix.Name)
				if ix.IPAddr4 != "" {
					fmt.Println(ix.Name, "IPv4")
					fmt.Printf("\n\n")
					fmt.Printf("      %d:\n", as.ASN)
					fmt.Printf("        description: %s\n", as.Name)
					fmt.Printf("        neighbor: %s\n", ix.IPAddr4)
					fmt.Printf("        source_address:\n")
					fmt.Printf("        prepend: 0\n")
					fmt.Printf("        prefix_limit: %d\n", as.InfoPrefixes4)
					fmt.Printf("        limit_exceeded_action: disable\n")
					if password != "" {
						fmt.Printf("        password: %s\n", password)
					}
					fmt.Printf("\n")
				}

				if ix.IPAddr6 != "" {
					fmt.Println(ix.Name, "IPv6")
					fmt.Printf("\n\n")
					fmt.Printf("      %d:\n", as.ASN)
					fmt.Printf("        description: %s\n", as.Name)
					fmt.Printf("        neighbor: %s\n", ix.IPAddr6)
					fmt.Printf("        source_address:\n")
					fmt.Printf("        prepend: 0\n")
					fmt.Printf("        prefix_limit: %d\n", as.InfoPrefixes6)
					fmt.Printf("        limit_exceeded_action: disable\n")
					if password != "" {
						fmt.Printf("        password: %s\n", password)
					}

					fmt.Printf("\n")
				}

			}
		} else {
			fmt.Println("Error whilst getting IX")
		}
	}
}

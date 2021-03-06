package cmd

import (
	"os"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var domainRecordListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: "civo domain record ls DOMAIN/DOMAIN_ID",
	Args:    cobra.MinimumNArgs(1),
	Short:   "List all domains records",
	Long: `List all current domain records.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* Value
	* Type
	* TTL
	* Priority`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		domain, err := client.FindDNSDomain(args[0])
		if err != nil {
			utility.Error("Unable to find domain for your search %s", err)
			os.Exit(1)
		}

		records, err := client.ListDNSRecords(domain.ID)
		if err != nil {
			utility.Error("Unable to list domains %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, record := range records {
			ow.StartLine()

			ow.AppendData("ID", record.ID)
			ow.AppendData("Name", record.Name)
			ow.AppendData("Value", record.Value)
			ow.AppendData("Type", string(record.Type))
			ow.AppendData("TTL", strconv.Itoa(record.TTL))
			ow.AppendData("Priority", strconv.Itoa(record.Priority))

		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveDefault
		}
		return getDomainList(toComplete), cobra.ShellCompDirectiveDefault
	},
}

func getDomainList(value string) []string {
	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	domain, err := client.FindDNSDomain(value)
	if err != nil {
		utility.Error("Unable to list domains %s", err)
		os.Exit(1)
	}

	var domainList []string
	domainList = append(domainList, domain.Name)

	return domainList

}

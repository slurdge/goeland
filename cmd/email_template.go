package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// emailTemplateCmd represents the email-template command
var emailTemplateCmd = &cobra.Command{
	Use:   "output-email-template",
	Short: "Output the default email template",
	Long: `Output the default email template that can be used for customization.

This command outputs the default email template to stdout, which you can redirect
to a file for customization. For example:

goeland email-template > my-template.html

You can then use this customized template by setting the email.template configuration
option to point to your file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Output the default email template
		fmt.Print(string(defaultEmailBytes))
	},
}

func init() {
	rootCmd.AddCommand(emailTemplateCmd)
}

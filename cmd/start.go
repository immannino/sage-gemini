package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/disgoorg/log"
	"github.com/immannino/sage-gemini/internal"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	userLicense string
	rootCmd     = &cobra.Command{
		Use:   "sage",
		Short: "Sage is a utility to crawl sitemaps and bulk email to an address.",
		Long:  `Sage is a utility to crawl sitemaps and bulk email articles to my new Sage colored Kindle, Sage Gemini.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "sage version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("sage version v0.0.1 -- HEAD")
		},
	}

	sitemapUrl = ""
	fetchCmd   = &cobra.Command{
		Use:   "fetch",
		Short: "fetch makes an http request to a sitemap and logs it's contents",
		Run: func(cmd *cobra.Command, args []string) {
			if sitemapUrl == "" {
				log.Fatalf("no sitemap provided, exiting.")
			}

			sitemap, err := internal.FetchSitemap(sitemapUrl)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(sitemap.XMLName.Local + "\n" + strings.Repeat("-", 32))
			for _, v := range sitemap.URL {
				fmt.Printf("%s %s %s %v\n", v.Loc, v.LastMod, v.ChangeFreq, v.Priority)
			}
		},
	}

	recipient = ""
	sendCmd   = &cobra.Command{
		Use:   "send",
		Short: "send fetches a sitemap, parses, and sends the contents to the recipient.",
		Run: func(cmd *cobra.Command, args []string) {
			if sitemapUrl == "" {
				log.Fatalf("no sitemap provided, exiting.")
			}

			var to string
			if to == "" {
				to = os.Getenv("RECIPIENT")
			}
			if to == "" {
				log.Fatalf("no email recipient provided.\n%s", cmd.Usage())
			}

			ctx := cmd.Context()
			sitemap, err := internal.FetchSitemap(sitemapUrl)
			if err != nil {
				log.Fatal(err)
			}

			for _, v := range sitemap.URL {
				contents, err := internal.FetchHTML(ctx, v.Loc)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Fetched contents for ", v.Loc, " content-length: ", len(contents))
			}
		},
	}

	testSendCmd = &cobra.Command{
		Use:   "test-send",
		Short: "test-send sends a test email to the recipient to ensure auth is working.",
		Run: func(cmd *cobra.Command, args []string) {
			var to string
			if to == "" {
				to = os.Getenv("RECIPIENT")
			}
			if to == "" {
				log.Fatalf("no email recipient provided.\n%s", cmd.Usage())
			}

			contents, err := internal.FetchHTML(cmd.Context(), "https://text.npr.org/")
			if err != nil {
				log.Fatal(err)
			}

			// sender := email.New(email.SenderOpts{
			// 	Host:       os.Getenv("EMAIL_HOST"),
			// 	Username:   os.Getenv("EMAIL_USERNAME"),
			// 	Password:   os.Getenv("EMAIL_PASSWORD"),
			// 	PortNumber: os.Getenv("EMAIL_PORT"),
			// })
			// m := email.NewMessage("sage gemini test email", "")
			// m.To = []string{to}
			f, err := ioutil.TempFile(path.Join("./"), "sage.*.html")
			if err != nil {
				log.Fatal("open temp file error ", err)
			}
			if _, err := f.WriteString(contents); err != nil {
				log.Fatal("write contents error ", err)
			}
			if err := f.Close(); err != nil {
				log.Fatal("file closer error ", err)
			}
			defer os.RemoveAll(f.Name())

			if err := internal.Send(f.Name()); err != nil {
				log.Error(err)
			} else {
				log.Info("Tee hee")
			}
		},
	}
)

func init() {
	viper.SetDefault("author", "Tony Mannino <tony@ope.cool>")
	viper.SetDefault("license", "mit")
}

func Start() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	fetchCmd.Flags().StringVarP(&sitemapUrl, "sitemap", "s", "", "The sitemap to fetch urls from")
	sendCmd.Flags().StringVarP(&sitemapUrl, "sitemap", "s", "", "The sitemap to fetch urls from")
	sendCmd.Flags().StringVarP(&recipient, "recipient", "r", "", "The email address to send content to")
	rootCmd.AddCommand(versionCmd, fetchCmd, sendCmd, testSendCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

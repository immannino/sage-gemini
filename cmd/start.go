package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

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

			var to string
			if to == "" {
				to = os.Getenv("RECIPIENT")
			}
			if to == "" {
				log.Fatalf("no email recipient provided.\n%s", cmd.Usage())
			}

			fmt.Println(sitemap.XMLName.Local + "\n" + strings.Repeat("-", 32))
			for _, v := range sitemap.URL {
				fmt.Printf("%s %s %s %f\n", v.Loc, v.ChangeFreq, v.LastMod, v.Priority)
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

			log.Debug("fetching sitemap")
			sitemap, err := internal.FetchSitemap(sitemapUrl)
			if err != nil {
				log.Fatal(err)
			}

			log.Debug("parsing entires")
			for _, v := range sitemap.URL {
				contents, err := internal.FetchHTML(cmd.Context(), v.Loc)
				if err != nil {
					log.Fatal(err)
				}

				title, err := internal.ParseTitle(contents)
				if err != nil {
					log.Error(err)
					title = v.Loc
				}

				subject := strings.TrimSpace(title)
				subject = strings.Trim(subject, "\n")
				subject = fmt.Sprintf("%s.html", subject)

				log.Debug("sending email")
				if err := internal.Send(to, subject, contents); err != nil {
					log.Error(err)
				} else {
					log.Infof("✧˖°. sage gemini email sent url=%s ✧˖°.", v.Loc)
				}

				time.Sleep(time.Second * 5)
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

			contents, err := internal.FetchHTML(cmd.Context(), "https://text.npr.org")
			if err != nil {
				log.Fatal(err)
			}

			subject, err := internal.ParseTitle(contents)
			if err != nil {
				log.Error(err)
				subject = "✧˖°. sage gemini test email ✧˖°."
			}

			tmp, err := ioutil.TempFile("./", "attachment.*.html")
			if err != nil {
				log.Fatal(err)
			}

			_, err = tmp.Write([]byte(contents))
			if err != nil {
				log.Fatal(err)
			}

			defer tmp.Close()
			defer os.RemoveAll(tmp.Name())

			if err := internal.SendWithAttachment(to, subject, contents, tmp.Name()); err != nil {
				log.Error(err)
			} else {
				log.Info("✧˖°. sage gemini test email sent ✧˖°.")
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

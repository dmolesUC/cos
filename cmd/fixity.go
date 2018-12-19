package cmd

import (
	"errors"
	"fmt"
	. "net/url"
	"strings"

	"github.com/spf13/cobra"
)

// ------------------------------------------------------------
// Constants: Help Text

const shortDescription = "Verify the digest of an object"

const longDescription = shortDescription + `
[TODO: long description]
`

// TODO: replace with fake regions, buckets, prefixes, objects
const example = `
  fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg
  fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg -x c99ad299fa53d5d9688909164cf25b386b33bea8d4247310d80f615be29978f5
  fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg -a md5 -x eac8a75e3b3023e98003f1c24137ebbd
  fixity s3://www.dmoles.net/images/fa/archive.svg -e https://s3.us-west-2.amazonaws.com/ -a md5 -x eac8a75e3b3023e98003f1c24137ebbd
`

// ------------------------------------------------------------
// Global variables

var objectUrl *URL
var verbose *bool
var expected *[]byte
var algorithm *string
var endpoint *string

// ------------------------------------------------------------
// Output

func maybePrintArgs(args []string) {
	if !*verbose {
		return
	}
	fmt.Printf("object URL: %v\n", args[0])
	fmt.Printf("verbose   : %v\n", *verbose)
	fmt.Printf("algorithm : %v\n", *algorithm)
	fmt.Printf("expected  : %x\n", *expected)
	fmt.Printf("endpoint  : %v\n", *endpoint)
}

// ------------------------------------------------------------
// Validators

// TODO: 1) decompose HTTP(S) URLs into endpoint URL + host
// TODO: 2) fail on S3 URLs w/o specified endpoint
// TODO:    2a) or look up default endpoint in S3 config / environment variables?
func initObjectUrl(objUrlStr string) error {
	objUrl, err := Parse(objUrlStr)
	if err != nil {
		return err
	}
	if !objUrl.IsAbs() {
		msg := fmt.Sprintf("object URL '%v' must have a scheme", objUrlStr)
		return errors.New(msg)
	}
	objectUrl = objUrl
	return nil
}

// ------------------------------------------------------------
// Command

func cmdMain(cmd *cobra.Command, args []string) error {
	if *verbose {
		maybePrintArgs(args)
	}

	err := initObjectUrl(args[0])
	if err != nil {
		return err
	}

	fmt.Println(objectUrl.Scheme)

	return nil
}

// ------------------------------------------------------------
// Command initialization

func init() {
	fixityCmd := &cobra.Command{
		Use:           "fixity <OBJECT-URL>",
		Short:         shortDescription,
		Long:          strings.TrimSpace(longDescription),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "  " + strings.TrimSpace(example),
		RunE:          cmdMain,
	}

	flags := fixityCmd.Flags()
	expected = flags.BytesHexP("expected", "x", nil, "Expected digest value")
	algorithm = flags.StringP("algorithm", "a", "sha256", "Algorithm (md5, sha256; default is sha256)")
	endpoint = flags.StringP("endpoint", "e", "", "S3 endpoint")
	verbose = flags.BoolP("verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(fixityCmd)
}

//func checkFixity(objUrlStr string) error {
//	objUrl, err := Parse(objUrlStr)
//	if err == nil {
//		return checkFixityUrl(objUrl)
//	}
//	return err
//}
//
//func checkFixityUrl(objUrl *URL) error {
//
//	//s3Config := aws.Config{
//	//	Endpoint: aws.String(objUrl.String()),
//	//}
//	//
//	//s3Opts := session.Options{
//	//	Config:            s3Config,
//	//	SharedConfigState: session.SharedConfigEnable,
//	//}
//	//
//	//sess, err := session.NewSessionWithOptions(s3Opts)
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//svc := s3.New(sess)
//	//result, err := svc.ListBuckets(nil)
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//for _, b := range result.Buckets {
//	//	fmt.Printf("* %s created on %s\n", aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
//	//}
//
//	return nil
//}

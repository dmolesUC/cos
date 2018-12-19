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

// TODO: format these with something like https://play.golang.org/p/EGaSR0TmMy-
const (
	shortDescription = "Verify the digest of an object"
	longDescription  = shortDescription + `
[TODO: long description]
`
	example = `
  coscheck fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg
  coscheck fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg -x c99ad299fa53d5d9688909164cf25b386b33bea8d4247310d80f615be29978f5
  coscheck fixity https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg -a md5 -x eac8a75e3b3023e98003f1c24137ebbd
  coscheck fixity s3://www.dmoles.net/images/fa/archive.svg -e https://s3.us-west-2.amazonaws.com/ -a md5 -x eac8a75e3b3023e98003f1c24137ebbd
`
)

// ------------------------------------------------------------
// Fixity type

type Fixity struct {
	Verbose   bool
	Expected  []byte
	Algorithm string
	Endpoint  string
	ObjectUrlStr string
	ObjectUrl URL
}

// ------------------------------
// Methods

// TODO:
//   1) decompose HTTP(S) URLs into endpoint URL + host
//   2) fail on S3 URLs w/o specified endpoint
//      2a) or look up default endpoint in S3 config / environment variables?
func (f *Fixity) runWith(_ *cobra.Command, args []string) error {
	f.ObjectUrlStr = args[0]
	f.maybePrintArgs()
	return f.initObjectUrl()
}

func (f *Fixity) initObjectUrl() error {
	objUrl, err := validUrl(f.ObjectUrlStr)
	if err != nil {
		return err
	}
	f.ObjectUrl = *objUrl
	return nil
}

func (f *Fixity) maybePrintArgs() {
	if f.Verbose {
		fmt.Printf("object URL: %v\n", f.ObjectUrlStr)
		fmt.Printf("verbose   : %v\n", f.Verbose)
		fmt.Printf("algorithm : %v\n", f.Algorithm)
		fmt.Printf("expected  : %x\n", f.Expected)
		fmt.Printf("endpoint  : %v\n", f.Endpoint)
	}
}

// ------------------------------------------------------------
// Global functions

func validUrl(objUrlStr string) (*URL, error) {
	objUrl, err := Parse(objUrlStr)
	if err != nil {
		return objUrl, err
	}
	if !objUrl.IsAbs() {
		msg := fmt.Sprintf("object URL '%v' must have a scheme", objUrlStr)
		return nil, errors.New(msg)
	}
	return objUrl, nil
}

// ------------------------------------------------------------
// Command initialization

func init() {
	fixity := Fixity{}

	cmd := &cobra.Command{
		Use:           "fixity <OBJECT-URL>",
		Short:         shortDescription,
		Long:          strings.TrimSpace(longDescription),
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       "  " + strings.TrimSpace(example),
		RunE:          fixity.runWith,
	}
	cmd.Flags().BoolVarP(&fixity.Verbose, "verbose", "v", false, "Verbose output")
	cmd.Flags().BytesHexVarP(&fixity.Expected, "expected", "x", nil, "Expected digest value")
	cmd.Flags().StringVarP(&fixity.Algorithm, "algorithm", "a", "sha256", "Algorithm (md5, sha256; default is sha256)")
	cmd.Flags().StringVarP(&fixity.Endpoint, "endpoint", "e", "", "S3 endpoint")

	rootCmd.AddCommand(cmd)
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

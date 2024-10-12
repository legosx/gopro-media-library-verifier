package verifyrun

import (
	"fmt"
	"github.com/legosx/gopro-media-library-verifier/buildclient"
	"github.com/legosx/gopro-media-library-verifier/client"
	"github.com/legosx/gopro-media-library-verifier/dirscan"
	"github.com/legosx/gopro-media-library-verifier/fetch"
	"github.com/legosx/gopro-media-library-verifier/verify"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strings"
)

const (
	configAuthTokenKey = "auth.token"
)

type TokenPromptMethod string

const (
	TokenPromptMethodInput TokenPromptMethod = "input"
	TokenPromptMethodCURL  TokenPromptMethod = "curl"
)

var tokenPromptMethodsAvailable = []TokenPromptMethod{TokenPromptMethodInput, TokenPromptMethodCURL}

type Runner struct {
	path              string
	outputFilePath    string
	tokenPromptMethod TokenPromptMethod
	buildClient       func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error)
	buildVerifier     func(fetcher fetch.Fetcher, scanner dirscan.Scanner) Verifier
}

type Verifier interface {
	IdentifyMissingFiles(path string) (filePaths []string, err error)
}

func NewRunner(path, outputFilePath string, tokenPromptMethod TokenPromptMethod, opts ...func(*Runner)) (runner Runner) {
	r := Runner{
		path:              path,
		outputFilePath:    outputFilePath,
		tokenPromptMethod: tokenPromptMethod,
		buildClient: func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error) {
			return buildclient.NewBuilder(opts...).Build()
		},
		buildVerifier: func(fetcher fetch.Fetcher, scanner dirscan.Scanner) Verifier {
			return verify.NewVerifier(fetcher, scanner)
		},
	}

	for _, opt := range opts {
		opt(&r)
	}

	return r
}

func WithBuildClient(buildClient func(opts ...func(builder *buildclient.Builder)) (c *client.Client, err error)) func(r *Runner) {
	return func(r *Runner) {
		r.buildClient = buildClient
	}
}

func WithBuildVerifier(buildVerifier func(fetcher fetch.Fetcher, scanner dirscan.Scanner) Verifier) func(r *Runner) {
	return func(r *Runner) {
		r.buildVerifier = buildVerifier
	}
}

func (r Runner) Run() (err error) {
	verifier, err := r.createVerifier()
	if err != nil {
		return err
	}

	filePaths, err := verifier.IdentifyMissingFiles(r.path)
	if err != nil {
		return err
	}

	if err = r.outputFilePaths(filePaths, r.outputFilePath); err != nil {
		return err
	}

	return nil
}

func (r Runner) outputFilePaths(filePaths []string, outputFilePath string) (err error) {
	if len(filePaths) == 0 {
		fmt.Println("\nAll files from specified local directory are already uploaded to Gopro Media Library.")
		return nil
	}

	sort.Strings(filePaths)

	filePathsInline := ""
	for _, filePath := range filePaths {
		filePathsInline = filePathsInline + fmt.Sprintln(filePath)
	}

	if outputFilePath != "" {
		if err = os.WriteFile(outputFilePath, []byte(filePathsInline), 0644); err != nil {
			return err
		}
		fmt.Printf("\nOutput written to %s\n\n", outputFilePath)
	} else {
		fmt.Printf("\nFiles that still can be uploaded to Gopro Media Library:\n%s\n\n", filePathsInline)
	}

	return nil
}

func (r Runner) getTokenPromptMethods() (methods []buildclient.TokenPromptMethod, err error) {
	if r.tokenPromptMethod == "" {
		return []buildclient.TokenPromptMethod{
			buildclient.NewTokenPromptMethodInput(),
			buildclient.NewTokenPromptMethodCURL(),
		}, nil
	}

	var tokenPromptMethod buildclient.TokenPromptMethod

	switch r.tokenPromptMethod {
	case TokenPromptMethodInput:
		tokenPromptMethod = buildclient.NewTokenPromptMethodInput()
	case TokenPromptMethodCURL:
		tokenPromptMethod = buildclient.NewTokenPromptMethodCURL()
	default:
		return []buildclient.TokenPromptMethod{}, fmt.Errorf("invalid token prompt method: %s", r.tokenPromptMethod)
	}

	return []buildclient.TokenPromptMethod{tokenPromptMethod}, nil
}

func (r Runner) createVerifier() (verifier Verifier, err error) {
	tokenPromptMethods, err := r.getTokenPromptMethods()
	if err != nil {
		return verify.Verifier{}, err
	}

	c, err := r.buildClient(
		buildclient.WithConfigAuthTokenKey(configAuthTokenKey),
		buildclient.WithTokenPromptMethods(tokenPromptMethods...),
		buildclient.WithPersistConfig(buildclient.PersistConfigIfChanged),
		buildclient.WithVerbose(buildclient.VerboseAll),
	)
	if err != nil {
		return verify.Verifier{}, err
	}

	scanner := dirscan.NewScanner(c.GetAllowedExtensions())

	fetcher := fetch.NewFetcher(*c)

	return r.buildVerifier(fetcher, scanner), nil
}

func Init(cmd *cobra.Command) error {
	cmd.Flags().StringP("path", "p", "", "path to the local directory to verify")

	if err := cmd.MarkFlagRequired("path"); err != nil {
		return err
	}

	cmd.Flags().StringP("outputFilePath", "o", "", "a path to a file where the result will be written to instead of stdout")

	var methods []string
	for _, tokenPromptMethod := range tokenPromptMethodsAvailable {
		methods = append(methods, string(tokenPromptMethod))
	}

	usage := fmt.Sprintf("method to use for token prompt (%s)", strings.Join(methods, ", "))
	cmd.Flags().StringP("tokenPromptMethod", "m", "", usage)

	return nil
}

package buildclient

import (
	"github.com/pkg/errors"
	"regexp"
)

type promptInputFunc func(label string, mask rune, hideEntered bool, validate func(value string) error) (value string, err error)

type TokenPromptMethodCURL struct {
	promptInput promptInputFunc
}

func NewTokenPromptMethodCURL(opts ...func(*TokenPromptMethodCURL)) TokenPromptMethodCURL {
	method := TokenPromptMethodCURL{
		promptInput: PromptInput,
	}

	for _, opt := range opts {
		opt(&method)
	}

	return method
}

func WithCURLPromptInput(promptInput promptInputFunc) func(t *TokenPromptMethodCURL) {
	return func(t *TokenPromptMethodCURL) {
		t.promptInput = promptInput
	}
}

func (t TokenPromptMethodCURL) GetName() (name string) {
	return "CURL request"
}

func (t TokenPromptMethodCURL) GetMessage() (message string) {
	return "Please provide CURL request"
}

func (t TokenPromptMethodCURL) GetToken() (token string, err error) {
	curlRequest, err := t.promptInput("CURL request", '*', true, t.Validate)
	if err != nil {
		return "", errors.Wrap(err, "curl prompt failed")
	}

	return t.getBearerValue(curlRequest)
}

func (t TokenPromptMethodCURL) Validate(value string) (err error) {
	if len(value) == 0 {
		return errors.New("curl request cannot be empty")
	}

	if _, err := t.getBearerValue(value); err != nil {
		return err
	}

	return nil
}

func (t TokenPromptMethodCURL) getBearerValue(input string) (value string, err error) {
	re := regexp.MustCompile(`Bearer ([\w-]+(\.[\w-]+)*(\.[\w-]+)*)`)
	if match := re.FindStringSubmatch(input); len(match) > 1 {
		return match[1], nil
	}

	return "", errors.New("no bearer token found in CURL request")
}

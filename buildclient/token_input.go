package buildclient

import (
	"github.com/pkg/errors"
)

type TokenPromptMethodInput struct {
	promptInput promptInputFunc
}

func NewTokenPromptMethodInput(opts ...func(*TokenPromptMethodInput)) TokenPromptMethodInput {
	method := TokenPromptMethodInput{
		promptInput: PromptInput,
	}

	for _, opt := range opts {
		opt(&method)
	}

	return method
}

func WithInputPromptInput(promptInput promptInputFunc) func(t *TokenPromptMethodInput) {
	return func(t *TokenPromptMethodInput) {
		t.promptInput = promptInput
	}
}

func (t TokenPromptMethodInput) GetName() (name string) {
	return "Direct input"
}

func (t TokenPromptMethodInput) GetMessage() (message string) {
	return "Please provide your Gopro Media Library token to authenticate:"
}

func (t TokenPromptMethodInput) GetToken() (token string, err error) {
	tokenValue, err := t.promptInput("Token", '*', true, t.Validate)
	if err != nil {
		return "", errors.Wrap(err, "input prompt failed")
	}

	return tokenValue, nil
}

func (t TokenPromptMethodInput) Validate(value string) (err error) {
	if len(value) == 0 {
		return errors.New("token cannot be empty")
	}

	return nil
}

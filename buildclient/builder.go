package buildclient

import (
	"fmt"
	"github.com/legosx/gopro-media-library-verifier/client"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type PersistConfig int

const (
	PersistConfigNever PersistConfig = iota
	PersistConfigIfChanged
	PersistConfigAlways
)

type Verbose int

const (
	VerboseNone Verbose = iota
	VerboseOnlyErrors
	VerboseAll
)

type TokenPromptMethod interface {
	GetToken() (token string, err error)
	GetName() string
	GetMessage() string
}

type createClientFunc func(token string, opts ...func(c *client.Client) error) (client *client.Client, err error)

type promptSelectFunc func(label string, items interface{}) (int, error)

type Builder struct {
	configAuthTokenKey *string
	persistConfig      PersistConfig
	verbose            Verbose
	tokenPromptMethods []TokenPromptMethod
	promptSelect       promptSelectFunc
	createClient       createClientFunc
}

func NewBuilder(opts ...func(builder *Builder)) (builder Builder) {
	builder = Builder{
		promptSelect: PromptSelect,
		createClient: func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
			return client.NewClient(token, opts...)
		},
	}

	for _, opt := range opts {
		opt(&builder)
	}

	return builder
}

func WithConfigAuthTokenKey(key string) func(builder *Builder) {
	return func(builder *Builder) {
		builder.configAuthTokenKey = &key
	}
}

func WithPersistConfig(persistConfig PersistConfig) func(builder *Builder) {
	return func(builder *Builder) {
		builder.persistConfig = persistConfig
	}
}

func WithVerbose(verbose Verbose) func(builder *Builder) {
	return func(builder *Builder) {
		builder.verbose = verbose
	}
}

func WithTokenPromptMethods(methods ...TokenPromptMethod) func(builder *Builder) {
	return func(builder *Builder) {
		builder.tokenPromptMethods = methods
	}
}

func WithPromptSelect(promptSelect promptSelectFunc) func(builder *Builder) {
	return func(builder *Builder) {
		builder.promptSelect = promptSelect
	}
}

func WithCreateClient(createClient createClientFunc) func(builder *Builder) {
	return func(builder *Builder) {
		builder.createClient = createClient
	}
}

func (b Builder) Build() (c *client.Client, err error) {
	if b.configAuthTokenKey == nil && len(b.tokenPromptMethods) == 0 {
		return nil, errors.New("no authentication method provided")
	}

	if b.configAuthTokenKey != nil {
		if c, err = b.fromConfigAuthToken(); err != nil {
			return nil, err
		} else if c != nil {
			return c, nil
		}
	}

	if tokenPromptMethod, err := b.askTokenPromptMethod(); err != nil {
		return nil, errors.Wrap(err, "failed to get chosen token prompt method")
	} else if tokenPromptMethod != nil {
		if c, err = b.fromUserPrompt(tokenPromptMethod); err != nil {
			return nil, errors.Wrap(err, "failed to get client from user prompt")
		}
	}

	if c == nil {
		return nil, errors.New("no client created")
	}

	return c, nil
}

func (b Builder) askTokenPromptMethod() (method TokenPromptMethod, err error) {
	if len(b.tokenPromptMethods) == 0 {
		return nil, nil
	}

	if len(b.tokenPromptMethods) == 1 {
		return b.tokenPromptMethods[0], nil
	}

	items := make([]string, len(b.tokenPromptMethods))

	for i, method := range b.tokenPromptMethods {
		items[i] = method.GetName()
	}

	if index, err := b.promptSelect("Select token prompt method", items); err != nil {
		return nil, errors.Wrap(err, "prompt failed")
	} else {
		method = b.tokenPromptMethods[index]
	}

	return method, nil
}

func (b Builder) fromConfigAuthToken() (c *client.Client, err error) {
	authToken := viper.GetString(*b.configAuthTokenKey)
	if authToken == "" {
		b.print(errors.New("No token found in config"))

		return nil, nil
	}

	b.print("Token found in config")

	if c, err = b.createClient(authToken, client.WithAuthCheck()); err == nil {
		b.print("Using stored token")
		b.handleValidToken(authToken)

		return c, nil
	}

	if errors.As(err, &client.ErrorResponse{}) {
		b.printErr(err, "stored token is not valid")

		return nil, nil
	}

	return nil, errors.Wrap(err, "failed to get token from config")
}

func (b Builder) fromUserPrompt(tokenPromptMethod TokenPromptMethod) (c *client.Client, err error) {
	b.print(fmt.Sprintf("\n%s", tokenPromptMethod.GetMessage()))

	for {
		tokenValue, err := tokenPromptMethod.GetToken()
		if err != nil {
			return nil, errors.Wrap(err, "token prompt failed")
		}

		if c, err = b.createClient(tokenValue, client.WithAuthCheck()); err != nil {
			b.printErr(err, "error checking client")
		} else {
			b.handleValidToken(tokenValue)
			break
		}
	}

	return c, nil
}

func (b Builder) print(a any) {
	if b.verbose == VerboseNone {
		return
	}

	if b.verbose == VerboseOnlyErrors {
		if _, ok := a.(error); !ok {
			return
		}
	}

	fmt.Println(a)
}

func (b Builder) printErr(err error, message string) {
	b.print(errors.Wrap(err, message))
}

func (b Builder) handleValidToken(authTokenValue string) {
	b.print("Token is valid")
	b.doPersistConfig(authTokenValue)
}

func (b Builder) doPersistConfig(authTokenValue string) {
	if b.persistConfig == PersistConfigNever {
		return
	}

	if b.persistConfig == PersistConfigIfChanged &&
		viper.GetString(*b.configAuthTokenKey) == authTokenValue {
		return
	}

	viper.Set(*b.configAuthTokenKey, authTokenValue)

	if viper.ConfigFileUsed() != "" {
		if err := viper.WriteConfig(); err != nil {
			b.printErr(err, "failed to write config")
		} else {
			b.print("Token saved to config file")
		}

		return
	}

	if err := viper.SafeWriteConfig(); err != nil {
		b.printErr(err, "failed to safe write config")
		return
	}

	b.print("Token saved to a new config file")
}

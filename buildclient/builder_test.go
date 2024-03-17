package buildclient_test

import (
	"crypto/rand"
	"github.com/legosx/gopro-media-library-verifier/buildclient"
	"github.com/legosx/gopro-media-library-verifier/client"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"os"
	"path/filepath"
	"testing"
)

func TestBuilder_Build(t *testing.T) {
	type fields struct {
		opts func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder)
	}

	type want struct {
		err error
	}

	createClient1try := 0

	tests := []struct {
		name string
		fields
		want
	}{
		{
			name: "sad path, no auth method provided",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					return []func(builder *buildclient.Builder){}
				},
			},
			want: want{
				err: errors.New("no authentication method provided"),
			},
		},
		{
			name: "sad path, token provided but it's empty",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey("emptyTokenKey"),
					}
				},
			},
			want: want{
				err: errors.New("no client created"),
			},
		},
		{
			name: "sad path, token provided but not authorized",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("invalid")

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
					}
				},
			},
			want: want{
				err: errors.New("no client created"),
			},
		},
		{
			name: "sad path, unexpected error from the client",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("invalid")

					createClient := func(token string, opts ...func(c *client.Client) error) (client *client.Client, err error) {
						return nil, assert.AnError
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
					}
				},
			},
			want: want{
				err: errors.Wrap(assert.AnError, "failed to get token from config"),
			},
		},
		{
			name: "happy path, token is valid",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
					}
				},
			},
		},
		{
			name: "happy path, persist config if changed",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")
					assert.NoError(t, setRandomViperConfigFile())

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
						buildclient.WithPersistConfig(buildclient.PersistConfigIfChanged),
					}
				},
			},
		},
		{
			name: "happy path, persist config if changed but it's not changed now",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")
					assert.NoError(t, setRandomViperConfigFile())

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
						buildclient.WithPersistConfig(buildclient.PersistConfigIfChanged),
					}
				},
			},
		},
		{
			name: "happy path, always persist config",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")
					assert.NoError(t, setRandomViperConfigFile())

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
						buildclient.WithPersistConfig(buildclient.PersistConfigAlways),
					}
				},
			},
		},
		{
			name: "happy path, but can't persist config into already existing file",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")
					viper.SetConfigFile("does-not-exist")

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
						buildclient.WithPersistConfig(buildclient.PersistConfigAlways),
					}
				},
			},
		},
		{
			name: "happy path, can persist config into already existing file",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")
					assert.NoError(t, setRandomViperConfigFile())

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
						buildclient.WithPersistConfig(buildclient.PersistConfigAlways),
					}
				},
			},
		},
		{
			name: "happy path, but can't persist config into not existent yet file",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")
					viper.SetConfigFile("")

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
						buildclient.WithPersistConfig(buildclient.PersistConfigAlways),
					}
				},
			},
		},
		{
			name: "happy path, can persist config into not existent yet file",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")
					viper.SetConfigFile("")

					configPath, err := createRandomConfigPath()
					assert.NoError(t, err)
					viper.SetConfigName("test")
					viper.SetConfigType("yaml")
					viper.AddConfigPath(configPath)

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
						buildclient.WithPersistConfig(buildclient.PersistConfigAlways),
					}
				},
			},
		},
		{
			name: "sad path, verbose only errors",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return nil, assert.AnError
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(""),
						buildclient.WithCreateClient(createClient),
						buildclient.WithVerbose(buildclient.VerboseOnlyErrors),
					}
				},
			},
			want: want{
				err: errors.New("no client created"),
			},
		},
		{
			name: "happy path, verbose only errors",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
						buildclient.WithVerbose(buildclient.VerboseOnlyErrors),
					}
				},
			},
		},
		{
			name: "happy path, verbose all",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					tokenKey := setRandomViperKey("valid")

					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithConfigAuthTokenKey(tokenKey),
						buildclient.WithCreateClient(createClient),
						buildclient.WithVerbose(buildclient.VerboseAll),
					}
				},
			},
		},
		{
			name: "happy path, token from user prompt, input method",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					promptSelect := func(label string, items interface{}) (int, error) {
						return 0, nil
					}

					promptInput := func(label string, mask rune, hideEntered bool, validate func(value string) error) (value string, err error) {
						return "userPromptValidToken", nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithCreateClient(createClient),
						buildclient.WithPromptSelect(promptSelect),
						buildclient.WithTokenPromptMethods(
							buildclient.NewTokenPromptMethodInput(
								buildclient.WithInputPromptInput(promptInput),
							),
							buildclient.NewTokenPromptMethodCURL(),
						),
					}
				},
			},
		},
		{
			name: "happy path, token from user prompt, curl method",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					promptSelect := func(label string, items interface{}) (int, error) {
						return 1, nil
					}

					promptInput := func(label string, mask rune, hideEntered bool, validate func(value string) error) (value string, err error) {
						return "Bearer valid", nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithCreateClient(createClient),
						buildclient.WithPromptSelect(promptSelect),
						buildclient.WithTokenPromptMethods(
							buildclient.NewTokenPromptMethodInput(),
							buildclient.NewTokenPromptMethodCURL(
								buildclient.WithCURLPromptInput(promptInput),
							),
						),
					}
				},
			},
		},
		{
			name: "happy path, token from user prompt, input method, no selects",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					promptInput := func(label string, mask rune, hideEntered bool, validate func(value string) error) (value string, err error) {
						return "userPromptValidToken", nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithCreateClient(createClient),
						buildclient.WithTokenPromptMethods(
							buildclient.NewTokenPromptMethodInput(
								buildclient.WithInputPromptInput(promptInput),
							),
						),
					}
				},
			},
		},
		{
			name: "sad path, token from user prompt, select failed",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					promptSelect := func(label string, items interface{}) (int, error) {
						return 0, assert.AnError
					}

					promptInput := func(label string, mask rune, hideEntered bool, validate func(value string) error) (value string, err error) {
						return "userPromptValidToken", nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithCreateClient(createClient),
						buildclient.WithPromptSelect(promptSelect),
						buildclient.WithTokenPromptMethods(
							buildclient.NewTokenPromptMethodInput(
								buildclient.WithInputPromptInput(promptInput),
							),
							buildclient.NewTokenPromptMethodCURL(),
						),
					}
				},
			},
			want: want{
				err: errors.Wrap(assert.AnError, "failed to get chosen token prompt method: "+
					"prompt failed"),
			},
		},
		{
			name: "sad path, token from user prompt failed, input method, no selects",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						return &client.Client{}, nil
					}

					promptInput := func(label string, mask rune, hideEntered bool, validate func(value string) error) (value string, err error) {
						return "", assert.AnError
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithCreateClient(createClient),
						buildclient.WithTokenPromptMethods(
							buildclient.NewTokenPromptMethodInput(
								buildclient.WithInputPromptInput(promptInput),
							),
						),
					}
				},
			},
			want: want{
				err: errors.Wrap(assert.AnError, "failed to get client from user prompt: "+
					"token prompt failed: "+
					"input prompt failed"),
			},
		},
		{
			name: "happy path, token from user prompt, error checking client first time but works second time",
			fields: fields{
				opts: func(mockCtrl *gomock.Controller) []func(builder *buildclient.Builder) {
					createClient := func(token string, opts ...func(c *client.Client) error) (c *client.Client, err error) {
						createClient1try++
						if createClient1try == 2 {
							return &client.Client{}, nil
						}

						return nil, assert.AnError
					}

					promptInput := func(label string, mask rune, hideEntered bool, validate func(value string) error) (value string, err error) {
						return "userPromptValidToken", nil
					}

					return []func(builder *buildclient.Builder){
						buildclient.WithCreateClient(createClient),
						buildclient.WithTokenPromptMethods(
							buildclient.NewTokenPromptMethodInput(
								buildclient.WithInputPromptInput(promptInput),
							),
						),
					}
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			viper.Reset()

			builder := buildclient.NewBuilder(tt.fields.opts(mockCtrl)...)

			got, err := builder.Build()
			if tt.want.err == nil {
				assert.NotNil(t, got)
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err.Error())
			}
		})
	}
}

func randomString() string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 30)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}

func setRandomViperKey(value string) (key string) {
	key = randomString()
	viper.Set(key, value)

	return key
}

func createRandomConfigPath() (configPath string, err error) {
	configPath = filepath.Join(os.TempDir(), randomString())
	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		return "", err
	}

	return configPath, nil
}

func createRandomConfigFile() (configFilePath string, err error) {
	configPath, err := createRandomConfigPath()
	if err != nil {
		return "", err
	}

	configFilePath = filepath.Join(configPath, ".test.yaml")

	file, err := os.Create(configFilePath)
	if err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return configFilePath, nil
}

func setRandomViperConfigFile() (err error) {
	configFilePath, err := createRandomConfigFile()
	if err != nil {
		return err
	}

	viper.SetConfigFile(configFilePath)

	return nil
}

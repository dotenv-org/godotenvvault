package godotenvvault

import (
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/joho/godotenv"
)

// TEMPORARY: JUST PROXY ALL FUNCTIONS TO godotenv PACKAGE.

func Parse(r io.Reader) (map[string]string, error) {
	return godotenv.Parse(r)
}

func Load(filenames ...string) (err error) {
	return godotenv.Load(filenames...)
}

func Overload(filenames ...string) (err error) {
	return godotenv.Overload(filenames...)
}

func Read(filenames ...string) (envMap map[string]string, err error) {
	return godotenv.Read(filenames...)
}

func Unmarshal(str string) (envMap map[string]string, err error) {
	return godotenv.Unmarshal(str)
}

func UnmarshalBytes(src []byte) (map[string]string, error) {
	return godotenv.UnmarshalBytes(src)
}

func Exec(filenames []string, cmd string, cmdArgs []string, overload bool) error {
	return godotenv.Exec(filenames, cmd, cmdArgs, overload)
}

func Write(envMap map[string]string, filename string) error {
	return godotenv.Write(envMap, filename)
}

func Marshal(envMap map[string]string) (string, error) {
	return godotenv.Marshal(envMap)
}

type DotEnvKey struct {
	Key            string
	Params         map[string][]string
	Environment    string
	EnvironmentKey string
}

const (
	errorURLScheme          = "INVALID_DOTENV_KEY: URL scheme must be 'dotenv'"
	errorKeyField           = "INVALID_DOTENV_KEY: Missing key part"
	errorMissingEnvironment = "INVALID_DOTENV_KEY: Missing environment part"
)

// ParseKey parses a URL into a DotEnvKey structure, checking that the
// format of the URL key input is correct.
//
// Valid example:
//
// dotenv://:key_1234@dotenv.org/vault/.env.vault?environment=production
func ParseKey(keyStr string) (*DotEnvKey, error) {
	uri, err := url.Parse(keyStr)
	if err != nil {
		return nil, err
	}

	if uri.Scheme != "dotenv" {
		return nil, errors.New(errorURLScheme)
	}
	if uri.User == nil {
		return nil, errors.New(errorKeyField)
	}
	password, ok := uri.User.Password()
	if !ok {
		return nil, errors.New(errorKeyField)
	}

	params := uri.Query()
	environment := params.Get("environment")
	if environment == "" {
		return nil, errors.New(errorMissingEnvironment)
	}

	return &DotEnvKey{
		Key:            password,
		Params:         params,
		Environment:    environment,
		EnvironmentKey: "DOTENV_VAULT_" + strings.ToUpper(environment),
	}, nil
}

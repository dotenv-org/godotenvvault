package godotenvvault

// Significant portions of this code are derived directly from John
// Barton's godotenv package, in order to provide an identical API to
// godotenv for encrypted .env.vault files.
//
// The gotodenv package is published under the following licence:
//
// ----------------------------------------------------------------------
// Copyright (c) 2013 John Barton
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
// ----------------------------------------------------------------------
//
// (Original package published at https://github.com/joho/godotenv/)

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

// Parse reads an encrypted .env.vault files from an io.Reader,
// returning a map of keys and values.
func Parse(r io.Reader) (map[string]string, error) {
	dotenvKey, ok := os.LookupEnv("DOTENV_KEY")
	if !ok {
		return nil, errors.New(errorMissingDotenvKey)
	}

	// Use original godotenv Parse function to retrieve the encrypted
	// environments from the vault file.
	envVault, err := godotenv.Parse(r)
	if err != nil {
		return nil, err
	}

	// Extract encrypted data for each provided key from the vault file.
	keys := []*keyData{}
	dotenvKeys := strings.Split(dotenvKey, ",")
	for _, dotenvKey := range dotenvKeys {
		dotenvKey = strings.TrimSpace(dotenvKey)

		key, err := parseKey(dotenvKey)
		if err != nil {
			return nil, err
		}

		cipherText, ok := envVault[key.environmentKey]
		if !ok {
			return nil, errors.New(fmt.Sprintf(errorMissingEnvInVault, key.environmentKey))
		}

		keys = append(keys, &keyData{key.key, cipherText})
	}

	// Attempt to decrypt the environments from the vault.
	plainText, err := keyRotation(keys)
	if err != nil {
		return nil, err
	}

	// Parse the resulting environment variable settings using
	// godotenv's Parse function.
	return godotenv.Parse(strings.NewReader(plainText))
}

// Load will read your encrypted env file(s) and load them into the
// environment for this process.
//
// Call this function as close as possible to the start of your
// program (ideally in main).
//
// If you call Load without any args it will default to loading
// .env.vault in the current path.
//
// You can otherwise tell it which files to load (there can be more
// than one) like:
//
//	godotenvvault.Load("fileone", "filetwo")
//
// It's important to note that it WILL NOT OVERRIDE an environment
// variable that already exists - consider the .env.vault file to set
// development variables or sensible defaults.
func Load(filenames ...string) error {
	// Fallback to godotenv if DOTENV_KEY environment variable isn't
	// set.
	if _, exists := os.LookupEnv("DOTENV_KEY"); !exists {
		return godotenv.Load(filenames...)
	}

	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err := loadFile(filename, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// Overload will read your encrypted env file(s) and load them into
// the environment for this process.
//
// Call this function as close as possible to the start of your
// program (ideally in main).
//
// If you call Overload without any args it will default to loading
// .env.vault in the current path.
//
// You can otherwise tell it which files to load (there can be more
// than one) like:
//
//	godotenvvault.Overload("fileone", "filetwo")
//
// It's important to note this WILL OVERRIDE an environment variable
// that already exists - consider the .env.vault file to forcefully
// set all environment variables.
func Overload(filenames ...string) error {
	// Fallback to godotenv if DOTENV_KEY environment variable isn't
	// set.
	if _, exists := os.LookupEnv("DOTENV_KEY"); !exists {
		return godotenv.Overload(filenames...)
	}

	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err := loadFile(filename, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// Read all encrypted environments (with the same file loading
// semantics as Load) but return values as a map rather than
// automatically writing values into the environment.
func Read(filenames ...string) (map[string]string, error) {
	// Fallback to godotenv if DOTENV_KEY environment variable isn't
	// set.
	if _, exists := os.LookupEnv("DOTENV_KEY"); !exists {
		return godotenv.Read(filenames...)
	}

	filenames = filenamesOrDefault(filenames)
	envMap := make(map[string]string)

	for _, filename := range filenames {
		individualEnvMap, err := readFile(filename)

		if err != nil {
			return nil, err
		}

		for key, value := range individualEnvMap {
			envMap[key] = value
		}
	}

	return envMap, nil
}

// Unmarshal reads an environment file from a string, returning a map
// of keys and values.
func Unmarshal(str string) (envMap map[string]string, err error) {
	return godotenv.Unmarshal(str)
}

// UnmarshalBytes parses an environment file from a byte slice of
// chars, returning a map of keys and values.
func UnmarshalBytes(src []byte) (map[string]string, error) {
	return godotenv.UnmarshalBytes(src)
}

// Exec loads environment variabless from the specified filenames
// (empty map falls back to default .env.vault file) then executes the
// specified command.
//
// This simply hooks os.Stdin/err/out up to the command and calls
// Run().
//
// If you want more fine grained control over your command, it's
// recommended that you use `Load()`, `Overload()` or `Read()` and the
// `os/exec` package yourself.
func Exec(filenames []string, cmd string, cmdArgs []string, overload bool) error {
	op := Load
	if overload {
		op = Overload
	}
	if err := op(filenames...); err != nil {
		return err
	}

	command := exec.Command(cmd, cmdArgs...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

// Write serializes the given environment and writes it to a file.
func Write(envMap map[string]string, filename string) error {
	return godotenv.Write(envMap, filename)
}

// Marshal outputs the given environment as a dotenv-formatted
// environment file. Each line is in the format: KEY="VALUE" where
// VALUE is backslash-escaped.
func Marshal(envMap map[string]string) (string, error) {
	return godotenv.Marshal(envMap)
}

// Key data extracted from individual key URLs in the DOTENV_KEY
// environment variable.
type dotEnvKey struct {
	key            string
	environmentKey string
}

// Error messages used during DOTENV_KEY processing.
const (
	errorInvalidKey        = "INVALID_DOTENV_KEY: Key must be valid."
	errorInvalidKeyLength  = "INVALID_DOTENV_KEY: Key part must be 64 characters long (or more)"
	errorKeyField          = "INVALID_DOTENV_KEY: Missing key part"
	errorMissingDotenvKey  = "NOT_FOUND_DOTENV_KEY: Cannot find environment variable 'DOTENV_KEY'"
	errorMissingEnvInKey   = "INVALID_DOTENV_KEY: Missing environment part"
	errorMissingEnvInVault = "NOT_FOUND_DOTENV_ENVIRONMENT: Cannot locate environment %s in your .env.vault file. Run 'npx dotenv-vault build' to include it."
	errorURLScheme         = "INVALID_DOTENV_KEY: URL scheme must be 'dotenv'"
)

// parseKey parses a URL into a dotEnvKey structure, checking that the
// format of the URL key input is correct.
//
// Valid example:
//
// dotenv://:key_1234@dotenv.org/vault/.env.vault?environment=production
func parseKey(keyStr string) (*dotEnvKey, error) {
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
		return nil, errors.New(errorMissingEnvInKey)
	}

	return &dotEnvKey{
		key:            password,
		environmentKey: "DOTENV_VAULT_" + strings.ToUpper(environment),
	}, nil
}

// Temporary data structure for passing around encryption keys and
// ciphertexts from vault files.
type keyData struct {
	encryptedKey string
	cipherText   string
}

// Attempt to decrypt encrypted environment strings one at a time,
// returning the first success.
func keyRotation(keys []*keyData) (string, error) {
	for _, k := range keys {
		plainText, err := decrypt(k.cipherText, k.encryptedKey)
		if err == nil {
			return string(plainText), nil
		}
	}
	return "", errors.New(errorInvalidKey)
}

// Decrypt a single encrypted environment string using the supplied
// key. The cipher is AES-GCM, and the first 12 bytes of the
// ciphertext are used as the nonce value.
func decrypt(cipherText string, key string) ([]byte, error) {
	key = key[4:]
	if len(key) < 64 {
		return nil, errors.New(errorInvalidKeyLength)
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	nonceBytes := cipherTextBytes[:12]
	cipherTextBytes = cipherTextBytes[12:]

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainText, err := aesgcm.Open(nil, nonceBytes, cipherTextBytes, nil)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

// ---------------------------------------------------------------------------
// Private functions taken from github.com/joho/godotenv start here:

func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env.vault"}
	}
	return filenames
}

func loadFile(filename string, overload bool) error {
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] || overload {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return Parse(file)
}

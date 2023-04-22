package godotenvvault

import (
	"io"

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

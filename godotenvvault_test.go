package godotenvvault

import (
	"strings"
	"testing"
)

type testKey struct {
	key            string
	ok             bool
	environmentKey string
}

var TEST_KEYS = []testKey{
	{"dotenv://:key_1234@dotenv.org/vault/.env.vault?environment=production", true, "DOTENV_VAULT_PRODUCTION"},
	{"dotenv://:key_0dec82bea24ada79a983dcc11b431e28838eae59a07a8f983247c7ca9027a925@dotenv.local/vault/.env.vault?environment=development", true, "DOTENV_VAULT_DEVELOPMENT"},

	// Missing key value.
	{"dotenv://dotenv.org/vault/.env.vault?environment=production", false, "DOTENV_VAULT_PRODUCTION"},

	// Missing environment.
	{"dotenv://:key_1234@dotenv.org/vault/.env.vault", false, ""},
}

func TestKeyParsing(t *testing.T) {
	for itest, test := range TEST_KEYS {
		key, err := parseKey(test.key)
		if !test.ok {
			if err == nil {
				t.Errorf("Parse should have failed but didn't! (test key #%d)", itest+1)
			}
			continue
		}

		if err != nil {
			t.Errorf("Parse failed (test key #%d): %v", itest+1, err)
			continue
		}
		if key.environmentKey != test.environmentKey {
			t.Errorf("Parse failed (test key #%d): bad environment key = '%s'", itest+1, key.environmentKey)
		}
	}
}

const PARSE_TEST_KEY = "dotenv://:key_0dec82bea24ada79a983dcc11b431e28838eae59a07a8f983247c7ca9027a925@dotenv.local/vault/.env.vault?environment=development"

const PARSE_TEST_VAULT = `# .env.vault (generated with npx dotenv-vault local build)
DOTENV_VAULT_DEVELOPMENT="H2A2wOUZU+bjKH3kTpeua9iIhtK/q7/VpAn+LLVNnms+CtQ/cwXqiw=="
`

func TestVaultParsing(t *testing.T) {
	vaultReader := strings.NewReader(PARSE_TEST_VAULT)
	t.Setenv("DOTENV_KEY", PARSE_TEST_KEY)
	vals, err := Parse(vaultReader)
	if err != nil {
		t.Errorf("Parse failed: %v", err)
	}
	check, ok := vals["HELLO"]
	if !ok || check != "world" {
		t.Errorf("Parse returned invalid contents: %v", vals)
	}
}

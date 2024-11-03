package security_test

import (
	"testing"

	"github.com/pablor21/goms/app/config"
	"github.com/pablor21/goms/pkg/security"
)

func TestSigner_Sign(t *testing.T) {
	config.InitConfig([]string{"config.yml"})
	// if err != nil {
	// 	t.Errorf("Error loading config: %v", err)
	// }
	config.GetConfig().Security.Encryption.Key = "6NJT65HjtLkbcJ1nDNuxNzsZz89v3riw"
	testString := "test"
	en, err := security.Encrypt(testString)
	if err != nil {
		t.Errorf("Error encrypting string: %v", err)
	}

	de, err := security.Decrypt(en)
	if err != nil {
		t.Errorf("Error decrypting string: %v", err)
	}

	if de != testString {
		t.Errorf("Decrypted string does not match original string")
	}

}

package message

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetSASLCredentials(t *testing.T) {
	t.Run("validEnabledSASLSCRAM", testCredentials(t, "user", "pass", "SCRAM", "user", "pass", SCRAM, true))
	t.Run("validEnabledSASLSPLAIN", testCredentials(t, "user", "pass", "PLAIN", "user", "pass", PLAIN, true))
	t.Run("validDisabledSASLSOAUTH2", testCredentials(t, "user", "pass", "OAUTH2", "user", "pass", OAUTH2, true))

	t.Run("invalidDisabledSASLSOAUTH2", testCredentials(t, "user", "pass", "garbo", "user", "pass", NONE, false))
	t.Run("invalidDisabledSASL", testCredentials(t, "", "", "", "", "", NONE, false))
	t.Run("invalidDisabledSASLNoPass", testCredentials(t, "user", "", "", "user", "", NONE, false))
	t.Run("invalidDisabledSASLNoUser", testCredentials(t, "", "pass", "", "", "pass", NONE, false))
}

func testCredentials(_ *testing.T, setUser, setPass, setAuth, expectedUser, expectedPass string, expectedMechanism Mechanism, expectedEnabled bool) func(t *testing.T) {
	return func(t *testing.T) {
		err := configEnv(setUser, setPass, setAuth)
		if err != nil {
			t.Errorf("error configuring environment. error %s", err)
			return
		}

		saslAuth := GetSASLAuthentication()
		assert.Assert(t, saslAuth.SASLEnabled() == expectedEnabled)
		assert.Assert(t, saslAuth.Username == expectedUser)
		assert.Assert(t, saslAuth.Password == expectedPass)
		assert.Assert(t, saslAuth.Mechanism == expectedMechanism)
	}
}

func configEnv(user, pass, mech string) error {
	err := os.Setenv("SASL_USERNAME", user)
	if err != nil {
		return err
	}
	err = os.Setenv("SASL_PASSWORD", pass)
	if err != nil {
		return err
	}
	err = os.Setenv("SASL_MECHANISM", mech)
	if err != nil {
		return err
	}
	return nil
}

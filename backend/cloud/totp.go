package cloud

import (
	"github.com/pquerna/otp/totp"
)

var userTOTPSecrets = map[uint]string{}

func SetUpTotp(userID uint, userEmail string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Prague",
		AccountName: userEmail,
		Period:      90,
	})
	if err != nil {
		return "", err
	}

	userTOTPSecrets[userID] = key.Secret()

	return key.URL(), nil
}

func ValidateTOTP(userID uint, code string) bool {
	secret := userTOTPSecrets[userID]
	if secret == "" {
		return false
	}

	return totp.Validate(code, secret)
}

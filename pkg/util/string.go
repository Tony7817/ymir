package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"ymir.com/pkg/xerr"
)

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func Mask(s string) (string, error) {
	if IsEmailValid(s) {
		s = strings.Split(s, "@")[0]
	}
	if s == "" {
		return "", errors.Wrapf(xerr.NewErrCode(xerr.ReuqestParamError), "mask empty string is not supported")
	}
	if len(s) < 5 {
		return string(s[0]) + "*****" + s[len(s)-1:], nil
	}

	return string(s[0]) + "*****" + s[len(s)-1:], nil
}

func GenerateCpatcha() (string, error) {
	const codeLength = 6
	const maxDigit = 10 // Digits are 0-9

	code := ""
	for i := 0; i < codeLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(maxDigit))
		if err != nil {
			return "", fmt.Errorf("failed to generate secure code: %w", err)
		}
		code += num.String()
	}
	return code, nil
}

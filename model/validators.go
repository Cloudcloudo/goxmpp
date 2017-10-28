package model

import (
	"strings"
)

func JidValidator(jid string) error {

	switch {
	case jid == "":
		return &ValidationError{ Field: "jid", Message: "Blank", }

	case len([]rune(jid)) > 256:
		return &ValidationError{ Field: "jid", Message: "Jid too long", }

	case strings.ContainsAny(jid," \x22\x26\x27\x2F\x3A\x3C\x3E\x40\x7F\xFFFE\xFFFF"):
		return &ValidationError{ Field: "jid", Message: "Invalid characters in jid", }
	}

	return nil
}

func PasswdValidator(pass string) error {
	switch {
	case pass == "":
		return &ValidationError{Field: "password", Message: "Blank",}

	case len([]rune(pass)) < 6:
		return &ValidationError{Field: "password", Message: "Password too short",}
	}
	return nil
}

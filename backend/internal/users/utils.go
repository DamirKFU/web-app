package users

import (
	"app/internal/core"
	"strings"
)

func RegisterValidators(s *core.Server) {
	aliases := map[string]string{
		"username": "required,min=3,max=32,alphanum",
		"password": "required,min=6,max=64",
	}
	s.RegisterValidators(nil, aliases)
}

func NormalizeEmail(email string) string {
	var canonicalDomains = map[string]string{
		"yandex.ru": "ya.ru",
	}

	var dots = map[string]string{
		"ya.ry":     "-",
		"gmail.com": "",
	}

	email = strings.ToLower(email)
	emailParts := strings.SplitN(email, "@", 2)
	if len(emailParts) != 2 {
		return email
	}

	emailName := emailParts[0]
	domainPart := emailParts[1]

	if plusIndex := strings.Index(emailName, "+"); plusIndex != -1 {
		emailName = emailName[:plusIndex]
	}

	if canon, ok := canonicalDomains[domainPart]; ok {
		domainPart = canon
	}

	if dotReplacement, ok := dots[domainPart]; ok {
		emailName = strings.ReplaceAll(emailName, ".", dotReplacement)
	}

	return emailName + "@" + domainPart
}

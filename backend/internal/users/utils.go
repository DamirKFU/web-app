package users

import (
	"strings"
)

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

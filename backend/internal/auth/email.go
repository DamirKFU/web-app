package auth

import (
	"app/internal/core"
	"app/internal/users"
	"fmt"
	"log"
	"net/http"
)

func SendRegistrationEmail(emailBackend *core.EmailSMTPEngine, user *users.User, frontURL, token string) error {
	confirmationURL := fmt.Sprintf("%s/confirm-registration?token=%s", frontURL, token)
	subject := "Подтверждение регистрации"
	data := struct {
		Username        string
		ConfirmationURL string
		ExpirationHours int
		ServiceName     string
	}{
		Username:        user.Username,
		ConfirmationURL: confirmationURL,
	}

	body, err := core.RenderTextTemplate("./templates/registration_confirmation.txt", data)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	go func() {
		if err := emailBackend.SendMail(
			subject,
			[]byte(body),
			[]string{user.Email},
		); err != nil {
			log.Printf("[Email Error] failed to send mail to %s: %v", user.Email, err)
		} else {
			log.Printf("[DEBUG] sent password reset email to %s", user.Email)
		}
	}()

	return nil
}

func SendForgotPasswordEmail(emailBackend *core.EmailSMTPEngine, user *users.User, frontURL, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontURL, token)
	subject := "Сброс пароля"
	data := struct {
		Username string
		ResetURL string
	}{
		Username: user.Username,
		ResetURL: resetURL,
	}

	body, err := core.RenderTextTemplate("./templates/reset_password.txt", data)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	go func() {
		if err := emailBackend.SendMail(
			subject,
			[]byte(body),
			[]string{user.Email},
		); err != nil {
			log.Printf("[Email Error] failed to send mail to %s: %v", user.Email, err)
		} else {
			log.Printf("[DEBUG] sent password reset email to %s", user.Email)
		}
	}()

	return nil
}

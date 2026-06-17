package validator

import (
	"errors"
	"strings"
)

func WebhookURL(url string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("webhook_url must be absolute http or https URL")
	}
	return nil
}

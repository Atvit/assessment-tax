package middleware

import (
	"github.com/Atvit/assessment-tax/config"
)

func Authenticate(u, p string, cfg *config.Configuration) (bool, error) {
	return u == cfg.AdminUsername && p == cfg.AdminPassword, nil
}

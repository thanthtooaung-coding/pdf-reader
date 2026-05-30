package database

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/thanthtooaung-coding/pdf-reader/backend/internal/models"
	"github.com/thanthtooaung-coding/pdf-reader/backend/pkg/utils"
)

func SeedDefaultAdmin(db *gorm.DB, log *logrus.Logger, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return errors.New("default admin email and password are required")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		adminRole, err := ensureRole(tx, "admin", "ADMIN", "Platform administrator")
		if err != nil {
			return err
		}
		if _, err := ensureRole(tx, "user", "USER", "Standard user"); err != nil {
			return err
		}

		var u models.User
		err = tx.Where("email = ?", email).First(&u).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			hashed, hErr := utils.HashPassword(password)
			if hErr != nil {
				return hErr
			}
			u = models.User{
				RoleID:   adminRole.ID,
				Username: "admin",
				Fullname: "PDF Reader Administrator",
				Email:    email,
				Password: hashed,
				Base:     models.Base{IsEnable: true},
			}
			if err := tx.Create(&u).Error; err != nil {
				return err
			}
			log.Infof("seeded default admin user: %s", email)
		case err != nil:
			return err
		default:
			if u.RoleID != adminRole.ID {
				u.RoleID = adminRole.ID
				if err := tx.Save(&u).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func ensureRole(tx *gorm.DB, name, code, _ string) (*models.Role, error) {
	var role models.Role
	err := tx.Where("code = ?", code).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		role = models.Role{
			Name: name,
			Code: code,
			Base: models.Base{IsEnable: true},
		}
		if err := tx.Create(&role).Error; err != nil {
			return nil, err
		}
		return &role, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

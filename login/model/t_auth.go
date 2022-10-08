package model

import (
	"context"
	"sync"

	"github.com/tirozhang/go-practise-demo/login/config"
	"gorm.io/gorm"
)

type AuthInstance struct {
	authDB *gorm.DB
}

var authInstance *AuthInstance
var authOnce sync.Once

func GetAuthInstance() *AuthInstance {
	authOnce.Do(func() {
		authInstance = &AuthInstance{
			authDB: config.DbConn,
		}
	})
	return authInstance
}

// ResolveUserID 根据openID获取UserID
func (t *AuthInstance) ResolveUserID(ctx context.Context, openID string) (int32, error) {
	var user TAuth
	result := t.authDB.WithContext(ctx).Model(&TAuth{}).Where("open_id = ?", openID).First(&user)
	if result.RowsAffected == 1 {
		return user.ID, nil
	}
	if result.Error != nil {
		return 0, result.Error
	}
	user.OpenID = openID
	result = t.authDB.WithContext(ctx).Model(&TAuth{}).Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

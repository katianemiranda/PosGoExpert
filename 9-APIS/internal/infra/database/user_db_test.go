package database

import (
	"katianemiranda/PosGoExpert/9-APIS/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})

	user, _ := entity.NewUser("John Doe", "j@j.com", "123456")
	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	var userfound entity.User
	err = db.First(&userfound, "id = ?", user.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, user.ID, userfound.ID)
	assert.Equal(t, user.Name, userfound.Name)
	assert.Equal(t, user.Email, userfound.Email)
	assert.NotNil(t, userfound.Password)
}

func TestFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})

	user, _ := entity.NewUser("John", "j@j.com", "123456")
	userDB := NewUser(db)

	err = userDB.Create(user)
	assert.Nil(t, err)

	userfound, err := userDB.FindByEmail(user.Email)
	assert.Nil(t, err)
	//assert.Equal(t, user.ID, userfound.ID)
	assert.Equal(t, user.Name, userfound.Name)
	assert.Equal(t, user.Email, userfound.Email)
	assert.NotNil(t, userfound.Password)
}

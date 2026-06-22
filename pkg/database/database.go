package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Username string `gorm:"type:varchar(100);not null"`
	Role     string `gorm:"type:varchar(20);default:'user'"`
}

var db *gorm.DB

// Init 初始化数据库
func Init(dbPath string) {
	var err error
	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 创建默认管理员用户
	createDefaultAdmin()
}

// createDefaultAdmin 创建默认管理员用户
func createDefaultAdmin() {
	var count int64
	db.Model(&User{}).Where("email = ?", "admin@gobox.com").Count(&count)

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("生成密码哈希失败: %v", err)
		}

		admin := User{
			Email:    "admin@gobox.com",
			Password: string(hashedPassword),
			Username: "管理员",
			Role:     "admin",
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("创建管理员用户失败: %v", err)
		}
		log.Println("已创建默认管理员用户: admin@gobox.com / admin123")
	}
}

// GetUserByEmail 通过邮箱查找用户
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 通过ID查找用户
func GetUserByID(id uint) (*User, error) {
	var user User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	return db.Create(user).Error
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}
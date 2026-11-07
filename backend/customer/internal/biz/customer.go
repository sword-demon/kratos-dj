package biz

import (
	"database/sql"

	"gorm.io/gorm"
)

// CustomerWork 业务逻辑部分
type CustomerWork struct {
	Telephone string `gorm:"type:varchar(15);uniqueIndex;" json:"telephone"`
	Name      string `gorm:"type:varchar(255);uniqueIndex;" json:"name"`
	Email     string `gorm:"type:varchar(255);uniqueIndex;" json:"email"`
	Wechat    string `gorm:"type:varchar(255);uniqueIndex;" json:"wechat"`
}

type CustomerToken struct {
	Token          string       `gorm:"type:varchar(4095);" json:"token"`
	TokenCreatedAt sql.NullTime `gorm:"" json:"token_created_at"`
}

// Customer 客户核心字段
type Customer struct {
	// 嵌入 4 个字段: ID, CreatedAt, UpdatedAt, DeletedAt
	// 这 4 个字段是 gorm 自动添加的, 用于记录数据的创建时间, 更新时间, 删除时间
	// 这 4 个字段是 gorm 自动添加的, 用于记录数据的创建时间, 更新时间, 删除时间
	gorm.Model
	// 业务逻辑部分
	CustomerWork
	// token 部分
	CustomerToken
}

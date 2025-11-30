package db

import (
	"complex-sql/internal/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Gormdb() (*gorm.DB, error) {

	dsn := "root:rootpassword@tcp(127.0.0.1:3307)/smiles?parseTime=true"

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("error connecting to MySQL with GORM: %v", err)
	}

	fmt.Println("connected to gorm mysql ")

	err=gormDB.AutoMigrate(&models.UserGorm{}, &models.OrderGorm{})

	if err != nil {
		return nil, fmt.Errorf("error during GORM automigration: %v", err)
	}

	return gormDB, nil

}

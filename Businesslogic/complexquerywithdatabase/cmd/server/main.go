package main

import (
	"complex-sql/internal/db"
	"complex-sql/internal/models"
	"complex-sql/internal/repository"
	"context"
	"fmt"
)

/*

 docker compose run only mysql servcies

 1. open terminal
 2. mysql -uroot -p.
    rootpassword
3. check databases
4. mismatch sql query harddisk and ram solved mismathc in column names


*/

func main() {

	err := db.InitMySQL()

	/*.  Very important for read and write database connections pooling

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	*/

	if err != nil {
		fmt.Printf("Error initializing MySQL: %v\n", err)
		return
	}

	/*
		err = db.Createtables()
		if err != nil {
			fmt.Printf("Error creating tables: %v\n", err)
			return
		}
	*/

	_, err = db.Gormdb()
	if err != nil {
		fmt.Printf("Error creating tables: %v\n", err)
		return
	}

	/*

		       Error handling in all db queries is very important
			   each error handling make the code robust and reliable

	*/

	/*    mysql table users

	repo := repository.NewUserRepo(db.GetDB())

	ctx := context.Background()

	repo.CreateUser(ctx, &models.User{
		Name:  "karix",
		Email: "dheerajsir@karix.com",
	})

	user, err := repo.GetUserByID(ctx, 1)
	if err != nil {
		fmt.Printf("Error retrieving user: %v\n", err)
		return
	}

	*/

	// gorm table creation user_gorms

	repo := repository.NewUserRepoGorm(db.GetDB())

	ctx := context.Background()

	repo.CreateUserGorm(ctx, &models.UserGorm{
		Name:  "karix",
		Email: "dheerajsir@karix.com",
	})

	user, err := repo.GetUserByIDGorm(ctx, 1)
	if err != nil {
		fmt.Printf("Error retrieving user: %v\n", err)
		return
	}

	fmt.Printf("Retrieved User: %+v\n", user)

}

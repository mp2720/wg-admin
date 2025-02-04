package main

import (
	"context"
	"mp2720/wg-admin/wg-admin/app/api"
	"mp2720/wg-admin/wg-admin/app/services"
	"mp2720/wg-admin/wg-admin/config"
	dbPkg "mp2720/wg-admin/wg-admin/db"
	"mp2720/wg-admin/wg-admin/transaction"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := dbPkg.NewDB(ctx, cfg.SQLiteDBPath)
	if err != nil {
		panic(err)
	}
	tm := transaction.NewManager(db)

	userRepo := dbPkg.NewUserRepo(db)

	userService := services.NewUserService(userRepo, nil, tm, nil)

	if err := api.RunHTTPServer(
		":8080",
		"http://localhost",
		userService,
	); err != nil {
		panic(err)
	}
}

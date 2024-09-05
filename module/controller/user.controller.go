package controller

import "github.com/labovector/vecsys-api/module/repository"

type userController struct {
	userRepo *repository.UserRepository
}

func NewUserController(userRepo *repository.UserRepository) userController {
	return userController{
		userRepo: userRepo,
	}
}

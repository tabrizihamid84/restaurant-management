package controllers

import "github.com/gin-gonic/gin"

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
func SignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func HashPassword(password string) string {
	return ""
}

func VerifyPassword(userPassword, providePassword string) (bool, string) {
	return false, ""
}

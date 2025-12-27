# Setup Golang

- mkdir todo_app
- cd todo_app
- go mod init todo_app

## Install Dependencies

Install semua dependencies yang diperlukan :

- go get -u gorm.io/gorm
- go get -u gorm.io/driver/mysql
- go get -u github.com/gin-gionic/gin
- go get -u github.com/golang-jwt/jwt/v5
- go get -u golang.org/x/crypto/bcrypt

Untuk menjalankan aplikasi :
- go run main.go


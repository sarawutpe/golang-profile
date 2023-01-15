# Installation
1 mysql
- mysql name ivypay_profile_db
2 build go
3 setup webp & webp-watchers
- guild https://www.digitalocean.com/community/tutorials/how-to-create-and-serve-webp-images-to-speed-up-your-website
- install webp
- install web-watchers
4. run

# go packages
go get github.com/gin-contrib/cors
go get github.com/gin-gonic/gin
go get github.com/google/uuid
go get gorm.io/driver/mysql
go get gorm.io/gorm

# deploy
- go build
- go build -ldflags "-s -w"
- using env: export GIN_MODE=release
- using code: gin.SetMode(gin.Re



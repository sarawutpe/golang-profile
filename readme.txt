go get "github.com/gin-gonic/gin"
go get "github.com/itsjamie/gin-cors"
go get "gorm.io/driver/sqlite"
go get "github.com/mattn/go-sqlite3"

# deploy
- go build
- go build -ldflags "-s -w"
- using env: export GIN_MODE=release
- using code: gin.SetMode(gin.Re



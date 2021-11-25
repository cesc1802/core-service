module example

go 1.17

require (
	github.com/cesc1802/core-service v0.0.0
	github.com/gin-gonic/gin v1.7.4
	github.com/spf13/cobra v1.2.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gorm.io/gorm v1.22.3
)

replace github.com/cesc1802/core-service v0.0.0 => ../../core-service

require (
	github.com/mattn/go-isatty v0.0.14 // indirect
	go.mongodb.org/mongo-driver v1.7.4 // indirect
	golang.org/x/sys v0.0.0-20211116061358-0a5406a5449c // indirect
)

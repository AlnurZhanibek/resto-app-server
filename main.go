package main

import (
	"resto-app-server/internal/handler"
	"resto-app-server/internal/repo"
)

// @title			Resto App Server
// @version		1.0
// @description	Resto App Server APIs
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:8080
// @BasePath		/
// @schemes		http https
func main() {
	r := repo.New()
	h := handler.New(r)
	server := h.Init()
	err := server.Run()
	if err != nil {
		panic(err)
	}
}

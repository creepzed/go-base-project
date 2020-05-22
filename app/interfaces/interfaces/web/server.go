package web

import (
	"bitbucket.org/walmartdigital/hermes/app/application/send_to_email"
	"bitbucket.org/walmartdigital/hermes/app/application/update_message_status"
	"bitbucket.org/walmartdigital/hermes/app/interfaces/web/middleware/json_validator"
	"bitbucket.org/walmartdigital/hermes/app/interfaces/web/routes"
	"bitbucket.org/walmartdigital/hermes/app/interfaces/web/routes/webhooks"
	"bitbucket.org/walmartdigital/hermes/app/shared/utils/log"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

var echoServer *echo.Echo

func NewWebServer() {
	echoServer = echo.New()
	echoServer.HideBanner = true
	//echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.CORS())
	echoServer.Use(middleware.RequestID())
	echoServer.Validator = json_validator.NewJsonValidator()

}

func InitRoutes(sendToEmailUseCase send_to_email.SendToEmailUseCase, updateMessageStatusUseCase update_message_status.UpdateMessageStatusUseCase) {
	routes.NewHealthHandler(echoServer)
	routes.NewMetricsHandler(echoServer)
	routes.NewPingHandler(echoServer)
	routes.NewLoginHandler(echoServer)
	routes.NewSendHandler(echoServer, sendToEmailUseCase)

	webhooks.NewSendHandler(echoServer, updateMessageStatusUseCase)
}

func Start(port string) {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  3 * time.Minute,
		WriteTimeout: 3 * time.Minute,
	}
	log.Info("Hermes ready to handle messages and listen in port %s", port)
	echoServer.Logger.Fatal(echoServer.StartServer(server))
}

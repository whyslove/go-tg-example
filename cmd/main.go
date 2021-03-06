package main

import (
	"os"
	"os/signal"
	"syscall"
	"test-tg/pkg/buisnes"
	"test-tg/pkg/transport"

	"test-tg/pkg/config"

	// "test-tg/pkg/user"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	// "github.com/sah4ez/tg-example/pkg/config"
	// "github.com/sah4ez/tg-example/pkg/errors"
	// "github.com/sah4ez/tg-example/pkg/storage"
	// "github.com/sah4ez/tg-example/pkg/user"
	// "github.com/sah4ez/tg-example/pkg/config"
	// "github.com/sah4ez/tg-example/pkg/errors"
	// "github.com/sah4ez/tg-example/pkg/storage"
	// "github.com/sah4ez/tg-example/pkg/transport"
	// "github.com/sah4ez/tg-example/pkg/user"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)

	log.Logger = config.Service().Logger()

	log.Log().Msg("hello world")
	defer log.Log().Msg("good bye")

	// var err error

	// svcUser := user.New(userStore)
	svc := buisnes.NewService(log.Logger)
	services := []transport.Option{
		transport.Mathematical(transport.NewMathematical(log.Logger, svc)),
	}
	// services := []transport.Option{
	// 	transport.User(transport.NewUser(log.Logger, svcUser)),
	// }

	srv := transport.New(log.Logger, services...).WithLog(log.Logger)

	srv.Fiber().Get("/api/healthcheck", func(ctx *fiber.Ctx) error {
		ctx.Status(fiber.StatusOK)
		return ctx.JSON(map[string]string{"status": "Ok"})
	})

	go func() {
		log.Info().Str("bind", ":9000").Msg("listen on") // Move to config
		if err := srv.Fiber().Listen(":9000"); err != nil {
			log.Panic().Err(err).Stack().Msg("server error")
		}
	}()

	<-shutdown
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"post_service/controllers"
	"post_service/repositories"
	service "post_service/services"
	"post_service/utils"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	port string
	l    *log.Logger
	g    *gin.Engine
}

func (a *App) NewServer(port string, l *log.Logger) (*App, error) {
	r := gin.Default()

	return &App{
		port: port,
		l:    l,
		g:    r,
	}, nil
}

func (a *App) Run() {

	srv := &http.Server{
		Addr:    a.port,
		Handler: a.g,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.l.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 2)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.l.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		a.l.Fatal("Server Shutdown:", err)
	}

	<-ctx.Done()
	a.l.Println("timeout of 5 seconds.")

	a.l.Println("Server exiting")

}

func main() {
	logger := log.New(os.Stdout, "post-service", log.LstdFlags)
	env := utils.LoadEnv(logger)
	db := env.ConnectToDB()

	app := new(App)
	app, err := app.NewServer(env.ServerPort, logger)
	if err != nil {
		logger.Fatalln(err)
	}

	postR := repositories.NewPostRepository(logger, *db.Collection("posts"))
	postS := service.NewPostService(logger, postR)
	postC := controllers.NewPostController(logger, postS)
	app.g.POST("/post", postC.CreatePost)
	app.g.DELETE("/post/:id", postC.DeletePost)

	app.Run()
}

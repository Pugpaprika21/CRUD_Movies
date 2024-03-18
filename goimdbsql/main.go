package main

import (
	"log"

	"github.com/Pugpaprika21/goimdb/controller"
	"github.com/Pugpaprika21/goimdb/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/proullon/ramsql/driver"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	movieRepository := repository.NewMovieRepository()
	movieController := controller.NewMovieController(movieRepository)

	movieRepository.CreateTable()

	e.GET("/movies", movieController.GetAllMovies)
	e.GET("/movies/:imdbID", movieController.GetMoviesByID)
	e.POST("/movies", movieController.CreateMovies)
	e.PUT("/movies/:imdbID", movieController.UpdateMoviesByID)
	e.DELETE("/movies/:imdbID", movieController.DeleteMoviesByID)

	port := "2565"
	log.Println("starting... port:", port)
	log.Fatal(e.Start(":" + port))
}

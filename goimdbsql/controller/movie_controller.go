package controller

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Pugpaprika21/goimdb/dto"
	"github.com/Pugpaprika21/goimdb/repository"
	"github.com/labstack/echo/v4"
)

type movieController struct {
	repository repository.IMovieRepository
}

func NewMovieController(repository *repository.MovieRepository) *movieController {
	return &movieController{
		repository: repository,
	}
}

func (m *movieController) GetAllMovies(c echo.Context) error {
	y := c.QueryParam("year")
	if y == "" {
		movies, _ := m.repository.GetAll()
		return c.JSON(http.StatusOK, movies)
	}

	year, err := strconv.Atoi(y)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	movies, _ := m.repository.GetByYear(year)
	return c.JSON(http.StatusOK, movies)
}

func (m *movieController) GetMoviesByID(c echo.Context) error {
	imdbID := c.Param("imdbID")
	movie, err := m.repository.GetByID(imdbID)
	switch err {
	case nil:
		return c.JSON(http.StatusOK, movie)
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, echo.Map{"message!": "not found"})
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (m *movieController) CreateMovies(c echo.Context) error {
	movie := &dto.Movie{}
	if err := c.Bind(movie); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	movieRes, err := m.repository.Create(movie)
	switch {
	case err == nil:
		return c.JSON(http.StatusCreated, movieRes)
	case err.Error() == "UNIQUE constraint violation":
		return c.JSON(http.StatusConflict, "movie already exists")
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (m *movieController) UpdateMoviesByID(c echo.Context) error {
	id := c.Param("imdbID")
	movie := &dto.Movie{}
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message!": "required imdbID"})
	}

	if err := c.Bind(movie); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	_, err := m.repository.Update(id, movie)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message!": "update Movie error " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message!": "update Movie success", "data": movie})
}

func (m *movieController) DeleteMoviesByID(c echo.Context) error {
	id := c.Param("imdbID")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message!": "required imdbID"})
	}

	_, err := m.repository.Delete(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message!": "delete Movie error " + err.Error()})
	}

	return c.JSON(http.StatusBadRequest, echo.Map{"message!": "delete Movie success"})
}

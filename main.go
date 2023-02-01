package main

import (
    "net/http"
    "strconv"
    "fmt"
    "log"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "github.com/yutanov/movies_api/docs"
    "github.com/yutanov/movies_api/models"
)

//	@title			Movies Example API
//	@version		1.0
//	@description	This is a sample server movie server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//  @BasePath  /api/v1
func main() {
    err := models.ConnectDatabase()
    checkErr(err)

    r := gin.Default()
    r.Use(cors.Default())

    v1 := r.Group("/api/v1")
    {
      v1.GET("/movies", getAllMovies)
      v1.GET("/movies/:id", getMovieById)
      v1.POST("/movies", addMovie)
      v1.PUT("/movies/:id", updateMovie)
      v1.DELETE("/movies/:id", deleteMovie)
    }
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.Run("localhost:8080")
}


// @Summary Get All Movies
// @Tags movies
// @Description get all movies
// @ID get-all-movies
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Movie
// @Router /movies [get]
func getAllMovies(c *gin.Context) {
	movies, err := models.GetMovies()
	checkErr(err)
	if movies == nil {
		c.JSON(404, gin.H{"error": "No Records Found"})
		return
	} else {
    c.IndentedJSON(http.StatusOK, movies)
	}
}

// @Summary get a movie by ID
// @Tags movies
// @Description get a movie by ID
// @ID get-movie-by-id
// @Produce json
// @Param id path string true "movie ID"
// @Success 200 {object} models.Movie
// @Router /movies/{id} [get]
func getMovieById(c *gin.Context) {
	id := c.Param("id")
	movie, err := models.GetMovieById(id)

	checkErr(err)
	if movie.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
    c.IndentedJSON(http.StatusOK, movie)
		// c.JSON(http.StatusOK, gin.H{"data": person})
	}
}

// @Summary Add New Movie
// @Tags movies
// @Description add new movie
// @ID add-movie
// @Produce json
// @Param data body models.Movie true "movie data"
// @Success 200 {object} models.Movie
// @Router /movies [post]
func addMovie(c *gin.Context) {
	var json models.Movie

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddMovie(json)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

// @Summary Update a movie
// @Description Update a movie
// @Tags movies
// @ID update-movie-by-id
// @Accept json
// @Produce json
// @Param id path int true	"Movie ID"
// @Param account body	models.Movie	true "Update movie"
// @Success 200	{object} models.Movie
// @Router /movies/{id} [put]
func updateMovie(c *gin.Context) {
	var json models.Movie
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movieId, err := strconv.Atoi(c.Param("id"))
	fmt.Printf("Updating id %d", movieId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdateMovie(json, movieId)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// @Summary Delete a movie
// @Description Delete a movie
// @Tags movies
// @ID delete-movie-by-id
// @Produce json
// @Param id path string true "movie ID"
// @Success 200 {object} models.Movie
// @Router /movies/{id} [delete]
func deleteMovie(c *gin.Context) {

	movieId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeleteMovie(movieId)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

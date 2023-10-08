package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/imgSample", "./resources/imgSample")
	router.Static("/js", "./resources/js")

	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"name": "Main website",
			"css":  "/styles/home.css",
		})
	})

	router.GET("/services/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			return
		}

		sample, err := a.repository.GetSampleByID(id)
		if err != nil {
			log.Println("Error with running\nServer down")
			return
		}
		fmt.Println(sample)

		c.HTML(http.StatusOK, "info.tmpl", gin.H{
			"css":    "/styles/info.css",
			"Sample": sample,
		})

	})

	router.GET("/search", func(c *gin.Context) {
		searchQuery := c.DefaultQuery("q", "")

		samples, err := a.repository.GetSampleByName(searchQuery)
		if err != nil {
			log.Println("Error with running\nServer down")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if searchQuery == "" {
			samples, err = a.repository.GetAllSamples()
			if err != nil {
				log.Println("Error with running\nServer down")
				return
			}
		}
		// for i := 0; i < len(samples); i++ {
		// 	samples[i].Date_Sealed = samples[i].Date_Sealed[:10]
		// }

		c.HTML(http.StatusOK, "services.tmpl", gin.H{
			"css":      "/styles/services.css",
			"Services": samples,
			"Search":   searchQuery,
		})

	})

	router.POST("/delete", func(c *gin.Context) {
		id, err := strconv.Atoi(c.DefaultQuery("q", ""))
		log.Print(c.DefaultQuery("q", ""))
		if err != nil {
			log.Print(err)
		}
		err = a.repository.DeleteSampleByID(id)
		if err != nil {
			log.Print(err)
		}
		sample, err := a.repository.GetAllSamples()
		if err != nil {
			log.Print(err)
		}
		data := gin.H{
			"css":      "/styles/employee_mode.css",
			"Services": sample,
		}
		c.HTML(http.StatusOK, "employee_mode.tmpl", data)
	})

	err := router.Run()
	if err != nil {
		log.Println("Error with running\nServer down")
		return
	}

	log.Println("Server down")
}

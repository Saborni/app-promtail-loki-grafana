package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

func helloWorld(c echo.Context) error {
	log.Println("helloWorld")
	return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
}

func health(c echo.Context) error {
	log.Println("health")
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func external(c echo.Context) error {
	start := time.Now()
	log.Println("calling external api")
	resp, err := http.Get("https://api.github.com")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	latency := time.Since(start).Milliseconds()
	log.Printf("external latency: %dms", latency)
	return c.String(http.StatusOK, string(body))
}

func main() {
	logFile, _ := os.OpenFile("/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(logFile)

	e := echo.New()

	e.GET("/", helloWorld)
	e.GET("/health", health)
	e.GET("/external", external)

	e.Logger.Fatal(e.Start(":8000"))
}

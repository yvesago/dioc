package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"sync"
	"time"

	. "./models"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		//	c.Writer.Header().Add("Access-Control-Expose-Headers", "X-myToken")
		c.Next()
	}
}

type Config struct {
	Port       string
	DBname     string
	Salt       string
	Token      string
	IPsAllowed []string
	MailServer string
	MailFrom   string
	MailTo     []string
}

// gin Middlware to set Config
func SetConfig(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Salt", config.Salt)
		c.Set("MailServer", config.MailServer)
		c.Set("MailTo", config.MailTo)
		c.Set("MailFrom", config.MailFrom)
		c.Next()
	}
}

func main() {
	confPtr := flag.String("c", "", "Json config file")
	//	debugPtr := flag.Bool("d", false, "Debug mode")
	flag.Parse()
	conf := *confPtr
	//	Debug := *debugPtr

	file, err := os.Open(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		flag.PrintDefaults()
		os.Exit(0)
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		flag.PrintDefaults()
		os.Exit(0)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go servermain(config)
	wg.Add(1)
	go checkOffline(config.DBname)
	wg.Wait()

}

func checkOffline(dbname string) {
	db := InitDb(dbname)
	for {
		time.Sleep(30 * time.Second)
		CheckAgentOffLine(db)
	}
}

func servermain(config Config) {
	//	gin.SetMode(gin.ReleaseMode)
	//r := gin.Default()
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(Logger())

	r.Use(SetConfig(config))
	r.Use(Database(config.DBname))
	r.Use(Cors())

	admin := r.Group("admin/api/v1")
	admin.Use(TokenAuthMiddleware(config))
	{
		admin.GET("/surveys", GetSurveys)
		admin.GET("/surveys/:id", GetSurvey)
		admin.POST("/surveys", PostSurvey)
		admin.PUT("/surveys/:id", UpdateSurvey)
		admin.DELETE("/surveys/:id", DeleteSurvey)
		admin.OPTIONS("/surveys", Options)     // POST
		admin.OPTIONS("/surveys/:id", Options) // PUT, DELETE

		admin.GET("/agents", GetAgents)
		admin.GET("/agents/:crca", GetAgent)
		admin.POST("/agents", PostAgent)
		admin.PUT("/agents/:crca", UpdateAgent)
		admin.DELETE("/agents/:crca", DeleteAgent)
		admin.OPTIONS("/agents", Options)       // POST
		admin.OPTIONS("/agents/:crca", Options) // PUT, DELETE

		admin.GET("/alertes", GetAlertes)
		admin.GET("/alertes/:id", GetAlerte)
		admin.POST("/alertes", PostAlerte)
		admin.PUT("/alertes/:id", UpdateAlerte)
		admin.DELETE("/alertes/:id", DeleteAlerte)
		admin.OPTIONS("/alertes", Options)     // POST
		admin.OPTIONS("/alertes/:id", Options) // PUT, DELETE
	}

	client := r.Group("client/api/v1")
	{
		client.POST("/agent", RegisterHandler)
		client.PUT("/agent/:crca", SendLinesHandler)
		client.GET("/agent/:crca", CMDHandler)
		client.POST("/alerte", PostNewAlerte)
		client.GET("/survey/:crcs", GetSurveyByCRCs)
	}

	r.Run(config.Port)
}

func Options(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,DELETE,POST,PUT")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TokenAuthMiddleware(config Config) gin.HandlerFunc {
	// some init
	return func(c *gin.Context) {
		//token := c.Request.Header.Get("x-mytoken")
		q := c.Request.URL.Query()
		if q["X-MyToken"] == nil {
			respondWithError(401, "API token required", c)
			return
		}
		token := q["X-MyToken"][0]

		//fmt.Println("token : ", token)
		if token != config.Token {
			respondWithError(401, "Invalid API token", c)
			return
		}

		//fmt.Println("clienIP : " + c.ClientIP())

		if contains(config.IPsAllowed, c.ClientIP()) == false {
			respondWithError(401, "Acces denied", c)
			return
		}

		c.Next()
	}
}

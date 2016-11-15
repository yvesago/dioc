package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"net/http"
	"os"
	"sync"
	"time"

	. "./models"
)

type Config struct {
	Port       string
	DBname     string
	Salt       string
	CorsOrigin string
	Token      string
	IPsAllowed []string
	MailServer string
	MailFrom   string
	MailTo     []string
	Debug      bool
	Verbose    bool
	OffLineMs  int64
	TLScert    string
	TLSkey     string
}

// gin Middlware to set Config
func SetConfig(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Salt", config.Salt)
		c.Set("CorsOrigin", config.CorsOrigin)
		c.Set("MailServer", config.MailServer)
		c.Set("MailTo", config.MailTo)
		c.Set("MailFrom", config.MailFrom)
		c.Set("Verbose", config.Verbose)
		c.Next()
	}
}

func main() {
	// Parameters
	confPtr := flag.String("c", "", "Json config file")
	debugPtr := flag.Bool("d", false, "Debug mode")
	verbosePtr := flag.Bool("v", false, "Verbose mode, need Debug mode")
	flag.Parse()
	conf := *confPtr
	Debug := *debugPtr
	Verbose := *verbosePtr
	if Debug == false { // Verbose need Debug
		Verbose = false
	}

	// Load config frome file
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
	config.Debug = Debug
	config.Verbose = Verbose
	offLineTest := int64(300000)
	if config.OffLineMs != 0 {
		offLineTest = config.OffLineMs
	}

	// Main Loop
	var wg sync.WaitGroup
	wg.Add(1)
	go servermain(config)
	wg.Add(1)
	go checkOffline(config.DBname, offLineTest)
	wg.Wait()

}

func checkOffline(dbname string, offLineMs int64) {
	db := InitDb(dbname)
	for {
		time.Sleep(30 * time.Second)
		CheckAgentOffLine(db, offLineMs)
	}
}

/*func addHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("X-Total-Count", "0")
		//	c.Writer.Header().Add("Access-Control-Expose-Headers", "X-myToken")
		c.Next()
	}
}*/

func servermain(config Config) {
	if config.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	if config.Debug == true {
		r.Use(gin.Logger())
	}

	r.Use(SetConfig(config))
	r.Use(Database(config.DBname))
	r.Use(cors.Middleware(cors.Config{
		Origins:         config.CorsOrigin,
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, X-MyToken",
		ExposedHeaders:  "x-total-count",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	//r.Use(addHeaders())

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

	if config.TLScert != "" && config.TLSkey != "" {
		if config.Debug == true {
			fmt.Println("Listening and serving HTTPS on ", config.Port)
		}
		err := http.ListenAndServeTLS(config.Port, config.TLScert, config.TLSkey, r)
		if err != nil {
			fmt.Println("ListenAndServe: ", err)
			os.Exit(0)
		}
	} else {
		r.Run(config.Port)
	}
}

func Options(c *gin.Context) {
	Origin := c.MustGet("CorsOrigin").(string)

	c.Writer.Header().Set("Access-Control-Allow-Origin", Origin)
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
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
		token := c.Request.Header.Get("X-MyToken")

		if config.Verbose == true {
			fmt.Println("token : ", token)
			fmt.Println("clienIP : " + c.ClientIP())
		}

		if config.Token != "" && token != config.Token {
			respondWithError(401, "Invalid API token", c)
			return
		}

		if config.IPsAllowed != nil && contains(config.IPsAllowed, c.ClientIP()) == false {
			respondWithError(401, "Acces denied", c)
			return
		}

		c.Next()
	}
}

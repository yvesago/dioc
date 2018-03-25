package main

import (
	"encoding/json"
	//	"flag"
	"fmt"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
	"sync"
	"time"

	. "./models"
)

// Config types read from config file
type Config struct {
	Port            string
	DBname          string
	Salt            string
	CorsOrigin      string
	Token           string
	IPsAllowed      []string
	CityDB          string
	AsnDB           string
	MailServer      string
	MailFrom        string
	MailTo          []string
	Debug           bool
	Verbose         bool
	OffLineMs       int64
	TLScert         string
	TLSkey          string
	AuthCASUrl      string   // CAS server
	AuthCASService  string   // Fix CAS service when a proxy need it
	AuthJWTTimeOut  int      // Hours for jwt timeout
	AuthJWTPassword string   // JWT secret password
	AuthJWTCallback string   // client url callback to validate and register jwt
	AuthValidLogins []string // valid cas users
}

// SetConfig gin Middlware to push some config values
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
	var Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage of %s\n\n  Default behaviour: start daemon\n\n", os.Args[0])
		//flag.SortFlags = false
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Usage = Usage

	// Parameters
	confPtr := flag.StringP("conf", "c", "", "Json config file")
	debugPtr := flag.BoolP("debug", "d", false, "Debug mode")
	verbosePtr := flag.BoolP("Verbose", "V", false, "Verbose mode, need Debug mode")
	flag.Parse()

	conf := *confPtr
	Debug := *debugPtr
	Verbose := *verbosePtr
	if Debug == false { // Verbose need Debug
		Verbose = false
	}

	// Load config from file
	file, err := os.Open(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError on mandatory config file:\n %s\n", err)
		Usage()
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		Usage()
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

	e := InitLocDbs(config.CityDB, config.AsnDB)
	if e != nil {
		fmt.Println(e)
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
		RequestHeaders:  "Origin, Authorization, Content-Type, X-MyToken, Bearer",
		ExposedHeaders:  "x-total-count",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	//r.Use(addHeaders())

	casHandler := setCasHandler(config)

	// add /auth/login /auth/logout to allow cas login to set jwt token
	auth := r.Group("auth")
	{
		auth.GET("/login", gin.WrapH(casHandler))
		auth.GET("/logout", gin.WrapH(casHandler))
	}

	admin := r.Group("admin/api/v1")
	//admin.Use(TokenAuthMiddleware(config))
	admin.Use(jwt.Auth(config.AuthJWTPassword))
	{
		admin.GET("/geojson", GetGeoJsonIPs)
		admin.GET("/actionextract", RestExtract)
		admin.GET("/actionfluship", RestFlushIP)

		admin.GET("/board", GetBoard)
		admin.GET("/board/:id", GetBoard)
		admin.PUT("/board/:id", UpdateBoard)

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

		admin.GET("/ips", GetIPs)
		admin.GET("/ips/:id", GetIP)
		admin.POST("/ips", PostIP)
		admin.PUT("/ips/:id", UpdateIP)
		admin.DELETE("/ips/:id", DeleteIP)
		admin.OPTIONS("/ips", Options)     // POST
		admin.OPTIONS("/ips/:id", Options) // PUT, DELETE

		admin.GET("/extracts", GetExtracts)
		admin.GET("/extracts/:id", GetExtract)
		admin.POST("/extracts", PostExtract)
		admin.PUT("/extracts/:id", UpdateExtract)
		admin.DELETE("/extracts/:id", DeleteExtract)
		admin.OPTIONS("/extracts", Options)    // POST
		admin.OPTIONS("/extracs/:id", Options) // PUT, DELETE

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

// Options common response for rest options
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

// TokenAuthMiddleware middleware with various auth options
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

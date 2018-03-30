package models

/*
  Shared functions for models tests
*/

import (
	"fmt"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
)

func deleteFile(file string) {
	// delete file
	var err = os.Remove(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
}

type Config struct {
	DBname     string
	Salt       string
	MailServer string
	MailFrom   string
	MailTo     []string
	Verbose    bool
	CityDB     string
	AsnDB      string
}

func SetConfig(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Salt", config.Salt)
		c.Set("MailServer", config.MailServer)
		c.Set("MailTo", config.MailTo)
		c.Set("MailFrom", config.MailFrom)
		c.Set("Verbose", config.Verbose)
		c.Set("claims", jwt_lib.MapClaims{"id": "test"})
		c.Next()
	}
}

// Set test config
var config = Config{
	DBname: "_test.sqlite3",
	Salt:   "xxxx",
	// CityDB: "./GeoLite2-City.mmdb",
	CityDB:     "./GeoIP2-City-Test.mmdb",
	AsnDB:      "./GeoLite2-ASN-Test.mmdb",
	MailServer: "smtp.my.test:25",
	MailFrom:   "No reply <noreply@my.test>",
	MailTo:     []string{"me@my.org", "other@my.test"},
	Verbose:    true, // Set with cmd line in release mode
}

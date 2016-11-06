package models

/*
  Shared functions for models tests
*/

import (
	"fmt"
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
}

func SetConfig(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Salt", config.Salt)
		c.Set("MailServer", config.MailServer)
		c.Set("MailTo", config.MailTo)
		c.Set("MailFrom", config.MailFrom)
		c.Next()
	}
}

// Set test config
var config = Config{
	DBname:     "_test.sqlite3",
	Salt:       "xxxx",
	MailServer: "smtp.my.test:25",
	MailFrom:   "No reply <noreply@my.test>",
	MailTo:     []string{"me@my.org", "other@my.test"},
}

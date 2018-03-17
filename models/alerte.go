package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gopkg.in/gorp.v2"
	//	"log"
	"strconv"
	"strings"
	"time"
)

/**
Search for XXX to fix fields mapping in Update handler, mandatory fields
or remove sqlite tricks

 vim search and replace cmd to customize struct, handler and instances
  :%s/Alerte/NewStruct/g
  :%s/alerte/newinst/g

**/

// XXX custom struct name and fields
type Alerte struct {
	Id         int64     `db:"id" json:"id"`
	CRCa       string    `db:"crca" json:"crca"`
	CRCs       string    `db:"crcs" json:"crcs"`
	IP         string    `db:"ip" json:"ip"`
	FileSurvey string    `db:"filesurvey" json:"filesurvey"`
	Role       string    `db:"role" json:"role"`
	Line       string    `db:"line" json:"line,size:16384"`
	Search     string    `db:"search" json:"search"`
	Level      string    `db:"level" json:"level"`
	Comment    string    `db:"comment" json:"comment,size:16384"`
	Created    time.Time `db:"created" json:"created"` // or int64
	Updated    time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

func (a *Alerte) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now() // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

func (a *Alerte) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now()
	return nil
}

// REST handlers

func GetAlertes(c *gin.Context) {
	verbose := c.MustGet("Verbose").(bool)
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	query := "SELECT * FROM alerte"

	// Parse query string
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	var count int64
	if s != "" {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM alerte  WHERE " + s)
		query = query + " WHERE " + s
	} else {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM alerte")
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	if verbose == true {
		fmt.Println(q)
		fmt.Println("query: " + query)
	}

	var alertes []Alerte
	_, err := dbmap.Select(&alertes, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.FormatInt(count, 10)) // float64 to string
		c.JSON(200, alertes)
	} else {
		c.JSON(404, gin.H{"error": "no alerte(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/alertes
}

func GetAlerte(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var alerte Alerte
	err := dbmap.SelectOne(&alerte, "SELECT * FROM alerte WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, alerte)
	} else {
		c.JSON(404, gin.H{"error": "alerte not found"})
	}

	// curl -i http://localhost:8080/api/v1/alertes/1
}

func PostAlerte(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)

	var alerte Alerte
	c.Bind(&alerte)

	//log.Println(alerte)

	if alerte.CRCa != "" && alerte.CRCs != "" && alerte.Line != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&alerte)
		if err == nil {
			c.JSON(201, alerte)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/alertes
}

func UpdateAlerte(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var alerte Alerte
	err := dbmap.SelectOne(&alerte, "SELECT * FROM alerte WHERE id=?", id)
	if err == nil {
		var json Alerte
		c.Bind(&json)

		//log.Println(json)
		alerteId, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		alerte := Alerte{
			Id:         alerteId,
			CRCa:       alerte.CRCa,
			CRCs:       alerte.CRCs,
			IP:         alerte.IP,
			FileSurvey: alerte.FileSurvey,
			Role:       alerte.Role,
			Line:       alerte.Line,
			Level:      alerte.Level,
			Search:     alerte.Search,
			Comment:    json.Comment,   // update only comment
			Created:    alerte.Created, //alerte read from previous select
		}

		if alerte.CRCa != "" { // XXX Check mandatory fields
			_, err = dbmap.Update(&alerte)
			if err == nil {
				c.JSON(200, alerte)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "alerte not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/alertes/1
}

func DeleteAlerte(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var alerte Alerte
	err := dbmap.SelectOne(&alerte, "SELECT * FROM alerte WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&alerte)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "alerte not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/alertes/1
}

/**



**/

func PostNewAlerte(c *gin.Context) {
	verbose := c.MustGet("Verbose").(bool)
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)

	var alerte Alerte
	c.Bind(&alerte)

	if verbose == true {
		fmt.Println(alerte)
	}

	if alerte.CRCa != "" && alerte.CRCs != "" && alerte.Line != "" { // XXX Check mandatory fields
		var agent Agent
		_ = dbmap.SelectOne(&agent, "SELECT ip,filesurvey,role,cmd FROM agent WHERE crca=?", alerte.CRCa)
		var survey Survey
		_ = dbmap.SelectOne(&survey, "SELECT * FROM survey WHERE crcs=?", alerte.CRCs)

		//TODO increase and update survey checked
		alerte.IP = agent.IP
		alerte.FileSurvey = agent.FileSurvey
		alerte.Role = agent.Role
		alerte.Search = survey.Search
		alerte.Level = survey.Level

		if alerte.Level != "" && agent.CMD == "" { // don't send mail for FullSearch
			mserver := c.MustGet("MailServer").(string)
			mto := c.MustGet("MailTo").([]string)
			mfrom := c.MustGet("MailFrom").(string)
			//SendMail
			if verbose == true {
				fmt.Printf("Sendmail : %s, %s %s %s\n", agent.CMD, mserver, mto, mfrom)
			}
			alerte.SendMail(mserver, mfrom, mto)
		}

		err := dbmap.Insert(&alerte)
		if err == nil {
			c.JSON(201, "OK")
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}
}

func (a *Alerte) SendMail(mserver string, from string, to []string) {
	s := strings.Split(mserver, ":")
	server, p := s[0], s[1]
	port, _ := strconv.Atoi(p)

	msg := "\nFrom " + a.IP + " " + a.FileSurvey +
		"\nMatch: " + a.Search +
		"\n\n" + a.Line + "\n\n"

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to[0])
	m.SetHeader("CC", strings.Join(to[:], " ,"))
	m.SetHeader("Subject", "[ALERTE] ("+strings.ToUpper(a.Level)+") "+a.Search)
	m.SetBody("text/plain", msg)

	d := gomail.Dialer{Host: server, Port: port}
	if server == "smtp.my.test" {
		return
	}

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Warn: send mail failed: %v\n", err)
	}

}

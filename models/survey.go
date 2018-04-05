package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v2"
	"hash/crc32"
	"log"
	"strconv"
	"time"
)

/**
Search for XXX to fix fields mapping in Update handler, mandatory fields
or remove sqlite tricks

 vim search and replace cmd to customize struct, handler and instances
  :%s/Survey/NewStruct/g
  :%s/survey/newinst/g

**/

// XXX custom struct name and fields
type Survey struct {
	Id      int64     `db:"id" json:"id"`
	CRCs    string    `db:"crcs" json:"crcs"`
	Role    string    `db:"role" json:"role"`
	Search  string    `db:"search" json:"search"`
	Level   string    `db:"level" json:"level"`
	Comment string    `db:"comment" json:"comment,size:65534"`
	Checked int       `db:"checked" json:"checked"`
	Created time.Time `db:"created" json:"created"` // or int64
	Updated time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

func hashSurvey(Role string, Search string) string {
	// Change CRCs on Role or Search update
	crc32q := crc32.MakeTable(0xD5828281)
	crcs := fmt.Sprintf("%08x", crc32.Checksum([]byte(Role+":"+Search), crc32q))
	return crcs
}

func (a *Survey) PreInsert(s gorp.SqlExecutor) error {
	a.Created = time.Now() // or time.Now().UnixNano()
	a.CRCs = hashSurvey(a.Role, a.Search)
	a.Updated = a.Created
	return nil
}

func (a *Survey) PreUpdate(s gorp.SqlExecutor) error {
	a.CRCs = hashSurvey(a.Role, a.Search)
	a.Updated = time.Now()
	return nil
}

// REST handlers

func GetSurveys(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	query := "SELECT * FROM survey"

	// Parse query string
	//  receive : map[_filters:[{"q":"wx"}] _sortField:[id] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	var count int64
	if s != "" {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM survey  WHERE " + s)
		query = query + " WHERE " + s
	} else {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM survey")
	}
	if o != "" {
		query = query + o
	}
	if l != "" {
		query = query + l
	}

	if verbose == true {
		fmt.Println(q)
		fmt.Println(" -- " + query)
	}

	var surveys []Survey
	_, err := dbmap.Select(&surveys, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.FormatInt(count, 10)) // float64 to string
		c.JSON(200, surveys)
	} else {
		c.JSON(404, gin.H{"error": "no survey(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/surveys
}

func GetSurvey(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var survey Survey
	err := dbmap.SelectOne(&survey, "SELECT * FROM survey WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, survey)
	} else {
		c.JSON(404, gin.H{"error": "survey not found"})
	}

	// curl -i http://localhost:8080/api/v1/surveys/1
}

func PostSurvey(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	claims := c.MustGet("claims").(jwt.MapClaims)
	log.Printf("%s [%s]: PostSurvey\n", c.ClientIP(), claims["id"])

	var survey Survey
	c.Bind(&survey)

	//log.Println(survey)

	if survey.Search != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&survey)
		if err == nil {
			c.JSON(201, survey)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/surveys
}

func UpdateSurvey(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")
	claims := c.MustGet("claims").(jwt.MapClaims)
	log.Printf("%s [%s]: UpdateSurvey %s\n", c.ClientIP(), claims["id"], id)

	var survey Survey
	err := dbmap.SelectOne(&survey, "SELECT * FROM survey WHERE id=?", id)
	if err == nil {
		var json Survey
		c.Bind(&json)

		//log.Println(json)
		surveyId, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		survey := Survey{
			Id:      surveyId,
			CRCs:    json.CRCs,
			Role:    json.Role,
			Search:  json.Search,
			Level:   json.Level,
			Comment: json.Comment,
			Checked: json.Checked,
			Created: survey.Created, //survey read from previous select
		}

		if survey.CRCs != "" { // XXX Check mandatory fields
			_, err = dbmap.Update(&survey)
			if err == nil {
				c.JSON(200, survey)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "survey not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/surveys/1
}

func DeleteSurvey(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")
	claims := c.MustGet("claims").(jwt.MapClaims)
	log.Printf("%s [%s]: DeleteSurvey %s\n", c.ClientIP(), claims["id"], id)

	var survey Survey
	err := dbmap.SelectOne(&survey, "SELECT * FROM survey WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&survey)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "survey not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/surveys/1
}

/**

Survey handlers


**/
func GetSurveyByCRCs(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	crcs := c.Params.ByName("crcs")

	type ShortSurvey struct {
		CRCs   string
		Search string
		Id     int64
	}
	var surveyResp ShortSurvey
	err := dbmap.SelectOne(&surveyResp, "SELECT crcs,search,id FROM survey WHERE crcs=? LIMIT 1", crcs)
	if verbose == true {
		fmt.Println(surveyResp)
	}

	if err == nil {
		c.JSON(200, surveyResp)
	} else {
		c.JSON(404, gin.H{"error": "survey not found"})
	}

	// curl -i http://localhost:8080/api/v1/surveys/1
}

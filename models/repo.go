package models

import (
	"database/sql"
	//"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
	"log"
	"regexp"
	"time"
)

// gin Middlware to select database
func Database(connString string) gin.HandlerFunc {
	dbmap := InitDb(connString)
	return func(c *gin.Context) {
		c.Set("DBmap", dbmap)
		c.Next()
	}
}

func InitDb(dbName string) *gorp.DbMap {
	db, err := sql.Open("sqlite3", dbName)
	checkErr(err, "sql.Open failed")
	//dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Agent{}, "Agent").SetKeys(false, "CRCa")
	dbmap.AddTableWithName(Survey{}, "Survey").SetKeys(true, "Id")
	dbmap.AddTableWithName(Alerte{}, "Alerte").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func CheckAgentOffLine(dbmap *gorp.DbMap) {
	var agents []Agent
	_, err := dbmap.Select(&agents, "SELECT * FROM agent WHERE status=\"OnLine\"")
	if err == nil {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		for _, a := range agents {
			t := a.Updated.UnixNano() / int64(time.Millisecond)
			//fmt.Printf("%s %d : %d \n", a.CRCa, t, now-t)
			if now-t > 300000 {
				a.Status = "OffLine"
				dbmap.Update(&a)
			}
		}
	}
}

func ParseQuery(q map[string][]string) string {
	query := " "
	if q["_filters"] != nil {
		re := regexp.MustCompile("{\"([a-zA-Z0-9_]+?)\":\"([a-zA-Z0-9_. ]+?)\"}")
		r := re.FindStringSubmatch(q["_filters"][0])
		// TODO: special col name for all fields via reflections
		col := r[1]
		search := r[2]
		if col != "" && search != "" {
			query = query + " WHERE " + col + " LIKE \"%" + search + "%\" "
		}
	}
	if q["_sortField"] != nil && q["_sortDir"] != nil {
		sortField := q["_sortField"][0]
		// prevent SQLi
		valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
		if !valid.MatchString(sortField) {
			sortField = ""
		}
		if sortField == "created" || sortField == "updated" { // XXX trick for sqlite
			sortField = "datetime(" + sortField + ")"
		}
		sortOrder := q["_sortDir"][0]
		if sortOrder != "ASC" {
			sortOrder = "DESC"
		}
		// _page, _perPage, _sortDir, _sortField
		if sortField != "" {
			query = query + " ORDER BY " + sortField + " " + sortOrder
		}
	}
	return query
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

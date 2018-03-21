package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v2"
	//"log"
	"strconv"
	"time"
)

/**
Search for XXX to fix fields mapping in Update handler, mandatory fields
or remove sqlite tricks

 vim search and replace cmd to customize struct, handler and instances
  :%s/IP/NewStruct/g
  :%s/ip/newinst/g

**/

// XXX custom struct name and fields
type IP struct {
	Id      int64     `db:"id" json:"id"`
	Name    string    `db:"name" json:"name"`
	Host    string    `db:"host" json:"host"`
	Lat     float64   `db:"lat" json:"lat"`
	Lon     float64   `db:"lon" json:"lon"`
	P       string    `db:"p" json:"p"` //Country
	R       string    `db:"r" json:"r"` //Region
	C       string    `db:"c" json:"c"` //City
	ASNnum  string    `db:"asnnum" json:"asnnum"`
	ASNname string    `db:"asnname" json:"asnname"`
	Count   int64     `db:"count" json:"count"`
	Comment string    `db:"comment" json:"comment,size:65534"`
	Created time.Time `db:"created" json:"created"` // or int64
	Updated time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

func (e *IP) PreInsert(s gorp.SqlExecutor) error {
	e.Created = time.Now() // or time.Now().UnixNano()
	e.Updated = e.Created
	return nil
}

func (e *IP) PreUpdate(s gorp.SqlExecutor) error {
	e.Updated = time.Now()
	return nil
}

// REST handlers

func GetIPs(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	query := "SELECT * FROM ip"

	// Parse query string
	//  receive : map[_filters:[{"q":"wx"}] _sortField:[id] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	var count int64
	if s != "" {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM ip  WHERE " + s)
		query = query + " WHERE " + s
	} else {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM ip")
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

	var ips []IP
	_, err := dbmap.Select(&ips, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.FormatInt(count, 10)) // float64 to string
		c.JSON(200, ips)
	} else {
		c.JSON(404, gin.H{"error": "no ip(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/ips
}

func GetIP(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var ip IP
	err := dbmap.SelectOne(&ip, "SELECT * FROM ip WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, ip)
	} else {
		c.JSON(404, gin.H{"error": "ip not found"})
	}

	// curl -i http://localhost:8080/api/v1/ips/1
}

func PostIP(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)

	var ip IP
	c.Bind(&ip)

	//log.Println(ip)

	if ip.Name != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&ip)
		if err == nil {
			c.JSON(201, ip)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/ips
}

func UpdateIP(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var ip IP
	err := dbmap.SelectOne(&ip, "SELECT * FROM ip WHERE id=?", id)
	if err == nil {
		var json IP
		c.Bind(&json)

		//log.Println(json)
		ipId, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		ip := IP{
			Id:      ipId,
			Name:    json.Name,
			Host:    json.Host,
			Lat:     json.Lat,
			Lon:     json.Lon,
			P:       json.P,
			R:       json.R,
			C:       json.C,
			ASNnum:  json.ASNnum,
			ASNname: json.ASNname,
			Count:   json.Count,
			Comment: json.Comment,
			Created: ip.Created, //ip read from previous select
		}

		if ip.Name != "" { // XXX Check mandatory fields
			_, err = dbmap.Update(&ip)
			if err == nil {
				c.JSON(200, ip)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "ip not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/ips/1
}

func DeleteIP(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var ip IP
	err := dbmap.SelectOne(&ip, "SELECT * FROM ip WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&ip)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "ip not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/ips/1
}

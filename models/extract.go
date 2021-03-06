package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v2"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/**
Search for XXX to fix fields mapping in Update handler, mandatory fields
or remove sqlite tricks

 vim search and replace cmd to customize struct, handler and instances
  :%s/Extract/NewStruct/g
  :%s/extract/newinst/g

**/

// Extract custom struct name and fields XXX
type Extract struct {
	Id       int64     `db:"id" json:"id"`
	Search   string    `db:"search" json:"search"`
	Role     string    `db:"role" json:"role"`
	Action   string    `db:"action" json:"action"`
	Comment  string    `db:"comment" json:"comment,size:65534"`
	Active   bool      `db:"active" json:"active"`
	FromDate time.Time `db:"fromdate" json:"fromdate"`
	ToDate   time.Time `db:"todate" json:"todate"`
	Created  time.Time `db:"created" json:"created"` // or int64
	Updated  time.Time `db:"updated" json:"updated"`
}

// Hooks
// PreInsert fix date hook
func (e *Extract) PreInsert(s gorp.SqlExecutor) error {
	e.Created = time.Now() // or time.Now().UnixNano()
	e.Updated = e.Created
	return nil
}

// PreUpdate fix date hook
func (e *Extract) PreUpdate(s gorp.SqlExecutor) error {
	e.Updated = time.Now()
	return nil
}

/**

REST handlers

**/

// GetExtracts all items
func GetExtracts(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	query := "SELECT * FROM extract"

	// Parse query string
	//  receive : map[_filters:[{"q":"wx"}] _sortField:[id] ...
	q := c.Request.URL.Query()
	s, o, l := ParseQuery(q)
	var count int64
	if s != "" {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM extract  WHERE " + s)
		query = query + " WHERE " + s
	} else {
		count, _ = dbmap.SelectInt("SELECT COUNT(*) FROM extract")
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

	var extracts []Extract
	_, err := dbmap.Select(&extracts, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.FormatInt(count, 10)) // float64 to string
		c.JSON(200, extracts)
	} else {
		c.JSON(404, gin.H{"error": "no extract(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/extracts
}

// GetExtract one item by id
func GetExtract(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")

	var extract Extract
	err := dbmap.SelectOne(&extract, "SELECT * FROM extract WHERE id=? LIMIT 1", id)

	if err == nil {
		c.JSON(200, extract)
	} else {
		c.JSON(404, gin.H{"error": "extract not found"})
	}

	// curl -i http://localhost:8080/api/v1/extracts/1
}

// PostExtract add item
func PostExtract(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	claims := c.MustGet("claims").(jwt.MapClaims)
	log.Printf("%s [%s]: PostExtract\n", c.ClientIP(), claims["id"])

	var extract Extract
	c.Bind(&extract)

	//log.Println(extract)

	if extract.Search != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&extract)
		if err == nil {
			c.JSON(201, extract)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory field Search is empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/extracts
}

// UpdateExtract update item by id
func UpdateExtract(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")
	claims := c.MustGet("claims").(jwt.MapClaims)
	log.Printf("%s [%s]: UpdateExtract %s\n", c.ClientIP(), claims["id"], id)

	var extract Extract
	err := dbmap.SelectOne(&extract, "SELECT * FROM extract WHERE id=?", id)
	if err == nil {
		var json Extract
		c.Bind(&json)

		//log.Println(json)
		extractId, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		extract := Extract{
			Id:       extractId,
			Search:   json.Search,
			Role:     json.Role,
			Action:   json.Action,
			Comment:  json.Comment,
			Active:   json.Active,
			FromDate: json.FromDate,
			ToDate:   json.ToDate,
			Created:  extract.Created, //extract read from previous select
		}

		if extract.Search != "" { // XXX Check mandatory fields
			_, err = dbmap.Update(&extract)
			if err == nil {
				c.JSON(200, extract)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "extract not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/extracts/1
}

// DeleteExtract delte item by id
func DeleteExtract(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	id := c.Params.ByName("id")
	claims := c.MustGet("claims").(jwt.MapClaims)
	log.Printf("%s [%s]: DeleteExtract %s\n", c.ClientIP(), claims["id"], id)

	var extract Extract
	err := dbmap.SelectOne(&extract, "SELECT * FROM extract WHERE id=?", id)

	if err == nil {
		_, err = dbmap.Delete(&extract)

		if err == nil {
			c.JSON(200, gin.H{"id #" + id: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "extract not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/extracts/1
}

// RestExtract run action via rest
func RestExtract(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	claims := c.MustGet("claims").(jwt.MapClaims)
	log.Printf("%s [%s]: RestExtract\n", c.ClientIP(), claims["id"])

	i := ExtractSearchs(dbmap, verbose)
	c.JSON(200, gin.H{"result": i})
}

/**

**/

// ExtractSearchs run action
func ExtractSearchs(dbmap *gorp.DbMap, verbose bool) int {
	var extracts []Extract
	dbmap.Select(&extracts, "SELECT * FROM extract WHERE active = 1")
	//Searchs  = make(map[string]*regexp.Regexp)
	count := 0

	for _, e := range extracts {
		//fmt.Println(e.Search)
		re, err := regexp.Compile(e.Search)
		if err != nil {
			continue
		}
		query := "SELECT id,line,comment,updated FROM alerte where role=? "
		var alertes []Alerte
		var compress []string
		firtAlerte := Alerte{} // to store compressed
		dbmap.Select(&alertes, query, e.Role)
		for _, a := range alertes {
			//fmt.Printf("%+v\n", a)
			// Continue if outside dates
			if e.FromDate.IsZero() == false {
				if a.Updated.Before(e.FromDate) {
					//fmt.Printf(" => continue before\n")
					continue
				}
			}
			if e.ToDate.IsZero() == false {
				if a.Updated.After(e.ToDate) {
					//fmt.Printf(" => continue after\n")
					continue
				}
			}

			res := re.FindStringSubmatch(a.Line)
			if res == nil {
				//fmt.Printf(" => continue no match\n")
				continue
			}
			count++
			ip := res[1]
			if verbose == true {
				fmt.Printf("%+v\n", a)
				fmt.Printf(" => <%s>\n", ip)
			}
			if e.Action == "Compress" && firtAlerte.Line == "" {
				compress = append(compress, a.Comment)
				compress = append(compress, a.Line)
				firtAlerte = a
				continue
			}

			switch e.Action {
			case "Delete":
				// TODO dbmap.Delete(&a)
				log.Println("Extract Delete")
			case "Compress":
				log.Println("Extract Compress")
				compress = append(compress, a.Line)
				count--
				dbmap.Delete(&a)
			case "AddIP":
				log.Println("Extract AddIP")
				i, _ := CreateOrUpdateIp(dbmap, ip)
				fmt.Println(i.totxt(false))
				if a.Comment == "" {
					// Update comment Alerte
					a.Comment = i.totxt(false)
					dbmap.Exec("UPDATE Alerte SET comment = ? WHERE id = ?", a.Comment, a.Id)
				}
			}
		}

		if e.Action == "Compress" && len(compress) > 0 && firtAlerte.Line != "" {
			newComment := strings.Join(compress[:], "\n")
			if len(newComment) > 60000 {
				newComment = newComment[:60000] + "\n ... \n" + compress[len(compress)-1]
			}
			firtAlerte.Comment = strings.Join(compress[:], "\n")
			dbmap.Exec("UPDATE Alerte SET comment = ? WHERE id = ?", newComment, firtAlerte.Id)
		}

	}

	return count
}

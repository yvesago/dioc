package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"
	"hash/crc32"
	"strconv"
	"strings"
	"time"
)

/**
Search for XXX to fix fields mapping in Update handler, mandatory fields
or remove sqlite tricks

 vim search and replace cmd to customize struct, handler and instances
  :%s/Agent/NewStruct/g
  :%s/agent/newinst/g

**/

// XXX custom struct name and fields
type Agent struct {
	//	Id         int64     `db:"id" json:"id"`
	CRCa       string    `db:"crca" json:"crca"` //primary Key
	IP         string    `db:"ip" json:"ip"`
	FileSurvey string    `db:"filesurvey" json:"filesurvey"`
	Role       string    `db:"role" json:"role"`
	Comment    string    `db:"comment" json:"comment"`
	Lines      string    `db:"lines" json:"lines"`
	Status     string    `db:"status" json:"status"`
	CMD        string    `db:"cmd" json:"cmd"`
	Salt       string    `db:"-" json:"-"`             // not registred in database
	Created    time.Time `db:"created" json:"created"` // or int64
	Updated    time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

func hashAgent(Salt string, IP string, File string) string {
	crc32q := crc32.MakeTable(0xD5828281)
	crca := fmt.Sprintf("%08x", crc32.Checksum([]byte(Salt+IP+":"+File), crc32q))
	return crca
}

func (a *Agent) PreInsert(s gorp.SqlExecutor) error {
	salt := a.Salt
	a.CRCa = hashAgent(salt, a.IP, a.FileSurvey)
	a.CMD = "SendLines"
	a.Created = time.Now() // or time.Now().UnixNano()
	a.Updated = a.Created
	return nil
}

func (a *Agent) PreUpdate(s gorp.SqlExecutor) error {
	a.Updated = time.Now()
	return nil
}

// REST handlers

func GetAgents(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	verbose := c.MustGet("Verbose").(bool)
	query := "SELECT * FROM agent"

	// Parse query string
	q := c.Request.URL.Query()
	tmpquery := query + ParseQuery(q)
	query = strings.Replace(tmpquery, "ORDER BY id ", "ORDER BY crca ", 1) // workaround for ng-admin bug
	if verbose == true {
		fmt.Println(q)
		fmt.Println("query: " + query)
	}

	var agents []Agent
	_, err := dbmap.Select(&agents, query)

	if err == nil {
		c.Header("X-Total-Count", strconv.Itoa(len(agents)))
		c.JSON(200, agents)
	} else {
		c.JSON(404, gin.H{"error": "no agent(s) into the table"})
	}

	// curl -i http://localhost:8080/api/v1/agents
}

func GetAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	//id := c.Params.ByName("id")
	crca := c.Params.ByName("crca")
	//fmt.Println(" -- " + id + " -- " + crca)

	var agent Agent
	//err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE id=? LIMIT 1", id)
	err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE crca=? LIMIT 1", crca)

	if err == nil {
		c.JSON(200, agent)
	} else {
		c.JSON(404, gin.H{"error": "agent not found"})
	}

	// curl -i http://localhost:8080/api/v1/agents/1
}

func PostAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)

	var agent Agent
	c.Bind(&agent)

	//log.Println(agent)

	if agent.IP != "" && agent.FileSurvey != "" { // XXX Check mandatory fields
		err := dbmap.Insert(&agent)
		if err == nil {
			c.JSON(201, agent)
		} else {
			checkErr(err, "Insert failed")
		}

	} else {
		c.JSON(400, gin.H{"error": "Mandatory fields are empty"})
	}

	// curl -i -X POST -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Queen\" }" http://localhost:8080/api/v1/agents
}

func UpdateAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	//id := c.Params.ByName("id")
	crca := c.Params.ByName("crca")

	var agent Agent
	//err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE id=?", id)
	err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE crca=?", crca)
	if err == nil {
		var json Agent
		c.Bind(&json)

		//log.Println(json)
		//agent_id, _ := strconv.ParseInt(id, 0, 64)

		//TODO : find fields via reflections
		//XXX custom fields mapping
		agent := Agent{
			CRCa:       crca,
			IP:         json.IP,
			FileSurvey: json.FileSurvey,
			Role:       json.Role,
			Lines:      json.Lines,
			CMD:        json.CMD,
			Comment:    json.Comment,
			Status:     json.Status,
			Created:    agent.Created, //agent read from previous select
		}

		if agent.CRCa != "" { // XXX Check mandatory fields
			_, err = dbmap.Update(&agent)
			if err == nil {
				c.JSON(200, agent)
			} else {
				checkErr(err, "Updated failed")
			}

		} else {
			c.JSON(400, gin.H{"error": "mandatory fields are empty"})
		}

	} else {
		c.JSON(404, gin.H{"error": "agent not found"})
	}

	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\": \"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/agents/1
}

func DeleteAgent(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	//id := c.Params.ByName("id")
	crca := c.Params.ByName("crca")

	var agent Agent
	//err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE id=?", id)
	err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE crca=?", crca)

	if err == nil {
		_, err = dbmap.Delete(&agent)

		if err == nil {
			//c.JSON(200, gin.H{"id #" + id: "deleted"})
			c.JSON(200, gin.H{"crca #" + crca: "deleted"})
		} else {
			checkErr(err, "Delete failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "agent not found"})
	}

	// curl -i -X DELETE http://localhost:8080/api/v1/agents/1
}

/**

Agent handlers


**/

func RegisterHandler(c *gin.Context) {
	verbose := c.MustGet("Verbose").(bool)
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)

	var json Agent
	c.Bind(&json)
	if verbose {
		fmt.Println(json)
	}
	if json.FileSurvey != "" {
		var agent Agent
		agent = Agent{
			IP:         c.ClientIP(),
			FileSurvey: json.FileSurvey,
		}
		salt := c.MustGet("Salt").(string)
		agent.Salt = salt
		crca := hashAgent(salt, agent.IP, agent.FileSurvey)
		// Check if exist
		err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE crca=?", crca)
		if err == nil {
			c.JSON(201, agent.CRCa)
		} else {
			err := dbmap.Insert(&agent)
			if err == nil {
				c.JSON(201, agent.CRCa)
			} else {
				checkErr(err, "Insert failed")
			}
		}
	} else {
		c.JSON(400, gin.H{"error": "missing mandatory file"})
	}
}

func SendLinesHandler(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	crca := c.Params.ByName("crca")
	//			fmt.Println("======>"+crca+"<==")
	var agent Agent
	err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE crca=?", crca)
	if err == nil {
		var json Agent
		c.Bind(&json)
		//		fmt.Println(json)
		if json.Lines != "" {
			agent.Lines = json.Lines
			agent.CMD = ""
			//fmt.Println("=======>"+agent.Lines)
			_, err = dbmap.Update(&agent)
			if err == nil {
				c.JSON(201, "0K")
			} else {
				checkErr(err, "Update failed")
			}
		} else {
			c.JSON(400, gin.H{"error": "missing mandatory lines"})
		}
	} else {
		c.JSON(400, gin.H{"error": "bad crca"})
	}

}

func CMDHandler(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	crca := c.Params.ByName("crca")
	var agent Agent
	err := dbmap.SelectOne(&agent, "SELECT * FROM agent WHERE crca=?", crca)
	if err == nil {
		var listeSurveys []string
		dbmap.Select(&listeSurveys, "SELECT crcs FROM survey WHERE role=?", agent.Role)
		agent.Status = "OnLine"
		_, err = dbmap.Update(&agent)
		if err == nil {
			c.JSON(200, gin.H{"CMD": agent.CMD, "ListeSurveys": listeSurveys})
		} else {
			checkErr(err, "Update failed")
		}
	} else {
		c.JSON(400, gin.H{"error": "bad crca"})
	}
}

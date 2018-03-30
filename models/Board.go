package models

import (
	//	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v2"
)

type m map[string]int64

type Board struct {
	Id         int64     `db:"id" json:"id"`
	Agents     []m       `json:"agents" db:"-"`  // ignore sql
	Surveys    []m       `json:"surveys" db:"-"` // ignore sql
	Alerts     []m       `json:"alerts" db:"-"`  // ignore sql
	DocAgents  string    `json:"docagents" db:"docagents,size:65534"`
	DocSearchs string    `json:"docsearchs" db:"docsearchs,size:65534"`
	Docs       string    `json:"docs" db:"docs,size:65534"`
	Created    time.Time `db:"created" json:"created"`
	Updated    time.Time `db:"updated" json:"updated"`
}

// Hooks : PreInsert and PreUpdate

func (b *Board) PreInsert(s gorp.SqlExecutor) error {
	b.Created = time.Now() // or time.Now().UnixNano()
	b.Updated = b.Created
	return nil
}

func (b *Board) PreUpdate(s gorp.SqlExecutor) error {
	b.Updated = time.Now()
	return nil
}

// Methods

func (b *Board) Load(tx *gorp.DbMap) error {

	var status []string
	_, err := tx.Select(&status, "select DISTINCT(status) from Agent")
	if err != nil {
		return err
	}
	for _, s := range status {
		c, _ := tx.SelectInt("select Count(*) from Agent where status=?", s)
		b.Agents = append(b.Agents, m{s: c})
	}

	var roles []string
	_, err = tx.Select(&roles, "select DISTINCT(role) from Survey")
	if err != nil {
		return err
	}
	for _, r := range roles {
		c, _ := tx.SelectInt("select Count(*) from Survey where role=?", r)
		b.Surveys = append(b.Surveys, m{r: c})
	}

	var rls []string
	_, err = tx.Select(&rls, "select DISTINCT(role) from Alerte")
	if err != nil {
		return err
	}
	for _, r := range rls {
		c, _ := tx.SelectInt("select Count(*) from Alerte where role=?", r)
		b.Alerts = append(b.Alerts, m{r: c})
	}
	//fmt.Printf("%+v\n", b)
	return err
}

// REST handlers

func GetBoard(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)

	var b Board
	dbmap.SelectOne(&b, "SELECT * FROM Board WHERE id=1")

	err := b.Load(dbmap)
	if err == nil {
		c.JSON(200, b)
	} else {
		c.JSON(404, gin.H{"error": "board not found"})
	}

	// curl -i http://localhost:8080/api/v1/board
}

func UpdateBoard(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)
	claims := c.MustGet("claims").(jwt.MapClaims)
	log.Printf("[%s] UpdateBoard\n", claims["id"])

	var b Board
	err := dbmap.SelectOne(&b, "SELECT * FROM Board WHERE id=1")
	if err == nil {
		var json Board
		c.Bind(&json)

		newb := Board{
			Id:         1,
			DocAgents:  json.DocAgents,
			DocSearchs: json.DocSearchs,
			Docs:       json.Docs,
		}
		_, err = dbmap.Update(&newb)
		if err == nil {
			c.JSON(200, newb)
		} else {
			checkErr(err, "Updated failed")
		}

	} else {
		c.JSON(404, gin.H{"error": "board not found"})
	}
}

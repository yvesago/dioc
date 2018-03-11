package models

import (
	//	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"
)

type m map[string]int64

type Board struct {
	Agents  []m `json:"agents" db:"-"`  // ignore sql
	Surveys []m `json:"surveys" db:"-"` // ignore sql
}

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
	//fmt.Printf("%+v\n", b)
	return err
}

func GetBoard(c *gin.Context) {
	dbmap := c.MustGet("DBmap").(*gorp.DbMap)

	var b Board
	err := b.Load(dbmap)

	if err == nil {
		c.JSON(200, b)
	} else {
		c.JSON(404, gin.H{"error": "board not found"})
	}

	// curl -i http://localhost:8080/api/v1/board
}

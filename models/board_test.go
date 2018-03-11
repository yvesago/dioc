package models

import (
	//	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBoardModel(t *testing.T) {
	defer deleteFile(config.DBname)

	dbmap := InitDb(config.DBname)
	dbmap.TraceOn("[gorp]", log.New(os.Stdout, "myapp:", log.Lmicroseconds))
	a1 := Agent{IP: "1", FileSurvey: "f", Status: "OffLine"}
	a2 := Agent{IP: "2", FileSurvey: "f", Status: "OffLine"}
	a3 := Agent{IP: "3", FileSurvey: "f", Status: "OnLine"}
	dbmap.Insert(&a1)
	dbmap.Insert(&a2)
	dbmap.Insert(&a3)

	s1 := Survey{Search: "a", Role: "test"}
	s2 := Survey{Search: "b", Role: "test"}
	s3 := Survey{Search: "c", Role: "test2"}
	dbmap.Insert(&s1)
	dbmap.Insert(&s2)
	dbmap.Insert(&s3)

	al := Alerte{CRCa: a1.CRCa, CRCs: s1.CRCs, Role: "test"}
	dbmap.Insert(&al)

	b := Board{}
	b.Load(dbmap)
	d, _ := json.Marshal(b)
	fmt.Println(string(d))
	assert.Equal(t, `{"agents":[{"OffLine":2},{"OnLine":1}],"surveys":[{"test":2},{"test2":1}],"alerts":[{"test":1}]}`, string(d), "todo")

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var urla = "/admin/api/v1/board"
	router.GET(urla, GetBoard)

	// Get Board
	log.Println("= http GET Board")
	req, err := http.NewRequest("GET", urla, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println("----", resp.Body)
	var b2 Board
	json.Unmarshal(resp.Body.Bytes(), &b2)
	//fmt.Println(b2)
	a := b2.Agents[0]
	assert.Equal(t, int64(2), a["OffLine"], "2 results")

}

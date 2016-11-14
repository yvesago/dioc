package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestModelAlerte(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var urla = "/admin/api/v1/alertes"
	router.POST(urla, PostAlerte)
	router.GET(urla, GetAlertes)
	router.GET(urla+"/:id", GetAlerte)
	router.DELETE(urla+"/:id", DeleteAlerte)
	router.PUT(urla+"/:id", UpdateAlerte)

	// Add
	log.Println("= http POST Alerte")
	var a = Alerte{CRCa: "xxx", CRCs: "yyyy", Line: "something"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(a)
	req, err := http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Add second alerte
	log.Println("= http POST more Alerte")
	var a2 = Alerte{CRCa: "xxx", CRCs: "yyyy", Line: "something"}
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("POST", urla, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Get all
	log.Println("= http GET all Alertes")
	req, err = http.NewRequest("GET", urla, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []Alerte
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 2, len(as), "2 results")

	log.Println("= Test parsing query")
	s := "http://127.0.0.1:8080/api?_filters={\"line\":\"t\"}&_sortDir=ASC&_sortField=created"
	u, _ := url.Parse(s)
	q, _ := url.ParseQuery(u.RawQuery)
	//fmt.Println(q)
	query := ParseQuery(q)
	//fmt.Println(query)
	assert.Equal(t, "  WHERE line LIKE \"%t%\" ORDER BY datetime(created) ASC", query, "Parse query")

	log.Println("= Test parsing page query")
	s = "http://127.0.0.1:8080/api?_perPage=5&_page=1"
	u, _ = url.Parse(s)
	q, _ = url.ParseQuery(u.RawQuery)
	//fmt.Println(q)
	query = ParseQuery(q)
	//fmt.Println(query)
	assert.Equal(t, "  LIMIT 5", query, "Parse query")

	log.Println("= Test parsing page query")
	s = "http://127.0.0.1:8080/api?_perPage=5&_page=2"
	u, _ = url.Parse(s)
	q, _ = url.ParseQuery(u.RawQuery)
	//fmt.Println(q)
	query = ParseQuery(q)
	//fmt.Println(query)
	assert.Equal(t, "  LIMIT 5 OFFSET 6", query, "Parse query")

	log.Println("= Test parsing multi filter query")
	s = "http://127.0.0.1:8080/api?_filters={\"line\":\"t\",\"line2\":\"t2\"}&_sortDir=DESC&_sortField=created"
	u, _ = url.Parse(s)
	q, _ = url.ParseQuery(u.RawQuery)
	//fmt.Println(q)
	query = ParseQuery(q)
	//fmt.Println(query)

	// Managed unsorted queries map
	res1 := "  WHERE line LIKE \"%t%\" AND line2 LIKE \"%t2%\" ORDER BY datetime(created) DESC"
	res2 := "  WHERE line2 LIKE \"%t2%\" AND line LIKE \"%t%\" ORDER BY datetime(created) DESC"
	if res1 != query && res2 != query {
		assert.Equal(t, res1+" -- OR -- "+res2, query, "Parse query")
	}

	// Get one
	log.Println("= http GET one Alerte")
	var a1 Alerte
	req, err = http.NewRequest("GET", urla+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a1.CRCa, a.CRCa, "a1 = a")

	// Delete one
	log.Println("= http DELETE one Alerte")
	req, err = http.NewRequest("DELETE", urla+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http DELETE success")
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	req, err = http.NewRequest("GET", urla, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all for count success")
	//fmt.Println(resp.Body)
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 1, len(as), "1 result")

	req, _ = http.NewRequest("GET", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "No more /1")
	req, _ = http.NewRequest("DELETE", urla+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "No more /1")

	// Update one
	log.Println("= http PUT one Alerte")
	a2.Comment = "Comment test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("PUT", urla+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	req, _ = http.NewRequest("PUT", urla+"/1", b)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "Can't update /1")

	var a3 Alerte
	req, err = http.NewRequest("GET", urla+"/2", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all updated success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a2.Comment, a3.Comment, "a2 Comment updated")

}

func TestAlerte(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var urlagent = "/agent/api/v1/agent"
	router.POST(urlagent, RegisterHandler)
	var urladminagent = "/admin/api/v1/agent"
	router.PUT(urladminagent+"/:crca", UpdateAgent)
	var urlsurvey = "/admin/api/v1/surveys"
	router.POST(urlsurvey, PostSurvey)
	var urlalerte = "/agent/api/v1/alerte"
	router.POST(urlalerte, PostNewAlerte)

	// Add agent
	log.Println("= POST Register Agent")
	var jsonStr = []byte(`{"FileSurvey":"/some/file","Hostname":"xxxx"}`)
	req, err := http.NewRequest("POST", urlagent, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Body)
	var crca string
	json.NewDecoder(resp.Body).Decode(&crca)
	assert.Equal(t, 201, resp.Code, "http Register success")
	fmt.Println("====", urladminagent)

	// remove cmd to test mail
	var a = Agent{CRCa: crca}
	b := new(bytes.Buffer)
	a.CMD = ""
	json.NewEncoder(b).Encode(a)
	req, err = http.NewRequest("PUT", urladminagent+"/"+crca, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	fmt.Println(resp)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	// Add survey
	log.Println("= http POST Survey")
	var s = Survey{Search: "something", Level: "warn"}
	b = new(bytes.Buffer)
	json.NewEncoder(b).Encode(s)
	req, err = http.NewRequest("POST", urlsurvey, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &s) //read CRCs
	assert.Equal(t, 201, resp.Code, "http POST success")
	crcs := s.CRCs

	// Add alerte
	log.Println("= http POST New Alerte")
	type NewAlerte struct {
		CRCa string
		CRCs string
		Line string
	}
	var newalerte = NewAlerte{CRCa: crca, CRCs: crcs, Line: "Some line"}
	json.NewEncoder(b).Encode(newalerte)
	req, err = http.NewRequest("POST", urlalerte, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Body)
	assert.Equal(t, 201, resp.Code, "http POST success")

}

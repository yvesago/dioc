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
	"os"
	"testing"
)

func TestExtractModel(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var url = "/admin/api/v1/extracts"
	router.POST(url, PostExtract)
	router.GET(url, GetExtracts)
	router.GET(url+"/:id", GetExtract)
	router.DELETE(url+"/:id", DeleteExtract)
	router.PUT(url+"/:id", UpdateExtract)

	// Add
	log.Println("= http POST Extract")
	var a = Extract{Search: "something"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(a)
	req, err := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &a)
	assert.Equal(t, 201, resp.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Add second extract
	log.Println("= http POST more Extract")
	var a2 = Extract{Search: "something else"}
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &a2)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Get all
	log.Println("= http GET all Extracts")
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	//fmt.Println(resp.Body)
	var as []Extract
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 2, len(as), "2 results")

	// Get one
	log.Println("= http GET one Extract")
	var a1 Extract
	req, err = http.NewRequest("GET", url+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(resp.Body)
	assert.Equal(t, a.Id, a1.Id, "a1 = a")

	// Delete one
	log.Println("= http DELETE one Extract")
	req, err = http.NewRequest("DELETE", url+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	req, err = http.NewRequest("GET", url+"?_start=1&_end=5&_sortField=id&_sortDir=ASC&_filters={\"search\":\"some\"}", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	//fmt.Println(resp.Body)
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 1, len(as), "1 result")

	req, _ = http.NewRequest("GET", url+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "No more /1")
	req, _ = http.NewRequest("DELETE", url+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "No more /1")

	// Update one
	log.Println("= http PUT one Extract")
	a2.Role = "Role test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("PUT", url+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	var a3 Extract
	req, err = http.NewRequest("GET", url+"/2", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a2.Role, a3.Role, "a2 Role updated")

	req, _ = http.NewRequest("PUT", url+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "Can't update /1")

}

func TestExtract(t *testing.T) {
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

	al := Alerte{CRCa: a1.CRCa, CRCs: s1.CRCs, Role: "test", Line: "some log with IP: 192.168.1.1"}
	al2 := Alerte{CRCa: a1.CRCa, CRCs: s1.CRCs, Role: "test", Line: "some log with IP: 192.168.1.2"}
	dbmap.Insert(&al)
	dbmap.Insert(&al2)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var url = "/admin/api/v1/extracts"
	router.POST(url, PostExtract)
	var url2 = "/admin/api/v1/actionextract"
	router.PUT(url2, RestExtract)

	// AddIP
	log.Println("= http POST Extract")
	var e = Extract{Search: "IP: (.*)", Active: true, Role: "test", Action: "AddIP"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(e)
	req, err := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &e)
	assert.Equal(t, 201, resp.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Searches with rest request
	r, _ := http.NewRequest("PUT", url2, nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, r)
	assert.Equal(t, 200, resp.Code, "http PUT success")
	fmt.Println(resp.Body)
	var jresp map[string]int
	json.Unmarshal(resp.Body.Bytes(), &jresp)
	assert.Equal(t, 2, jresp["result"], "2 alerte line match")
	e.Active = false
	dbmap.Update(&e)

	ip, _ := CreateOrUpdateIp(dbmap, "192.168.1.1")
	assert.Equal(t, int64(2), ip.Count, "create IP")

	// Compress
	log.Println("= http POST Compress")
	var e2 = Extract{Search: "IP: (.*)", Active: true, Role: "test", Action: "Compress"}
	b2 := new(bytes.Buffer)
	json.NewEncoder(b2).Encode(e2)
	req2, err2 := http.NewRequest("POST", url, b2)
	req2.Header.Set("Content-Type", "application/json")
	if err2 != nil {
		fmt.Println(err2)
	}
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)
	json.Unmarshal(resp2.Body.Bytes(), &e2)
	assert.Equal(t, 201, resp2.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Searches
	log.Println("= ExtractSearchs")
	res2 := ExtractSearchs(dbmap, true)
	assert.Equal(t, 1, res2, "1 alerte line remain")
	e2.Active = false
	dbmap.Update(&e2)

}

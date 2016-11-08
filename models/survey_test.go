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
	"testing"
)

func TestSurveyModel(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var url = "/admin/api/v1/surveys"
	router.POST(url, PostSurvey)
	router.GET(url, GetSurveys)
	router.GET(url+"/:id", GetSurvey)
	router.DELETE(url+"/:id", DeleteSurvey)
	router.PUT(url+"/:id", UpdateSurvey)

	// Add
	log.Println("= http POST Survey")
	var a = Survey{Search: "something"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(a)
	req, err := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &a) //read CRCs
	assert.Equal(t, 201, resp.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Add second survey
	log.Println("= http POST more Survey")
	var a2 = Survey{Search: "something else"}
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &a2) //read CRCs
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Get all
	log.Println("= http GET all Surveys")
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	//fmt.Println(resp.Body)
	var as []Survey
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 2, len(as), "2 results")

	// Get one
	log.Println("= http GET one Survey")
	var a1 Survey
	req, err = http.NewRequest("GET", url+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(resp.Body)
	assert.Equal(t, a.CRCs, a1.CRCs, "a1 = a")

	// Delete one
	log.Println("= http DELETE one Survey")
	req, err = http.NewRequest("DELETE", url+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	req, err = http.NewRequest("GET", url, nil)
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

	// Update one
	log.Println("= http PUT one Survey")
	a2.Role = "Role test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("PUT", url+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	var a3 Survey
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
	assert.NotEqual(t, a2.CRCs, a3.CRCs, "Change CRCs on role update")

}

func TestSurvey(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var url = "/admin/api/v1/surveys"
	router.POST(url, PostSurvey)

	var urlagent = "/agent/api/v1/surveys"
	router.GET(urlagent+"/:crcs", GetSurveyByCRCs)

	// Add
	log.Println("= http POST Survey")
	var a = Survey{Search: "something"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(a)
	req, err := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &a) //read CRCs
	assert.Equal(t, 201, resp.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Get by CRCs
	log.Println("= http GET one Survey by CRCs")
	var a1 Survey
	req, err = http.NewRequest("GET", urlagent+"/"+a.CRCs, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(resp.Body)
	assert.Equal(t, a.CRCs, a1.CRCs, "a1 = a")
}

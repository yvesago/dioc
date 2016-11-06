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


func TestAgentModel(t *testing.T) {
	defer deleteFile(config.DBname)


	gin.SetMode(gin.TestMode)
	router := gin.New()
    router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var url = "/admin/api/v1/agents"
	router.POST(url, PostAgent)
	router.GET(url, GetAgents)
	router.GET(url+"/:crca", GetAgent)
	router.DELETE(url+"/:crca", DeleteAgent)
	router.PUT(url+"/:crca", UpdateAgent)

	// Add
	log.Println("= http POST Agent")
	var a = Agent{IP: "192.168.1.1", FileSurvey: "/some/file"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(a)
	req, err := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &a) //read CRCa
	assert.Equal(t, 201, resp.Code, "http POST success")
	//fmt.Println(resp.Body)

	// Add second agent
	log.Println("= http POST more Agent")
	var a2 = Agent{IP: "192.168.1.2", FileSurvey: "/some/file"}
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &a2) //read CRCa
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Get all
	log.Println("= http GET all Agents")
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET all success")
	//fmt.Println(resp.Body)
	var as []Agent
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 2, len(as), "2 results")

	// Get one
	log.Println("= http GET one Agent")
	var a1 Agent
	req, err = http.NewRequest("GET", url+"/"+a.CRCa, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	//fmt.Println(resp.Body)
	assert.Equal(t, a.CRCa, a1.CRCa, "a1 = a")

	// Delete one
	log.Println("= http DELETE one Agent")
	req, err = http.NewRequest("DELETE", url+"/"+a.CRCa, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http DELETE success")
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	req, err = http.NewRequest("GET", url, nil)
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

	// Update one
	log.Println("= http PUT one Agent")
	a2.Role = "Role test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("PUT", url+"/"+a2.CRCa, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	var a3 Agent
	req, err = http.NewRequest("GET", url+"/"+a2.CRCa, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http GET one updated success")
	json.Unmarshal(resp.Body.Bytes(), &a3)
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	assert.Equal(t, a2.Role, a3.Role, "a2 Role updated")

}

func TestAgent(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
    router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var url = "/agent/api/v1/agent"
	router.POST(url, RegisterHandler)
	router.PUT(url+"/:crca", SendLinesHandler)
	router.GET(url+"/:crca", CMDHandler)

	// Add
	log.Println("= POST Register Agent")
	var jsonStr = []byte(`{"FileSurvey":"/some/file"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
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

	// CMD
	log.Println("= GET CMD")
	req, err = http.NewRequest("GET", url+"/"+crca, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Body)
	type CMDresp struct {
		CMD          string
		ListeSurveys []string
	}
	var cmd CMDresp
	json.Unmarshal(resp.Body.Bytes(), &cmd)
	assert.Equal(t, "SendLines", cmd.CMD, "CMD SendLines after Register")

	// SendLines
	log.Println("= PUT SendLines")
	jsonStr = []byte(`{"Lines":"line1\nline2"}`)
	req, err = http.NewRequest("PUT", url+"/"+crca, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Body)
	assert.Equal(t, 201, resp.Code, "http Update lines success")
	// next CMD
	log.Println("= GET next CMD")
	req, err = http.NewRequest("GET", url+"/"+crca, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	//fmt.Println(resp.Body)
	json.Unmarshal(resp.Body.Bytes(), &cmd)
	assert.Equal(t, "", cmd.CMD, "no more CMD SendLines after lines sended")
}

package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	geojson "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIPModel(t *testing.T) {
	defer deleteFile(config.DBname)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var url = "/admin/api/v1/ips"
	router.POST(url, PostIP)
	router.GET(url, GetIPs)
	router.GET(url+"/:id", GetIP)
	router.DELETE(url+"/:id", DeleteIP)
	router.PUT(url+"/:id", UpdateIP)

	// Add
	log.Println("= http POST IP")
	var a = IP{Name: "2.125.160.216"}
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

	// Add second ip
	log.Println("= http POST more IP")
	var a2 = IP{Name: "something else"}
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &a2)
	assert.Equal(t, 201, resp.Code, "http POST success")

	// Get all
	log.Println("= http GET all IPs")
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	//fmt.Println(resp.Body)
	var as []IP
	json.Unmarshal(resp.Body.Bytes(), &as)
	//fmt.Println(len(as))
	assert.Equal(t, 2, len(as), "2 results")

	// Get one
	log.Println("= http GET one IP")
	var a1 IP
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
	log.Println("= http DELETE one IP")
	req, err = http.NewRequest("DELETE", url+"/1", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	//fmt.Println(a1.Name)
	//fmt.Println(resp.Body)
	req, err = http.NewRequest("GET", url+"?_start=1&_end=5&_sortField=id&_sortDir=ASC&_filters={\"name\":\"some\"}", nil)
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
	log.Println("= http PUT one IP")
	a2.Comment = "Comment test2 updated"
	json.NewEncoder(b).Encode(a2)
	req, err = http.NewRequest("PUT", url+"/2", b)
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http PUT success")

	var a3 IP
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
	assert.Equal(t, a2.Comment, a3.Comment, "a2 Comment updated")

	req, _ = http.NewRequest("PUT", url+"/1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 404, resp.Code, "Can't update /1")

}

func TestIP(t *testing.T) {
	defer deleteFile(config.DBname)

	e := InitLocDbs(config.CityDB, config.AsnDB)
	if e != nil {
		fmt.Println(e)
		return
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(SetConfig(config))
	router.Use(Database(config.DBname))

	var url = "/admin/api/v1/ips"
	router.POST(url, PostIP)
	router.GET("/geojson", GetGeoJsonIPs)
	router.GET("/actionfluship", RestFlushIP)

	// Add
	log.Println("= http POST IP")
	var a = IP{Name: "81.2.69.142"}
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
	assert.Equal(t, "Londres", a.C, "Geoloc London City")
	fmt.Println(resp.Body)

	log.Println("= http GET geojson")
	var a1 IP
	req, err = http.NewRequest("GET", "/geojson", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
	json.Unmarshal(resp.Body.Bytes(), &a1)
	fmt.Println(resp.Body)
	a2, _ := geojson.UnmarshalFeatureCollection(resp.Body.Bytes())
	//fmt.Println(a2)
	assert.Equal(t, "GB", a2.Features[0].Properties["Loc"], "IP loc country in GeoJson")

	log.Println("= http GET actionfluship")
	req, err = http.NewRequest("GET", "/actionfluship", nil)
	if err != nil {
		fmt.Println(err)
	}
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code, "http success")
}

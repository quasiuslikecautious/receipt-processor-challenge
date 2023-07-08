package handler_test

import (
  "encoding/json"
  "fmt"
  "math"
	"net/http"
	"net/http/httptest"
  "strings"
	"testing"

  "quasiuslikecautious/receipt-processor-challenge/app"
  "quasiuslikecautious/receipt-processor-challenge/config"
)

func TestPostReceiptEmpty(t *testing.T) {
  config := config.DefaultConfig()

  app := &app.App{}
  app.Initialize(config)

	// Create a new HTTP request
	request, err := http.NewRequest("POST", "/receipts/process", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create an HTTP response recorder
	recorder := httptest.NewRecorder()

	// Perform the request using the router
	app.Router.ServeHTTP(recorder, request)

	// Check the response status code
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d but got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestPostReceiptExample1(t *testing.T) {
  config := config.DefaultConfig()

  app := &app.App{}
  app.Initialize(config)

  // check valid receipt is stored
  jsonPayload := `{
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
      {
        "shortDescription": "Mountain Dew 12PK",
        "price": "6.49"
      },{
        "shortDescription": "Emils Cheese Pizza",
        "price": "12.25"
      },{
        "shortDescription": "Knorr Creamy Chicken",
        "price": "1.26"
      },{
        "shortDescription": "Doritos Nacho Cheese",
        "price": "3.35"
      },{
        "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
        "price": "12.00"
      }
    ],
    "total": "35.35"
  }`

  requestBody := strings.NewReader(jsonPayload)

  request, err := http.NewRequest("POST", "/receipts/process", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("POST Expected status %d but got %d", http.StatusOK, recorder.Code)
	}

  // check points calculation
  var responseBody map[string]interface{}
  err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
  if err != nil {
    t.Fatal(err)
  }

  id := responseBody["id"].(string)
  pointsRoute := fmt.Sprintf("/receipts/%s/points", id)

  request, err = http.NewRequest("GET", pointsRoute, nil)
  if err != nil {
    t.Fatal(err)
  }

  recorder = httptest.NewRecorder()

  app.Router.ServeHTTP(recorder, request)

  if recorder.Code != http.StatusOK {
    t.Errorf("GET Expected status %d but got %d", http.StatusOK, recorder.Code)
  }

  err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
  if err != nil {
    fmt.Println("response Body:", string(recorder.Body.Bytes()))
    t.Fatal(err)
  }

  points := responseBody["points"].(float64)
  parsedPoints := int(math.Round(points))
  if err != nil {
    t.Fatal(err)
  }

  if parsedPoints != 28 {
    t.Errorf("Expected %d points but got %d points", 28, parsedPoints)
  }
}

func TestPostReceiptExample2(t *testing.T) {
  config := config.DefaultConfig()

  app := &app.App{}
  app.Initialize(config)

  // check valid receipt is stored
  jsonPayload := `{
    "retailer": "M&M-Corner-Market",
    "purchaseDate": "2022-03-20",
    "purchaseTime": "14:33",
    "items": [
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      },{
        "shortDescription": "Gatorade",
        "price": "2.25"
      },{
        "shortDescription": "Gatorade",
        "price": "2.25"
      },{
        "shortDescription": "Gatorade",
        "price": "2.25"
      }
    ],
    "total": "9.00"
  }`

  requestBody := strings.NewReader(jsonPayload)

  request, err := http.NewRequest("POST", "/receipts/process", requestBody)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	app.Router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("POST Expected status %d but got %d", http.StatusOK, recorder.Code)
	}

  // check points calculation
  var responseBody map[string]interface{}
  err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
  if err != nil {
    t.Fatal(err)
  }

  id := responseBody["id"].(string)
  pointsRoute := fmt.Sprintf("/receipts/%s/points", id)

  request, err = http.NewRequest("GET", pointsRoute, nil)
  if err != nil {
    t.Fatal(err)
  }

  recorder = httptest.NewRecorder()

  app.Router.ServeHTTP(recorder, request)

  if recorder.Code != http.StatusOK {
    t.Errorf("GET Expected status %d but got %d", http.StatusOK, recorder.Code)
  }

  err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
  if err != nil {
    fmt.Println("response Body:", string(recorder.Body.Bytes()))
    t.Fatal(err)
  }

  points := responseBody["points"].(float64)
  parsedPoints := int(math.Round(points))
  if err != nil {
    t.Fatal(err)
  }

  if parsedPoints != 109 {
    t.Errorf("Expected %d points but got %d points", 109, parsedPoints)
  }
}

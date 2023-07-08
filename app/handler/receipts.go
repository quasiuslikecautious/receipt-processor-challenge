package handler

import (
  "math"
  "net/http"
  "strings"
  "time"
  "unicode"

  "github.com/gin-gonic/gin"
  "github.com/google/uuid"
  log "github.com/sirupsen/logrus"
  "quasiuslikecautious/receipt-processor-challenge/app/model"
)

// processReceipt adds a receipt from JSON received in the request body.
func PostReceipt(c *gin.Context, db *model.MemoryDB) {
  log.Debug("Received call on post receipt.")
  var newReceipt model.Receipt

  // Call BindJSON to bind the received JSON to
  // newReceipt
  if err := c.ShouldBindJSON(&newReceipt); err != nil {
    log.Error(err)
    c.Status(http.StatusBadRequest)
    return
  }

  if err := newReceipt.Validate(); err != nil {
    log.Error(err)
    c.Status(http.StatusBadRequest)
    return
  }

  id := uuid.New()
  (*db)[id] = newReceipt
  log.Debug("Added ", id, " to db.")

  c.IndentedJSON(http.StatusOK, gin.H{"id": id.String()})
}

// getReceiptPoints locates the receipt by matching id parameter sent
// by the client, and the returns corresponding points as a response.
func GetReceiptPoints(c *gin.Context, db *model.MemoryDB) {
  log.Debug("Received call on get receipt.")
  id := c.Param("id")

  parsedUUID, err := uuid.Parse(id)
  if err != nil {
    log.Error(err)
    c.Status(http.StatusNotFound)
    return
  }

  if r, ok := (*db)[parsedUUID]; ok {
    c.IndentedJSON(http.StatusOK, gin.H{"points": calculateReceiptPoints(&r)})
    return
  }
  
  log.Error("Failed to find receipt ", id, " from db.")
  c.Status(http.StatusNotFound)
}

// Calculate the cummulative points from a receipt
func calculateReceiptPoints(r *model.Receipt) int {
  receiptPoints := alphanumericPoints(r) + roundTotalPoints(r) + modQuarterPoints(r) + pairsPoints(r) + itemDescriptionLengthPoints(r) + oddDayPoints(r) + purchaseTimePoints(r)
 log.Debug("Total points: ", receiptPoints)
 return receiptPoints
}

// A point is gained for every alphanumeric character in the retailer name
func alphanumericPoints(r *model.Receipt) int {
  count := 0
  for _, char := range r.Retailer {
    if unicode.IsLetter(char) || unicode.IsDigit(char) {
      count++
    }
  }

  log.Debug("alphanumeric points: ", count)
  return count 
}

// if the total is a round number (i.e. ends in .00), 50 points are awarded
func roundTotalPoints(r *model.Receipt) int {
  rounded := math.Round(float64(r.Total))

  if float64(r.Total) - rounded < 0.001 {
    log.Debug("round total points: 50")
    return 50
  }
  
  log.Debug("round total points: 0")
  return 0
}

// if the total is a multiple of 0.25, 25 points are awarded
func modQuarterPoints(r *model.Receipt) int {
  cents := int(float64(r.Total) * 100)
  if cents % 25 == 0 {
    log.Debug("quarter total points: 25")
    return 25
  }

  log.Debug("quarter total points: 0")
  return 0
}

// for every two items, 5 points are awarded
func pairsPoints(r *model.Receipt) int {
  numItems := len(r.Items)
  points := 5 * (numItems / 2)
  log.Debug("pairs points: ", points)
  return points
}

// for every item with a trimmed whitespace character count that
// is divisible by three, add 0.2x the price of the item, rounded
// up to the nearest whol point
func itemDescriptionLengthPoints(r *model.Receipt) int {
  totalPoints := 0

  for _, item := range r.Items {
    trimmed := strings.TrimSpace(item.ShortDescription)
    if len(trimmed) % 3 == 0 {
      totalPoints += int(math.Ceil(float64(item.Price) * 0.2))
    }
  }

  log.Debug("item desc points: ", totalPoints)
  return totalPoints
}

// if the day of te purchase was odd, add 6 points
func oddDayPoints(r *model.Receipt) int {
  if time.Time(r.PurchaseDate).Day() % 2 == 1 {
    log.Debug("day points: 6")
    return 6
  }

  log.Debug("day points: 0")
  return 0
}

// if the purchase was between 14:00 and 16:00, add 10 points
func purchaseTimePoints(r *model.Receipt) int {
  hour := time.Time(r.PurchaseTime).Hour() 

  if 14 <= hour && hour <= 15 {
    log.Debug("hour points: 10")
    return 10
  }

  log.Debug("hour points: 0")
  return 0
}

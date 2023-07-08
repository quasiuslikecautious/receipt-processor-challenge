package model

import (
  "encoding/json"
  "errors"
  "regexp"
  "strconv"
  "time"

  "github.com/gin-gonic/gin/binding"
  "github.com/google/uuid"
  "github.com/go-playground/validator/v10"
)

type MemoryDB map[uuid.UUID]Receipt

type Item struct {
  ShortDescription  string  `json:"shortDescription" binding:"required,item_desc"`
  Price             PriceT  `json:"price" binding:"required"`
}

type Receipt struct {
  Retailer      string    `json:"retailer" binding:"required,retailer"`
  PurchaseDate  DateOnly  `json:"purchaseDate" binding:"required"`
  PurchaseTime  TimeOnly  `json:"purchaseTime" binding:"required"`
  Items         []Item    `json:"items" binding:"required,dive"`
  Total         PriceT    `json:"total" binding:"required"`
}

func RegisterModelValidators() error {
  v, ok := binding.Validator.Engine().(*validator.Validate)
  if !ok {
    return errors.New("failed to register validators")
  }

  v.RegisterValidation("retailer", retailerValidator)
  v.RegisterValidation("item_desc", itemDescriptionValidator)

  return nil
}

// ------------------------------------------------------
//      FIELD FORMAT VALIDATION
// ------------------------------------------------------

// Validate validates conditions that rely on multiple fields within the receipt
func (r *Receipt) Validate() error {
  pd := time.Time(r.PurchaseDate)
  pt := time.Time(r.PurchaseTime)

  fullPurchaseDateTime := time.Date(pd.Year(), pd.Month(), pd.Day(), pt.Hour(), pt.Minute(), 59, 99, time.UTC)

  now := time.Now()

  if fullPurchaseDateTime.After(now) {
    return errors.New("purchase datetime after current datetime")
  }

  totalPrice := 0.00
  for _, item := range r.Items {
    totalPrice += float64(item.Price)
  }

  if totalPrice - float64(r.Total) > 0.001 {
    return errors.New("total price does not match total price of all items")
  }

  return nil
}

var retailerValidator validator.Func = func(fl validator.FieldLevel) bool {
  retailer := fl.Field().String()
  // No whitespace
  return regexp.MustCompile(`^\S+$`).MatchString(retailer)
}

var itemDescriptionValidator validator.Func = func(fl validator.FieldLevel) bool {
  itemDesc := fl.Field().String()
  // Can contain alphanumeric, '_' and '-'
  return regexp.MustCompile(`^[\w\s\-]+$`).MatchString(itemDesc)
}
type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(data []byte) error {
  var input string
  if err := json.Unmarshal(data, &input); err != nil {
    return err
  }

  format := "2006-01-02"
  parsedDate, err := time.Parse(format, input)
  if err != nil {
    return err
  }

  *d = DateOnly(parsedDate)
  return nil
}

type TimeOnly time.Time

func (t* TimeOnly) UnmarshalJSON(data []byte) error {
  var input string
  if err := json.Unmarshal(data, &input); err != nil {
    return err
  }
  
  format := "15:04"
  parsedTime, err := time.Parse(format, input)
  if err != nil {
    return err
  }

  *t = TimeOnly(parsedTime)
  return nil
}

type PriceT float64

func (p* PriceT) UnmarshalJSON(data []byte) error {
  var input string
  if err := json.Unmarshal(data, &input); err != nil {
    return err
  }

  // any number of digits, followed by with exactly two digits after a '.'
  match, err := regexp.MatchString(`^\d+\.\d{2}$`, input)
  if err != nil {
    return err
  }

  if !match {
    return errors.New("invalid price format")
  }

  value, err := strconv.ParseFloat(input, 64)
  if err != nil {
    return err
  }

  if value < 0.00 {
    return errors.New("negative price not allowed")
  }

  *p = PriceT(value)
  return nil
}

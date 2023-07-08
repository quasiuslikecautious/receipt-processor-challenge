package app

import (
  "os"

  "github.com/gin-gonic/gin"
  "github.com/google/uuid"
  log "github.com/sirupsen/logrus"
  "quasiuslikecautious/receipt-processor-challenge/app/handler"
  "quasiuslikecautious/receipt-processor-challenge/app/model"
  "quasiuslikecautious/receipt-processor-challenge/config"
)

type App struct {
  Config *config.Config
  Db     model.MemoryDB
  Router *gin.Engine
}

// Initialize initializes the app db and router with the config
func (a *App) Initialize(config *config.Config) {
  a.Config = config
  a.Db = make(map[uuid.UUID]model.Receipt)
  a.Router = gin.Default()

  if a.Config.Debug {
    log.SetLevel(log.DebugLevel)
    log.SetOutput(os.Stdout)
  }

  if err := model.RegisterModelValidators(); err != nil {
    panic(err)
  }

  a.Post("/receipts/process", handler.PostReceipt)
  a.Get("/receipts/:id/points", handler.GetReceiptPoints)
}

// Get wraps the router for GET methods, adding a db parameter 
func (a *App) Get(path string, f func(c *gin.Context, db *model.MemoryDB)) {
  a.Router.GET(path, func(c *gin.Context) {
    f(c, &a.Db)
  })
}

// Post wraps the router for POST methods, adding a db parameter 
func (a *App) Post(path string, f func(c *gin.Context, db *model.MemoryDB)) {
  a.Router.POST(path, func(c *gin.Context) {
    f(c, &a.Db)
  })
}

// Run the app on its router
func (a *App) Run() {
  address := a.Config.Host + ":" + a.Config.Port
  a.Router.Run(address)
}

package api

import (
	"log"
	"time"

	"github.com/Mohsen20031203/blockchain-insight/config"
	_ "github.com/Mohsen20031203/blockchain-insight/docs"
	"github.com/Mohsen20031203/blockchain-insight/internal/enthblock"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	client *ethclient.Client
	config config.Config
	router *gin.Engine
	cach   *cache.Cache
}

const LastBlock = "last_block"

func NewServer(config config.Config) *Server {

	client, err := enthblock.NewClient(config.RPCURL)
	if err != nil {
		log.Fatal(err)
	}

	cach := cache.New(cache.NoExpiration, 1*time.Hour)
	server := &Server{
		client: client,
		config: config,
		cach:   cach,
	}

	server.setupRouter()
	return server

}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/balace/:address", s.GetAddressBalance)

	router.GET("/block/:id", s.GetBlockById)
	router.Use(s.Cache())
	router.GET("/last/block", s.GetLastBlock)
	s.router = router
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

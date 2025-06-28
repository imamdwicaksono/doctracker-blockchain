package main

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/routes"
	"doc-tracker/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	fmt.Println("[Init] Starting Doc-Tracker Node...")

	blockchain.LoadChainFromStorage()
	fmt.Println("[Blockchain] Chain loaded")

	mempool.LoadFromFile()
	fmt.Println("[Mempool] Mempool loaded")

	services.StartMinerWorker()
	fmt.Println("[Miner] Worker started")

	services.StartSyncWorker()
	fmt.Println("[Sync] Worker started")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	routes.P2PRoutes(app)
	routes.SyncRoutes(app)

	routes.MinerRoutes(app)

	api := app.Group("/api")
	routes.TrackerRoutes(api)
	routes.SetupAuthRoutes(api)
	routes.RegisterDecryptRoutes(api)
	routes.RegisterEvidenceRoutes(api)
	routes.RegisterCheckpointRoutes(api)
	routes.BlockRoutes(api)

	fmt.Println("[Server] Listening on :3001")

	err := app.Listen(":3001")
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

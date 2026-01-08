package main

import (
	"context"
	"doc-tracker/blockchain"
	"doc-tracker/blockchain/adapter"
	"doc-tracker/grpc"
	"doc-tracker/mempool"
	"doc-tracker/middlewares"
	"doc-tracker/routes"
	"doc-tracker/services"
	"doc-tracker/storage"
	"doc-tracker/storage/redis"
	"doc-tracker/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "doc-tracker/docs"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// ===================== MAIN =====================
func main() {

	// ===== ENV =====
	loadEnv()

	// ===== CRYPTO KEYS =====
	fmt.Println("✅ Checking and creating ECDSA keys...")
	utils.CreatePemIfNotExists("data/private.pem")
	utils.CreatePemIfNotExists("data/public.pem")

	// ===== REDIS =====
	redis.InitRedis()
	fmt.Println("[Redis] Initialized")

	// ===== DATABASE =====
	storage.InitDB()
	fmt.Println("[DB] Connected")

	// ===== BLOCKCHAIN CORE (LEGACY) =====
	blockchain.InitChain()
	mempoolInit()

	// ===== BLOCKCHAIN LISTENER (EVM) =====
	startBlockchainListener()

	// ===== WORKERS =====
	services.StartMinerWorker()
	services.StartSyncWorker()
	fmt.Println("[Workers] Miner & Sync started")

	// ===== GRPC =====
	killProcessOnPort(3003)
	go grpc.StartGRPCServer("3003")
	fmt.Println("[GRPC] Server started on :3003")

	// ===== HTTP API =====
	startHTTPServer()
}

// ===================== BLOCKCHAIN LISTENER =====================
func startBlockchainListener() {

	client, err := adapter.NewClient()
	if err != nil {
		log.Fatal("[Blockchain] client error:", err)
	}

	contractAddr := common.HexToAddress(
		os.Getenv("AUDIT_CONTRACT_ADDRESS"),
	)

	auditContract, err := adapter.NewDocumentAudit(
		contractAddr,
		client.Eth,
	)
	if err != nil {
		log.Fatal("[Blockchain] contract bind error:", err)
	}

	listener := &adapter.Listener{
		Client:   client,
		Contract: auditContract,
		Address:  contractAddr,
	}

	ctx := context.Background()

	// 1️⃣ catch-up dari last processed block
	if err := listener.SyncFromLastBlock(ctx, storage.DB); err != nil {
		log.Fatal("[Blockchain] sync failed:", err)
	}

	// 2️⃣ realtime listener
	go func() {
		listener.StartPolling(ctx, storage.DB)
	}()

	fmt.Println("[Blockchain] Event listener started (resume-safe)")
}

// ===================== MEMPOOL INIT =====================
func mempoolInit() {

	if _, err := mempool.InitKeys(); err != nil {
		fmt.Printf("Failed init keys: %v\n", err)
	}

	if err := mempool.InitEncryptMempool(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	if err := mempool.LoadFromFile(); err != nil {
		fmt.Printf("Failed load mempool: %v\n", err)
	}

	mempool.RemoveDuplicateEntries()
	blockchain.RemoveDuplicateBlocks()

	fmt.Println("[Mempool] Loaded & cleaned")
}

// ===================== HTTP SERVER =====================
func startHTTPServer() {

	app := fiber.New()

	if os.Getenv("ALLOWED_ORIGIN") != "" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     os.Getenv("ALLOWED_ORIGIN"),
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
			AllowCredentials: true,
			MaxAge:           12 * 3600,
		}))
	}

	app.Use(limiter.New(limiter.Config{Max: 100, Expiration: time.Minute}))

	app.Use(func(c *fiber.Ctx) error {
		fmt.Printf("👉 [%s] %s\n", c.Method(), c.Path())
		return c.Next()
	})

	// ROUTES
	routes.P2PRoutes(app)
	routes.SyncRoutes(app)
	routes.MinerRoutes(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	routes.SetupAuthRoutes(api)

	protected := app.Group("/api", middlewares.JWTMiddleware)
	routes.TrackerRoutes(protected)
	routes.SetupAuthProtectedRoutes(protected)
	routes.RegisterDecryptRoutes(protected)
	routes.RegisterEvidenceRoutes(protected)
	routes.RegisterCheckpointRoutes(protected)
	routes.BlockRoutes(protected)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}

	portInt, _ := strconv.Atoi(port)
	killProcessOnPort(portInt)

	fmt.Println("[HTTP] Listening on :", port)
	log.Fatal(app.Listen(":" + port))
}

// ===================== ENV =====================
func loadEnv() {

	wd, _ := os.Getwd()
	envPath := filepath.Join(wd, ".env")

	if _, err := os.Stat(envPath); err == nil {
		_ = godotenv.Overload(envPath)
		fmt.Println("✅ .env loaded")
	}
}

// ===================== KILL PORT =====================
func killProcessOnPort(port int) error {

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		findCmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t")
		out, _ := findCmd.Output()
		pid := strings.TrimSpace(string(out))
		if pid != "" {
			cmd = exec.Command("kill", "-9", pid)
		}
	default:
		return nil
	}

	if cmd != nil {
		return cmd.Run()
	}
	return nil
}

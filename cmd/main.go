package main

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/middlewares"
	"doc-tracker/routes"
	"doc-tracker/services"
	"doc-tracker/storage/redis"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "doc-tracker/docs"

	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// @title           Document Tracker API
// @version         1.0
// @description     REST API for internal document blockchain tracking system
// @host            localhost:8080
// @BasePath        /
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	fmt.Println("[Init] Starting Doc-Tracker Node...")

	redis.InitRedis()
	fmt.Println("[Redis] Redis initialized")

	blockchain.InitChain()
	fmt.Println("[Blockchain] Chain loaded")

	// Inisialisasi kunci
	if _, err := mempool.InitKeys(); err != nil {
		fmt.Printf("Failed to initialize keys: %v", err)
	}
	// Inisialisasi mempool
	if err := mempool.InitEncryptMempool(); err != nil {
		fmt.Printf("Warning: %v", err)
	}
	// Load data dari file
	if err := mempool.LoadFromFile(); err != nil {
		fmt.Printf("Failed to load mempool: %v", err)
	}
	fmt.Println("[Mempool] Mempool loaded")

	services.StartMinerWorker()
	fmt.Println("[Miner] Worker started")

	services.StartSyncWorker()
	fmt.Println("[Sync] Worker started")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "*",
	}))

	routes.P2PRoutes(app)
	routes.SyncRoutes(app)
	routes.MinerRoutes(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use(limiter.New(limiter.Config{Max: 100, Expiration: time.Minute}))
	// api
	HandlerApiRoute(app)

	HandlerWebRoute(app)
	//api protected
	HandlerApiProtectedRoute(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
		portInt, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("Invalid port:", port)
			return
		}
		killProcessOnPort(portInt)
		err = app.Listen(":" + port)
		fmt.Println("[Server] Listening on :", port)
		if err != nil {
			fmt.Println("Server error:", err)

			return
		}
	}

	app.Use(func(c *fiber.Ctx) error {
		fmt.Println("Request Origin:", c.Get("Origin"))
		return c.Next()
	})
}

func HandlerWebRoute(app *fiber.App) {
	protectedWeb := app.Group("", middlewares.JWTMiddleware)
	routes.RegisterEvidenceRoutesWeb(protectedWeb)
}

func HandlerApiProtectedRoute(app *fiber.App) {
	// api protected
	protected := app.Group("/api", middlewares.JWTMiddleware)
	routes.TrackerRoutes(protected)
	routes.SetupAuthProtectedRoutes(protected)
	routes.RegisterDecryptRoutes(protected)
	routes.RegisterEvidenceRoutes(protected)
	routes.RegisterCheckpointRoutes(protected)
	routes.BlockRoutes(protected)

}

func HandlerApiRoute(app *fiber.App) {
	api := app.Group("/api")
	routes.SetupAuthRoutes(api)
}

func killProcessOnPort(port int) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Find PID using netstat and kill with taskkill
		findCmd := exec.Command("netstat", "-ano")
		findOut, err := findCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to find process: %v", err)
		}

		lines := strings.Split(string(findOut), "\n")
		for _, line := range lines {
			if strings.Contains(line, fmt.Sprintf(":%d", port)) {
				parts := strings.Fields(line)
				if len(parts) > 4 {
					pid := parts[len(parts)-1]
					cmd = exec.Command("taskkill", "/F", "/PID", pid)
					break
				}
			}
		}
	case "darwin", "linux", "freebsd", "openbsd":
		// Find PID using lsof and kill with kill
		findCmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t")
		findOut, err := findCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to find process: %v", err)
		}

		pid := strings.TrimSpace(string(findOut))
		if pid != "" {
			cmd = exec.Command("kill", "-9", pid)
		}
	default:
		return fmt.Errorf("unsupported platform")
	}

	if cmd == nil {
		return fmt.Errorf("no process found on port %d", port)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

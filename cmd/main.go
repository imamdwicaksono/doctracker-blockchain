package main

import (
	"doc-tracker/blockchain"
	"doc-tracker/grpc"
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

	mempool.RemoveDuplicateEntries()
	fmt.Println("[Mempool] Duplicate entries removed")
	blockchain.RemoveDuplicateBlocks()
	fmt.Println("[Blockchain] Duplicate blocks removed")

	services.StartMinerWorker()
	fmt.Println("[Miner] Worker started")

	services.StartSyncWorker()
	fmt.Println("[Sync] Worker started")

	go grpc.StartGRPCServer("3003")
	fmt.Println("[GRPC] Server started on port 3003")

	app := fiber.New()

	// allowedOrigins := ""
	// if os.Getenv("ENV") == "development" {
	// 	allowedOrigins = "http://172.24.4.25:3000,http://localhost:3000"
	// } else {
	// 	allowedOrigins = "https://production.com"
	// }

	// app.Use(cors.New(cors.Config{
	// 	AllowOriginsFunc: func(origin string) bool {
	// 		return origin == "http://docutrack.test" ||
	// 			origin == "http://172.24.4.25:3000" ||
	// 			origin == "http://localhost:3000" ||
	// 			origin == "http://blockchain.mmsgroup.co.id:3000"
	// 	},
	// 	AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	// 	AllowHeaders:     "Origin, Content-Type, Accept, Authorization", // ⚠️ jangan tambahkan `credentials`
	// 	ExposeHeaders:    "Content-Length",
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * 3600,
	// }))

	app.Use(func(c *fiber.Ctx) error {
		fmt.Printf("👉 [%s] %s from %s\n", c.Method(), c.Path(), c.Get("Origin"))
		return c.Next()
	})

	app.Use(limiter.New(limiter.Config{Max: 100, Expiration: time.Minute}))

	routes.P2PRoutes(app)
	routes.SyncRoutes(app)
	routes.MinerRoutes(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

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

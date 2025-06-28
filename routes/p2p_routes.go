package routes

import (
	"doc-tracker/p2p"

	"github.com/gofiber/fiber/v2"
)

func P2PRoutes(app *fiber.App) {
	app.Get("/p2p/latest-block", p2p.GetLatestBlock)
	app.Get("/p2p/mempool", p2p.GetMempool)
	app.Post("/p2p/block", p2p.ReceiveBlock)
	app.Post("/p2p/mempool", p2p.ReceiveMempool)
}

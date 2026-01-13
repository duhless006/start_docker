package main

import (
	httppack "copr/httpPack"
	"copr/worker"
	"fmt"
)

func main() {
	// ctx := context.Background()

	// conn, err := database.ConnectDataBase(ctx)
	// if err != nil {
	// 	fmt.Println("error connect database", err)
	// }

	// if err := conn.Ping(ctx); err != nil {
	// 	fmt.Println("no signal", err)
	// 	return
	// }
	// fmt.Println("database connect complete")

	fmt.Println("server start")

	worker := worker.NewWorkers()
	httpHandlers := httppack.NewHTTPHandlers(worker)
	httpServer := httppack.NewHTTPServer(*httpHandlers)

	if err := httpServer.ConnectServer(); err != nil {
		fmt.Println("failed to start http server:", err)
	}
}

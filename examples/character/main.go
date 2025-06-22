package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 命令行參數
	var (
		port   = flag.String("port", "8080", "Server port")
		dbPath = flag.String("db", "characters.db", "SQLite database path")
	)
	flag.Parse()

	// 建立伺服器
	server, err := NewServer(*dbPath, *port)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	// 處理優雅關閉
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down server...")
		server.Close()
		os.Exit(0)
	}()

	// 啟動伺服器
	log.Fatal(server.Start())
}

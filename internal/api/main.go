package main

import (
    "github.com/gin-gonic/gin"
    "github.com/rahulguha/promptly/api"
    "github.com/rahulguha/promptly/routes"
    "github.com/rahulguha/promptly/storage"
)

func main() {
    store := storage.NewFileStorage("data/prompts.json")
    handler := &api.Handler{Store: store}

    r := gin.Default()
    routes.RegisterRoutes(r, handler)

    r.Run(":8080")
}

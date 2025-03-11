package main

import (
    "log"
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "todo-app/backend/internal/handlers"
    "todo-app/backend/internal/redis"
)

func main() {
    // inicializar cliente Redis
    redisClient := redis.NewRedisClient()
    defer redisClient.Close()

    r := mux.NewRouter()
    
    // iunicializar handler
    todoHandler := handlers.NewTodoHandler(redisClient)

    r.HandleFunc("/api/todos", todoHandler.GetTodos).Methods("GET")
    r.HandleFunc("/api/todos", todoHandler.CreateTodo).Methods("POST")
    r.HandleFunc("/ws", todoHandler.HandleWebSocket)
    r.HandleFunc("/api/todos/{id}/toggle", todoHandler.ToggleTodo).Methods("PUT")

    // configurar CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Content-Type"},
    })

    // aplicar middleware CORS ao router
    handler := c.Handler(r)

    log.Println("Servidor iniciado na porta 8080")
    if err := http.ListenAndServe(":8080", handler); err != nil {
        log.Fatal(err)
    }
}
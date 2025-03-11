package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "time"
    "github.com/google/uuid"
    "github.com/gorilla/websocket"
    "github.com/go-redis/redis/v8"
    "todo-app/backend/internal/models"
    "github.com/gorilla/mux"
)

type TodoHandler struct {
    rdb      *redis.Client
    upgrader websocket.Upgrader
}

func NewTodoHandler(rdb *redis.Client) *TodoHandler {
    return &TodoHandler{
        rdb: rdb,
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                return true 
            },
        },
    }
}

func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    
    todos := []models.Todo{}
    keys, err := h.rdb.Keys(ctx, "todo:*").Result()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    for _, key := range keys {
        todoJSON, err := h.rdb.Get(ctx, key).Result()
        if err != nil {
            continue
        }

        var todo models.Todo
        if err := json.Unmarshal([]byte(todoJSON), &todo); err != nil {
            continue
        }
        todos = append(todos, todo)
    }

    json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    
    var todoCreate models.TodoCreate
    if err := json.NewDecoder(r.Body).Decode(&todoCreate); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    todo := models.Todo{
        ID:        uuid.New().String(),
        Title:     todoCreate.Title,
        Completed: false,
        CreatedAt: time.Now(),
    }

    todoJSON, _ := json.Marshal(todo)
    err := h.rdb.Set(ctx, "todo:"+todo.ID, todoJSON, 0).Err()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    h.rdb.Publish(ctx, "todos", todoJSON)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := h.upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer conn.Close()

    // Inscrever no canal de todos
    pubsub := h.rdb.Subscribe(context.Background(), "todos")
    defer pubsub.Close()

    // Enviar mensagens do Redis para o WebSocket
    for {
        msg, err := pubsub.ReceiveMessage(context.Background())
        if err != nil {
            break
        }
        
        if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
            break
        }
    }
}

func (h *TodoHandler) ToggleTodo(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    vars := mux.Vars(r)
    todoID := vars["id"]

    // Buscar todo atual
    todoJSON, err := h.rdb.Get(ctx, "todo:"+todoID).Result()
    if err != nil {
        http.Error(w, "Todo não encontrada", http.StatusNotFound)
        return
    }

    var todo models.Todo
    if err := json.Unmarshal([]byte(todoJSON), &todo); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Toggle completed
    todo.Completed = !todo.Completed

    // Salvar atualização
    updatedJSON, _ := json.Marshal(todo)
    err = h.rdb.Set(ctx, "todo:"+todoID, updatedJSON, 0).Err()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Publicar evento de atualização
    h.rdb.Publish(ctx, "todos", updatedJSON)

    json.NewEncoder(w).Encode(todo)
}
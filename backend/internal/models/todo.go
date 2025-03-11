package models

import (
    "time"
)

type Todo struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
}

type TodoCreate struct {
    Title string `json:"title"`
}
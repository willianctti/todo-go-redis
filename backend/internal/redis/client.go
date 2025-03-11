package redis

import (
    "context"
    "github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", 
        DB:       0,
    })

    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        panic("Não foi possível conectar ao Redis: " + err.Error())
    }

    return client
}
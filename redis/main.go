package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	//连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()

	//定义有序集合的键名
	key := "myZSet"

	//添加元素到有序集合，使用分数表示条件范围
	scores := []int{0, 10, 100, 200, 300, 400, 0, 10, 20}
	members := []string{"A1", "A2", "A3", "A4", "A5", "B1", "B2", "B3"}
	for i := 0; i < len(scores); i++ {
		err := rdb.ZAdd(ctx, key, redis.Z{Member: members[i], Score: float64(scores[i])}).Err()
		if err != nil {
			fmt.Println("Failed to add members to sorted set:", err)
			return
		}
	}

	//创建管道并添加命令
	pipe := rdb.Pipeline()

	// 添加条件查询命令到管道
	_, _ = pipe.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: "10", Max: "100"}).Result()

	_, _ = pipe.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: "10", Max: "20"}).Result()

	_, _ = pipe.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: "20"}).Result()

	// 执行管道中的命令并获取结果
	results, err := pipe.Exec(ctx)
	if err != nil {
		fmt.Println("Failed to execute pipeline:", err)
		return
	}

	// 处理结果，检查满足条件的情况并输出对应的结果
	intersectCount := len(results[0].String()) // 第一个查询的结果作为交集处理后的结果集长度
	fmt.Printf("Intersection count: %d\n", intersectCount)

}

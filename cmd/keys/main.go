package main

import "github.com/PoolHealth/PoolHealthServer/internal/repo/keys"

func main() {
	print(keys.NewBuilder().UsersPools())
}

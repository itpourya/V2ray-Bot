package cache

func NewCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "haze_cache:6379",
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
		DB:       0, // use default DB
	})

	return rdb
}

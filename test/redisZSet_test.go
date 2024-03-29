package test

//
//func Test_redisZSet(t *testing.T) {
//
//	client := NewClient("default")
//	defer func() {
//		_, _ = client.Del("key_set_z")
//	}()
//	//测试
//	add, err := client.ZSet.SetAdd("key_set_z", &redisZ{1, "小猫"}, &redisZ{2, "小狗"}, &redisZ{3, "小鸟"})
//	flog.Info("添加返回结果：", add, err)
//
//	score, err := client.ZSet.Score("key_set_z", "小狗")
//	flog.Info("返回指定成员的score:", score, err)
//
//	strings, err := client.ZSet.Range("key_set_z", 0, 1)
//	flog.Info("获取所有集合：", strings, err)
//
//	revRange, err := client.ZSet.RevRange("key_set_z", 0, 3)
//	flog.Info("有序集合指定区间内的成员分数从高到低：", revRange, err)
//
//	byScore, err := client.ZSet.RangeByScore("key_set_z", &redisZRangeBy{Max: "3", Min: "1", Offset: 1, Count: 3})
//	flog.Info("获取指定分数区间的成员列表：", byScore, err)
//
//}

//func Test_redisZSet_Add(t *testing.T) {
//
//	client := redis.NewClient("default")
//	defer client.Del("key_set_z")
//	add, _ := client.ZSet.SetAdd("key_set_z", &redis.redisZ{1, "小猫"}, &redis.redisZ{2, "小狗"}, &redis.redisZ{3, "小鸟"})
//	assert.Equal(t, add, true)
//}
//
//func Test_redisZSet_Range(t *testing.T) {
//
//	client := redis.NewClient("default")
//	defer client.Del("key_set_z")
//	client.ZSet.SetAdd("key_set_z", &redis.redisZ{1, "小猫"}, &redis.redisZ{2, "小狗"}, &redis.redisZ{3, "小鸟"})
//	strings, _ := client.ZSet.Range("key_set_z", 0, 1)
//	assert.Equal(t, strings[0], "小猫")
//}
//
//func Test_redisZSet_RangeByScore(t *testing.T) {
//
//	client := redis.NewClient("default")
//	defer client.Del("key_set_z")
//	client.ZSet.SetAdd("key_set_z", &redis.redisZ{1, "小猫"}, &redis.redisZ{2, "小狗"}, &redis.redisZ{3, "小鸟"})
//	byScore, _ := client.ZSet.RangeByScore("key_set_z", &redis.redisZRangeBy{Max: "3", Min: "1", Offset: 1, Count: 3})
//	assert.Equal(t, byScore[0], "小狗")
//}
//
//func Test_redisZSet_RevRange(t *testing.T) {
//
//	client := redis.NewClient("default")
//	defer client.Del("key_set_z")
//	client.ZSet.SetAdd("key_set_z", &redis.redisZ{1, "小猫"}, &redis.redisZ{2, "小狗"}, &redis.redisZ{3, "小鸟"})
//	revRange, _ := client.ZSet.RevRange("key_set_z", 0, 3)
//	assert.Equal(t, revRange[0], "小鸟")
//}
//
//func Test_redisZSet_Score(t *testing.T) {
//
//	client := redis.NewClient("default")
//	defer client.Del("key_set_z")
//	client.ZSet.SetAdd("key_set_z", &redis.redisZ{1, "小猫"}, &redis.redisZ{2, "小狗"}, &redis.redisZ{3, "小鸟"})
//	score, _ := client.ZSet.Score("key_set_z", "小狗")
//	assert.Equal(t, score, float64(2))
//}

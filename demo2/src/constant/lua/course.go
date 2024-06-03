package lua

// 1. 是否存在时间冲突
// 2. 用户是否已经选择了

const (
	CourseSelectOK = iota
	CourseSelected
	CourseFull
	/*
		CourseSelectLuaScript lua
			key1 = 用户key
			key2 = 课程id
			key3 = 课程key
			key4 = capacity key
	*/
	CourseSelectLuaScript = `
	-- 1. 用户是否已经选择了
	if redis.call("sismember", KEYS[1], KEYS[2]) == 1 then
		return 1 
	end	
	-- 2. 选课操作
	local count =tonumber(redis.call("hget", KEYS[3], KEYS[4]))
	if count and count > 0 then
		-- 课程人数减 一
		redis.call("hincrby", KEYS[3], KEYS[4], -1)
		-- 选课，添加到用户集合
		redis.call("sadd", KEYS[1], KEYS[2])
		return 0
	else
		-- 容量满了
		return 2
	end
`
	CourseBackOK      = 0
	CourseNotSelected = 1
	/*
		CourseBackLuaScript
		key1 = 用户key
		key2 = 课程id
		key3 = 课程key
		key4 = capacity key

	*/

	CourseBackLuaScript = `
	-- 1. 用户是否已经选择了
	if redis.call("sismember", KEYS[1], KEYS[2]) == 1 then
		-- 课程人数加一
		redis.call("hincrby", KEYS[3], KEYS[4], 1)
		-- 退课，从用户集合中删除
		redis.call("srem", KEYS[1], KEYS[2])
		return 0
	else
		return 1
	end
`
)

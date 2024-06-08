package lua

// 1. 是否存在时间冲突
// 2. 用户是否已经选择了

const (
	CourseSelectOK = iota
	CourseSelected
	CourseTimeConflict
	CourseFull
	/*
		CourseSelectLuaScript lua
			key1 = 用户key
			key2 = 课程id
			key3 = 课程key
			key4 = capacity key
			key5 = user course schedule key
			key6 = offset
	*/
	CourseSelectLuaScript = `

	-- 1. 用户是否已经选择了
	if redis.call("sismember", KEYS[1], KEYS[2]) == 1 then
		return 1 
	end	

    -- 2. 是否存在时间冲突 判断某个时间段是否为1
    local bitmap = tonumber(redis.call("getbit", KEYS[5], KEYS[6]))
    if bitmap and bitmap == 1 then
        return 2
    end
	
	-- 3. 选课操作
	local count =tonumber(redis.call("hget", KEYS[3], KEYS[4]))
	if count and count > 0 then
		-- 课程人数减 一
		redis.call("hincrby", KEYS[3], KEYS[4], -1)
		-- 选课，添加到用户集合
		redis.call("sadd", KEYS[1], KEYS[2])
		-- 课程时间段设置为1
		redis.call("setbit", KEYS[5], KEYS[6], 1) 
		return 0
	else
		-- 容量满了
		return 3
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
		key5 = user course schedule key
		key6 = offset
	*/

	CourseBackLuaScript = `
	-- 1. 用户是否已经选择了
	if redis.call("sismember", KEYS[1], KEYS[2]) == 1 then
		-- 课程人数加一
		redis.call("hincrby", KEYS[3], KEYS[4], 1)
		-- 退课，从用户集合中删除
		redis.call("srem", KEYS[1], KEYS[2])
		-- 课程时间段设置为0
		redis.call("setbit", KEYS[5], KEYS[6], 0)
		return 0
	else
		return 1
	end
`
)

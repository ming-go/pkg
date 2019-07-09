package ratelimiting

const (
	LUA_SCRIPT_LEAKY_BUCKET = `
		local ret = redis.call("INCRBY", KEYS[1], "1")
		if ret == 1 then
			redis.call("PEXPIRE", KEYS[1], KEYS[2])
		end
		return ret
	`
)

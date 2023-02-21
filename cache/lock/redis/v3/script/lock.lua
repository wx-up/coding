local val = redis.call("get", KEYS[1])
if not val then
    -- key 不存在的时候
    return redis.call("set", KEYS[1], ARGV[1],"EX",ARGV[2])
elseif val == ARGV[1] then
    -- 重新刷新有效期
    redis.call("expire",KEYS[1],ARGV[2] )
    return "OK"
else
    -- 别人持有锁
    return ""
end
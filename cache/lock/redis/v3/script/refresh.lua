-- 检查是不是自己的锁，是的话就续约
if redis.call("get",KEYS[1]) == ARGV[1]
then
    return redis.call("expire", KEYS[1], ARGV[2])
else
    return 0
end
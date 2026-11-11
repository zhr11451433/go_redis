1.keys命令

    DEL key
    该命令用于在 key 存在时删除 key。

    EXISTS key
    检查给定 key 是否存在。

    EXPIRE key seconds
    为给定 key 设置过期时间，以秒计。

    PERSIST key
    移除 key 的过期时间，key 将持久保持

    TTL key
    以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。

    RANDOMKEY
    从当前数据库中随机返回一个 key 。

    RENAME key newkey
    修改 key 的名称

    TYPE key
    返回 key 所储存的值的类型
2.字符串命令
    
    SET key value
    设置指定 key 的值

    GET key
    获取指定 key 的值。
    
    GETSET key value
    将给定 key 的值设为 value ，并返回 key 的旧值(old value)。

    MGET key1 [key2..]
    获取所有(一个或多个)给定 key 的值。

    SETNX key value
    只有在 key 不存在时设置 key 的值。

    STRLEN key
    返回 key 所储存的字符串值的长度。

    MSET key value [key value ...]
    同时设置一个或多个 key-value 对。

    INCR key
    将 key 中储存的数字值增一。
    
    INCRBY key increment
    将 key 所储存的值加上给定的增量值（increment） 。

    DECR key
    将 key 中储存的数字值减一。

    DECRBY key decrement
    key 所储存的值减去给定的减量值（decrement） 。

    APPEND key value
    如果 key 已经存在并且是一个字符串， APPEND 命令将指定的 value 追加到该 key 原来值（value）的末尾。
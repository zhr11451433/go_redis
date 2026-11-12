package main

import (
	"bufio"
	"fmt"
	"go_redis/db"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const dataPath = "./redis.json"

func main() {
	redis, err := db.Load(dataPath)
	if err != nil {
		log.Fatalf("加载数据库失败：%v", err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		//strings.Fields() 会按照连续空白字符拆分字符串
		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}
		command := strings.ToUpper(args[0])
		switch command {
		case "FLUSHDB":
			if len(args) != 1 {
				fmt.Printf("用法：%s\n", command)
				continue
			}
			if err := db.FlushDB(dataPath, redis); err != nil {
				log.Printf("清空数据库失败：%v", err)
				continue
			} else {
				fmt.Println("清空数据库成功")
			}
		case "STRLEN":
			if len(args) != 2 {
				fmt.Printf("用法：%s key\n", command)
				continue
			}
			fmt.Println(redis.StrLen(args[1]))
		case "MSET":
			if len(args)%2 == 0 || len(args) == 1 {
				fmt.Printf("用法：%s key value [key value]\n", command)
				continue
			}
			pair := make([]db.Pair, 0, (len(args)-1)/2)
			for i := 1; i < len(args); i += 2 {
				pair = append(pair, db.Pair{
					Key:   args[i],
					Value: args[i+1],
				})
			}
			redis.MSet(pair...)
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println("OK")
		case "INCRBY":
			if len(args) != 3 {
				fmt.Printf("用法：%s key increment\n", command)
				continue
			}
			value, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("必须为整数")
				continue
			}
			n, err := redis.IncrBy(args[1], value)
			if err != nil {
				log.Println("key中值非整数")
				continue
			}
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println(n)
		case "APPEND":
			if len(args) != 3 {
				fmt.Printf("用法：%s key value\n", command)
				continue
			}
			n, err := redis.Append(args[1], args[2])
			if err != nil {
				log.Println("key不是string类型")
				continue
			}
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println(n)
		case "DECRBY":
			if len(args) != 3 {
				fmt.Printf("用法：%s key increment\n", command)
				continue
			}
			value, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Println("必须为整数")
				continue
			}
			n, err := redis.DecrBY(args[1], value)
			if err != nil {
				log.Println("key中值非整数")
				continue
			}
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println(n)
		case "INCR":
			if len(args) != 2 {
				fmt.Printf("用法：%s key\n", command)
				continue
			}
			n, err := redis.Incr(args[1])
			if err != nil {
				log.Println("key中值非整数")
				continue
			}
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println(n)
		case "DECR":
			if len(args) != 2 {
				fmt.Printf("用法：%s key\n", command)
				continue
			}
			n, err := redis.Decr(args[1])
			if err != nil {
				log.Println("key中值非整数")
				continue
			}
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println(n)
		case "SETNX":
			if len(args) != 3 {
				fmt.Printf("用法：%s key value\n", command)
				continue
			}
			n := redis.SetNx(args[1], args[2])
			if n == 0 {
				fmt.Println(0)
				continue
			}
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println(n)
		case "MGET":
			if len(args) < 2 {
				fmt.Printf("用法：%s key [key...]\n", command)
				continue
			}
			keys := args[1:]
			value := redis.MGet(keys...)
			for _, v := range value {
				fmt.Println(v)
			}
		case "DEL":
			if len(args) != 2 {
				fmt.Printf("用法：%s <键>\n", command)
				continue
			}
			if n := redis.Del(args[1]); n == 0 {
				log.Println("键不存在")
			} else {
				if err := db.Save(dataPath, redis); err != nil {
					log.Printf("保存数据库失败：%v", err)
					continue
				}
				fmt.Println("删除成功")
			}
		case "EXISTS":
			if len(args) != 2 {
				fmt.Printf("用法：%s <键>\n", command)
				continue
			}
			n := redis.Exists(args[1])
			fmt.Println(n)
		case "EXPIRE":
			if len(args) != 3 {
				fmt.Printf("用法：%s <键> <秒数>\n", command)
				continue
			}
			key := args[1]
			seconds, err := strconv.Atoi(args[2])
			if err != nil {
				log.Println("过期时间必须是整数秒")
				continue
			}
			if n := redis.Expire(key, time.Duration(seconds)*time.Second); n != 1 {
				log.Println("过期时间设置失败")
				continue
			}
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println("过期时间设置成功")
		case "PERSIST":
			if len(args) != 2 {
				fmt.Printf("用法：%s <键>\n", command)
				continue
			}
			if n := redis.Persist(args[1]); n != 1 {
				log.Println("移除过期时间失败")
				continue
			}
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println("移除过期时间成功")
		case "TTL":
			if len(args) != 2 {
				fmt.Printf("用法：%s <键>\n", command)
				continue
			}
			keyTime := redis.Ttl(args[1])
			if keyTime == -1 {
				fmt.Println("键永不过期")
			} else if keyTime == -2 {
				fmt.Println("键不存在")
			} else {
				fmt.Println(keyTime)
			}
			fmt.Println("查询过期时间成功")
		case "RANDOMKEY":
			if len(args) != 1 {
				fmt.Printf("用法：%s\n", command)
				continue
			}
			name, ok := redis.RandOmKey()
			if !ok {
				fmt.Println("没有可返回的键")
			} else {
				fmt.Println(name)
			}
		case "RENAME":
			if len(args) != 3 {
				fmt.Printf("用法：%s <旧键> <新键>\n", command)
				continue
			}
			oldKey, newKey := args[1], args[2]
			if err := redis.ReName(oldKey, newKey); err != nil {
				log.Printf("重命名失败:%v\n", err)
				continue
			}
			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println("重命名成功")
		case "TYPE":
			if len(args) != 2 {
				fmt.Printf("用法：%s <键>\n", command)
				continue
			}
			typeString, ok := redis.Type(args[1])
			if !ok {
				log.Println("键不存在")
			} else {
				fmt.Println(typeString)
			}
		case "SET", "GETSET":
			if len(args) != 3 {
				fmt.Printf("用法：%s <键> <值>\n", command)
				continue
			}

			reply := "设置成功"
			if command == "SET" {
				redis.Set(args[1], args[2])
			} else {
				old, existed := redis.GetSet(args[1], args[2])
				reply = "键不存在"
				if existed {
					reply = fmt.Sprintf("%q", old)
				}
			}

			if err := db.Save(dataPath, redis); err != nil {
				log.Printf("保存数据库失败：%v", err)
				continue
			}
			fmt.Println(reply)
		case "GET":
			if len(args) != 2 {
				fmt.Println("用法:GET key")
				continue
			}
			keyValue, err := redis.Get(args[1])
			if err != nil {
				log.Println(nil)
			} else {
				fmt.Println(keyValue)
			}
		default:
			fmt.Printf("命令 %s 尚未实现\n", command)
		}
	}
	if err := scanner.Err(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "读取输入失败：", err)
	}
}

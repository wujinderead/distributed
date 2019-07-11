package main

import (
	"encoding/hex"
	"fmt"
	"github.com/go-redis/redis"
	"io"
	"math/rand"
	"time"
)

var (
	// run redis-server on docker
	// d run -d --name redis -h redis -p 6379:6379 redis:5.0.5
	redisServerAddr = "10.19.138.22:6379"
	counties = []string{"Norway", "Australia", "Switzerland", "Denmark", "Netherlands", "Germany", "Ireland",
		"United States", "Canada", "New Zealand", "Singapore", "Hong Kong", "Liechtenstein", "Sweden", "United Kingdom",
		"Iceland", "Korea, South", "Israel", "Luxembourg", "Japan", "Belgium", "France", "Austria", "Finland",
		"Slovenia", "Spain", "Italy", "Czech Republic", "Greece", "Estonia", "Brunei", "Cyprus", "Qatar", "Andorra",
		"Slovakia", "Poland", "Lithuania", "Malta", "Saudi Arabia", "Argentina", "United Arab Emirates", "Chile",
		"Portugal", "Hungary", "Bahrain", "Latvia", "Croatia", "Kuwait", "Montenegro", "Belarus", "Russia", "Oman",
		"Romania", "Uruguay", "Bahamas", "Kazakhstan", "Barbados", "Antigua and Barbuda", "Bulgaria", "Palau", "Panama",
		"Malaysia", "Mauritius", "Seychelles", "Trinidad and Tobago", "Serbia", "Cuba", "Lebanon", "Costa Rica", "Iran",
		"Venezuela", "Turkey", "Sri Lanka", "Mexico", "Brazil", "Georgia", "Saint Kitts and Nevis", "Azerbaijan",
		"Grenada", "Jordan", "Macedonia", "Ukraine", "Algeria", "Peru", "Albania", "Armenia", "Bosnia and Herzegovina",
		"Ecuador", "Saint Lucia", "China", "Fiji", "Mongolia", "Thailand", "Dominica", "Libya", "Tunisia", "Colombia",
		"Saint Vincent and the Grenadines", "Jamaica", "Tonga", "Belize", "Dominican Republic", "Suriname", "World",
		"Maldives", "Samoa", }
	hdis = []float64{0.944, 0.935, 0.93, 0.923, 0.922, 0.916, 0.916, 0.915, 0.913, 0.913, 0.912, 0.91, 0.908, 0.907,
		0.907, 0.899, 0.898, 0.894, 0.892, 0.891, 0.89, 0.888, 0.885, 0.883, 0.88, 0.876, 0.873, 0.87, 0.865, 0.861,
		0.856, 0.85, 0.85, 0.845, 0.844, 0.843, 0.839, 0.839, 0.837, 0.836, 0.835, 0.832, 0.83, 0.828, 0.824, 0.819,
		0.818, 0.816, 0.802, 0.798, 0.798, 0.793, 0.793, 0.793, 0.79, 0.788, 0.785, 0.783, 0.782, 0.78, 0.78, 0.779,
		0.777, 0.772, 0.772, 0.771, 0.769, 0.769, 0.766, 0.766, 0.762, 0.761, 0.757, 0.756, 0.755, 0.754, 0.752, 0.751,
		0.75, 0.748, 0.747, 0.747, 0.736, 0.734, 0.733, 0.733, 0.733, 0.732, 0.729, 0.727, 0.727, 0.727, 0.726, 0.724,
		0.724, 0.721, 0.72, 0.72, 0.719, 0.717, 0.715, 0.715, 0.714, 0.711, 0.706, 0.702, }
)

func main() {
	testString()
	testList()
	testSet()
	testZset()
}

func toClose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		fmt.Println("close err:", err.Error())
	}
}

func getRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: redisServerAddr,
		DB: 0,
		MaxRetries: 0,
		DialTimeout: 5*time.Second,
		ReadTimeout: 5*time.Second,
		WriteTimeout: 5*time.Second,
		//TLSConfig *tls.Config
	})
}

func testString() {
	client := getRedisClient()
	defer toClose(client)

	// set command: set value for key
	kvs := [][]interface{}{
		{"lgq", "programmer"},
		{"杜兰特", "篮球运动员"},
		// integer, float are deemed as string, though the encoding is different
		{"floater", 1.23e45},
		{"inter", int32(1193046)},
		{"int64er", int64(0xabcdef098765432)},
		{"uint64er", uint64(0xabcdef0987654321)},
		{"runer", 'Ⅲ'},
		{"byter", []byte("⅐⅑⅒⅓⅔⅕⅖⅗⅘⅙⅚⅛⅜⅝⅞")},
	}
	for i := range kvs {
		statusCmd := client.Set(kvs[i][0].(string), kvs[i][1], 0)
		str, err := statusCmd.Result()
		fmt.Println("re:", str, "err:", err)
	}
	fmt.Println()

	// get command: get value
	for i := range kvs {
		key := kvs[i][0].(string)
		strCmd := client.Get(key)
		fmt.Println(key, "=")
		str, err := strCmd.Result()
		if err != nil {
			fmt.Println("str:", err)
		} else {
			fmt.Println("str:", str)
		}
		byter, err := strCmd.Bytes()
		if err != nil {
			fmt.Println("bytes:", err)
		} else {
			fmt.Println("bytes:", hex.EncodeToString(byter))
		}
		f64, err := strCmd.Float64()      // parse string to float64
		if err != nil {
			fmt.Println("f64:", err)
		} else {
			fmt.Println("f64:", f64)
		}
		i64, err := strCmd.Int64()        // return int value of type is integer
		if err != nil {
			fmt.Println("i64:", err)
		} else {
			fmt.Println("i64:", i64)
		}
		fmt.Println()
	}

	// append command: append content to value
	intCmd := client.Append("lgq", "李哥")
	afterlen, err := intCmd.Result()
	fmt.Println(afterlen, err)   // return value len after append

	// setnx command: set value only when key not exists
	boolCmd := client.SetNX("lgq", "李哥", 0)
	hasset, err := boolCmd.Result()
	fmt.Println(hasset, err)     // return if preform set

	// `SET key value [expiration] XX`, set value only when key exists
	boolCmd = client.SetXX("lgq", "李哥", 0)
	hasset, err = boolCmd.Result()
	fmt.Println(hasset, err)     // return if preform set

	// exists command: check if key exists
	intCmd = client.Exists("lgq", "杜兰特", "durant")
	has, err := intCmd.Result()
	fmt.Println(has, err)        // one key exists, return 1; none key exists, return 0
	fmt.Println()

	// keys command: find keys
	strSliceCmd := client.Keys("i*")
	strs, err := strSliceCmd.Result()
	fmt.Println(strs, err)                // return all keys follow the queried pattern
	fmt.Println()

	// type command: get value type
	for i := range kvs {
		key := kvs[i][0].(string)
		statusCmd := client.Type(key)  // all type is string
		re, err := statusCmd.Result()
		fmt.Println(key, re, err)
		strCmd := client.ObjectEncoding(key)
		str, err := strCmd.Result()
		fmt.Println(key, str, err)
		/*
		| key       |  value                      | encoding  | 
		+-----------|-----------------------------+-----------|
		| "lgq"     | "programmer"                | embstr    |
		| "杜兰特"   | "篮球运动员"                  | embstr    |
		| "floater" | 1.23e45                     | raw       |
		| "inter"   | int32(1193046)              | int       |
		| "int64er" | int64(0xabcdef098765432)    | int       |
		| "uint64er"| uint64(0xabcdef0987654321)  | embstr    |
		| "runer"   | 'Ⅲ',                        | int      |
		| "byter"   | []byte("⅐⅑⅒⅓⅔⅕⅖⅗⅘⅙⅚⅛⅜⅝⅞")  | raw       |      */
	}
	fmt.Println()

	// mget, mset command: multi get and set
	statusCmd := client.MSet("lgq", "good luck", "runer", 123)
	re, err := statusCmd.Result()
	fmt.Println(re, err)
	sliceCmd := client.MGet("lgq", "runer", "aa", "byter")
	values, err := sliceCmd.Result()    // return all values for keys, nil for non-existed key
	fmt.Println("err:", err, len(values))
	for i := range values {
		if values[i] != nil {
			fmt.Println(i, values[i].(string))
		}
	}
	fmt.Println()

	// incr command: increment value if it is a number
	for i := range kvs {
		key := kvs[i][0].(string)
		intCmd := client.Incr(key)
		re, err := intCmd.Result()
		fmt.Println(key, re, err)   // return error if value is not integer, return incremented value
	}
}

func testList() {
	client := getRedisClient()
	defer toClose(client)

	// lpush command: push to list head, if list not exists, create first
	list := "lister"
	intCmd := client.LPush(list, "lgq", "墨西哥",)
	re, err := intCmd.Result()
	fmt.Println(re, err)        // return list length after push

	// rpush command: push to list tail
	intCmd = client.RPush(list, 'Ⅲ', []byte("⅐⅑⅒⅓⅔⅕⅖⅗⅘⅙⅚⅛⅜⅝⅞"))
	re, err = intCmd.Result()
	fmt.Println(re, err)        // return list length after add

	// lset command: set list item based on index
	statusCmd := client.LSet(list, 1, "秘鲁")
	status, err := statusCmd.Result()
	fmt.Println(status, err)    // return if lset ok
	fmt.Println()

	// lindex command: get list item based on index
	for i:=int64(0); i<4; i++ {
		strCmd := client.LIndex(list, i)
		str, err := strCmd.Result()
		fmt.Println(i, str, err)
	}
	fmt.Println()

	// lrange command: iterate list based on index range
	strSliceCmd := client.LRange(list, 0, 3)
	strs, err := strSliceCmd.Result()
	fmt.Println(strs, err)
	fmt.Println()

	// lrem command: remove item that equal to specific value
	// count > 0 : search from list head, remove |count| items that equals to value
	// count < 0 : search from list tail, remove |count| items that equals to value
	// count = 0 : remove all items equal to value
	intCmd = client.LRem(list, 0, int('Ⅲ'))
	re, err = intCmd.Result()
	fmt.Println(re, err)        // return how many items are removed
	fmt.Println()

	// lpop command: pop list head
	strCmd := client.LPop(list)
	str, err := strCmd.Result()
	fmt.Println(str, err)        // return popped value
	fmt.Println()

	// llen command: return list length
	intCmd = client.LLen(list)
	llen, err := intCmd.Result()
	fmt.Println(llen, err)        // return list length
	fmt.Println()

	// del command: delete key
	intCmd = client.Del(list)
	re, err = intCmd.Result()
	fmt.Println(re, err)        // return deleted key number
}

func testSet() {
	client := getRedisClient()
	defer toClose(client)

	// sadd command: add an element to set
	cities := "cities"
	intCmd := client.SAdd(cities, "los angeles", "são paulo", "quito", "ក្រុងសៀមរាប", "Нур-Султан", "თბილისი")
	re, err := intCmd.Result()
	fmt.Println(re, err)        // return the successfully added number

	// sismember command: check if an element is in set
	boolCmd := client.SIsMember(cities, []byte("ក្រុងសៀមរាប"))
	in, err := boolCmd.Result()
	fmt.Println(in, err)        // return true if in set else false

	// srandmember key [count]: return random members of given count
	strSliceCmd := client.SRandMemberN(cities, 2)
	strs, err := strSliceCmd.Result()
	fmt.Println("rand member:", strs, err)      // return members in string slice

	// spop key [count]: return random members of given count and remove them from set
	strSliceCmd = client.SPopN(cities, 2)
	strs, err = strSliceCmd.Result()
	fmt.Println("rand pop:", strs, err)         // return members in string slice

	// srem command: remove elements from set
	intCmd = client.SRem(cities, "Нур-Султан", "თბილისი")
	re, err = intCmd.Result()
	fmt.Println("deleted number:", re, err)     // return successfully deleted number

	// scard command: return set cardinality (elements number)
	intCmd = client.SCard(cities)
	re, err = intCmd.Result()
	fmt.Println("set cardinality:", re, err)    // return set elements number

	// smembers command: return all elements of a set
	strSliceCmd = client.SMembers(cities)
	strs, err = strSliceCmd.Result()
	fmt.Println("all members:", strs, err)      // return members in string slice
	fmt.Println()

	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0; i<500; i++ {
		client.SAdd("seta", rander.Intn(1000))
		client.SAdd("setb", rander.Intn(1000))
	}
	sa := client.SMembers("seta")
	saa, err := sa.Result()
	fmt.Println("sa:", saa, err)
	sb := client.SMembers("setb")
	sbb, err := sb.Result()
	fmt.Println("sb:", sbb, err)

	// sinter command: get intersection of sets
	strSliceCmd = client.SInter("seta", "setb")
	strs, err = strSliceCmd.Result()
	fmt.Println("intersect:", strs, err)       // return intersection set members

	// sdiff command: get difference of set, A-B = A-(A∩B)
	strSliceCmd = client.SDiff("seta", "setb")
	strs, err = strSliceCmd.Result()
	fmt.Println("diff:", strs, err)            // return difference set members

	// sunionstore command: get union of sets, and store the results
	intCmd = client.SUnionStore("union", "seta", "setb")
	re, err = intCmd.Result()
	fmt.Println("union size:", re, err)        // return union set size

	// sscan command: iterate set elements in multiple batches
	scanCmd := client.SScan("union", 0, "", 0)
	keys, cursor, err := scanCmd.Result()
	fmt.Println("keys:", keys, "cursor:", cursor, "err:", err)
	for cursor!=0 {
		scanCmd = client.SScan("union", cursor, "", 0)
		keys, cursor, err = scanCmd.Result()
		fmt.Println("keys:", keys, "cursor:", cursor, "err:", err)
	}
	fmt.Println()

	// del command: delete key
	intCmd = client.Del(cities, "seta", "setb")
	re, err = intCmd.Result()
	fmt.Println(re, err)       // return deleted key number
}

func testZset()  {
	client := getRedisClient()
	defer toClose(client)
	hdi := "hdi"

	// zadd command: add elements to sorted set
	// ZADD key [NX|XX] [CH] [INCR] score member [score member ...]
	// XX: Only update elements that already exist. Never add elements.
	// NX: Don't update already existing elements. Always add new elements.
	// CH: Modify the return value from the number of new elements added, to the total number of elements changed
	perm := rand.Perm(len(counties))
	for i := range perm {
		client.ZAdd(hdi, redis.Z{Score: hdis[perm[i]], Member: counties[perm[i]]})
	}

	// zscore command: get zset element score
	floatCmd := client.ZScore(hdi, "Lithuania")
	f, err := floatCmd.Result()
	fmt.Println("zscore:", f, err)    // return score of elements

	// zrank command: get zset element rank
	intCmd := client.ZRank(hdi, "Lithuania")
	i, err := intCmd.Result()
	fmt.Println("zrank:", i, err)     // return rank of elements

	// zincby command: increase score for an element
	floatCmd = client.ZIncrBy(hdi, -0.1, "Lithuania")
	f, err = floatCmd.Result()
	fmt.Println("zincby:", f, err)    // return modified score

	// zcard command: get zset cardinality
	intCmd = client.ZCard(hdi)
	i, err = intCmd.Result()
	fmt.Println("zcard:", i, err)     // return cardinality of zset

	// zcount command: count elements number within range of score
	intCmd = client.ZCount(hdi, "0.75", "0.85")
	i, err = intCmd.Result()
	fmt.Println("zcount:", i, err)     // return count of elements within score range
	fmt.Println()

	// zpopmax, zpopmin command: pop given count elements with max or min scores
	zSliceCmd := client.ZPopMax(hdi, 3)
	zs, err := zSliceCmd.Result()
	fmt.Println("zpopmax 3:", zs, err)   // return elements within range

	// zrange command: get elements based on rank
	strSliceCmd := client.ZRange(hdi, 0, 9)
	strs, err := strSliceCmd.Result()
	fmt.Println("zrange 0 9:", strs, err)   // return elements within range

	// zrevrange command: get elements based on descent rank
	zSliceCmd = client.ZRevRangeWithScores(hdi, 0, 9)
	zs, err = zSliceCmd.Result()
	fmt.Println("zrevrange 0 9 withscores:", zs, err)   // return elements within range

	// zrangebyscore command: get elements based on score range
	strSliceCmd = client.ZRangeByScore(hdi, redis.ZRangeBy{Min: "0.75", Max: "(0.85", Offset: 0, Count:10})
	strs, err = strSliceCmd.Result()
	fmt.Println("zrangebyscore 0.75 (0.85 limit 0 10:", strs, err)   // return elements within range

	// zrem command: remove elements from zset
	fmt.Println("to remove", counties[perm[0]], counties[perm[1]])
	intCmd = client.ZRem(hdi, counties[perm[0]], counties[perm[1]])
	i, err = intCmd.Result()
	fmt.Println("zrem:", i, err)     // return deleted elements count

	// zremrangebyscore command: remove elements from zset by score
	intCmd = client.ZRemRangeByScore(hdi, "(0.75", "0.78")
	i, err = intCmd.Result()
	fmt.Println("zremrangebyscore (0.75 0.78:", i, err)     // return deleted elements count
	fmt.Println()

	// zscan command: iterate zset elements in multiple batches
	scanCmd := client.ZScan(hdi, 0, "", 5)
	keys, cursor, err := scanCmd.Result()
	fmt.Println("keys:", keys, "cursor:", cursor, "err:", err)
	for cursor!=0 {
		scanCmd = client.SScan(hdi, cursor, "", 5)
		keys, cursor, err = scanCmd.Result()
		fmt.Println("keys:", keys, "cursor:", cursor, "err:", err)
	}
	fmt.Println()

	// del command: delete key
	intCmd = client.Del(hdi)
	i, err = intCmd.Result()
	fmt.Println(i, err)       // return deleted key number
}
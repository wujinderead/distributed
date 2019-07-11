package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"etcder"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	clientKeyFile  = "/home/xzy/etcd-keys/client.key"
	clientCertFile = "/home/xzy/etcd-keys/client.crt"
	caCertFile     = "/home/xzy/etcd-keys/ca.crt"
	endPoints      = []string{"https://10.19.138.135:12379", "https://10.19.138.22:12379", "https://10.19.138.127:12379"}
	timeout        = 5 * time.Second
)

func main() {
	//testGrpcConn()
	//testMemberList()
	//testMaintainState()
	//testSingleGetPut()
	//testRangeGet()
	//testKvcDoTxn()
	//testWatch()
	//testLease()
	//testLock()
	//testElection()
	//testElectionContext()
	testStm()
}

func testGrpcConn() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}
	conn := cli.ActiveConnection()
	fmt.Println("grpc conn target:", conn.Target())
	fmt.Println("grpc conn state:", conn.GetState().String())
	// method is typically called by this grpc conn
	// err := conn.Invoke(ctx context.Context, method string, args, reply interface{}, opts ...CallOption)
}

func testMemberList() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}

	// create a cluster
	cluster := clientv3.NewCluster(cli)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	// get cluster member list
	memberList, err := cluster.MemberList(ctx)
	cancel()
	if err != nil {
		fmt.Println("get member list err:", err)
		return
	}
	fmt.Println("clusterid:", memberList.Header.ClusterId, strconv.FormatUint(memberList.Header.ClusterId, 16))
	fmt.Println("memberid:", memberList.Header.MemberId)
	fmt.Println("raftterm:", memberList.Header.RaftTerm)
	fmt.Println("revision:", memberList.Header.Revision)
	for i, member := range memberList.Members {
		fmt.Println(i, "id:", member.ID, strconv.FormatUint(member.ID, 16))
		fmt.Println(i, "name:", member.Name)
		fmt.Println(i, "client:", member.ClientURLs)
		fmt.Println(i, "peer:", member.PeerURLs)
	}
}

func testMaintainState() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}
	maintain := clientv3.NewMaintenance(cli)
	fmt.Println("endpoints:", cli.Endpoints())
	for _, ep := range cli.Endpoints() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		state, err := maintain.Status(ctx, ep)
		cancel()
		if err != nil {
			fmt.Println("get status err:", err, ". endpoint:", ep)
			continue
		}
		fmt.Println("clusterid:", state.Header.ClusterId, strconv.FormatUint(state.Header.ClusterId, 16))
		fmt.Println("memberid:", state.Header.MemberId)
		fmt.Println("raftterm:", state.Header.RaftTerm)
		fmt.Println("revision:", state.Header.Revision)
		fmt.Println(ep, "version:", state.Version)
		fmt.Println(ep, "raft term:", state.RaftTerm)
		fmt.Println(ep, "raft index:", state.RaftIndex)
		fmt.Println(ep, "dbsize:", state.DbSize)
		fmt.Println(ep, "leader:", state.Leader)
	}
}

func testSingleGetPut() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}
	testdata := etcder.GetRandomKv(200, 5, 10, 8, 10)

	// kv client
	kvc := clientv3.NewKV(cli)

	// put
	for i := range testdata {
		putResp, err := kvc.Put(context.TODO(), testdata[i][0], testdata[i][1])
		if err != nil {
			fmt.Println("put err:", err)
			continue
		}
		fmt.Println("memberid:", putResp.Header.MemberId)
		fmt.Println("revision:", putResp.Header.Revision)
	}

	// get
	for i := range testdata {
		getResp, err := kvc.Get(context.TODO(), testdata[i][0])
		if err != nil {
			fmt.Println("get err:", err)
			continue
		}
		fmt.Println("memberid:", getResp.Header.MemberId)
		fmt.Println("revision:", getResp.Header.Revision)
		if len(getResp.Kvs) > 0 {
			fmt.Println(string(getResp.Kvs[0].Key), "=", string(getResp.Kvs[0].Value),
				", ", getResp.Kvs[0].Version, getResp.Kvs[0].CreateRevision, getResp.Kvs[0].ModRevision)
		}
		// Header.revision is the current global revision of all kvs
		// kv.createRevision is the revision when the kv is created
		// kv.modifyRevision is the latest revision when the kv is modified
		// kv.Version is how many times the kv has been modified
	}
}

func testRangeGet() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}
	kvc := clientv3.NewKV(cli)

	// with range
	ranged, err := kvc.Get(context.TODO(), "M", clientv3.WithRange("Y"))
	if err != nil {
		fmt.Println("get err:", err)
		return
	}
	fmt.Println("memberId:", ranged.Header.MemberId)
	fmt.Println("revision:", ranged.Header.Revision)
	fmt.Println("count:", ranged.Count, len(ranged.Kvs))
	for i := range ranged.Kvs {
		fmt.Println(string(ranged.Kvs[i].Key), "=", string(ranged.Kvs[i].Value), ", ",
			", ", ranged.Kvs[i].Version, ranged.Kvs[i].CreateRevision, ranged.Kvs[i].ModRevision)
	}
	fmt.Println()

	// with from key
	from, err := kvc.Get(context.TODO(), "z", clientv3.WithFromKey())
	if err != nil {
		fmt.Println("get err:", err)
		return
	}
	fmt.Println("memberId:", from.Header.MemberId)
	fmt.Println("revision:", from.Header.Revision)
	fmt.Println("count:", from.Count, len(from.Kvs))
	for i := range from.Kvs {
		fmt.Println(string(from.Kvs[i].Key), "=", string(from.Kvs[i].Value), ", ",
			", ", from.Kvs[i].Version, from.Kvs[i].CreateRevision, from.Kvs[i].ModRevision)
	}
	fmt.Println()

	// with rev, with limit, with sort
	rev, err := kvc.Get(context.TODO(), "O", clientv3.WithRange("z"), clientv3.WithRev(180),
		clientv3.WithLimit(20), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		fmt.Println("get err:", err)
		return
	}
	fmt.Println("memberId:", rev.Header.MemberId)
	fmt.Println("revision:", rev.Header.Revision)
	// count is 41, len is 20. it means there is 41 kvs met the queried condition,
	// but since limit is 20, on 20 kvs returned
	fmt.Println("count:", rev.Count, len(rev.Kvs))
	for i := range rev.Kvs {
		fmt.Println(string(rev.Kvs[i].Key), "=", string(rev.Kvs[i].Value), ", ",
			", ", rev.Kvs[i].Version, rev.Kvs[i].CreateRevision, rev.Kvs[i].ModRevision)
	}
}

func testKvcDoTxn() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}
	kvc := clientv3.NewKV(cli)

	// do
	/*doPut, err := kvc.Do(context.TODO(), clientv3.OpPut("libai", "侠客行"))
	if err != nil {
		fmt.Println("do put err:", err)
		return
	}
	fmt.Println("do put memberId:", doPut.Put().Header.MemberId)
	fmt.Println("do put revision:", doPut.Put().Header.Revision)
	fmt.Println() */
	doGet, err := kvc.Do(context.TODO(), clientv3.OpGet("libai"))
	if err != nil {
		fmt.Println("do put err:", err)
		return
	}
	fmt.Println("do get memberId:", doGet.Get().Header.MemberId)
	fmt.Println("do get revision:", doGet.Get().Header.Revision)
	fmt.Println("do get:", string(doGet.Get().Kvs[0].Key), "=", string(doGet.Get().Kvs[0].Value),
		", ", doGet.Get().Kvs[0].Version, doGet.Get().Kvs[0].CreateRevision, doGet.Get().Kvs[0].ModRevision)
	fmt.Println()

	// transaction, one txn can has an 'if-then-else' condition statement,
	// to do all get, set update in one single txn.
	txn := kvc.Txn(context.TODO())
	commit, err := txn.If(clientv3.Compare(
		clientv3.Value("libai"), "=", "侠客行"),
	).Then(
		clientv3.OpPut("libai", "梦游天姥吟留别"),
		clientv3.OpPut("dufu", "茅屋为秋风所破歌"),
	).Else(
		clientv3.OpPut("libai", "侠客行"),
		clientv3.OpPut("dufu", "qian"),
	).Commit()
	if err != nil {
		fmt.Println("txn commit err:", err)
		return
	}
	fmt.Println("commit memberId:", commit.Header.MemberId)
	fmt.Println("commit revision:", commit.Header.Revision)
	fmt.Println("commit succeeded:", commit.Succeeded)
	for i := range commit.Responses {
		fmt.Println("commit put:", commit.Responses[i].GetResponsePut().Header.Revision)
	}
	fmt.Println()

	doGet, _ = kvc.Do(context.TODO(), clientv3.OpGet("libai"))
	fmt.Println("do get memberId:", doGet.Get().Header.MemberId)
	fmt.Println("do get revision:", doGet.Get().Header.Revision)
	fmt.Println("do get:", string(doGet.Get().Kvs[0].Key), "=", string(doGet.Get().Kvs[0].Value),
		", ", doGet.Get().Kvs[0].Version, doGet.Get().Kvs[0].CreateRevision, doGet.Get().Kvs[0].ModRevision)
	doGet, _ = kvc.Do(context.TODO(), clientv3.OpGet("dufu"))
	fmt.Println("do get memberId:", doGet.Get().Header.MemberId)
	fmt.Println("do get revision:", doGet.Get().Header.Revision)
	fmt.Println("do get:", string(doGet.Get().Kvs[0].Key), "=", string(doGet.Get().Kvs[0].Value),
		", ", doGet.Get().Kvs[0].Version, doGet.Get().Kvs[0].CreateRevision, doGet.Get().Kvs[0].ModRevision)
}

func testWatch() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}

	// watch, put and delete event
	watcher := clientv3.NewWatcher(cli)
	ctx, cancel := context.WithCancel(context.Background())
	watchChan := watcher.Watch(ctx, "foo", clientv3.WithPrefix())
	go func() { // to put in a separate goroutine
		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			_, err = cli.Put(context.TODO(), "foo"+strconv.Itoa(rand.Intn(5)), "bar"+strconv.Itoa(rand.Intn(10000)))
			if err != nil {
				fmt.Println("put err:", err)
			}
			if rand.Uint32()&0x80000001 == 0 { // 1/4 rate
				_, err = cli.Delete(context.TODO(), "foo"+strconv.Itoa(rand.Intn(5)))
				if err != nil {
					fmt.Println("delete err:", err)
					continue
				}
			}
		}
		err := watcher.Close()
		if err != nil {
			fmt.Println("close watcher err:", err)
		}
	}()
	for change := range watchChan {
		fmt.Println("revision:", change.Header.Revision)
		fmt.Println("canceled:", change.Canceled)
		fmt.Println("compact:", change.CompactRevision)
		fmt.Println("created:", change.Created)
		fmt.Println("err:", change.Err())
		if err != nil {
			fmt.Println("watch err:", err)
			break
		}
		for i := range change.Events {
			fmt.Println(string(change.Events[i].Kv.Key), "=", string(change.Events[i].Kv.Value),
				", ", change.Events[i].Kv.Version, change.Events[i].Kv.CreateRevision, change.Events[i].Kv.ModRevision)
			fmt.Println("iscreate:", change.Events[i].IsCreate())
			fmt.Println("ismodify:", change.Events[i].IsModify())
			fmt.Println("type:", change.Events[i].Type.String())
			fmt.Println()
		}
	}
	cancel()
}

func testLease() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}

	// lease is used for key ttl in api v3
	leaser := clientv3.NewLease(cli)
	leaseResp, err := leaser.Grant(context.TODO(), 8) // 0
	if err != nil {
		fmt.Println("grant lease err:", err)
		return
	}
	id := leaseResp.ID
	fmt.Println("lease id:", id)
	fmt.Println("revision:", leaseResp.Revision)
	fmt.Println("ttl:", leaseResp.TTL)
	fmt.Println("err:", leaseResp.Error)
	fmt.Println()

	time.Sleep(2 * time.Second) // 2
	ttlResp, err := leaser.TimeToLive(context.TODO(), id)
	if err != nil {
		fmt.Println("lease ttl err:", err)
		return
	}
	fmt.Println("revision:", ttlResp.Revision)
	fmt.Println("ttl:", ttlResp.TTL)
	fmt.Println("granted ttl:", ttlResp.GrantedTTL)
	fmt.Println("keys:", ttlResp.Keys)
	fmt.Println()

	time.Sleep(2 * time.Second) // 4
	kvc := clientv3.NewKV(cli)
	putResp, _ := kvc.Put(context.TODO(), "李白", "蜀道难", clientv3.WithLease(id))
	fmt.Println("put revision:", putResp.Header.Revision)
	putResp, _ = kvc.Put(context.TODO(), "杜甫", "兵车行", clientv3.WithLease(id))
	fmt.Println("put revision:", putResp.Header.Revision)

	time.Sleep(2 * time.Second) // 6
	ttlResp, err = leaser.TimeToLive(context.TODO(), id)
	if err != nil {
		fmt.Println("lease ttl err:", err)
		return
	}
	fmt.Println("revision:", ttlResp.Revision)
	fmt.Println("ttl:", ttlResp.TTL)
	fmt.Println("granted ttl:", ttlResp.GrantedTTL)
	fmt.Println("keys:", ttlResp.Keys)
	getResp, _ := kvc.Get(context.TODO(), "李白")
	fmt.Println("get revision:", putResp.Header.Revision)
	if len(getResp.Kvs) > 0 {
		fmt.Println("get key:", string(getResp.Kvs[0].Key), "=", string(getResp.Kvs[0].Value))
	}
	getResp, _ = kvc.Get(context.TODO(), "杜甫")
	fmt.Println("get revision:", putResp.Header.Revision)
	if len(getResp.Kvs) > 0 {
		fmt.Println("get key:", string(getResp.Kvs[0].Key), "=", string(getResp.Kvs[0].Value))
	}
	fmt.Println()

	time.Sleep(4 * time.Second) // 10
	ttlResp, err = leaser.TimeToLive(context.TODO(), id)
	if err != nil {
		fmt.Println("lease ttl err:", err)
		return
	}
	fmt.Println("revision:", ttlResp.Revision)
	fmt.Println("ttl:", ttlResp.TTL)
	fmt.Println("granted ttl:", ttlResp.GrantedTTL)
	fmt.Println("keys:", ttlResp.Keys)
	// ttl is over, get null keys
	getResp, _ = kvc.Get(context.TODO(), "李白")
	fmt.Println("get revision:", putResp.Header.Revision)
	if len(getResp.Kvs) > 0 {
		fmt.Println("get key:", string(getResp.Kvs[0].Key), "=", string(getResp.Kvs[0].Value))
	}
	getResp, _ = kvc.Get(context.TODO(), "杜甫")
	fmt.Println("get revision:", putResp.Header.Revision)
	if len(getResp.Kvs) > 0 {
		fmt.Println("get key:", string(getResp.Kvs[0].Key), "=", string(getResp.Kvs[0].Value))
	}
}

func testLock() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}

	// different goroutine manipulate distributed processes
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(i int) {
			session, err := concurrency.NewSession(cli, concurrency.WithTTL(3))
			if err != nil {
				fmt.Println("new session err:", err)
				return
			}
			locker := concurrency.NewLocker(session, "TestLock")
			locker.Lock()
			fmt.Println("role", i, "acquire lock at", time.Now())
			time.Sleep(2 * time.Second)
			locker.Unlock()
			fmt.Println("role", i, "release lock at", time.Now())
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func testElection() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}

	// different goroutine manipulate distributed processes
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(i int) {
			session, err := concurrency.NewSession(cli, concurrency.WithTTL(3))
			defer etcder.ToClose(session)
			if err != nil {
				fmt.Println("new session err:", err)
				return
			}
			elect := concurrency.NewElection(session, "TestElect")
			err = elect.Campaign(context.TODO(), "cand"+strconv.Itoa(i))
			if err != nil {
				fmt.Println("campaign err:", err)
				return
			}
			fmt.Println("role", i, "elected as leader at", time.Now())
			time.Sleep(2 * time.Second)
			err = elect.Resign(context.TODO())
			if err != nil {
				fmt.Println("resign err:", err)
				return
			}
			fmt.Println("role", i, "resigned leader at", time.Now())
			wg.Done()
		}(i)
	}
	// main session to observe
	session, err := concurrency.NewSession(cli, concurrency.WithTTL(3))
	defer etcder.ToClose(session)
	if err != nil {
		fmt.Println("new session err:", err)
		return
	}
	elect := concurrency.NewElection(session, "TestElect")
	time.Sleep(1 * time.Second)
	for i := 0; i < 4; i++ {
		leaderResp, err := elect.Leader(context.TODO())
		if err != nil {
			fmt.Println("get leader err:", err)
			continue
		}
		if len(leaderResp.Kvs) > 0 {
			fmt.Println("main get", string(leaderResp.Kvs[0].Key), "=", string(leaderResp.Kvs[0].Value), "at", time.Now())
		}
		time.Sleep(2 * time.Second)
	}
	wg.Wait()
}

func testElectionContext() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}

	// different goroutine manipulate distributed processes
	var wg sync.WaitGroup
	numMember := 3
	wg.Add(numMember)
	for i := 0; i < numMember; i++ {
		go func(i int) {
			defer wg.Done()
			session, err := concurrency.NewSession(cli, concurrency.WithTTL(20))
			defer etcder.ToClose(session)
			if err != nil {
				fmt.Println("new session err:", err)
				return
			}
			elect := concurrency.NewElection(session, "TestElect")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			myKey := "cand" + strconv.Itoa(i)
			err = elect.Campaign(ctx, myKey)
			if err != nil {
				fmt.Println("campaign err:", err, "at", time.Now())
			}
			cancel()
			leaderResp, err := elect.Leader(context.TODO())
			if err != nil {
				fmt.Println("get leader err:", err)
				return
			}
			if string(leaderResp.Kvs[0].Value) == myKey {
				fmt.Println(myKey, "is leader")
			} else {
				fmt.Println(myKey, "is follower")
			}
		}(i)
	}
	wg.Wait()
}

func testStm() {
	cli, err := getTlsClient()
	defer etcder.ToClose(cli)
	if err != nil {
		fmt.Println("create tls client err:", err)
		return
	}

	// stm: software transactional memory, used for distributed transaction
	numAccount := 5
	// make accounts
	for i := 0; i < numAccount; i++ {
		_, _ = cli.Put(context.TODO(), fmt.Sprintf("account%d", i), "100")
	}

	// exchange func
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	exchange := func(stm concurrency.STM) error {
		from, to := rander.Intn(numAccount), rander.Intn(numAccount)
		if from == to {
			return nil
		}

		fromStr := stm.Get(fmt.Sprintf("account%d", from))
		toStr := stm.Get(fmt.Sprintf("account%d", to))
		fromBalance, _ := strconv.Atoi(fromStr)
		toBalance, _ := strconv.Atoi(toStr)

		transfer := fromBalance / 2
		stm.Put(fmt.Sprintf("account%d", from), strconv.Itoa(fromBalance-transfer))
		stm.Put(fmt.Sprintf("account%d", to), strconv.Itoa(toBalance+transfer))
		return nil
	}

	var wg sync.WaitGroup
	concurrencyLevel := 20
	wg.Add(concurrencyLevel)
	for i := 0; i < concurrencyLevel; i++ {
		go func() {
			defer wg.Done()
			txnResp, err := concurrency.NewSTM(cli, exchange)
			if err != nil {
				fmt.Println("stm err:", err)
				return
			}
			fmt.Println("success:", txnResp.Succeeded)
		}()
	}

	// get sum
	wg.Wait()
	all := 0
	for i := 0; i < numAccount; i++ {
		getResp, _ := cli.Get(context.TODO(), fmt.Sprintf("account%d", i))
		cur, _ := strconv.Atoi(string(getResp.Kvs[0].Value))
		fmt.Println(i, "has", cur)
		all += cur
	}
	fmt.Println("all:", all)
}

func getTlsClient() (cli *clientv3.Client, err error) {
	// get client cert and key
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		fmt.Println("get client key cert pair err:", err)
		return nil, err
	}

	// set ca cert pool
	pool := x509.NewCertPool()
	cacert, err := etcder.GetCertFromFile(caCertFile)
	if err != nil {
		fmt.Println("get key cert pair err:", err)
		return nil, err
	}
	pool.AddCert(cacert)

	// create client
	return clientv3.New(clientv3.Config{
		Endpoints:   endPoints,
		DialTimeout: timeout,
		TLS: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		},
	})
}

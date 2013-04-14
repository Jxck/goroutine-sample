// http://talks.golang.org/2012/concurrency.slide
// http://www.slideshare.net/jgrahamc/go-oncurrency
package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func randomTime() time.Duration {
	return time.Duration(rand.Intn(1e3)) * time.Millisecond
}

func message(msg string) {
	for i := 0; i < 5; i++ {
		log.Println(msg)
		time.Sleep(randomTime())
	}
}

func main1() {
	// メッセージを止めるまで出し続ける
	message("main1")
}

func main2() {
	// goroutine で実行
	go message("main2")
	log.Println("start")
	time.Sleep(1 * time.Second)
	log.Println("end")
	// goroutine は終わってなくても
	// この関数を抜けると止まる。
}

///////////////////////////////////////////////////////////////////////////////

func message3(msg string, receiver chan string) {
	for i := 0; i < 5; i++ {
		// ちょっと加工して五回返す
		receiver <- fmt.Sprintf("%d %s", i, msg)
		time.Sleep(2 * time.Second) // ここはあえて長く
	}
}

func main3() {
	// 受信目的でチャネルを生成
	receiver := make(chan string)
	// 一緒に渡して goroutine にする
	go message3("main3", receiver)
	log.Println(<-receiver) // ここの受信はそれぞれブロックしてる
	log.Println(<-receiver) // もしメソッドが 3 回未満しか戻さなければ
	log.Println(<-receiver) // deadlock になる
	// 3 回だけ受け取って終わり。
}

///////////////////////////////////////////////////////////////////////////////

func message4(msg string) <-chan string { // 受信専用のチャネルを返す
	// チャネルを生成して返す
	receiver := make(chan string)
	go func() { // 無名関数をその場で goroutine にして実行
		for i := 0; i < 5; i++ {
			// チャネルへはクロージャでアクセス
			receiver <- fmt.Sprintf("%d %s", i, msg)
			time.Sleep(2 * time.Second) // ここはあえて長く
		}
	}()
	return receiver
}

func main4() {
	// channel は受信目的だけだったので channel の生成を委譲し
	// 受信用のチャネルだけを戻り値として取得する。
	// ただし、戻り値をとる場合は go で実行できないので、
	// goroutine の生成も関数の中になる。
	receiver := message4("main4")
	log.Println(<-receiver) // ここの受信はそれぞれブロックしてる
	log.Println(<-receiver)
	log.Println(<-receiver)
	// 3 回だけ受け取って終わり。
}

///////////////////////////////////////////////////////////////////////////////

func message5(msg string) <-chan string {
	receiver := make(chan string)
	go func() {
		for i := 0; i < 5; i++ {
			receiver <- fmt.Sprintf("%d %s", i, msg)
			time.Sleep(randomTime()) // ランダムに戻しただけ
		}
	}()
	return receiver
}

func main5() {
	// message5 を別メッセージで呼ぶ
	receiver1 := message5("main5-1")
	receiver2 := message5("main5-2")
	for i := 0; i < 5; i++ {
		log.Println(<-receiver1) // ここの受信はそれぞれブロックしてる
		log.Println(<-receiver2) // なので、ランダムに戻っても順序が保たれる
	}
}

///////////////////////////////////////////////////////////////////////////////

func message6(msg string) <-chan string {
	receiver := make(chan string)
	go func() {
		for i := 0; i < 5; i++ {
			receiver <- fmt.Sprintf("%d %s", i, msg)
			time.Sleep(time.Millisecond / 2) // ごく短時間に設定
		}
	}()
	return receiver
}

func mux(c1, c2 <-chan string) <-chan string {
	// 二つのチャネルをまとめるためのチャネル
	c := make(chan string)

	// それぞれに対し goroutine を作る
	go func() {
		for { // 終わらないようにループ
			c <- <-c1 // 右から左に流す
		}
	}()
	go func() {
		for {
			c <- <-c2
		}
	}()
	return c // このチャネルに両方が流れる
}

func main6() {
	// message5, message6 を別メッセージで呼ぶ
	receiver1 := message5("main5-1")
	receiver2 := message6("main5-2")
	receiver := mux(receiver1, receiver2)
	for i := 0; i < 10; i++ {
		log.Println(<-receiver) // 先着順になる
	}
}

///////////////////////////////////////////////////////////////////////////////

type Message struct {
	str  string
	wait chan bool // HL
}

// 先の mux と同じ、チャネルをまとめる
// 型を Message にしただけ
func muxMessage(c1, c2 <-chan Message) <-chan Message {
	c := make(chan Message)
	go func() {
		for {
			c <- <-c1
		}
	}()
	go func() {
		for {
			c <- <-c2
		}
	}()
	return c
}

// 同期をとるために Message 全体で 1 つの channel を共有
func message7(msg string) <-chan Message {
	c := make(chan Message)
	waitForIt := make(chan bool) // 全ての Message で共有する
	go func() {
		for i := 0; ; i++ {
			c <- Message{ // メッセージを投げる (1) (5)
				fmt.Sprintf("%s: %d", msg, i),
				waitForIt,
			}
			time.Sleep(randomTime()) // 遅延 (4) (8)
			// ここで同期をとっている。
			// メッセージが送られてくるまで先に進めない
			<-waitForIt // (10)	(12)
		}
	}()
	return c
}

func main7() {
	receiver1 := message7("main7-1")
	receiver2 := message7("main7-2")
	receiver := muxMessage(receiver1, receiver2)
	// 1, 2 両方の表示が終わったら
	// wait にシグナルを送って再開する。
	// mux した意味があるのか。。
	for i := 0; i < 5; i++ {
		msg1 := <-receiver    // 受信 (2)
		log.Println(msg1.str) // 表示 (3)
		msg2 := <-receiver    // 受信 (6)
		log.Println(msg2.str) // 表示 (7)
		msg1.wait <- true     // 再開 (9)
		msg2.wait <- true     // 再開 (11)
	}
	log.Println("end")
}

///////////////////////////////////////////////////////////////////////////////

func muxSelect(c1, c2 <-chan string) <-chan string {
	c := make(chan string) // このチャネルに束ねる
	go func() {
		// １つのチャネルに束ねる
		for {
			select {
			case s := <-c1:
				c <- s // c はクロージャで参照
			case s := <-c2:
				c <- s
			}
		}
	}()
	return c
}

func main8() {
	// message5, message6 を別メッセージで呼ぶ
	receiver1 := message5("main8-1")
	receiver2 := message6("main8-2")
	// チャネルを１つに束ねる
	receiver := muxSelect(receiver1, receiver2)
	for i := 0; i < 10; i++ {
		log.Println(<-receiver) // 先着順になる
	}
}

///////////////////////////////////////////////////////////////////////////////

func message9(msg string) <-chan string {
	receiver := make(chan string)
	go func() {
		for i := 0; i < 50; i++ { // 長めにデータを返す
			receiver <- fmt.Sprintf("%d %s", i, msg)
			time.Sleep(time.Second / 2)
		}
	}()
	return receiver
}

func main9() {
	c := message9("main9")
	timeout := time.After(3 * time.Second) // 3 秒後にメッセージを投げる
	for {
		select {
		case s := <-c:
			log.Println(s)
		case s := <-timeout: // 3 秒後に実行
			log.Println("end", s)
			return // goroutine を抜けて止める
		}
	}
}

// これだとループの度に channel が生成されるので
// 確実なタイムアウトが出来ない(時間が短ければ起こるときも)
// func main9() {
// 	c := message9("main9")
// 	for {
// 		select {
// 		case s := <-c:
// 			log.Println(s)
// 		case s := <-time.After(3 * time.Second)
// 			log.Println("end", s)
// 			return
// 		}
// 	}
// }

///////////////////////////////////////////////////////////////////////////////

// quit channel
// 呼び出し側から、停止要求を投げる
func message10(msg string, quit <-chan bool) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Second / 2)
			select {
			case c <- fmt.Sprintf("%d %s", i, msg):
				// noop
			case <-quit: // このメッセージで止める
				return
			}
		}
	}()
	return c
}

func main10() {
	quit := make(chan bool)               // 停止用の channel
	receiver := message10("main10", quit) // ここで一緒に渡す
	for i := 0; i < 10; i++ {             // 10 個メッセージを取得
		log.Println(<-receiver)
	}
	quit <- true // 終わったら止める
}

///////////////////////////////////////////////////////////////////////////////

// quit channel
// 呼び出し側からの停止要求で終了処理をさせて
// 終了処理が完了した事を呼び出し側に知らせる。
// クリーンアップ用の同期パターン

func message11(msg string, quit chan string) <-chan string {
	receiver := make(chan string)
	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Second / 2)
			select {
			case receiver <- fmt.Sprintf("%d %s", i, msg):
				// noop
			case s := <-quit:
				log.Println(s)
				quit <- "done"
				return
			}
		}
	}()
	return receiver
}

func main11() {
	quit := make(chan string)
	receiver := message11("main11", quit)
	for i := 0; i < 10; i++ {
		log.Println(<-receiver)
	}
	quit <- "stop"
	log.Println(<-quit)
}

///////////////////////////////////////////////////////////////////////////////

// [leftmost] - ([left] - [right]) <- 1
func pipe(left, right chan int) {
	left <- 1 + <-right // 右から左に 1 加えて渡す
}

func main12() {
	// 最初は全て左端に初期化
	leftmost := make(chan int)
	right := leftmost
	left := leftmost
	for i := 0; i < 10; i++ {
		right = make(chan int) // 新しい右 channel
		go pipe(left, right)   // 左と右を繋ぐ
		left = right           // もう一個右に移動
	}
	go func() {
		right <- 1 // 右端にデータを送る
	}()
	log.Println(<-leftmost) // 左端から取り出す
}

///////////////////////////////////////////////////////////////////////////////
// Google Search

type Search func(query string) Result
type Result string

// fake
func fakeSearch(kind string) Search {
	// 検索結果を返す関数を返す関数
	return func(query string) Result {
		time.Sleep(randomTime())
		return Result(fmt.Sprintf("%s result for %q", kind, query))
	}
}

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

// search
func Google(query string) (results []Result) {
	// 各サーチを実行する
	// ブロックするので同期になる
	results = append(results, Web(query))
	results = append(results, Image(query))
	results = append(results, Video(query))
	return
}

// test
func main13() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	log.Println(results)
	log.Println(elapsed)
}

///////////////////////////////////////////////////////////////////////////////
// 各々を並行に実行する

func GoogleAsync(query string) (results []Result) {
	c := make(chan Result)
	go func() {
		c <- Web(query)
	}()
	go func() {
		c <- Image(query)
	}()
	go func() {
		c <- Video(query)
	}()
	for i := 0; i < 3; i++ {
		results = append(results, <-c)
	}
	return results
}

func main14() {
	results := GoogleAsync("golang")
	log.Println(results)
}

///////////////////////////////////////////////////////////////////////////////
// タイムアウトを設定して、遅いのを捨てる

func GoogleTimeout(query string) (results []Result) {
	c := make(chan Result)
	go func() {
		c <- Web(query)
	}()
	go func() {
		c <- Image(query)
	}()
	go func() {
		c <- Video(query)
	}()
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-time.After(600 * time.Millisecond):
			log.Println("timeout")
			return
		}
	}
	return
}

func main15() {
	results := GoogleTimeout("golang")
	log.Println(results)
}

///////////////////////////////////////////////////////////////////////////////
// 一番早いのだけ使う

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	for i := range replicas {
		go func() {
			c <- replicas[i](query)
		}()
	}
	// 一番最初だけ返す
	return <-c
}

func main16() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	result := First("golang",
		fakeSearch("replica 1"),
		fakeSearch("replica 2"))
	elapsed := time.Since(start)
	log.Println(result)
	log.Println(elapsed)
}

///////////////////////////////////////////////////////////////////////////////
//

var (
	Web1   = fakeSearch("web")
	Web2   = fakeSearch("web")
	Image1 = fakeSearch("image")
	Image2 = fakeSearch("image")
	Video1 = fakeSearch("video")
	Video2 = fakeSearch("video")
)

func GoogleReduceTail(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- First(query, Web1, Web2) }()
	go func() { c <- First(query, Image1, Image2) }()
	go func() { c <- First(query, Video1, Video2) }()
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-time.After(80 * time.Second):
			log.Println("time out")
			return
		}
	}
	return
}

func main17() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	result := First("golang",
		fakeSearch("replica 1"),
		fakeSearch("replica 2"))
	elapsed := time.Since(start)
	log.Println(result)
	log.Println(elapsed)
}

func main() {
	/**
	 * Don't communicate by sharing memory, share memory by communicating.
	 */

	log.SetFlags(log.Lmicroseconds)
	main1()
	main2()
	main3()
	main4()
	main5()
	main6()
	main7()
	main8()
	main9()
	main10()
	main11()
	main12()
	main13()
	main14()
	main15()
	main16()
	main17()

	// Channel の buffer の話が無い
	// rob 先生のだと main3 の次
}

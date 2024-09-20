package storage

import (
	"IpProxyPool/middleware/database"
	"context"
	"fmt"
	"testing"
	"time"
)

//func init() {
//	config.ConfigFile = "/home/ycd/self_data/source_code/IpProxyPool/conf/config.yaml"
//	config.InitConfig()
//	setting := config.ServerSetting
//	log.Init(true)
//	log.SetLogLevel(setting.Log.Level)
//	database.InitDB(&setting.Database)
//}
//func TestCheckProxyDB(t *testing.T) {
//	CheckProxyDB()
//}
//
//func TestCheckIp(t *testing.T) {
//	ip1 := &database.IP{
//		ProxyId:       80,
//		ProxyHost:     "47.115.219.60",
//		ProxyPort:     7890,
//		ProxyType:     "HTTPS",
//		ProxyLocation: "SSL高匿_中国阿里云",
//		ProxySpeed:    1046,
//		ProxySource:   "http://www.ip3366.net",
//		CreateTime:    time.Now(),
//		UpdateTime:    time.Now(),
//	}
//	ip2 := &database.IP{
//		ProxyId:       4842,
//		ProxyHost:     "81.69.33.240",
//		ProxyPort:     7890,
//		ProxyType:     "HTTPS",
//		ProxyLocation: "SSL高匿_上海市腾讯云",
//		ProxySpeed:    1046,
//		ProxySource:   "http://www.ip3366.net",
//		CreateTime:    time.Now(),
//		UpdateTime:    time.Now(),
//	}
//	ip3 := &database.IP{
//		ProxyId:       4842,
//		ProxyHost:     "222.190.173.176",
//		ProxyPort:     8089,
//		ProxyType:     "https",
//		ProxyLocation: "SSL高匿_上海市腾讯云",
//		ProxySpeed:    1046,
//		ProxySource:   "http://www.ip3366.net",
//		CreateTime:    time.Now(),
//		UpdateTime:    time.Now(),
//	}
//
//	fmt.Println(ip2)
//	fmt.Println(ip1)
//	fmt.Println(ip3)
//	fmt.Println(CheckIP(ip3))
//
//}
//
//func Test1(t *testing.T) {
//	aChan := make(chan int)
//	go func() {
//		for {
//			if aChan != nil {
//				aChan <- 1
//				fmt.Println("来了")
//				time.Sleep(1 * time.Second)
//			}
//		}
//	}()
//	go func() {
//		for {
//			<-aChan
//		}
//	}()
//	time.Sleep(10 * time.Second)
//	fmt.Println("结束")
//	close(aChan)
//	select {}
//}
//
//func contextPart4() {
//	ctx, cancel := context.WithCancel(context.Background())
//	go watch(ctx, "task 1")
//	go watch(ctx, "task 2")
//	go watch(ctx, "task 3")
//
//	time.Sleep(10 * time.Second)
//	fmt.Println("可以通知任务结束")
//	cancel()
//	time.Sleep(5 * time.Second)
//}
//
//func watch(ctx context.Context, name string) {
//	for {
//		select {
//		case <-ctx.Done():
//			fmt.Println(name + " 任务即将要退出了...")
//			return
//		default:
//			fmt.Println(name + " goroutine 继续处理任务中...")
//			time.Sleep(2 * time.Second)
//		}
//	}
//}
//
//func Test_checkIP(t *testing.T) {
//	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancelFunc()
//	resultChan := make(chan *database.IP, 1)
//	go func() {
//		checkIP(ctx, &database.IP{ProxyHost: "205.178.186.112", ProxyPort: 8443}, resultChan)
//	}()
//	d := <-resultChan
//	if d == nil {
//		fmt.Println("nil")
//		return
//	}
//	fmt.Println("来了  ", d)
//}

func Benchmark_checkIP(b *testing.B) {

	for n := 0; n < b.N; n++ {
		aa()
	}

}
func aa() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer func() {
		cancelFunc()
		fmt.Println("结束")
	}()

	d := checkIP(ctx, &database.IP{ProxyHost: "205.178.186.112", ProxyPort: 8443})

	fmt.Println("来了  ", d)
}

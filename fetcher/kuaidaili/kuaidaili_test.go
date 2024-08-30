package kuaidaili

import (
	"IpProxyPool/middleware/database"
	"fmt"
	"github.com/youcd/toolkit/log"
	"reflect"
	"testing"
)

func init() {
	log.Init(true)
}
func TestKuaiDaiLi(t *testing.T) {
	type args struct {
		proxyType string
	}
	tests := []struct {
		name string
		args args
		want []*database.IP
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KuaiDaiLi(tt.args.proxyType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KuaiDaiLi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKuaiDaiLiInha(t *testing.T) {
	got := KuaiDaiLiInha()
	for _, ip := range got {
		fmt.Printf("%#v\n", ip)
	}
}

func TestKuaiDaiLiIntr(t *testing.T) {
	got := KuaiDaiLiIntr()
	for _, ip := range got {
		fmt.Printf("%#v\n", ip)
	}
}

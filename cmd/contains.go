package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/learn-go/cidr/internal/core"
	"github.com/spf13/cobra"
)

const (
	containsExample = "# Check whether an IPv4 CIDR range contains a given IPv4 address\n" +
		"cidr contains 10.0.0.0/16 10.0.14.5\n" +
		"\n" +
		"# Check whether an IPv6 CIDR range contains a given IPv6 address\n" +
		"cidr contains 2001:db8:1234:1a00::/106 2001:db8:1234:1a00::"
)

var (
	containsCmd = &cobra.Command{
		// [填空 9] Use / Short / Example，参考 count.go 的写法自己写一遍
		Run: func(cmd *cobra.Command, args []string) {
			// [填空 10] 参数校验：必须刚好 2 个参数（CIDR + IP 地址），否则：
			//   fmt.Println("error: provide a CIDR range and an IP address")
			//   fmt.Println("See 'cidr contains -h' for help and examples")
			//   os.Exit(1)
			if len(args) == 0 {
				// TODO: 占位代码，完成填空 10 之后把它删掉
				fmt.Println("error: 还没有做参数校验（填空 10）")
				os.Exit(1)
			}

			// [填空 11] 解析第 1 个参数：
			//   network, err := core.ParseCIDR(args[0])
			//   出错时 fmt.Printf("error: %s\n", err) 加帮助提示，然后 os.Exit(1)
			network, err := core.ParseCIDR("0.0.0.0/0") // TODO: 改成 args[0]（填空 11）
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// [填空 12] 解析第 2 个参数：
			//   ip := net.ParseIP(args[1])
			//   ip == nil 时打印 "error: invalid IP address: <用户输入>" 加帮助提示，然后 os.Exit(1)
			ip := net.ParseIP("0.0.0.0") // TODO: 改成 args[1]（填空 12）

			// [填空 13] 打印结果：fmt.Println(contains(network, ip))
			_ = contains(network, ip) // TODO: 改成把它打印出来（填空 13）
		},
	}
)

func init() {
	// [填空 14] rootCmd.AddCommand(containsCmd)
}

// contains 只是把 core.ContainsAddress 包了一层。
func contains(network *net.IPNet, ip net.IP) bool {
	return core.ContainsAddress(network, ip)
}

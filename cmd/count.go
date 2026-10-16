package cmd

import (
	"fmt"
	"math/big"
	"net"
	"os"

	"github.com/learn-go/cidr/internal/core"
	"github.com/spf13/cobra"
)

const (
	countExample = "# Return the count of all addresses within a given IPv4 CIDR range\n" +
		"cidr count 10.0.0.0/16\n" +
		"\n" +
		"# Return the count of all addresses within a given IPv6 CIDR range\n" +
		"cidr count 2001:db8:1234:1a00::/106"
)

var (
	countCmd = &cobra.Command{
		// [填空 3] 把这个子命令描述清楚（注意三个字段各自会出现在哪里）：
		//   Use:     "count"          // 命令行里输入的名字，会出现在父命令的帮助里
		//   Short:   "..."            // 一句话说明，同样出现在父命令的帮助里
		//   Example: countExample     // 用上面已经写好的常量，会出现在 `cidr count -h` 里
		Run: func(cmd *cobra.Command, args []string) {
			// [填空 4] 参数校验：必须刚好 1 个参数（一个 CIDR），不是就报错退出：
			//   fmt.Println("error: provide a CIDR range")
			//   fmt.Println("See 'cidr count -h' for help and examples")
			//   os.Exit(1)
			if len(args) == 0 {
				// TODO: 占位代码，完成填空 4 之后把它删掉
				fmt.Println("error: 还没有做参数校验（填空 4）")
				os.Exit(1)
			}

			// [填空 5] 用 core.ParseCIDR 解析用户输入的网段（现在是写死的字符串）：
			//   network, err := core.ParseCIDR(args[0])
			//   出错时打印 "error: invalid CIDR range: <用户输入>" 和帮助提示，然后 os.Exit(1)
			network, err := core.ParseCIDR("0.0.0.0/0") // TODO: 改成 args[0]（填空 5）
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// [填空 6] 把地址总数打印出来：fmt.Println(count(network))
			_ = count(network) // TODO: 改成把它打印出来（填空 6）
		},
	}
)

func init() {
	// [填空 7] 把 countCmd 挂到根命令下，`cidr count ...` 才会生效：
	//   rootCmd.AddCommand(countCmd)
}

// count 是 cmd 层到 core 层的一层薄封装。
// cmd 层只负责“解析命令行 + 打印”，具体算法都放在 internal/core 里。
func count(network *net.IPNet) *big.Int {
	// [填空 8] 提示：调用 core.GetAddressCount(network)
	return big.NewInt(0)
}

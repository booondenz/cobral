package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/learn-go/cidr/internal/core"
	"github.com/spf13/cobra"
)

const (
	overlapsExample = "# Check whether 2 IPv4 CIDR ranges overlap\n" +
		"cidr overlaps 10.0.0.0/16 10.0.14.0/22\n" +
		"\n" +
		"# Check whether 2 IPv6 CIDR ranges overlap\n" +
		"cidr overlaps 2001:db8:1111:2222:1::/80 2001:db8:1111:2222:1:1::/96"
)

var overlapsCmd = &cobra.Command{
	// [填空 15] Use / Short / Example
	Run: func(cmd *cobra.Command, args []string) {
		// [填空 16] 参数校验：必须刚好 2 个参数（两个 CIDR），否则：
		//   fmt.Println("error: provide 2 CIDR ranges")
		//   fmt.Println("See 'cidr overlaps -h' for help and examples")
		//   os.Exit(1)
		if len(args) == 0 {
			// TODO: 占位代码，完成填空 16 之后把它删掉
			fmt.Println("error: 还没有做参数校验（填空 16）")
			os.Exit(1)
		}

		// [填空 17] 依次解析两个 CIDR。注意第二个参数要用 `:=` 重新赋值 err，
		//           因为 network2 是新变量，err 是已存在的变量（Go 允许这样写）。
		//   network1, err := core.ParseCIDR(args[0])
		//   network2, err := core.ParseCIDR(args[1])
		//   任意一个出错都要 fmt.Printf("error: %s\n", err) 加帮助提示，然后 os.Exit(1)
		network1, err := core.ParseCIDR("0.0.0.0/0") // TODO: 改成 args[0]（填空 17）
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		network2, err := core.ParseCIDR("0.0.0.0/0") // TODO: 改成 args[1]（填空 17）
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// [填空 18] 打印结果：fmt.Println(overlaps(network1, network2))
		_ = overlaps(network1, network2) // TODO: 改成把它打印出来（填空 18）
	},
}

func init() {
	// [填空 19] rootCmd.AddCommand(overlapsCmd)
}

// overlaps 只是把 core.Overlaps 包了一层。
func overlaps(network1, network2 *net.IPNet) bool {
	return core.Overlaps(network1, network2)
}

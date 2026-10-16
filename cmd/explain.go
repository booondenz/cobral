package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
	"github.com/learn-go/cidr/internal/core"
	"github.com/learn-go/cidr/internal/helper"
	"github.com/spf13/cobra"
)

const (
	explainExample = "# Explain the details of a given IPv4 CIDR range\n" +
		"cidr explain 10.1.0.0/16\n" +
		"\n" +
		"# Explain the details of a given IPv6 CIDR range\n" +
		"cidr explain 2001:db8:1234:1a00::/106"
)

var (
	explainCmd = &cobra.Command{
		// [填空 20] Use / Short / Example
		Run: func(cmd *cobra.Command, args []string) {
			// [填空 21] 参数校验：必须刚好 1 个参数（一个 CIDR），否则：
			//   fmt.Println("error: provide a CIDR range")
			//   fmt.Println("See 'cidr explain -h' for help and examples")
			//   os.Exit(1)
			if len(args) == 0 {
				// TODO: 占位代码，完成填空 21 之后把它删掉
				fmt.Println("error: 还没有做参数校验（填空 21）")
				os.Exit(1)
			}

			// [填空 22] 解析参数并输出详情：
			//   network, err := core.ParseCIDR(args[0])
			//   出错时 fmt.Printf("error: %s\n", err) 加帮助提示，然后 os.Exit(1)
			//   然后 details := getNetworkDetails(network)，再 explain(details)
			network, err := core.ParseCIDR("0.0.0.0/0") // TODO: 改成 args[0]（填空 22）
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			_ = getNetworkDetails(network) // TODO: 拿到 details 后调用 explain(details)（填空 22）
		},
	}
)

func init() {
	// [填空 23] rootCmd.AddCommand(explainCmd)
}

// networkDetailsToDisplay 是「要打印的东西」的中间结构。
// 先把数据整理到这个结构里，再统一打印，是个很好用的习惯。
type networkDetailsToDisplay struct {
	IsIPV4Network              bool
	IsIPV6Network              bool
	BroadcastAddress           string
	BroadcastAddressHasError   bool
	Netmask                    net.IP
	PrefixLength               int
	BaseAddress                net.IP
	Count                      string
	HostCount                  string
	UsableAddressRangeHasError bool
	FirstUsableIPAddress       string
	LastUsableIPAddress        string
}

func getNetworkDetails(network *net.IPNet) *networkDetailsToDisplay {
	details := &networkDetailsToDisplay{}

	// [填空 24] 判断这是 IPv4 还是 IPv6 网络，给 details.IsIPV4Network 赋值：
	//   提示：helper.IsIPv4Network(network) / helper.IsIPv6Network(network)

	// [填空 25] 取广播地址（IPv6 没有广播地址，会返回 error）：
	//   ipBroadcast, err := core.GetBroadcastAddress(network)
	//   出错时：details.BroadcastAddressHasError = true，并把 err.Error() 存进 details.BroadcastAddress
	//   成功时：把 ipBroadcast.String() 存进 details.BroadcastAddress

	// ---- 下面这一段已经写好了，注意它是怎么把 core / helper 组合起来的 ----

	// 取子网掩码和前缀长度
	netmask := core.GetNetmask(network)
	details.Netmask = core.NetMaskToIPAddress(netmask)
	details.PrefixLength = core.GetPrefixLength(details.Netmask)

	// 取网络的基础地址（网络地址）
	details.BaseAddress = core.GetBaseAddress(network)

	// 总地址数，用 helper.FormatNumber 加上千分位逗号
	count := core.GetAddressCount(network)
	details.Count = helper.FormatNumber(count.String())

	// 可用主机地址数
	hostCount := core.GetHostAddressCount(network)
	details.HostCount = helper.FormatNumber(hostCount.String())

	// [填空 26] 取第一个和最后一个可用地址：
	//   firstUsableIP, err := core.GetFirstUsableIPAddress(network)  // 出错则 details.UsableAddressRangeHasError = true
	//                                                               // 否则 details.FirstUsableIPAddress = firstUsableIP.String()
	//   lastUsableIP, err := core.GetLastUsableIPAddress(network)    // 同样处理 LastUsableIPAddress

	return details
}

//nolint:goconst
func explain(details *networkDetailsToDisplay) {
	var lengthIndicator string

	// [填空 27] 打印基础地址，标题用蓝色：
	//   fmt.Printf(color.BlueString("Base Address:\t\t ")+"%s\n", details.BaseAddress)

	// [填空 28] 打印可用地址范围：
	//   正常时：蓝色标题 "Usable Address Range:\t " + "%s to %s (%s)\n"
	//           （依次是 FirstUsableIPAddress、LastUsableIPAddress、HostCount）
	//   出错时：红色标题 + "unable to calculate usable address range"
	//   提示：color.RedString(...)

	if !details.BroadcastAddressHasError && details.IsIPV4Network {
		fmt.Printf(color.BlueString("Broadcast Address:\t ")+"%s\n", details.BroadcastAddress)
	} else if details.BroadcastAddressHasError && details.IsIPV4Network {
		fmt.Printf(color.RedString("Broadcast Address:\t ")+"%s\n", details.BroadcastAddress)
	}

	fmt.Printf(color.BlueString("Addresses:\t\t ")+"%s\n", details.Count)

	if details.PrefixLength > 1 {
		lengthIndicator = "bits"
	} else {
		lengthIndicator = "bit"
	}

	fmt.Printf(color.BlueString("Netmask:\t\t ")+"%s (/%d %s)\n", details.Netmask, details.PrefixLength, lengthIndicator)
}

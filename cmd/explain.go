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
		Use:     "explain",
		Short:   "Provides information about a CIDR range",
		Example: explainExample,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("error: provide a CIDR range")
				fmt.Println("See 'cidr explain -h' for help and examples")
				os.Exit(1)
			}
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				fmt.Printf("error: %s\n", err)
				fmt.Println("See 'cidr explain -h' for help and examples")
				os.Exit(1)
			}
			details := getNetworkDetails(network)
			explain(details)
		},
	}
)

func init() {
	rootCmd.AddCommand(explainCmd)
}

// networkDetailsToDisplay 是要打印的内容的中间结构：先把数据整理进来，再统一打印。
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

// getNetworkDetails 把一个网段的各种信息整理成 networkDetailsToDisplay。
func getNetworkDetails(network *net.IPNet) *networkDetailsToDisplay {
	details := &networkDetailsToDisplay{}

	// 判断这是 IPv4 还是 IPv6 网络。
	if helper.IsIPv4Network(network) {
		details.IsIPV4Network = true
	} else if helper.IsIPv6Network(network) {
		details.IsIPV6Network = true
	}

	// 取广播地址（IPv6 会返回错误）。错误信息存下来，等打印时再显示。
	ipBroadcast, err := core.GetBroadcastAddress(network)
	if err != nil {
		details.BroadcastAddressHasError = true
		details.BroadcastAddress = err.Error()
	} else {
		details.BroadcastAddress = ipBroadcast.String()
	}

	// 取子网掩码和前缀长度，掩码会以人类可读的形式打印。
	netmask := core.GetNetmask(network)
	details.Netmask = core.NetMaskToIPAddress(netmask)
	details.PrefixLength = core.GetPrefixLength(details.Netmask)

	// 网络的基础地址（网段起点）。
	details.BaseAddress = core.GetBaseAddress(network)

	// 总地址数，格式化成带千分位的字符串。
	count := core.GetAddressCount(network)
	details.Count = helper.FormatNumber(count.String())

	// 可分配的主机地址数。
	hostCount := core.GetHostAddressCount(network)
	details.HostCount = helper.FormatNumber(hostCount.String())

	// 取第一个和最后一个可用地址，任一出错就标记 UsableAddressRangeHasError。
	firstUsableIP, err := core.GetFirstUsableIPAddress(network)
	if err != nil {
		details.UsableAddressRangeHasError = true
	} else {
		details.FirstUsableIPAddress = firstUsableIP.String()
	}

	lastUsableIP, err := core.GetLastUsableIPAddress(network)
	if err != nil {
		details.UsableAddressRangeHasError = true
	} else {
		details.LastUsableIPAddress = lastUsableIP.String()
	}

	return details
}

//nolint:goconst
func explain(details *networkDetailsToDisplay) {
	var lengthIndicator string

	fmt.Printf(color.BlueString("Base Address:\t\t ")+"%s\n", details.BaseAddress)

	if !details.UsableAddressRangeHasError {
		fmt.Printf(color.BlueString("Usable Address Range:\t ")+"%s to %s (%s)\n", details.FirstUsableIPAddress, details.LastUsableIPAddress, details.HostCount)
	} else {
		fmt.Printf(color.RedString("Usable Address Range:\t ")+"%s\n", "unable to calculate usable address range")
	}

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

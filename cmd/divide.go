package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/fatih/color"
	"github.com/learn-go/cidr/internal/core"
	"github.com/learn-go/cidr/internal/helper"
	"github.com/spf13/cobra"
)

const (
	divideExample = "# Divides the given CIDR range into N distinct networks *Truncates output to 50\n" +
		"$ cidr divide 10.0.0.0/16 9\n" +
		"10.0.0.0/20\n" +
		"10.0.16.0/20\n" +
		"10.0.32.0/20\n" +
		"10.0.48.0/20\n" +
		"10.0.64.0/20\n" +
		"10.0.80.0/20\n" +
		"10.0.96.0/20\n" +
		"10.0.112.0/20\n" +
		"10.0.128.0/20\n"
)

// divideCmd 用的是 cobra 的“推荐写法”，和前面几个命令都不一样：
//   - Args：交给 cobra 做参数个数校验（比在 Run 里手写 if len(args) != 2 更规范）
//   - PreRunE：在 RunE 之前做参数内容校验，返回 error 就会中止
//   - RunE：返回 error 而不是自己 os.Exit(1)，把错误交给 rootCmd.Execute() 统一处理
//
// 注意命名里多出来的那个 E，就是 Error 的意思。
var divideCmd = &cobra.Command{
	// [填空 29] 补全这个命令的字段：
	//   Use:     "divide"
	//   Short:   "Divides the given CIDR range into N distinct networks"
	//   Args:    cobra.MinimumNArgs(2)
	//   PreRunE: validateDivideArguments
	//   RunE:    executeDivide
	Example: divideExample,
}

func init() {
	rootCmd.AddCommand(divideCmd)
}

// validateDivideArguments 在真正干活之前，先把参数格式检查一遍。
// 这个函数已经写好了，读一读它是怎么用返回 error 来表达“参数不合法”的。
func validateDivideArguments(cmd *cobra.Command, args []string) error {
	// 确认第 1 个参数是合法的 CIDR
	_, _, err := net.ParseCIDR(args[0])
	if err != nil {
		return fmt.Errorf("invalid network: %s", args[0])
	}

	// 确认第 2 个参数是合法的整数
	_, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid divisor: %s", args[1])
	}

	return nil
}

func executeDivide(cmd *cobra.Command, args []string) error {
	// [填空 30] 实现它，一共 5 步：
	//   1) network, err := core.ParseCIDR(args[0])，出错返回 fmt.Errorf("invalid network: %s", args[0])
	//   2) maskSize, _ := network.Mask.Size()
	//      如果 (helper.IsIPv4Network(network) && maskSize == 32) || maskSize >= 128，
	//      返回 fmt.Errorf("invalid network mask size: %s", args[0])
	//   3) divisor, err := strconv.ParseInt(args[1], 10, 64)，出错返回 fmt.Errorf("invalid divisor: %s", args[1])
	//   4) networks, err := core.DivideCIDR(network, divisor)，出错直接 return err
	//   5) printNetworkPartitions(networks) 然后 return nil
	network, err := core.ParseCIDR("10.0.0.0/16") // TODO: 改用 args[0]（填空 30）
	if err != nil {
		return err
	}
	_ = helper.IsIPv4Network(network) // TODO: 占位代码，完成填空 30 之后删掉
	fmt.Println("TODO: divide 还没实现（填空 30）")
	return nil
}

// printNetworkPartitions 打印切分结果：最多 50 行，中间用 "......" 省略。
// 这个函数也已经写好了，可以直接读。
func printNetworkPartitions(networks []net.IPNet) {
	const truncateLimit = 50
	networkCount := len(networks)

	if networkCount <= truncateLimit {
		for _, network := range networks {
			fmt.Println(network.String())
		}
	} else {
		for i := 0; i < truncateLimit/2; i++ {
			fmt.Println(networks[i].String())
		}
		fmt.Println(color.BlueString("......"))
		for i := networkCount - truncateLimit/2; i < networkCount; i++ {
			fmt.Println(networks[i].String())
		}
	}
}

# 参考答案

> 用法：先自己写，卡住时来看。每个空都给了「填在哪」和代码。
> 所有代码都来自原项目 [bschaatsbergen/cidr](https://github.com/bschaatsbergen/cidr)，导入路径改成了 `github.com/learn-go/cidr`。
> 建议填完之后再看，或者只看某一关的答案。

填空总览：共 **39** 个，难度从「一行 metadata」到「CIDR 位运算」。

---

## 第 1 关：命令树（填空 1-3, 7, 9, 14, 15, 19, 20, 23, 29）

### 填空 1 — `cmd/root.go`：根命令的名字和描述

```go
	rootCmd = &cobra.Command{
		Use:   "cidr",
		Short: "cidr - CLI to perform various actions on CIDR ranges",
		Version: version,
```

`Use` 的第一个词就是命令名，`Short` 会出现在父命令的 `--help` 里。根命令没有父命令，所以 `Short` 显示在 `cidr --help` 的第一行。

### 填空 2 — `cmd/root.go`：不带子命令时打印帮助

```go
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
```

不加 `Run` 的话，用户敲 `cidr` 时 cobra 只会返回 `flag.ErrHelp`；加了才能在帮助里同时看到 `Usage:` 和所有子命令。

### 填空 3 — `cmd/count.go`：子命令的三个描述字段

```go
	countCmd = &cobra.Command{
		Use:     "count",
		Short:   "Return the count of all addresses in a given CIDR range",
		Example: countExample,
		Run: func(cmd *cobra.Command, args []string) {
```

### 填空 7 — `cmd/count.go`：注册到根命令

```go
func init() {
	rootCmd.AddCommand(countCmd)
}
```

`init()` 在 `main()` 之前自动执行，cobra 项目全靠它把命令树拼起来。

### 填空 9 — `cmd/contains.go`

```go
		Use:     "contains",
		Short:   "Checks whether an IP address belongs to a CIDR range",
		Example: containsExample,
```

### 填空 14 — `cmd/contains.go`

```go
func init() {
	rootCmd.AddCommand(containsCmd)
}
```

### 填空 15 — `cmd/overlaps.go`

```go
var overlapsCmd = &cobra.Command{
	Use:     "overlaps",
	Short:   "Checks if a CIDR range overlaps with another CIDR range",
	Example: overlapsExample,
```

### 填空 19 — `cmd/overlaps.go`

```go
func init() {
	rootCmd.AddCommand(overlapsCmd)
}
```

### 填空 20 — `cmd/explain.go`

```go
		Use:     "explain",
		Short:   "Provides information about a CIDR range",
		Example: explainExample,
```

### 填空 23 — `cmd/explain.go`

```go
func init() {
	rootCmd.AddCommand(explainCmd)
}
```

### 填空 29 — `cmd/divide.go`：cobra 推荐写法

```go
var divideCmd = &cobra.Command{
	Use:     "divide",
	Short:   "Divides the given CIDR range into N distinct networks",
	Args:    cobra.MinimumNArgs(2),
	Example: divideExample,
	PreRunE: validateDivideArguments,
	RunE:    executeDivide,
}
```

区别记住三句话：

- `Args` 交给 cobra 校验参数个数，比自己写 `if len(args) != 2` 规范；
- `PreRunE` 在干活前校验参数**内容**，返回 error 就中止，`RunE` 不会被执行；
- `RunE` 返回 error 而不是 `os.Exit(1)`，错误由 `rootCmd.Execute()` 统一处理。

---

## 第 2 关：core 基础（填空 31, 32）

### 填空 31 — `internal/core/core.go`：解析 CIDR

```go
func ParseCIDR(network string) (*net.IPNet, error) {
	_, ip, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	return ip, err
}
```

`net.ParseCIDR` 返回三个值：`(IP, *IPNet, error)`，那个 `IP` 是原始输入（可能带主机位，比如 `10.0.14.5/16`），我们要的是中间的 `*IPNet`（已经按掩码规整过的网络）。所以用 `_` 丢掉第一个。

### 填空 32 — `internal/core/core.go`：地址总数

```go
func GetAddressCount(network *net.IPNet) *big.Int {
	prefixLen, bits := network.Mask.Size()

	// Handle specific cases for IPv4 prefix lengths.
	if network.IP.To4() != nil {
		switch prefixLen {
		case 32:
			return big.NewInt(1)
		case 31:
			return big.NewInt(2)
		}
	}

	return new(big.Int).Lsh(big.NewInt(1), uint(bits-prefixLen))
}
```

要点：

- IPv6 地址数量是 2 的 128 次方量级，`int64` 装不下，所以用 `math/big`。`Lsh(x, n)` 就是左移，等价于乘 2 的 n 次方；
- `/32` 是单个主机地址，`/31` 是点对点链路（只有两个地址，没有网络地址和广播地址），这两种是特例，不能用 2^n − 2 那套公式；
- `network.IP.To4() != nil` 是判断「这是不是 IPv4」的惯用写法（IPv6 地址调用 `To4()` 返回 nil）。

### 填空 39 — `internal/helper/contains.go`

```go
func ContainsInt(ints []int, specifiedInt int) bool {
	for _, i := range ints {
		if i == specifiedInt {
			return true
		}
	}
	return false
}
```

Go 1.21+ 的标准库有 `slices.Contains`，但这里手写一遍正好熟悉 `range` 和 `for` 里 `return` 的行为。

---

## 第 3 关：`count` 跑通（填空 4, 5, 6, 8）

### 填空 4 — 参数个数校验

```go
			if len(args) != 1 {
				fmt.Println("error: provide a CIDR range")
				fmt.Println("See 'cidr count -h' for help and examples")
				os.Exit(1)
			}
```

（记得把 `if len(args) == 0 { ... }` 那段占位代码删掉。）

### 填空 5 — 解析用户输入的网段

```go
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				fmt.Printf("error: invalid CIDR range: %s\n", args[0])
				fmt.Println("See 'cidr count -h' for help and examples")
				os.Exit(1)
			}
```

如果这里报 `undefined: core`，就在文件顶部的 import 里加：

```go
	"github.com/learn-go/cidr/internal/core"
```

### 填空 6 — 打印结果

```go
			fmt.Println(count(network))
```

### 填空 8 — cmd 层调用 core 层

```go
func count(network *net.IPNet) *big.Int {
	count := core.GetAddressCount(network)
	return count
}
```

这层「薄封装」看着多余，但它是分层的关键：`cmd` 负责解析命令行和打印，`core` 负责算。想在别处复用算法时，直接调 `core` 就行，不用碰 cobra。

---

## 第 4 关：`contains` 跑通（填空 10, 11, 12, 13）

### 填空 10 — 两个参数

```go
			if len(args) != 2 {
				fmt.Println("error: provide a CIDR range and an IP address")
				fmt.Println("See 'cidr contains -h' for help and examples")
				os.Exit(1)
			}
```

### 填空 11 — 解析 CIDR（第 1 个参数）

```go
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				fmt.Printf("error: %s\n", err)
				fmt.Println("See 'cidr contains -h' for help and examples")
				os.Exit(1)
			}
```

### 填空 12 — 解析 IP（第 2 个参数）

```go
			ip := net.ParseIP(args[1])
			if ip == nil {
				fmt.Printf("error: invalid IP address: %s\n", args[1])
				fmt.Println("See 'cidr contains -h' for help and examples")
				os.Exit(1)
			}
```

`net.ParseIP` 解析失败时返回 `nil` 而不是 error，所以判断方式是不等于 nil。

### 填空 13 — 打印结果

```go
			fmt.Println(contains(network, ip))
```

`fmt.Println` 打印 bool 就是 `true` / `false`，正好是期望的输出格式。

---

## 第 5 关：`overlaps` 跑通（填空 16, 17, 18, 34）

### 填空 16 — 两个 CIDR

```go
		if len(args) != 2 {
			fmt.Println("error: provide 2 CIDR ranges")
			fmt.Println("See 'cidr overlaps -h' for help and examples")
			os.Exit(1)
		}
```

### 填空 17 — 依次解析两个网段

```go
		network1, err := core.ParseCIDR(args[0])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			fmt.Println("See 'cidr overlaps -h' for help and examples")
			os.Exit(1)
		}
		network2, err := core.ParseCIDR(args[1])
		if err != nil {
			fmt.Printf("error: %s\n", err)
			fmt.Println("See 'cidr overlaps -h' for help and examples")
			os.Exit(1)
		}
```

第二行用 `:=` 不会报「重复声明」，因为 `network2` 是新变量，Go 允许「至少有一个新变量」的短声明。

### 填空 18 — 打印结果

```go
		fmt.Println(overlaps(network1, network2))
```

### 填空 34 — `internal/core/core.go`：判断重叠

```go
func Overlaps(network1, network2 *net.IPNet) bool {
	return network1.Contains(network2.IP) || network2.Contains(network1.IP)
}
```

判断两个网段是否重叠，不需要比较首尾地址：只要**一方的网络地址落在另一方内**就是重叠（两个网段要么完全包含，要么完全不相交，不存在部分交叉）。

---

## 第 6 关：`explain` 跑通（填空 21, 22, 24, 25, 26, 27, 28, 33, 35, 36）

### 填空 21 — 参数个数校验

```go
			if len(args) != 1 {
				fmt.Println("error: provide a CIDR range")
				fmt.Println("See 'cidr explain -h' for help and examples")
				os.Exit(1)
			}
```

### 填空 22 — 解析参数并输出详情

```go
			network, err := core.ParseCIDR(args[0])
			if err != nil {
				fmt.Printf("error: %s\n", err)
				fmt.Println("See 'cidr explain -h' for help and examples")
				os.Exit(1)
			}
			details := getNetworkDetails(network)
			explain(details)
```

这里体现了 `explain.go` 的组织方式：`getNetworkDetails` 只负责「把数据凑齐」，`explain` 只负责「打印」，职责分开后两边都好改。

### 填空 24 — 判断 IPv4 / IPv6

```go
	if helper.IsIPv4Network(network) {
		details.IsIPV4Network = true
	} else if helper.IsIPv6Network(network) {
		details.IsIPV6Network = true
	}
```

为什么不用 `else`？因为「既不是 IPv4 也不是 IPv6」理论上不该发生，显式判断能在出问题时留下线索（两个字段都是 false）。

### 填空 25 — 取广播地址（IPv6 会报错）

```go
	ipBroadcast, err := core.GetBroadcastAddress(network)
	if err != nil {
		details.BroadcastAddressHasError = true
		details.BroadcastAddress = err.Error()
	} else {
		details.BroadcastAddress = ipBroadcast.String()
	}
```

注意这里**没有直接退出**：IPv6 没有广播地址是正常现象，不是致命错误，所以只记一个 `HasError` 标记，等打印时用红字提示。什么时候「返回错误」、什么时候「记个标记继续跑」，是设计要拿捏的地方。

### 填空 26 — 首 / 末可用地址

```go
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
```

### 填空 27 — 打印基础地址

```go
	fmt.Printf(color.BlueString("Base Address:\t\t ")+"%s\n", details.BaseAddress)
```

`color.BlueString("...")` 返回的是**带 ANSI 转义色的字符串**，所以把它和 `"%s\n"` 拼起来传给 `fmt.Printf`。注意 `\t\t` 用来对齐——终端里对不齐就别纠结，原版就是这样。

### 填空 28 — 打印可用地址范围

```go
	if !details.UsableAddressRangeHasError {
		fmt.Printf(color.BlueString("Usable Address Range:\t ")+"%s to %s (%s)\n", details.FirstUsableIPAddress, details.LastUsableIPAddress, details.HostCount)
	} else {
		fmt.Printf(color.RedString("Usable Address Range:\t ")+"%s\n", "unable to calculate usable address range")
	}
```

### 填空 33 — `internal/core/core.go`：可用主机数

```go
func GetHostAddressCount(network *net.IPNet) *big.Int {
	prefixLen, bits := network.Mask.Size()

	if network.IP.To4() != nil {
		switch prefixLen {
		case 32:
			return big.NewInt(1)
		case 31:
			return big.NewInt(2)
		}
	}

	totalAddresses := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits-prefixLen)), nil)
	return totalAddresses.Sub(totalAddresses, big.NewInt(2))
}
```

`Exp(x, y, nil)` 是 x 的 y 次方（第三个参数是取模，nil 表示不取模）。
普通网段里「可用主机数 = 总地址数 − 2」（扣掉网络地址和广播地址），但 `/31` 和 `/32` 没有这个概念，所以是特例：`/31` 两个地址都可用，`/32` 就一个地址。

### 填空 35 — `internal/core/core.go`：第一个可用地址

```go
func GetFirstUsableIPAddress(network *net.IPNet) (net.IP, error) {
	// If it's an IPv6 network
	if network.IP.To4() == nil {
		ones, bits := network.Mask.Size()
		if ones == bits {
			return nil, errors.New(IPv6NetworkHasNoFirstUsableAddressError)
		}

		firstIP := make(net.IP, len(network.IP))
		copy(firstIP, network.IP)

		return firstIP, nil
	}

	// If it's an IPv4 network, first handle edge cases
	switch ones, _ := network.Mask.Size(); ones {
	case 32:
		return nil, errors.New(IPv4NetworkHasNoFirstUsableAddressError)
	case 31:
		firstIP := make(net.IP, len(network.IP))
		copy(firstIP, network.IP)
		return firstIP, nil
	default:
		ip := make(net.IP, len(network.IP))
		copy(ip, network.IP)
		ip[3]++ // Add 1 to the last octet

		return ip, nil
	}
}
```

三个容易踩的点：

1. **IPv6 和 IPv4 的规则不一样**：IPv6 没有广播地址，所以网络地址本身就是第一个可用地址；IPv4 要从网络地址 +1 开始。
2. **`net.IP` 是 `[]byte`（切片）**，`make` + `copy` 是在复制一份，否则直接改 `network.IP` 会污染调用者的数据。
3. `switch ones, _ := network.Mask.Size(); ones {` 是「在 switch 里初始化变量」的写法，`ones` 只在 switch 内可见。
4. `ip[3]++` 只对 IPv4 成立（IPv4 是 4 字节，最后一段是第 4 个字节）。注意这里没处理「最后一段溢出」——原版就是这样，因为默认分支下网络地址的最后一段必然小于 255。

### 填空 36 — `internal/core/core.go`：最后一个可用地址

```go
func GetLastUsableIPAddress(network *net.IPNet) (net.IP, error) {
	if network.IP.To4() == nil {
		ones, bits := network.Mask.Size()
		if ones == bits {
			return nil, errors.New(IPv6NetworkHasNoLastUsableAddressError)
		}

		lastIP := make(net.IP, len(network.IP))
		copy(lastIP, network.IP)
		for i := range lastIP {
			lastIP[i] |= ^network.Mask[i]
		}

		return lastIP, nil
	}

	switch ones, _ := network.Mask.Size(); ones {
	case 32:
		return nil, errors.New(IPv4NetworkHasNoLastUsableAddressError)
	case 31:
		lastIP := make(net.IP, len(network.IP))
		copy(lastIP, network.IP)
		lastIP[3] |= 1
		return lastIP, nil
	default:
		ip := make(net.IP, len(network.IP))
		for i := range ip {
			ip[i] = network.IP[i] | ^network.Mask[i]
		}
		ip[3]--

		return ip, nil
	}
}
```

`ip[i] |= ^network.Mask[i]` 是位运算的核心：把「掩码为 0 的那些位」（也就是主机位）全部置 1，得到的就是广播地址；再 −1 就是最后一个可用地址。`^` 在 Go 里是**按位取反**。

想彻底搞懂，拿 `172.16.18.0/31` 手算一遍：掩码是 `255.255.255.254`，网络位 31 位，主机位 1 位，所以只有两个地址，第一个可用是 `172.16.18.0`，最后一个可用是 `172.16.18.1`。

---

## 第 7 关：`divide` 跑通（填空 30, 37, 38）

### 填空 30 — `cmd/divide.go`：`executeDivide`

```go
func executeDivide(cmd *cobra.Command, args []string) error {
	network, err := core.ParseCIDR(args[0])
	if err != nil {
		return fmt.Errorf("invalid network: %s", args[0])
	}

	maskSize, _ := network.Mask.Size()
	if (helper.IsIPv4Network(network) && maskSize == 32) || maskSize >= 128 {
		return fmt.Errorf("invalid network mask size: %s", args[0])
	}

	divisor, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid divisor: %s", args[1])
	}

	networks, err := core.DivideCIDR(network, divisor)
	if err != nil {
		return err
	}

	printNetworkPartitions(networks)
	return nil
}
```

对比 `count.go`：这里全程 `return fmt.Errorf(...)`，一点没有 `os.Exit`。同一个项目里两种风格都有，正好体会一下哪种更舒服。

### 填空 37 — `internal/core/core.go`：算出新掩码

```go
func GetMaskWithDivisor(divisor int64, addressCount *big.Int, IPv4 bool) (net.IPMask, error) {
	div := big.NewInt(divisor)
	if addressCount.Cmp(div) == -1 || div.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("cannot divide %d addresses into %d divisions", addressCount, div)
	}

	addressPartition := new(big.Int).Div(addressCount, div)
	two := big.NewInt(2)
	exponent := big.NewInt(0)
	for two.Cmp(addressPartition) <= 0 {
		two.Lsh(two, 1)
		exponent.Add(exponent, big.NewInt(1))
	}
	subnetPrefix := int(exponent.Int64())
	bits := net.IPv6len * 8
	if IPv4 {
		bits = net.IPv4len * 8
		if subnetPrefix > 30 {
			return nil, fmt.Errorf("address Space is insufficient for %d subnets", div)
		}
	}
	if subnetPrefix > 126 {
		return nil, fmt.Errorf("address Space is insufficient for %d subnets", div)
	}
	return net.CIDRMask(bits-subnetPrefix, bits), nil
}
```

思路：每个子网要装 `addressCount / divisor` 个地址，那就找一个最小的 2 的幂次 `2^exponent` 不小于它——这个 `exponent` 就是新掩码多出来的位数。
`big.Int` 的 `Cmp` 返回 −1/0/1（小于/等于/大于），别用 `==` 比较。
上限：IPv4 最多 30 位、IPv6 最多 126 位，再分下去就没地址可用了。

### 填空 38 — `internal/core/core.go`：切分网段

```go
func DivideCIDR(network *net.IPNet, divisor int64) ([]net.IPNet, error) {
	isIPv4 := helper.IsIPv4Network(network)

	addressCount := GetAddressCount(network)
	newSubnetMask, err := GetMaskWithDivisor(divisor, addressCount, isIPv4)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	networks := make([]net.IPNet, divisor)
	nextAddress := new(net.IPNet)
	nextAddress.IP = network.IP
	nextAddress.Mask = newSubnetMask
	subnetSize := GetAddressCount(nextAddress)
	for i := int64(0); i < divisor; i++ {
		networks[i] = *nextAddress
		ipAsInt := new(big.Int).SetBytes(nextAddress.IP)
		nextAddress.IP = new(big.Int).Add(ipAsInt, subnetSize).Bytes()
	}
	return networks, nil
}
```

关键技巧：**把 IP 当整数来加**。`big.Int.SetBytes(ip)` 把 4/16 字节的 IP 变成一个大整数，加上子网大小后再 `Bytes()` 转回字节切片。
`networks[i] = *nextAddress` 是「解引用赋值」，把结构体的值拷贝进切片（所以后面改 `nextAddress` 不影响已存的元素）。

注意 `make([]net.IPNet, divisor)` 的 `divisor` 是 `int64`——Go 允许用任何整数类型做长度（只要不是负数），这一点和 C 不一样。

---

## 通关验收

```powershell
powershell -ExecutionPolicy Bypass -File scripts\check.ps1
```

期望看到：

```
  22 passed, 0 failed
```

有任何 FAIL 都会打印 `expected` / `actual` 的逐行对比，照着改就行。

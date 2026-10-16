package core

import (
	"errors"
	"fmt"
	"math/big"
	"net"

	"github.com/learn-go/cidr/internal/helper"
)

// 说明：core 包不打印任何东西，只负责计算。所有输出都由 cmd 包完成。

// ParseCIDR parses the given CIDR notation string and returns the corresponding IP network.
func ParseCIDR(network string) (*net.IPNet, error) {
	// [填空 31] 提示：net.ParseCIDR 返回 3 个值 (IP, *IPNet, error)，
	//           把中间那个 *net.IPNet 和 error 返回出去。用 `_` 忽略第一个返回值。
	return nil, fmt.Errorf("ParseCIDR 还没实现（填空 31）")
}

// GetAddressCount returns the total number of addresses in the given IP network.
// It accounts for both IPv4 and IPv6 networks, and handles specific cases for certain prefix lengths.
func GetAddressCount(network *net.IPNet) *big.Int {
	// [填空 32] 提示：
	//   prefixLen, bits := network.Mask.Size()
	//   IPv4 的 /32 只有 1 个地址，/31 有 2 个地址，这两种情况要单独返回
	//   其余情况就是 2 的 (bits - prefixLen) 次方，用 big.Int 来算：
	//     new(big.Int).Lsh(big.NewInt(1), uint(bits-prefixLen))
	return big.NewInt(0)
}

// GetHostAddressCount returns the number of distinct host addresses in the given IP network.
// It considers the network type (IPv4 or IPv6) and handles edge cases for specific prefix lengths.
// The result excludes the network address and the broadcast address, if applicable.
func GetHostAddressCount(network *net.IPNet) *big.Int {
	// [填空 33] 提示：和填空 32 一样先处理 IPv4 的 /32、/31 特例，
	//           然后 2^(bits-prefixLen) 再减去 2（网络地址和广播地址）。
	return big.NewInt(0)
}

// ContainsAddress checks if the given IP network contains the specified IP address.
// It returns true if the address is within the network, otherwise false.
func ContainsAddress(network *net.IPNet, ip net.IP) bool {
	return network.Contains(ip)
}

// Overlaps checks if there is an overlap between two IP networks.
// It returns true if there is any overlap, otherwise false.
func Overlaps(network1, network2 *net.IPNet) bool {
	// [填空 34] 提示：只要任一网络的网络地址落在对方网络里，就算有重叠。
	//           用 network1.Contains(network2.IP) 和 network2.Contains(network1.IP)。
	return false
}

// GetNetmask retrieves the netmask associated with the provided IP network.
func GetNetmask(network *net.IPNet) net.IPMask {
	return network.Mask
}

// NetMaskToIPAddress converts a netmask represented as a sequence of bytes
// to its corresponding IP address representation.
func NetMaskToIPAddress(netmask net.IPMask) net.IP {
	return net.IP(netmask)
}

// GetPrefixLength returns the prefix length from the given netmask.
func GetPrefixLength(netmask net.IP) int {
	ones, _ := net.IPMask(netmask).Size()
	return ones
}

// GetBaseAddress returns the base address of the given IP network.
func GetBaseAddress(network *net.IPNet) net.IP {
	return network.IP
}

// GetFirstUsableIPAddress returns the first usable IP address in the given IP network.
func GetFirstUsableIPAddress(network *net.IPNet) (net.IP, error) {
	// [填空 35] 提示：
	//   - IPv6（network.IP.To4() == nil）：如果 ones == bits（也就是 /128），
	//     返回 errors.New(IPv6NetworkHasNoFirstUsableAddressError)；否则复制 network.IP 直接返回
	//     （IPv6 没有广播地址，网络地址本身就算第一个可用地址）
	//   - IPv4：/32 返回 errors.New(IPv4NetworkHasNoFirstUsableAddressError)；
	//     /31 直接返回 network.IP 的副本；其余情况把最后一段 +1
	//   小技巧：net.IP 是 []byte，想改就先用 make + copy 复制一份，别改到调用者的数据。
	return nil, errors.New("GetFirstUsableIPAddress 还没实现（填空 35）")
}

// GetLastUsableIPAddress returns the last usable IP address in the given IP network.
func GetLastUsableIPAddress(network *net.IPNet) (net.IP, error) {
	// [填空 36] 提示：
	//   - IPv6：/128 返回 errors.New(IPv6NetworkHasNoLastUsableAddressError)；
	//     其余情况在网络地址上把「网络位之后的所有位」置 1：
	//       for i := range lastIP { lastIP[i] |= ^network.Mask[i] }
	//   - IPv4：/32 返回 errors.New(IPv4NetworkHasNoLastUsableAddressError)；
	//     /31 把最后一位翻转（lastIP[3] |= 1）；其余情况先算出广播地址再减 1
	return nil, errors.New("GetLastUsableIPAddress 还没实现（填空 36）")
}

// GetBroadcastAddress returns the broadcast address of the given IPv4 network, or an error if the IP network is IPv6.
// 这个函数已经写好了：注意「网络地址 | 掩码取反」就等于广播地址，
// 以及 helper.ContainsInt 是怎么用来处理 /31、/32 这两种没有广播地址的特例的。
func GetBroadcastAddress(network *net.IPNet) (net.IP, error) {
	if network.IP.To4() == nil {
		// IPv6 networks do not have broadcast addresses.
		return nil, errors.New(IPv6HasNoBroadcastAddressError)
	}

	// Handle edge case for /31 and /32 networks as they have no broadcast address.
	if prefixLen, _ := network.Mask.Size(); helper.ContainsInt([]int{31, 32}, prefixLen) {
		return nil, errors.New(IPv4HasNoBroadcastAddressError)
	}

	ip := make(net.IP, len(network.IP))
	for i := range ip {
		ip[i] = network.IP[i] | ^network.Mask[i]
	}

	return ip, nil
}

// GetMaskWithDivisor calculates the subnet mask for the given divisor and address count.
func GetMaskWithDivisor(divisor int64, addressCount *big.Int, IPv4 bool) (net.IPMask, error) {
	// [填空 37] 这是本项目的“硬骨头”，可以放到最后再做。提示：
	//   1) 用 big.Int 计算每个子网的大小：addressPartition = addressCount / divisor
	//      如果 addressCount < divisor 或 divisor == 0，返回
	//      fmt.Errorf("cannot divide %d addresses into %d divisions", addressCount, divisor)
	//   2) 找到最小的 2 的幂次（即主机位位数）使得 2^exponent >= addressPartition
	//      （源码里用循环 two.Lsh(two, 1) 不断翻倍来找）
	//   3) bits 按 IPv4 取 32、IPv6 取 128；IPv4 时 subnetPrefix > 30、
	//      IPv6 时 subnetPrefix > 126 都算地址空间不够，返回
	//      fmt.Errorf("address Space is insufficient for %d subnets", divisor)
	//   4) 最后返回 net.CIDRMask(bits-subnetPrefix, bits)
	return nil, errors.New("GetMaskWithDivisor 还没实现（填空 37）")
}

// DivideCIDR divides the given IP network into the specified number of subnets.
func DivideCIDR(network *net.IPNet, divisor int64) ([]net.IPNet, error) {
	// [填空 38] 提示：
	//   1) isIPv4 := helper.IsIPv4Network(network)
	//   2) 用上面的 GetMaskWithDivisor 算出新掩码
	//   3) 造一个 []net.IPNet，长度 divisor；用一个 nextAddress 变量从 network.IP 开始，
	//      每轮把 nextAddress 存进结果，然后 nextAddress.IP += 子网里的地址数
	//      （用 big.Int 做加法：new(big.Int).Add(new(big.Int).SetBytes(nextAddress.IP), subnetSize).Bytes()）
	return nil, errors.New("DivideCIDR 还没实现（填空 38）")
}

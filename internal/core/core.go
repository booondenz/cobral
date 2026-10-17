// Package core 只负责计算，不打印任何东西；所有输出都由 cmd 包完成。
package core

import (
	"errors"
	"fmt"
	"math/big"
	"net"

	"github.com/learn-go/cidr/internal/helper"
)

// ParseCIDR 解析 CIDR 写法的字符串，返回它对应的 IP 网络（*net.IPNet）。
func ParseCIDR(network string) (*net.IPNet, error) {
	_, ip, err := net.ParseCIDR(network)
	if err != nil {
		return nil, err
	}
	return ip, err
}

// GetAddressCount 返回这个网络里的地址总数。
// IPv4 / IPv6 都适用，并且处理了几种特殊的前缀长度（/31、/32）。
func GetAddressCount(network *net.IPNet) *big.Int {
	prefixLen, bits := network.Mask.Size()

	// IPv4 的这几种前缀长度是特例，不能用下面的通用公式。
	if network.IP.To4() != nil {
		switch prefixLen {
		case 32:
			// /32 只有一个地址。
			return big.NewInt(1)
		case 31:
			// /31 用于点对点链路，有两个地址。
			return big.NewInt(2)
		}
	}

	// 其余情况：2 的 (bits - prefixLen) 次方。
	return new(big.Int).Lsh(big.NewInt(1), uint(bits-prefixLen))
}

// GetHostAddressCount 返回这个网络里「可分配的主机地址」数量。
// 会区分网络类型（IPv4 / IPv6），并处理特殊的前缀长度。
// 结果里扣掉了网络地址和广播地址（如果这个网络有的话）。
func GetHostAddressCount(network *net.IPNet) *big.Int {
	prefixLen, bits := network.Mask.Size()

	// IPv4 的 /32、/31 是特例。
	if network.IP.To4() != nil {
		switch prefixLen {
		case 32:
			// /32 就是单个 IP（比如点对点链路）。
			return big.NewInt(1)
		case 31:
			// /31 是两个 IP（点对点链路）。
			return big.NewInt(2)
		}
	}

	// 其余情况：总地址数减去 2（网络地址和广播地址）。
	totalAddresses := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bits-prefixLen)), nil)
	return totalAddresses.Sub(totalAddresses, big.NewInt(2))
}

// ContainsAddress 判断这个网络是否包含指定的 IP 地址。
// 在范围内返回 true，否则返回 false。
func ContainsAddress(network *net.IPNet, ip net.IP) bool {
	return network.Contains(ip)
}

// Overlaps 判断两个 IP 网络是否有重叠。
// 只要有重叠就返回 true，否则返回 false。
func Overlaps(network1, network2 *net.IPNet) bool {
	return network1.Contains(network2.IP) || network2.Contains(network1.IP)
}

// GetNetmask 取出这个网络对应的子网掩码（net.IPMask 就是一串字节）。
func GetNetmask(network *net.IPNet) net.IPMask {
	return network.Mask
}

// NetMaskToIPAddress 把「字节形式的掩码」转成对应的 IP 地址写法。
// 255.255.0.0 这种看着像 IP 的掩码，本质就是 4 个字节的 net.IPMask。
func NetMaskToIPAddress(netmask net.IPMask) net.IP {
	return net.IP(netmask)
}

// GetPrefixLength 从掩码里算出前缀长度（比如 255.255.0.0 对应 /16 里的 16）。
func GetPrefixLength(netmask net.IP) int {
	ones, _ := net.IPMask(netmask).Size()
	return ones
}

// GetBaseAddress 返回网络地址，也就是整个网段的起点。
func GetBaseAddress(network *net.IPNet) net.IP {
	return network.IP
}

// GetFirstUsableIPAddress 返回这个网络里第一个可以分配给设备的 IP 地址。
func GetFirstUsableIPAddress(network *net.IPNet) (net.IP, error) {
	// 如果这是 IPv6 网络
	if network.IP.To4() == nil {
		ones, bits := network.Mask.Size()
		if ones == bits {
			return nil, errors.New(IPv6NetworkHasNoFirstUsableAddressError)
		}

		// IPv6 没有广播地址，网络地址本身就是第一个可用地址。
		firstIP := make(net.IP, len(network.IP))
		copy(firstIP, network.IP)

		return firstIP, nil
	}

	// 如果这是 IPv4 网络，先把几种特例处理掉
	switch ones, _ := network.Mask.Size(); ones {
	case 32:
		return nil, errors.New(IPv4NetworkHasNoFirstUsableAddressError)
	case 31:
		// /31 只有两个地址，网络地址就是第一个可用地址。
		firstIP := make(net.IP, len(network.IP))
		copy(firstIP, network.IP)
		return firstIP, nil
	default:
		// 网络地址 +1 就是第一个可用地址。
		ip := make(net.IP, len(network.IP))
		copy(ip, network.IP)
		ip[3]++ // 最后一段 +1

		return ip, nil
	}
}

// GetLastUsableIPAddress 返回这个网络里最后一个可以分配给设备的 IP 地址。
func GetLastUsableIPAddress(network *net.IPNet) (net.IP, error) {
	// 如果这是 IPv6 网络
	if network.IP.To4() == nil {
		ones, bits := network.Mask.Size()
		if ones == bits {
			return nil, errors.New(IPv6NetworkHasNoLastUsableAddressError)
		}

		// IPv6 没有广播地址，把主机位全置 1 就是最后一个可用地址。
		lastIP := make(net.IP, len(network.IP))
		copy(lastIP, network.IP)
		for i := range lastIP {
			lastIP[i] |= ^network.Mask[i]
		}

		return lastIP, nil
	}

	// 如果这是 IPv4 网络，先把几种特例处理掉
	switch ones, _ := network.Mask.Size(); ones {
	case 32:
		return nil, errors.New(IPv4NetworkHasNoLastUsableAddressError)
	case 31:
		// /31 只有两个地址，另一个地址就是最后一个可用地址。
		lastIP := make(net.IP, len(network.IP))
		copy(lastIP, network.IP)
		lastIP[3] |= 1 // 把最后一位翻过来，得到另一个地址
		return lastIP, nil
	default:
		// 广播地址 -1 就是最后一个可用地址。
		ip := make(net.IP, len(network.IP))
		for i := range ip {
			ip[i] = network.IP[i] | ^network.Mask[i]
		}
		ip[3]-- // 最后一段 -1

		return ip, nil
	}
}

// GetBroadcastAddress 返回 IPv4 网络的广播地址；如果传进来的是 IPv6 网络，则返回错误。
// 注意：「网络地址 | 掩码取反」就是广播地址；/31、/32 这两种没有广播地址，属于特例。
func GetBroadcastAddress(network *net.IPNet) (net.IP, error) {
	if network.IP.To4() == nil {
		// IPv6 网络没有广播地址这个概念，直接返回错误。
		return nil, errors.New(IPv6HasNoBroadcastAddressError)
	}

	// /31 和 /32 这两种长度的 IPv4 网络也没有广播地址，属于特例。
	if prefixLen, _ := network.Mask.Size(); helper.ContainsInt([]int{31, 32}, prefixLen) {
		return nil, errors.New(IPv4HasNoBroadcastAddressError)
	}

	ip := make(net.IP, len(network.IP))
	for i := range ip {
		ip[i] = network.IP[i] | ^network.Mask[i]
	}

	return ip, nil
}

// GetMaskWithDivisor 根据「要切成几份（divisor）」和「总地址数」，算出每个子网该用哪个掩码。
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

// DivideCIDR 把给定的 IP 网络切成指定数量的子网（divisor 就是份数）。
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

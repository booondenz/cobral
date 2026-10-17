package helper

import "net"

// IsIPv4Network 判断这个网络是不是 IPv4 网络。
// 是就返回 true，否则返回 false。
func IsIPv4Network(network *net.IPNet) bool {
	return network.IP.To4() != nil
}

// IsIPv6Network 判断这个网络是不是 IPv6 网络。
// 是就返回 true，否则返回 false。
func IsIPv6Network(network *net.IPNet) bool {
	return network.IP.To16() != nil
}

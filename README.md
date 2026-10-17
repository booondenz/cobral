# cidr

用 Go + [cobra](https://github.com/spf13/cobra) 写的 CIDR 网段计算工具：数地址、判断 IP 归属、检查网段重叠、解释网段、切分子网。

## 来源与许可

本项目是 **[bschaatsbergen/cidr](https://github.com/bschaatsbergen/cidr)**（MIT，见 `LICENSE`）的学习版：
功能和实现与上游一致，差异只有三处——模块路径换成 `github.com/learn-go/cidr`、注释翻译成中文、修掉了上游 `explain` 里一处 `-h` 提示写错的笔误。

## 快速开始

```powershell
go mod tidy                 # 拉依赖（GOPROXY 已指向 goproxy.cn）
go build -o cidr.exe .      # 编译成 cidr.exe
.\cidr.exe --help
```

想在任何目录下直接敲 `cidr`，就装进 PATH（`C:\Users\TTTT\go\bin` 已在 PATH 里）：

```powershell
go install .                # 生成 C:\Users\TTTT\go\bin\cidr.exe
cidr explain 10.0.0.0/16
```

## 命令一览

| 命令 | 作用 | 示例 |
| --- | --- | --- |
| `count <cidr>` | 网段里的地址总数 | `cidr count 10.0.0.0/16` → `65536` |
| `contains <cidr> <ip>` | IP 是否属于该网段 | `cidr contains 10.0.0.0/16 10.0.14.5` → `true` |
| `overlaps <cidr1> <cidr2>` | 两个网段是否重叠 | `cidr overlaps 10.0.0.0/16 10.0.14.0/22` → `true` |
| `explain <cidr>` | 网段详情（网络地址、可用范围、广播地址、掩码） | `cidr explain 10.1.0.0/16` |
| `divide <cidr> <n>` | 把网段切成 n 个 | `cidr divide 10.0.0.0/16 9` |

```
$ cidr explain 10.1.0.0/16
Base Address:            10.1.0.0
Usable Address Range:    10.1.0.1 to 10.1.255.254 (65,534)
Broadcast Address:       10.1.255.255
Addresses:               65,536
Netmask:                 255.255.0.0 (/16 bits)

$ cidr divide 10.0.0.0/16 9
10.0.0.0/20
10.0.16.0/20
...
10.0.128.0/20
```

IPv4 / IPv6 都支持，也处理了这些容易算错的边界情况：

- `/31`（点对点链路，两个地址）和 `/32`（单个地址）没有网络地址、广播地址的概念；
- IPv6 没有广播地址（`explain` 不会显示这一行）；
- `divide` 每个子网的大小取「不超过 总量÷N 的最大 2 的幂」，所以**不是精确等分**，会有富余；
  输出超过 50 行会截断成「前 25 行 + `......` + 后 25 行」。

## 目录结构

```
main.go                     入口，只调用 cmd.Execute()
cmd/                        命令行层（cobra 相关全在这里，一个命令一个文件）
├── root.go                 根命令、--help 模板、Execute()
├── count.go                cidr count <cidr>          ← Run + 手写参数校验
├── contains.go             cidr contains <cidr> <ip>
├── overlaps.go             cidr overlaps <cidr1> <cidr2>
├── explain.go              cidr explain <cidr>        ← 输出最丰富
└── divide.go               cidr divide <cidr> <n>     ← Args + PreRunE + RunE 写法
internal/core/              算法层，只算不打印
├── core.go                 CIDR 的各种计算
└── const.go                错误信息常量
internal/helper/            小工具：整数切片包含判断、数字加千分位
scripts/check.ps1           自测脚本：跑 22 个用例跟上游输出逐字对比
```

## 一些实现说明（cobra）

一个子命令就是一个 `*cobra.Command`，最常用的字段：

| 字段 | 作用 |
| --- | --- |
| `Use` | 命令名，也是 `--help` 里显示的名字 |
| `Short` | 一句话描述，出现在**父命令**的帮助里 |
| `Example` | 示例，出现在自己的 `--help` 的 Examples 里 |
| `Args` | 参数校验器，如 `cobra.MinimumNArgs(2)`、`cobra.ExactArgs(1)` |
| `Run` | 干活函数，自己处理错误（本项目里 `count`/`contains`/`overlaps`/`explain` 用它 + `os.Exit(1)`） |
| `RunE` | 干活函数，返回 `error`（`divide` 用它，更推荐） |
| `PreRunE` | `RunE` 之前的参数内容校验，返回 error 就中止 |

命令树靠各文件里的 `init()` 拼起来（`rootCmd.AddCommand(xxxCmd)`），
`RunE`/`PreRunE` 返回的 error 会冒泡到 `rootCmd.Execute()`，由 `cmd/root.go` 的 `Execute()` 统一打印并退出。

`--help` / `-h` / `help` 子命令由 cobra 自动生成；`--version` 需要 `Version` 字段非空（本项目在编译期注入）：

```powershell
go run -ldflags "-X github.com/learn-go/cidr/cmd.version=0.1.0" . --version
```

## 自测

```powershell
powershell -ExecutionPolicy Bypass -File scripts\check.ps1
```

脚本会自动 `go build`，然后跑 22 个用例（含 IPv6、`/31`、`/32`、`divide` 截断、各种错误参数）
跟上游原版的输出逐字对比，最后打印 `22 passed, 0 failed`。改完代码跑一遍，就是个简易回归测试。

## 还能继续练的几个点

1. 加持久 flag `--no-color`（`rootCmd.PersistentFlags().BoolVar` + `color.NoColor`）。
2. 给 `explain` 加 `--json`，用 `encoding/json` 输出 `networkDetailsToDisplay` 结构体。
3. 加 `SilenceUsage = true`，让 `divide` 参数报错时不喷一大段 Usage。
4. 给 `internal/core` 写 table-driven 单测，`go test ./...`。
5. 把 `divide` 的"每份取 2 的幂"改成精确等分，并补上对应的测试。

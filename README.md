# cidr 填空练习 —— 用 cobra 写 Go CLI

这是一个「**填空骨架**」：代码逻辑是完整的，但关键地方被挖空了（标着 `[填空 N]`）。
你把 39 个空填完，就得到了一个真的能用的命令行工具。

## 这个架子是哪来的

- 原项目：**[bschaatsbergen/cidr](https://github.com/bschaatsbergen/cidr)**（MIT，见 `LICENSE`）——一个对 IPv4/IPv6 网段做运算的小工具，纯 Go + cobra，全部源码约 30KB。
- 选它的原因：**小**（一共 11 个 Go 文件）、**分层干净**（`cmd/` 只管命令行，`internal/core/` 管算法，`internal/helper/` 管小工具）、**cobra 用法齐全**（子命令、参数校验、`Run`/`RunE`、`Args`/`PreRunE`、自定义 help 模板，全都有）。
- 本项目只是把它做成了填空版：所有 import 路径改成 `github.com/learn-go/cidr`，逻辑与上游一致，只是挖了空。

## 3 分钟跑起来

```powershell
go mod tidy        # 第一次会下载 cobra / fatih-color（GOPROXY 已设为 goproxy.cn）
go build ./...     # 开局就能编译通过
go run . --help    # 空壳状态：看不到任何子命令，这是正常的
go run .           # 看到 Usage 里命令名是空的，因为 [填空 1] 还没填
```

> 提示：本项目 `go run .` 每次都要重新编译，慢的话先 `go build -o cidr.exe .`，之后跑 `.\cidr.exe count 10.0.0.0/16`。

## 目录结构

```
main.go                     入口，只调用 cmd.Execute()
cmd/                        ★ 命令行层（cobra 相关全在这里）
├── root.go                 根命令：命令树、--help 模板、Execute()
├── count.go                cidr count <cidr>
├── contains.go             cidr contains <cidr> <ip>
├── overlaps.go             cidr overlaps <cidr1> <cidr2>
├── explain.go              cidr explain <cidr>     ← 输出最丰富，填空最多
└── divide.go               cidr divide <cidr> <n>  ← cobra「推荐写法」示范
internal/core/              算法层：只算不打印
├── core.go                 CIDR 的各种计算
└── const.go                错误信息常量
internal/helper/            小工具：整数切片包含判断、数字加千分位
scripts/check.ps1           自动验收脚本（拿原版输出来对答案）
ANSWERS.md                  ★ 参考答案（先自己写，卡住再看）
```

## 填空怎么用

1. 全项目搜索 `[填空` 就能找到所有空（VSCode 里 `Ctrl+Shift+F` 搜 `[填空`，或搜 `TODO`）。
2. 每个空上面都写了「要做什么」，多数还给了函数名提示。旁边标 `TODO: 占位代码` 的那几行是**假的实现**，填完记得删掉。
3. **注意 import**：骨架里有一些文件故意没有 `import "github.com/learn-go/cidr/internal/core"`，你用到 `core.` 的时候要自己加进去。`go build ./...` 的报错会明确告诉你缺什么。装了 gopls / goimports 的话会自动帮你补。
4. Go 的两个硬规矩：**import 了就必须用**，**局部变量声明了就必须用**。填到一半编译报错很正常，按报错改就行。
5. 拿不准某个空要填什么？先看 `ANSWERS.md` 的提示部分，再不行看答案。原版代码也在 `ANSWERS.md` 里，抄完请读一遍理解为什么。

## 关卡路线（建议按顺序）

| 关卡 | 目标 | 填空编号 | 验收 |
| --- | --- | --- | --- |
| 0 | 热身：读懂 helper 包 | 39 | `go build ./...` |
| 1 | **命令树搭起来** | 1, 2, 3, 7, 9, 14, 15, 19, 20, 23, 29 | `go run . --help` 能看到 5 个子命令和一句话描述 |
| 2 | core 基础 | 31, 32 | 这一步还看不到效果，先过 |
| 3 | `count` 跑通 | 4, 5, 6, 8 | `go run . count 10.0.0.0/16` → `65536` |
| 4 | `contains` 跑通 | 10, 11, 12, 13 | `go run . contains 10.0.0.0/16 10.0.14.5` → `true` |
| 5 | `overlaps` 跑通 | 16, 17, 18, 34 | `go run . overlaps 10.0.0.0/16 10.0.14.0/22` → `true` |
| 6 | `explain` 跑通 | 21, 22, 24, 25, 26, 27, 28, 33, 35, 36 | 见下面「四条命令的效果」 |
| 7 | `divide` 跑通 | 30, 37, 38 | `go run . divide 10.0.0.0/16 9` |
| 通关 | 22 项自动验收全绿 | 全部 | `powershell -ExecutionPolicy Bypass -File scripts\check.ps1` |

依赖关系：第 3 关要用到第 2 关的 `31`/`32`；第 4 关只需要 `31`；第 5/6/7 关都依赖第 2 关。
第 1 关只填 `Use`/`Short` 就能让 `--help` 好看，`Example` 也可以在后面的关卡里补。

## 四条命令的效果（原版输出，也就是验收标准）

```
$ cidr count 10.0.0.0/16
65536

$ cidr contains 10.0.0.0/16 10.0.14.5
true

$ cidr overlaps 10.0.0.0/16 10.0.14.0/22
true

$ cidr explain 10.1.0.0/16
Base Address:            10.1.0.0
Usable Address Range:    10.1.0.1 to 10.1.255.254 (65,534)
Broadcast Address:       10.1.255.255
Addresses:               65,536
Netmask:                 255.255.0.0 (/16 bits)

$ cidr divide 10.0.0.0/16 9
10.0.0.0/20
10.0.16.0/20
10.0.32.0/20
10.0.48.0/20
10.0.64.0/20
10.0.80.0/20
10.0.96.0/20
10.0.112.0/20
10.0.128.0/20
```

## cobra 关键概念速查

一个子命令就是一个 `*cobra.Command` 结构体，最常填的字段：

| 字段 | 作用 |
| --- | --- |
| `Use` | 命令名（第一个空格前的部分），也是 `--help` 里显示的名字 |
| `Short` | 一句话描述，出现在**父命令**的帮助里 |
| `Long` | 长描述，出现在自己的 `--help` 顶部（本项目没用） |
| `Example` | 示例，出现在自己的 `--help` 的 Examples 里 |
| `Args` | 参数个数/内容的校验器，如 `cobra.MinimumNArgs(2)`、`cobra.ExactArgs(1)` |
| `Run` | 干活函数 `func(cmd, args)`，自己处理错误（本项目手写 `os.Exit(1)`） |
| `RunE` | 干活函数 `func(cmd, args) error`，**推荐** |
| `PreRunE` | 在 `RunE` 之前跑的校验，返回 error 就中止 |
| `PersistentFlags()` | 挂在自己和所有子命令上都生效的 flag |

命令树是怎么长起来的：每个 `xxx.go` 里有一个 `init()`，里面一句 `rootCmd.AddCommand(xxxCmd)`。
所以**新增一个子命令 = 新建一个文件 + 在 init 里 AddCommand**。

`PreRunE`/`RunE` 返回的 error 会一路冒泡到 `rootCmd.Execute()`，由根命令统一打印并退出（`cmd/root.go` 的 `Execute()`），
这就是为什么 `divide` 的写法比 `count` 更规范——`count` 是在 `Run` 里写了一堆 `os.Exit(1)`，测试和复用都很别扭。

`--help` / `-h` / `--version` / `help` 子命令都是 cobra 自动加的（`Version` 字段非空才会有 `--version`），
本项目只关掉了自动生成的 `completion` 子命令，见 `cmd/root.go` 的 `init()`。

## 自动验收

```powershell
powershell -ExecutionPolicy Bypass -File scripts\check.ps1
```

它会自动 `go build`，然后跑 22 个用例，把真实输出和原版输出逐字对比：

```
  PASS  L1  root help lists every subcommand
  FAIL  L3  count 10.0.0.0/16
        output mismatch
            expected : 65536
            actual   : unknown command "count" for "" | ...
  ...
  22 passed, 0 failed
```

注意几个用例是**故意的边界情况**，实现了才算真的懂：
`count 172.16.18.0/31`（点对点链路，2 个地址）、`count 172.16.18.0/32`（单个地址）、
`explain 172.16.18.0/31`（IPv4 的 /31 没有广播地址，要显示红字提示）、
`divide 10.0.0.0/16 60`（输出超过 50 行要截断成「前 25 行 + `......` + 后 25 行」）。

## 进阶练习（填完空之后）

1. **加持久 flag**：给根命令加 `--no-color`，用它把 `color.NoColor` 置为 true，让输出不带颜色。
2. **换成 cobra 的校验器**：把 `count`/`contains` 里手写的 `if len(args) != 1` 改成 `Args: cobra.ExactArgs(1)`，感受一下 cobra 报错信息的不同（也顺便理解为什么 `divide` 的写法更好）。
3. **关掉 Usage**：随便给 `divide` 传个错参数，你会看到它把一大段 Usage + Examples 都喷出来了。在 `rootCmd` 上加 `SilenceUsage = true`（以及与 `SilenceErrors` 的区别），让报错干净一点。
4. **给 explain 加 `--json`**：用 `cmd.Flags().BoolVar` 加一个开关，打开时用 `encoding/json` 输出 `networkDetailsToDisplay` 结构体。这是最贴近真实工作的一项。
5. **写单测**：给 `internal/core` 写 table-driven 测试，`go test ./...`。上游的 `internal/core/core_test.go` 有 15KB 用例，可以直接参考（它用 testify）。
6. **注入版本号**：
   ```powershell
   go run -ldflags "-X github.com/learn-go/cidr/cmd.version=0.1.0" . --version
   ```
   （`Version` 字段已经在 `root.go` 里接好了，这也是一个「编译期注入变量」的常见套路。）
7. **新增一个子命令**：比如 `cidr netmask <cidr>`，只打印掩码和前戳长度，练一遍「新建文件 + init 注册」。

## 想对照原版 / 网络受限

GitHub 直连大概率不通。查原版单个文件可以用 jsDelivr（国内可访问）：

```
https://cdn.jsdelivr.net/gh/bschaatsbergen/cidr@main/cmd/explain.go
```

`GOPROXY` 已经指向 `https://goproxy.cn,direct`，所以 `go get` / `go mod tidy` 是正常的。

## License

`LICENSE` 是原项目的 MIT 协议（Copyright (c) 2023 Bruno Schaatsbergen），本项目是其衍生练习版，同样保留该协议。

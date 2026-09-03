# ev1 格式解码器

`ev1` 是一个命令行工具,用于把 **EV1** 视频解码还原成可播放的 **FLV** 文件。

EV1 并不是一种真正的编码格式 —— 它其实就是一个普通的 FLV 视频,只是把**开头 100 个字节用 `XOR 0xFF` 做了混淆**(目的是隐藏 FLV 文件头签名)。文件其余部分本来就是明文 FLV 数据,所以解码只需要还原开头这段字节并给文件改名即可。

## 工作原理

对找到的每一个 `*.ev1` 文件:

1. 以读写方式打开,读取开头 100 个字节。
2. 将每个字节与 `0xFF` 做 `XOR`,还原出原始的 FLV 文件头。
3. 把还原后的字节就地写回文件。
4. 将 `name.ev1` 改名为 `name.ev1.flv`。

## 项目结构

```text
.
├── main.go            # 程序入口,调用 cmd.Execute()
├── cmd/
│   ├── root.go        # cobra 根命令:ev1
│   └── decoder.go     # decoder 子命令,带 -d/--dir 参数
└── core/
    └── decrypt.go     # 核心单文件解码逻辑(DecryptFile)
```

## 环境要求

- Go 1.27 或更高版本

## 构建

```bash
go build -o ev1 .
```

## 使用方法

递归解码某个目录下的所有 `.ev1` 文件:

```bash
./ev1 decoder -d /path/to/videos
# 或者
./ev1 decoder --dir /path/to/videos
```

`-d/--dir` 为必填参数。命令会遍历整个目录树,对找到的每个 `.ev1` 文件进行解码,并打印汇总信息,例如:

```text
[ok]   /path/to/videos/a.ev1 -> /path/to/videos/a.ev1.flv
[ok]   /path/to/videos/sub/b.ev1 -> /path/to/videos/sub/b.ev1.flv

Done. 2 decoded, 0 failed.
```

查看帮助:

```bash
./ev1 --help
./ev1 decoder --help
```

## 注意事项

- **递归扫描**:子目录同样会被处理。
- **就地修改**:原始的 `x.ev1` 会被修改并改名为 `x.ev1.flv`,处理后原 `.ev1` 文件将不复存在。对同一文件重复运行会导致其损坏。
- **大小写敏感**:只匹配小写的 `.ev1` 扩展名。
- 单个文件解码失败不会中断整体流程;失败项会被打印出来,并计入最终的汇总统计。

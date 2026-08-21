# weir-manning

weir-manning 是明渠均匀流与薄壁堰流量核算的 Go 命令行内核：输入矩形或梯形断面（底宽 b、边坡 m，矩形 m=0）、曼宁糙率 n、底坡 S 与流量 Q，用曼宁公式 Q=(1/n)·A·R^(2/3)·S^(1/2)（R=A/P）对水深做单调二分迭代解出正常水深 yn，并输出断面平均流速 v、弗劳德数 Fr=v/√(g·A/T)（T 为水面宽）与水力半径 R；输入堰宽 b、流量系数 Cd 与堰上水头 H，按矩形薄壁堰 Q=Cd·b·√(2g)·H^(3/2) 输出泄量。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run . uniform example/rect-canal.json     # 矩形渠正常水深：yn≈0.81 m，Fr<1
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```

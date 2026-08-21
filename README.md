# weir-manning

weir-manning 是明渠均匀流与薄壁堰流量核算的 Go 命令行内核：输入矩形或梯形断面（底宽 b、边坡 m，矩形 m=0）、曼宁糙率 n、底坡 S 与流量 Q，用曼宁公式 Q=(1/n)·A·R^(2/3)·S^(1/2)（R=A/P）对水深做单调二分迭代解出正常水深 yn，并输出断面平均流速 v、弗劳德数 Fr=v/√(g·A/T)（T 为水面宽）与水力半径 R；输入堰宽 b、流量系数 Cd 与堰上水头 H，按矩形薄壁堰 Q=Cd·b·√(2g)·H^(3/2) 输出泄量。湿周、过水面积、水面宽、水力半径与弗劳德数中的水深尺度全部来自同一个断面函数，不允许曼宁一套面积、Fr 又用另一套水面宽。边界：n≤0、S≤0、b≤0、Q≤0、H≤0 一律 stderr 明文报错并以非零退出码结束；迭代不收敛、水深括区间达不到目标 Q 同样报错；堰的 H 与渠的 yn 是两种物理量，`compare` 只并排提示、绝不改公式。本仓只做明渠断面几何、曼宁迭代与堰的 H–Q，不做有压管网、水锤或流域汇流。

- 输入：一个 JSON 算例文件（字段见 `main.go` 的 usage）；`example/rect-canal.json`（矩形渠，手算 yn 量级约 0.81 m）、`example/trapz-canal.json`（梯形渠）、`example/rect-weir.json`（薄壁堰）
- 输出：`uniform` 打印 `yn`、`v`、`Fr`（含亚临界/超临界/临界判定）、`R`，并把 yn 代回曼宁收回原 Q 做自检；`profile` 打印正常水深附近的水位–流量曲线；`crit` 打印临界水深（矩形用闭式 yc=(Q²/(g·b²))^(1/3)）与给定底坡下的流态；`weir` 打印 Q，也支持 `-q` 反解水头
- 边界：仅矩形/梯形棱柱断面；堰为无侧收缩基本式，Cd 恒定取 0.62（可显式给出）；同一 Cd 下 H 增大 Q 严格按 3/2 次幂升；底坡加大同一 Q 的 yn 下降；不做非均匀渐变流、不做管道

## 钉死的约定

- **断面几何**：A(y)=y·(b+m·y)、P(y)=b+2y√(1+m²)、T(y)=b+2m·y、R=A/P，矩形即 m=0 的特例，全部走同一函数；验证 A、P、R 随 y 严格单调升、T 不降（`internal/cross`）。
- **正常水深**：对给定 Q，Q(y) 单调升，先用倍增括区间再用二分收敛到相对容差 1e-8，迭代上限 200；n≤0、S≤0、Q≤0 error。
- **弗劳德数**：Fr=v/√(g·A/T)，v=Q/A；矩形临界水深闭式 yc=(Q²/(g·b²))^(1/3)，一般断面用 Fr=1 求根，与比能最小水深互验（`internal/flow`）。
- **交叉规则**：yn 代回曼宁的相对误差 ≤1e-6；底坡加大、Q 不变，yn 必须下降；堰 H 与渠 yn 并排输出只提示（`compare`）。
- **薄壁堰**：Q=Cd·b·√(2g)·H^(3/2)，H≤0、b≤0、Cd≤0 error；H 翻倍 Q 变为 2^(3/2)≈2.828 倍。

## 构建 / 运行 / 测试

```text
go build ./...                       # 编译（纯标准库）
go test ./...                        # 全部测试（cross / flow / weir / main）
go run . uniform example/rect-canal.json     # 矩形渠正常水深：yn≈0.81 m，Fr<1
go run . weir example/rect-weir.json         # 薄壁堰泄量：Q≈0.515 m3/s
go run . compare example/rect-canal.json -b 1.5 -h 0.25
go run . crit example/rect-canal.json
go run . profile example/rect-canal.json
```

命令反例（应 stderr 报错并非零退出）：

```text
go run . uniform example/trapz-canal.json    # 先改成 n=0 再跑，糙率必须为正
go run . weir -b 1.5 -cd 0.62 -h -0.1        # 堰上水头必须为正
```

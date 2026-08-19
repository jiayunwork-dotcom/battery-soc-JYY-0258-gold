# battery-soc — 题354 电池储能系统 SOC/SOH 估计

经典数值方法实现的电池储能系统状态估计（**题源 / gold-state，无 bug**）。
不使用任何机器学习、训练或外部依赖，仅依赖 Go 标准库。

## 功能

- **库仑计数（coulomb counting）**：基于电流积分估计 SOC，单步与整段电流日志两种接口。
- **扩展卡尔曼滤波（EKF）**：标量 EKF，用电压派生 SOC 测量值修正库仑计数预测。
- **健康状态（SOH）**：容量比健康度、循环寿命剩余比例、线性退化模型。

## 模块结构

```
battery-soc/
├── main.go                 # CLI: -input <path> 读取电流日志，输出 SOC；缺省读 stdin
├── internal/coulomb/       # 库仑计数 SOC
├── internal/kalman/        # 标量 EKF SOC 估计
├── internal/soh/           # SOH / 循环寿命 / 退化模型
├── example/current_log.csv # 可运行的电流日志样例
├── Dockerfile
└── .dockerignore
```

## CLI 用法

```bash
# 处理样例电流日志（退出码 0）
go run . -input example/current_log.csv

# 通过 stdin 管道输入
cat example/current_log.csv | go run .

# 自定义容量（Ah）
go run . -input example/current_log.csv -capacity 100
```

电流日志为 CSV，两列 `dt,current`：

- `dt`：步长，单位小时（h）
- `current`：电流，单位安培（A），正值充电、负值放电

## 导出 API

- `coulomb.CurrentSample` / `CoulombSOC` / `CoulombFromLog`
- `kalman.EKF` / `NewEKF` / `(*EKF).Step`
- `soh.SOH` / `soh.CycleLife` / `soh.DegradationModel`

## 构建与测试

```bash
go vet ./...
go build ./...
go test ./...
```

## 许可证

MIT — 见 `LICENSE`。

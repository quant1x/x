# 📊 algorithms - 时间序列趋势治理与交易机会发现引擎

> 一个基于支撑/压力分析的 Go 模块，用于识别主趋势、检测突破、发现交易机会。

[![GoDoc](https://pkg.go.dev/badge/github.com/quant1x/x/algorithms)](https://pkg.go.dev/github.com/quant1x/x/algorithms)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## ️ 功能概览

- ✅ **主趋势识别**：以“全局最值”为中心分片，识别主波峰/主波谷
- ✅ **支撑/压力建模**：基于 `High` 序列找压力，`Low` 序列找支撑
- ✅ **突破检测**：自动识别价格是否突破主压力或跌破主支撑
- ✅ **交易机会发现**：从突破点递归分析，发现近期局部趋势
- ✅ **趋势合规性审查**：支持 `PreserveTrend`（保趋势）和 `FindInflection`（找拐点）

---

## 🧩 核心思想

本模块采用三段式趋势分析框架：

1. **主趋势分析**
    - 在 `High` 序列中找主波峰（压力）
    - 在 `Low` 序列中找主波谷（支撑）

2. **突破检测**
    - 检查价格是否突破主压力或跌破主支撑

3. **递归机会发现**
    - 若发生突破：
        - 从**最早突破点**开始
        - 使用 `PreserveTrend` 模式分析后续数据
        - 发现近期局部趋势（如回调结束、反弹结束）

4. **输出交易信号**
    - 返回 `TradeOpportunity`：包含类型、位置、价格、趋势方向


> 本模块不是简单的“找极值”，而是**趋势治理 + 机会发现**的智能系统。

---

## 📦 安装

```bash
go get github.com/quant1x/x/algorithms
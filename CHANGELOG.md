# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.0.5] - 2024-08-24
### Changed
- 新增测试性c代码
- 新增bool util函数集
- 新增地址、指针运算和转换的实验性功能函数
- 调整部分创建切片的代码, 不指定最大长度

## [0.0.4] - 2024-08-22
### Changed
- 修订布尔类型切片的avx2版本的汇编指令

## [0.0.3] - 2024-08-21
### Changed
- 修订__m256编写的汇编函数的前缀给b32x8
- update changelog

## [0.0.2] - 2024-08-21
### Changed
- 新增一个内存逃逸的测试代码
- 修订代码注释
- 修订go sdk最低版本为1.21.12
- 删除废弃package
- 迁移chart和freetype到gitee.com/quant1x/pkg
- 梳理simd组件包路径
- 调整部分包路径
- 调整mul代码
- 新增mul测试代码
- 新增mul测试代码
- 修订c源代码
- 修订c源代码
- 调整函数签名
- 调整函数签名
- 新增float32切片计算avx2加速
- 删除废弃的代码
- 调整sse2代码包路径
- 约束avo生成plan9汇编代码的关键字
- 新增切片判断泛型函数
- 新增布尔类型切片的and函数
- 新增布尔类型切片的and函数
- 调整部分代码
- update changelog

## [0.0.1] - 2024-07-05
### Changed
- Initial commit
- add: github.com/wcharczuk/go-chart/v2 v2.1.1
- add: github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
- delete: chart _examples
- add: .gitignore
- update go.mod
- update changelog


[Unreleased]: https://gitee.com/quant1x/x-go.git/compare/v0.0.5...HEAD
[0.0.5]: https://gitee.com/quant1x/x-go.git/compare/v0.0.4...v0.0.5
[0.0.4]: https://gitee.com/quant1x/x-go.git/compare/v0.0.3...v0.0.4
[0.0.3]: https://gitee.com/quant1x/x-go.git/compare/v0.0.2...v0.0.3
[0.0.2]: https://gitee.com/quant1x/x-go.git/compare/v0.0.1...v0.0.2

[0.0.1]: https://gitee.com/quant1x/x-go.git/releases/tag/v0.0.1

# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.0.15] - 2025-03-13
### Changed
- 修复部分存在警告信息的代码

## [0.0.14] - 2025-03-13
### Changed
- 获取类型尺寸的另一个实现
- 补充函数注释
- update changelog

## [0.0.13] - 2025-03-13
### Changed
- 删除源文件头部不规范的注释信息
- 消除unsafe.pointer的告警信息
- 预备一个指针类型功能的源文件
- update changelog

## [0.0.12] - 2025-03-12
### Changed
- 新增内存方向的内存页尺寸和类型尺寸的小功能函数
- mmap方式的cache改成泛型
- update changelog

## [0.0.11] - 2025-03-11
### Changed
- 新增跨平台的mmap功能
- update changelog

## [0.0.10] - 2025-03-11
### Changed
- 新增16进制字符串和字节数据互转的单元测试
- update changelog

## [0.0.9] - 2025-03-11
### Changed
- 新增一个sync.pool的封装
- update changelog

## [0.0.8] - 2025-03-11
### Changed
- 删除部分废弃的源文件头部的注释
- 删除部分废弃的源文件头部的注释
- update changelog

## [0.0.7] - 2025-03-11
### Changed
- 封装一层字节切片和字符串互转
- 调整测试代码
- 调整切片类型转换的测试代码
- 调整number泛型, 新增bool类型
- 新增切片强制转换函数
- 调整部分测试代码
- 调整go版本最低支持1.24.0
- 增加一个简易的评分系统
- update changelog

## [0.0.6] - 2024-08-26
### Changed
- 分解输出路径, 预备相对路径的处理方式
- makefile增加清理.o文件
- makefile增加清理.o文件
- makefile 新增float32.c
- 去掉无用的参数n
- 调整参数offset类型位int_t
- 调整参数offset类型为int64_t,int offset goat编译失败, 不能识别int
- 删除废弃的中间汇编文件
- 调整测试代码
- 修复n的指针
- 调整n的观测代码
- 调整n的观测代码
- 更新add plan9汇编代码
- 调整测试代码
- 优化浮点随机数的生成
- 调整加速接口
- 调整方法名称
- 新增f32x8 add函数内不申请内存的
- 调整float32基准测试函数名前缀
- 调整实验性load1单元测试代码

## [0.0.5] - 2024-08-24
### Changed
- 新增测试性c代码
- 新增bool util函数集
- 新增地址、指针运算和转换的实验性功能函数
- 调整部分创建切片的代码, 不指定最大长度
- update changelog

## [0.0.4] - 2024-08-22
### Changed
- 修订布尔类型切片的avx2版本的汇编指令
- update changelog

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


[Unreleased]: https://gitee.com/quant1x/x-go.git/compare/v0.0.15...HEAD
[0.0.15]: https://gitee.com/quant1x/x-go.git/compare/v0.0.14...v0.0.15
[0.0.14]: https://gitee.com/quant1x/x-go.git/compare/v0.0.13...v0.0.14
[0.0.13]: https://gitee.com/quant1x/x-go.git/compare/v0.0.12...v0.0.13
[0.0.12]: https://gitee.com/quant1x/x-go.git/compare/v0.0.11...v0.0.12
[0.0.11]: https://gitee.com/quant1x/x-go.git/compare/v0.0.10...v0.0.11
[0.0.10]: https://gitee.com/quant1x/x-go.git/compare/v0.0.9...v0.0.10
[0.0.9]: https://gitee.com/quant1x/x-go.git/compare/v0.0.8...v0.0.9
[0.0.8]: https://gitee.com/quant1x/x-go.git/compare/v0.0.7...v0.0.8
[0.0.7]: https://gitee.com/quant1x/x-go.git/compare/v0.0.6...v0.0.7
[0.0.6]: https://gitee.com/quant1x/x-go.git/compare/v0.0.5...v0.0.6
[0.0.5]: https://gitee.com/quant1x/x-go.git/compare/v0.0.4...v0.0.5
[0.0.4]: https://gitee.com/quant1x/x-go.git/compare/v0.0.3...v0.0.4
[0.0.3]: https://gitee.com/quant1x/x-go.git/compare/v0.0.2...v0.0.3
[0.0.2]: https://gitee.com/quant1x/x-go.git/compare/v0.0.1...v0.0.2

[0.0.1]: https://gitee.com/quant1x/x-go.git/releases/tag/v0.0.1

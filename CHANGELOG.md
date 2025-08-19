# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.1.1] - 2025-08-19
### Changed
- 新增基础极值检测函数
- 新增从右侧波峰单调递减的检测功能
- 新增新版本的波峰检测, 支持力竭和全局两种模式
- 新增波谷极值的检测方法, 与波峰对称
- 可定制左右两侧的极值检测方式
- 优化部分代码
- 优化部分代码
- 新增波浪检测文档
- 新增信号绘图
- 删除废弃的波浪检测方法
- 明确防止重叠作业

## [0.1.0] - 2025-08-14
### Changed
- go版本要求最低1.25
- update changelog

## [0.0.21] - 2025-08-11
### Changed
- sort imports
- sort imports
- update changelog

## [0.0.20] - 2025-03-16
### Changed
- 剔除关闭全局调度器函数中的终端输出信息
- update changelog

## [0.0.19] - 2025-03-16
### Changed
- 系统底层的组件暂时去掉日志的输出
- 消除Unused告警信息
- 补充测试代码
- update changelog

## [0.0.18] - 2025-03-16
### Changed
- 新增系统信号监控封装
- 新增系统内部的全局上下文
- 增加一个临时的logger, 用于观察调试
- 增加专用于内置功能的任务调度器
- 增加每日定时初始化的功能，第一版
- 增加按照cron定时规则的周期性sync.Once功能
- 增加按照周期性Once组件的测试应用
- 新增gocron的依赖, 版本号2.16.1
- update changelog

## [0.0.17] - 2025-03-15
### Changed
- 将该项目打包成一个控制台简单可以使用的工具, 并计划支持跨平台的守护进程
- 更新依赖库版本
- 引入github.com/GavinClarke0/lockless-generic-ring-buffer作为实时环形队列组件, 并尝试修改为动态的消费模式
- 新增求不小于一个整型的最小2的幂次方函数
- 新增用区分并行和并发的源文件
- 新增获取gid的功能函数
- 补充RTRB代码的注释
- 新增一个RingBuffer的spsc的实现
- 优化无锁环形队列为单生产多消费模式(spmc)
- 单生产单消费, 确定没有问题
- 去掉冗余注释
- 调整RB应用测试代码
- 调整部分代码逻辑
- 继续优化了，RingBuffer支持了单生产多消费模式
- 继续优化了，RingBuffer支持了单生产多消费模式
- 继续优化了，RingBuffer支持了多生产多消费模式
- 继续优化了，RingBuffer支持了多生产单消费模式，存在小概率死循环的问题
- 继续优化了，RingBuffer支持了多生产单消费模式，修复小概率死循环的问题。至此无锁环形队列spsc、spmc、mpsc、mpmc四种使用模式全部支持。
- 修订部分测试代码
- 新增cpu缓存行尺寸的常量
- 更新依赖库版本
- update changelog

## [0.0.16] - 2025-03-13
### Changed
- 新增文件路径判断函数
- 文件系统功能归于fs包
- update changelog

## [0.0.15] - 2025-03-13
### Changed
- 修复部分存在警告信息的代码
- update changelog

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


[Unreleased]: https://gitee.com/quant1x/x-go.git/compare/v0.1.1...HEAD
[0.1.1]: https://gitee.com/quant1x/x-go.git/compare/v0.1.0...v0.1.1
[0.1.0]: https://gitee.com/quant1x/x-go.git/compare/v0.0.21...v0.1.0
[0.0.21]: https://gitee.com/quant1x/x-go.git/compare/v0.0.20...v0.0.21
[0.0.20]: https://gitee.com/quant1x/x-go.git/compare/v0.0.19...v0.0.20
[0.0.19]: https://gitee.com/quant1x/x-go.git/compare/v0.0.18...v0.0.19
[0.0.18]: https://gitee.com/quant1x/x-go.git/compare/v0.0.17...v0.0.18
[0.0.17]: https://gitee.com/quant1x/x-go.git/compare/v0.0.16...v0.0.17
[0.0.16]: https://gitee.com/quant1x/x-go.git/compare/v0.0.15...v0.0.16
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

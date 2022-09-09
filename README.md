# gommon
一些平时项目中使用到的库。

欢迎issues，pr！

注：仓库代码测试用例比较简单，如需生产环境使用请充分测试。

# cache
泛型LRU、LFU、FIFO、ARC、Random、NearlyLRU算法

# cmd
命令执行

# consistenthash
一致性哈希，参考groupcache的实现，进行了一点点修改

# container
泛型容器

# conv 
类型转换

# crypto
加密算法

# env
获取环境变量的工具

# hash
泛型哈希函数

# limiter
限流器

# math
一些数值工具

# pool 
对`sync.Pool`的泛型改造，`channel+select`实现的固定长度pool，分级对象池，以及`[]byte`和`bytes.Buffer`的字节对象池

# qps
基于滑动窗口的QPS统计

# slices
泛型slice工具

# validate
基于函数的参数校验，包含数值、字符串和slice类型校验函数





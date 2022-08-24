# pool

各种对象池

# 泛型sync.Pool
对`sync.Pool`进行泛型包装，添加`ClearFunc()`用于清理回收的对象

# FixPool
基于`channel+select`实现的固定长度缓冲区对象池

# LevelPool
多个等级的pool，一般用于字节池但是对象大小跨度比较大的场景

比如指定5个level，分别为1KB、2KB、4KB、8KB、16KB和更大

`Get(cap)`的时候指定长度，这样就会优先获取小的但是满足条件的对象，更加节约资源
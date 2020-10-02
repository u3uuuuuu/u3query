#简单实现的一个B+树的数据查询系统

生成的测试数据是KeySize，Key，ValueSize， Value... 的二进制文件：test/test1.binary 

程序读取该测试文件后生成B+树，并同时删除原本的二进制文件

B+树每十万条做一次分割存储，

以json序列化的二进制文件后保存到文件中（暂时未作压缩）：data/0-100000.binary


API通过下面链接可以看到

http://localhost:8080/swagger


具体API传参参见swagger内容

每个unit对象表示一个{KeySize,Key,ValueSize, Value}


在我的VPS云服务器上已经部署了一个容器化的服务可以直接使用

http://39.99.221.239:8080/swagger
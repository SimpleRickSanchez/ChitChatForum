# ChitChat Forum

### 都2024年了，居然还有人手写论坛

*参考：go web programming - sau sheong chang*

使用golang gin gorm搭建

数据库原作用的是postgresql，这里用的mysql+redis（redis只用于储存session）

- 对所有下一层的子目录`go mod tidy`
```bash
> find . -maxdepth 1 -type d \( ! -name '.*' \) -exec sh -c 'cd "{}" && go mod tidy' \;
```
- 需要在mysql上手动
```bash
mysql> source /path/to/chitchat_forum/models/setup.sql;
```  
- 需要在workspace目录运行（build之后也是）
```bash
> go run ./main
> ./main/chitchat-forum
```
- 启动后访问`/user/forge`自动创建100个用户，并用这些用户发帖发评论，自行修改创建数量的话需谨慎（根据机器性能最好也同时修改连接池数量和`randomUsers`里面规定的用于限制goroutine并发数量的信号量名为sem的channel），数据库操作写的比较随意，如果一次性创建太多同时连接数太多或goroutine太多，比如创建一百万个用户，可能会爆内存（并且slow query到几十秒才能执行完一次`SELECT COUNT(*)`，可能就是连接池分配的问题）。
```bash
browser> http://localhost:8080/user/forge
```

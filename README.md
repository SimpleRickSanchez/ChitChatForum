# ChitChat Forum

### 都2024年了，居然还有人手写论坛

*参考：go web programming - sau sheong chang*

并不完整，只有CR，没有UD。勉强能看。

使用golang gin gorm搭建

数据库原作用的是postgresql，这里用的mysql+redis（redis只用于储存session）

- 对所有下一层的子目录`go mod tidy`
```bash
> find . -maxdepth 1 -type d \( ! -name '.*' \) -exec sh -c 'cd "{}" && go mod tidy' \;
```
- 需要在workspace目录运行（build之后也是）
```bash
go run ./main
./main/chitchat-forum
```
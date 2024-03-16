# ChitChat Forum

### 都2024年了，居然还有人手写论坛

参考：go web programming - sau sheong chang

并不完整，只有CR，没有UD。勉强能看。

对所有下一层的子目录`go mod tidy`
```bash
> find . -maxdepth 1 -type d \( ! -name '.*' \) -exec sh -c 'cd "{}" && go mod tidy' \;
```
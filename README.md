# replay
a send pkg tool for stress testing version 1 <br>
发包压测工具版本1
- 单进程
- 支持seq修改 //需用户定制
- 支持简单压测报告(耗时\发包数\失败数)

## how to use 
```shell
./replay -h 

Usage of ./replay:
  -c int
        pkg count to send (default 10)
  -f string
        pkg file to send (default "req.bin")
  -o int
        seq offset (default 2000)
  -p int
        pkg cnt per period (default 10000)
  -r string
        remote addr to send pkg to (default "113.60.48.140:80")
  -t int
        period, in ms (default 100)
```

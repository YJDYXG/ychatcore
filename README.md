# ychatcore
后端核心进程
1：一级目录为src,两个一级package分别为monitor_cloud和ychatcore
2：main.go作为参数接收启动函数，ychatcore.go则是实际的核心文件。
3: 之后延伸出的package均放置在二级目录ychatcore下，成为其子目录
4: 编译出的进程可执行文件放在编译时生成的build目录下

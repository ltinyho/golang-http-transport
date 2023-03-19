
# 背景
A 服务提供外网访问能力

B 能连接A，A不能连接B

用户访问A时，希望通过A发送http请求到B

这个有什么办法，除了B连接A用websocket转


# 问题
使用 websocket 可以转，但是需要将http转换成 websocket，每个http方法都要转，比较繁琐，有什么办法不转？

# 解决思路
通过 A/B 之间用 TCP 连接，自己解析 TCP 成 HTTP 来处理

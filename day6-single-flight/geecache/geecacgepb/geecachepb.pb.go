syntax = "proto3";

package geecachepb;
//设置响应请求题消息格式

message Request {
string group = 1;
string key = 2;
}

message Response {
bytes value = 1;
}

service GroupCache {
rpc Get(Request) returns (Response);
}
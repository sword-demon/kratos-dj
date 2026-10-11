package service

import (
	"context"
	"customer/api/verifyCode"
	"regexp"
	"time"

	pb "customer/api/customer"

	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

// GetVerifyCode 获取验证码 - 业务逻辑实现
// 1. 校验号码,保证 telephone 满足手机的格式,使用正则验证
// 2. 请求验证码服务获取验证码, 和验证码服务进行通信
func (s *CustomerService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeReq) (*pb.GetVerifyCodeResp, error) {
	// 1. 校验手机号
	pattern := `^1[3-9]\d{9}$`
	regPattern := regexp.MustCompile(pattern)
	if !regPattern.MatchString(req.Telephone) {
		// 校验错误
		return &pb.GetVerifyCodeResp{
			Code:    1,
			Message: "电话号格式不正确",
		}, nil
	}

	// 2. 通过验证码服务生成验证码 (服务间通信,使用 gRPC)

	// 2.1 连接 gRPC 服务
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("localhost:9000"), // verifyCode 服务的 gRPC 的地址
	)
	if err != nil {
		return &pb.GetVerifyCodeResp{
			Code:    1,
			Message: "验证码服务不可用",
		}, nil
	}
	// 关闭连接
	defer func() {
		_ = conn.Close()
	}()

	// 2.2 发送获取验证码请求 必须有 存根 stub 基于 proto 生成的 client 的代码
	// 需要在 customer 服务中使用 verifyCode 服务中的 .proto 文件,生成客户端代码 stub 存根代码,才可以完成 gRPC 的远程调用
	// 拷贝 verifyCode.proto 到 customer 里
	// 执行 kratos proto client 来生成,只需要 client 的代码
	client := verifyCode.NewVerifyCodeClient(conn)
	reply, err := client.GetVerifyCode(context.Background(), &verifyCode.GetVerifyCodeRequest{
		Length: 6,
		Type:   1,
	})
	if err != nil {
		return &pb.GetVerifyCodeResp{
			Code:    1,
			Message: "验证码获取失败",
		}, err
	}

	return &pb.GetVerifyCodeResp{
		Code:           0,
		VerifyCode:     reply.Code,
		VerifyCodeTime: time.Now().Unix(),
		VerifyCodeLife: 60, // 给前端进行倒计时的
	}, nil
}

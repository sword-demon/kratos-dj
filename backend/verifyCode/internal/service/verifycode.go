package service

import (
	"context"
	"math/rand"
	"strings"
	pb "verifyCode/api/verifyCode"
)

type VerifyCodeService struct {
	pb.UnimplementedVerifyCodeServer
}

func NewVerifyCodeService() *VerifyCodeService {
	return &VerifyCodeService{}
}

func (s *VerifyCodeService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequest) (*pb.GetVerifyCodeReply, error) {
	return &pb.GetVerifyCodeReply{
		Code: RandCode(req.Length, req.Type),
	}, nil
}

func RandCode(length uint32, t pb.Type) string {
	switch t {
	default:
	case pb.Type_DEFAULT:
		fallthrough
	case pb.Type_DIGIT:
		return randCodeNew("0123456789", length, 4)
	case pb.Type_LEFTER:
		return randCodeNew("qwertyuiopasdfghjklzxcvbnm", length, 5)
	case pb.Type_MIXED:
		return randCodeNew("0123456789qwertyuiopasdfghjklzxcvbnm", length, 6)
	}
	return ""
}

func randCode(chars string, length uint32) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// randCodeNew 生成随机的唯一性的码
func randCodeNew(chars string, length uint32, idxBits int) string {

	// 形成掩码 mask
	idxMask := int64(1<<idxBits - 1)
	// 63 位 可以使用多少次
	idxMax := 63 / idxBits

	var result strings.Builder
	// 配合 Grom 预分配,整个过程只有一次内存分配
	result.Grow(int(length))

	// 生成随机字符
	for i, cache, remain := 0, rand.Int63(), idxMax; uint32(i) < length; {
		// 如果使用的剩余次数为 0，则重新获取随机数
		if remain == 0 {
			cache, remain = rand.Int63(), idxMax
		}

		// 利用掩码获取有效部位的随机数位
		if randIndex := int(cache & idxMask); randIndex < len(chars) {
			result.WriteByte(chars[randIndex])
			i++
		}

		// 使用下一组随机位
		// 右移消耗已使用的位，并减少剩余次数
		cache >>= idxBits
		remain--
	}

	// 零拷贝返回字符串
	// 避免了一次内存拷贝
	return result.String()
}

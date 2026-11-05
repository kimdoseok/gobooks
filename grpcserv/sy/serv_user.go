package sy

import (
	"context"
	"grpcserv/pb"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	UserService struct {
		pb.UnimplementedUserServiceServer
		repo *UserRepo
	}
)

func NewUserService(repo *UserRepo) pb.UserServiceServer {
	var userserv pb.UserServiceServer
	userserv = &UserService{repo: repo}
	log.Printf("NewUserService: %+v", userserv)
	return userserv
}

func convTimestamp(ts *time.Time) *timestamppb.Timestamp {
	if ts == nil {
		return nil
	}
	return timestamppb.New(*ts)
}

func (s *UserService) GetUserById(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := s.repo.GetUserById(ctx, in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		Id:        user.ID,
		YusCode:   user.Yus_code,
		YusName:   user.Yus_name,
		YusMemo:   user.Yus_memo,
		YusActive: user.Yus_active,
		CreatedAt: convTimestamp(user.CreatedAt),
		UpdatedAt: convTimestamp(user.UpdatedAt),
	}, nil
}

func (u UserService) GetUserByCode(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := u.repo.GetUserByCode(ctx, in.GetYusCode())
	if err != nil {
		return nil, err
	}
	log.Printf("UserService.GetUserByCode: %+v, %+v", user, err)

	return &pb.UserResponse{
		Id:        user.ID,
		YusCode:   user.Yus_code,
		YusName:   user.Yus_name,
		YusMemo:   user.Yus_memo,
		YusActive: user.Yus_active,
		CreatedAt: convTimestamp(user.CreatedAt),
		UpdatedAt: convTimestamp(user.UpdatedAt),
	}, nil
}
func (u UserService) GetUserByLogin(ctx context.Context, in *pb.UserLoginRequest) (*pb.UserResponse, error) {
	user, err := u.repo.GetUserByLogin(ctx, in.GetYusCode(), in.GetYusPasswd())
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{
		YusCode:   in.GetYusCode(),
		YusName:   in.GetYusPasswd(),
		YusMemo:   user.Yus_memo,
		YusActive: user.Yus_active,
		CreatedAt: convTimestamp(user.CreatedAt),
		UpdatedAt: convTimestamp(user.UpdatedAt),
	}, nil
}
func (u UserService) AddUser(ctx context.Context, in *pb.UserAddRequest) (*pb.UserResultResponse, error) {
	id, err := u.repo.AddUser(ctx, in)
	if err != nil {
		return &pb.UserResultResponse{
			Id:      id,
			Success: false,
		}, err
	}
	return &pb.UserResultResponse{
		Id:      id,
		Success: true,
	}, nil
}

func (u UserService) SetUser(ctx context.Context, in *pb.UserSetRequest) (*pb.UserResultResponse, error) {
	err := u.repo.SetUser(ctx, in)
	if err != nil {
		return &pb.UserResultResponse{
			Success: false,
		}, err
	}
	return &pb.UserResultResponse{
		Success: true,
	}, nil
}

func (u UserService) GetUserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	lpp, err := strconv.ParseInt(os.Getenv("LINE_PER_PAGE"), 10, 64)
	if lpp <= 0 || err != nil {
		lpp = 100 // Default value if not set or invalid
	}
	users, err := u.repo.GetUserList(ctx, in.GetFilter(), in.GetPage(), lpp)
	if err != nil {
		return &pb.UserListResponse{
			Usergroups: []*pb.UserResponse{},
		}, err
	}
	outusers := make([]*pb.UserResponse, len(users))
	for i, user := range users {
		outusers[i] = &pb.UserResponse{
			Id:        user.ID,
			YusCode:   user.Yus_code,
			YusName:   user.Yus_name,
			YusMemo:   user.Yus_memo,
			YusActive: user.Yus_active,
			CreatedAt: convTimestamp(user.CreatedAt),
			UpdatedAt: convTimestamp(user.UpdatedAt),
		}
	}
	return &pb.UserListResponse{
		Usergroups: outusers,
	}, nil
}
func (u UserService) GetUserCount(ctx context.Context, in *pb.UserListRequest) (*pb.UserCountResponse, error) {
	count, err := u.repo.GetUserCount(ctx, in.GetFilter())
	if err != nil {
		return &pb.UserCountResponse{}, err
	}
	return &pb.UserCountResponse{
		Count: count,
	}, nil
}
func (u UserService) DeleteUserById(ctx context.Context, in *pb.UserRequest) (*pb.UserResultResponse, error) {
	err := u.repo.DeleteUserById(ctx, in.GetId())
	if err != nil {
		return &pb.UserResultResponse{
			Success: false,
		}, err
	}
	return &pb.UserResultResponse{
		Success: true,
	}, nil
}
func (u UserService) DeleteUserByCode(ctx context.Context, in *pb.UserRequest) (*pb.UserResultResponse, error) {
	err := u.repo.DeleteUserByCode(ctx, in.GetYusCode())
	if err != nil {
		return &pb.UserResultResponse{
			Success: false,
		}, err
	}
	return &pb.UserResultResponse{
		Success: true,
	}, nil
}

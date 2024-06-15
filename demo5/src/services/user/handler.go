package main

import (
	"context"
	"fmt"
	"select-course/demo5/src/constant/services"
	"select-course/demo5/src/rpc/course"
	"select-course/demo5/src/rpc/user"
	grpc2 "select-course/demo5/src/utils/grpc"
)

type User struct {
	user.UnimplementedUserServiceServer
}

var courseClient course.CourseServiceClient

func (u *User) New() {
	conn := grpc2.Connect(context.Background(), services.CourseRpcServerName)
	courseClient = course.NewCourseServiceClient(conn)

}
func (u *User) GetUserInfo(ctx context.Context, res *user.UserRequest) (*user.UserResponse, error) {
	courseRes, err := courseClient.GetAllCourses(ctx, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println(courseRes)
	return nil, nil
}

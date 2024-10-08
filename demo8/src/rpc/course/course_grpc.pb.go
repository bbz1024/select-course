// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.19.4
// source: course.proto

package course

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	CourseService_GetAllCourses_FullMethodName            = "/course.CourseService/GetAllCourses"
	CourseService_GetMyCourses_FullMethodName             = "/course.CourseService/GetMyCourses"
	CourseService_SelectCourse_FullMethodName             = "/course.CourseService/SelectCourse"
	CourseService_BackCourse_FullMethodName               = "/course.CourseService/BackCourse"
	CourseService_EnQueueCourse_FullMethodName            = "/course.CourseService/EnQueueCourse"
	CourseService_TryTryDeductCourse_FullMethodName       = "/course.CourseService/TryTryDeductCourse"
	CourseService_TryConfirmDeductCourse_FullMethodName   = "/course.CourseService/TryConfirmDeductCourse"
	CourseService_TryCancelDeductCourse_FullMethodName    = "/course.CourseService/TryCancelDeductCourse"
	CourseService_TryTryEnqueueMessage_FullMethodName     = "/course.CourseService/TryTryEnqueueMessage"
	CourseService_TryConfirmEnqueueMessage_FullMethodName = "/course.CourseService/TryConfirmEnqueueMessage"
	CourseService_TryCancelEnqueueMessage_FullMethodName  = "/course.CourseService/TryCancelEnqueueMessage"
)

// CourseServiceClient is the client API for CourseService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CourseServiceClient interface {
	// 获取所有课程列表
	GetAllCourses(ctx context.Context, in *GetAllCoursesRequest, opts ...grpc.CallOption) (*GetAllCoursesResponse, error)
	// 获取指定用户的课程列表（我的课程）
	GetMyCourses(ctx context.Context, in *GetMyCoursesRequest, opts ...grpc.CallOption) (*GetMyCoursesResponse, error)
	// 选课操作
	SelectCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error)
	// 退课操作
	BackCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error)
	// 其中选课操作和退课操作都存在在redis预创建，如何丢入到消息队列进行处理操作，消息队列处理操作为：扣减课程，修改用户课程表，添加用户选课记录操作。
	// 丢入消息队列操作
	EnQueueCourse(ctx context.Context, in *EnQueueCourseRequest, opts ...grpc.CallOption) (*CourseOptResponse, error)
	// 1. 定义TCC事务的Try、Confirm和Cancel方法
	TryTryDeductCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error)
	TryConfirmDeductCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error)
	TryCancelDeductCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error)
	TryTryEnqueueMessage(ctx context.Context, in *EnQueueCourseRequest, opts ...grpc.CallOption) (*CourseOptResponse, error)
	TryConfirmEnqueueMessage(ctx context.Context, in *EnQueueCourseRequest, opts ...grpc.CallOption) (*CourseOptResponse, error)
	TryCancelEnqueueMessage(ctx context.Context, in *EnQueueCourseRequest, opts ...grpc.CallOption) (*CourseOptResponse, error)
}

type courseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCourseServiceClient(cc grpc.ClientConnInterface) CourseServiceClient {
	return &courseServiceClient{cc}
}

func (c *courseServiceClient) GetAllCourses(ctx context.Context, in *GetAllCoursesRequest, opts ...grpc.CallOption) (*GetAllCoursesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAllCoursesResponse)
	err := c.cc.Invoke(ctx, CourseService_GetAllCourses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) GetMyCourses(ctx context.Context, in *GetMyCoursesRequest, opts ...grpc.CallOption) (*GetMyCoursesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetMyCoursesResponse)
	err := c.cc.Invoke(ctx, CourseService_GetMyCourses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) SelectCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseOptResponse)
	err := c.cc.Invoke(ctx, CourseService_SelectCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) BackCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseOptResponse)
	err := c.cc.Invoke(ctx, CourseService_BackCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) EnQueueCourse(ctx context.Context, in *EnQueueCourseRequest, opts ...grpc.CallOption) (*CourseOptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseOptResponse)
	err := c.cc.Invoke(ctx, CourseService_EnQueueCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) TryTryDeductCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseOptResponse)
	err := c.cc.Invoke(ctx, CourseService_TryTryDeductCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) TryConfirmDeductCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseOptResponse)
	err := c.cc.Invoke(ctx, CourseService_TryConfirmDeductCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) TryCancelDeductCourse(ctx context.Context, in *CourseOptRequest, opts ...grpc.CallOption) (*CourseOptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseOptResponse)
	err := c.cc.Invoke(ctx, CourseService_TryCancelDeductCourse_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) TryTryEnqueueMessage(ctx context.Context, in *EnQueueCourseRequest, opts ...grpc.CallOption) (*CourseOptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseOptResponse)
	err := c.cc.Invoke(ctx, CourseService_TryTryEnqueueMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) TryConfirmEnqueueMessage(ctx context.Context, in *EnQueueCourseRequest, opts ...grpc.CallOption) (*CourseOptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseOptResponse)
	err := c.cc.Invoke(ctx, CourseService_TryConfirmEnqueueMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) TryCancelEnqueueMessage(ctx context.Context, in *EnQueueCourseRequest, opts ...grpc.CallOption) (*CourseOptResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseOptResponse)
	err := c.cc.Invoke(ctx, CourseService_TryCancelEnqueueMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CourseServiceServer is the server API for CourseService service.
// All implementations must embed UnimplementedCourseServiceServer
// for forward compatibility
type CourseServiceServer interface {
	// 获取所有课程列表
	GetAllCourses(context.Context, *GetAllCoursesRequest) (*GetAllCoursesResponse, error)
	// 获取指定用户的课程列表（我的课程）
	GetMyCourses(context.Context, *GetMyCoursesRequest) (*GetMyCoursesResponse, error)
	// 选课操作
	SelectCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error)
	// 退课操作
	BackCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error)
	// 其中选课操作和退课操作都存在在redis预创建，如何丢入到消息队列进行处理操作，消息队列处理操作为：扣减课程，修改用户课程表，添加用户选课记录操作。
	// 丢入消息队列操作
	EnQueueCourse(context.Context, *EnQueueCourseRequest) (*CourseOptResponse, error)
	// 1. 定义TCC事务的Try、Confirm和Cancel方法
	TryTryDeductCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error)
	TryConfirmDeductCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error)
	TryCancelDeductCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error)
	TryTryEnqueueMessage(context.Context, *EnQueueCourseRequest) (*CourseOptResponse, error)
	TryConfirmEnqueueMessage(context.Context, *EnQueueCourseRequest) (*CourseOptResponse, error)
	TryCancelEnqueueMessage(context.Context, *EnQueueCourseRequest) (*CourseOptResponse, error)
	mustEmbedUnimplementedCourseServiceServer()
}

// UnimplementedCourseServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCourseServiceServer struct {
}

func (UnimplementedCourseServiceServer) GetAllCourses(context.Context, *GetAllCoursesRequest) (*GetAllCoursesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllCourses not implemented")
}
func (UnimplementedCourseServiceServer) GetMyCourses(context.Context, *GetMyCoursesRequest) (*GetMyCoursesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMyCourses not implemented")
}
func (UnimplementedCourseServiceServer) SelectCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SelectCourse not implemented")
}
func (UnimplementedCourseServiceServer) BackCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BackCourse not implemented")
}
func (UnimplementedCourseServiceServer) EnQueueCourse(context.Context, *EnQueueCourseRequest) (*CourseOptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnQueueCourse not implemented")
}
func (UnimplementedCourseServiceServer) TryTryDeductCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryTryDeductCourse not implemented")
}
func (UnimplementedCourseServiceServer) TryConfirmDeductCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryConfirmDeductCourse not implemented")
}
func (UnimplementedCourseServiceServer) TryCancelDeductCourse(context.Context, *CourseOptRequest) (*CourseOptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryCancelDeductCourse not implemented")
}
func (UnimplementedCourseServiceServer) TryTryEnqueueMessage(context.Context, *EnQueueCourseRequest) (*CourseOptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryTryEnqueueMessage not implemented")
}
func (UnimplementedCourseServiceServer) TryConfirmEnqueueMessage(context.Context, *EnQueueCourseRequest) (*CourseOptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryConfirmEnqueueMessage not implemented")
}
func (UnimplementedCourseServiceServer) TryCancelEnqueueMessage(context.Context, *EnQueueCourseRequest) (*CourseOptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryCancelEnqueueMessage not implemented")
}
func (UnimplementedCourseServiceServer) mustEmbedUnimplementedCourseServiceServer() {}

// UnsafeCourseServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CourseServiceServer will
// result in compilation errors.
type UnsafeCourseServiceServer interface {
	mustEmbedUnimplementedCourseServiceServer()
}

func RegisterCourseServiceServer(s grpc.ServiceRegistrar, srv CourseServiceServer) {
	s.RegisterService(&CourseService_ServiceDesc, srv)
}

func _CourseService_GetAllCourses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllCoursesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).GetAllCourses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_GetAllCourses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).GetAllCourses(ctx, req.(*GetAllCoursesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_GetMyCourses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMyCoursesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).GetMyCourses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_GetMyCourses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).GetMyCourses(ctx, req.(*GetMyCoursesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_SelectCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseOptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).SelectCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_SelectCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).SelectCourse(ctx, req.(*CourseOptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_BackCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseOptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).BackCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_BackCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).BackCourse(ctx, req.(*CourseOptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_EnQueueCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnQueueCourseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).EnQueueCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_EnQueueCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).EnQueueCourse(ctx, req.(*EnQueueCourseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_TryTryDeductCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseOptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).TryTryDeductCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_TryTryDeductCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).TryTryDeductCourse(ctx, req.(*CourseOptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_TryConfirmDeductCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseOptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).TryConfirmDeductCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_TryConfirmDeductCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).TryConfirmDeductCourse(ctx, req.(*CourseOptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_TryCancelDeductCourse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseOptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).TryCancelDeductCourse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_TryCancelDeductCourse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).TryCancelDeductCourse(ctx, req.(*CourseOptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_TryTryEnqueueMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnQueueCourseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).TryTryEnqueueMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_TryTryEnqueueMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).TryTryEnqueueMessage(ctx, req.(*EnQueueCourseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_TryConfirmEnqueueMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnQueueCourseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).TryConfirmEnqueueMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_TryConfirmEnqueueMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).TryConfirmEnqueueMessage(ctx, req.(*EnQueueCourseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_TryCancelEnqueueMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnQueueCourseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).TryCancelEnqueueMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_TryCancelEnqueueMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).TryCancelEnqueueMessage(ctx, req.(*EnQueueCourseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CourseService_ServiceDesc is the grpc.ServiceDesc for CourseService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CourseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "course.CourseService",
	HandlerType: (*CourseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllCourses",
			Handler:    _CourseService_GetAllCourses_Handler,
		},
		{
			MethodName: "GetMyCourses",
			Handler:    _CourseService_GetMyCourses_Handler,
		},
		{
			MethodName: "SelectCourse",
			Handler:    _CourseService_SelectCourse_Handler,
		},
		{
			MethodName: "BackCourse",
			Handler:    _CourseService_BackCourse_Handler,
		},
		{
			MethodName: "EnQueueCourse",
			Handler:    _CourseService_EnQueueCourse_Handler,
		},
		{
			MethodName: "TryTryDeductCourse",
			Handler:    _CourseService_TryTryDeductCourse_Handler,
		},
		{
			MethodName: "TryConfirmDeductCourse",
			Handler:    _CourseService_TryConfirmDeductCourse_Handler,
		},
		{
			MethodName: "TryCancelDeductCourse",
			Handler:    _CourseService_TryCancelDeductCourse_Handler,
		},
		{
			MethodName: "TryTryEnqueueMessage",
			Handler:    _CourseService_TryTryEnqueueMessage_Handler,
		},
		{
			MethodName: "TryConfirmEnqueueMessage",
			Handler:    _CourseService_TryConfirmEnqueueMessage_Handler,
		},
		{
			MethodName: "TryCancelEnqueueMessage",
			Handler:    _CourseService_TryCancelEnqueueMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "course.proto",
}

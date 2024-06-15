package main

import (
	"context"
	"fmt"
	"select-course/demo5/src/rpc/course"
)

type Course struct {
	course.UnimplementedCourseServiceServer
}

func (c Course) GetAllCourses(ctx context.Context, request *course.GetAllCoursesRequest) (*course.GetAllCoursesResponse, error) {
	fmt.Println(1111)
	return nil, nil
}

func (c Course) GetMyCourses(ctx context.Context, request *course.GetMyCoursesRequest) (*course.GetMyCoursesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c Course) SelectCourse(ctx context.Context, request *course.CourseOptRequest) (*course.CourseOptResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c Course) BackCourse(ctx context.Context, request *course.CourseOptRequest) (*course.CourseOptResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c Course) EnQueueCourse(ctx context.Context, request *course.EnQueueCourseRequest) (*course.CourseOptResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c Course) mustEmbedUnimplementedCourseServiceServer() {
	//TODO implement me
	panic("implement me")
}

package test

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
)

const BaseUrl = "http://localhost:8888/api/v1/"

func TestSelectCourse(t *testing.T) {
	// 模拟100个用户同时对课程为2的id进行选课
	users := 100
	courseID := 2

	var wg sync.WaitGroup
	wg.Add(users)
	for i := 0; i < users; i++ {
		go func(i int) {
			defer wg.Done()
			// 模拟用户i进行选课
			SelectCourse(i, courseID)
		}(i)
	}
	wg.Wait()
}
func SelectCourse(userID int, courseID int) {
	// 发送请求
	resp, err := http.Post(BaseUrl+"course/select/", "application/json",
		strings.NewReader(fmt.Sprintf(`{"user_id":%d,"course_id":%d}`, userID, courseID)),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		fmt.Println("用户", userID, "选课成功")
	} else {
		fmt.Println("用户", userID, "选课失败")
	}

}

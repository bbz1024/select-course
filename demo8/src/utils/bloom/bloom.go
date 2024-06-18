package bloom

import (
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"select-course/demo8/src/models"
	"select-course/demo8/src/storage/database"
	"sync"
)

var UserBloom *bloom.BloomFilter
var CourseBloom *bloom.BloomFilter

func InitializeBloom() error {
	// 初始化操作
	CourseBloom = bloom.NewWithEstimates(100000, 0.01)
	UserBloom = bloom.NewWithEstimates(100000, 0.01)
	var wg sync.WaitGroup
	var err error
	wg.Add(2)
	go func() {
		defer wg.Done()
		var userList []*models.User
		if err = database.Client.Find(&userList).Error; err != nil {
			return
		}
		for _, user := range userList {
			// 加载id进入布隆过滤器
			UserBloom.AddString(fmt.Sprintf("%d", user.ID))
		}

	}()
	go func() {
		defer wg.Done()
		var courseList []*models.Course
		if err = database.Client.Find(&courseList).Error; err != nil {
			return
		}
		for _, course := range courseList {
			// 加载id进入布隆过滤器
			CourseBloom.AddString(fmt.Sprintf("%d", course.ID))
		}
	}()
	wg.Wait()

	return err
}

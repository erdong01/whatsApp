package visitor

import (
	"fmt"
	"whatsApp/models"

	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	var visitor models.Visitor
	err := c.BindJSON(&visitor)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 417,
			"msg":  "参数不正确！" + err.Error(),
		})
		return
	}
	err = visitor.Create()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 417,
			"msg":  "数据添加失败！",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "添加成功！",
	})
	return
}
func Count(c *gin.Context) {
	// var userIdStr, b = c.GetQuery("user_id")
	// if b {
	// 	c.JSON(200, gin.H{
	// 		"code": 417,
	// 		"msg":  "用户错误",
	// 	})
	// 	return
	// }
	// userId, err := strconv.Atoi(userIdStr)
	// if err != nil {
	// 	c.JSON(200, gin.H{
	// 		"code": 417,
	// 		"msg":  "用户错误",
	// 	})
	// 	return
	// }
	var visitor = models.Visitor{}
	visitorArr, err := visitor.Find()
	if err != nil {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "暂时没有数据",
		})
		return
	}
	var x []int
	var data []int
	var date = visitorArr[0].RegisterTime.Day()
	var count int
	var pieces []map[string]interface{}
	var threshold = 10
	var gt = -1
	var lte = 0
	var greenLte = 0
	for index, v := range visitorArr {
		NewDate := v.RegisterTime.Day()
		fmt.Println("date", date)
		fmt.Println("NewDate", NewDate)
		if NewDate != date {
			x = append(x, date)
			data = append(data, count)
			date = NewDate
			if count <= threshold {
				if greenLte > 0 {
					pieces = append(pieces, map[string]interface{}{"gt": gt, "lte": gt + greenLte + 1, "color": "green"})
					fmt.Println("green 前", gt)
					gt = gt + greenLte + 1
					fmt.Println("green 后", gt)
				}
				lte++
				greenLte = 0
			}
			if count > threshold {
				if lte > 0 {
					pieces = append(pieces, map[string]interface{}{"gt": gt, "lte": gt + lte - 1, "color": "red"})
					fmt.Println("red 前", gt)
					gt = gt + lte - 1
					fmt.Println("red 后", gt)
				}
				greenLte++
				lte = 0
				fmt.Println(" count > threshold gt", gt)
			}
			count = 0
		}
		count++
		fmt.Println("index >= (len(visitor) - 1)", index, (len(visitorArr) - 1))
		if index == (len(visitorArr) - 1) {
			fmt.Println("gt", gt)
			if count <= threshold {
				lte++
			}
			if count > threshold {
				greenLte++
			}
		}

	}
	x = append(x, date)
	data = append(data, count)
	fmt.Println("greenLte", greenLte)
	if lte > 0 {
		if greenLte > 1 {
			pieces = append(pieces, map[string]interface{}{"gt": gt, "lte": gt + greenLte, "color": "green"})
			fmt.Println("green 前", gt)
			gt = gt + greenLte
			fmt.Println("green 后", gt)
		}
		pieces = append(pieces, map[string]interface{}{"gt": gt, "lte": gt + lte, "color": "red"})
	}

	if greenLte > 0 {
		if lte > 1 {
			pieces = append(pieces, map[string]interface{}{"gt": gt, "lte": gt + lte, "color": "red"})
			fmt.Println("red 前", gt)
			gt = gt + lte
			fmt.Println("red 后", gt)
		}
		pieces = append(pieces, map[string]interface{}{"gt": gt, "lte": gt + greenLte, "color": "green"})
	}

	fmt.Println(visitor)
	c.JSON(200, gin.H{
		"code": 400,
		"msg":  "粉丝手机号不能为空",
		"data": gin.H{
			"x":      x,
			"data":   data,
			"pieces": pieces,
		},
	})
}

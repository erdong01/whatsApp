package m007

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func Login() {
	c := colly.NewCollector(colly.MaxDepth(1), colly.DetectCharset(), colly.Async(true), colly.AllowURLRevisit())
	err := c.Post("https://007.mn/java-api/system/user/login", map[string]string{
		"env_code": "46ce9cd791a5c6bacae06a637f5efc66",
		"password": "Tm587973",
		"sub_pwd":  "",
		"type":     "0",
		"user_id":  "erdong88155745@outlook.com",
	})

	if err != nil {
		fmt.Println("err", err)
	}

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("response received", r.StatusCode)
		fmt.Println("response received", string(r.Body))
	})

}

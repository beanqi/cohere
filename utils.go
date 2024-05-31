/*
 * @Author: Vincent Yang
 * @Date: 2024-04-18 03:50:36
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2024-04-18 03:50:56
 * @FilePath: /cohere2openai/utils.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2024 by Vincent, All Rights Reserved.
 */

package main

import "github.com/gin-gonic/gin"

func isInSlice(str string, list []string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

func stringPtr(s string) *string {
	return &s
}

// SetMessageChan sets the message channel in the context
func SetMessageChan(ch chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("chan", ch)
	}

}

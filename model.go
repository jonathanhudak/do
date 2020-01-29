package main

import "github.com/jinzhu/gorm"

type Entry struct {
	gorm.Model
	Title string
}

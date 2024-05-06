package service

type Event struct {
}

type Evener interface {
	Add() error
	Del(id int64) error
}

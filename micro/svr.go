package main

import "fmt"

type MySvc interface {
	Get(name string) error
	Put(name string) error
}

var _ MySvc = &MySvcImpl{}

type MySvcImpl struct {
}

func (s *MySvcImpl) Get(name string) error {
	fmt.Println("print get from svc")
	return nil
}

func (s *MySvcImpl) Put(name string) error {
	fmt.Println("print put from svc")
	return nil
}

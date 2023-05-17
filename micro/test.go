package main

type Other struct {
	S MySvc `inject:""` //依赖接口，注入不同的实现，达到不同效果
}

// Start func of Other
func (o *Other) Start() {
	o.S.Put("test")
}

package main

type Configuration struct {
	Listen  string            `json:"listen"`
	Program string            `json:"program"`
	EnvMap  map[string]string `json:"env"`
}

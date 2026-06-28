package cmd

import "github.com/jcyamacho/memo/internal/memory"

var service *memory.Service

func SetService(s *memory.Service) {
	service = s
}

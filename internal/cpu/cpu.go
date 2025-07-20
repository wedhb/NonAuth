package cpu

type x86 struct {
    HasAES        bool
    HasPCLMULQDQ  bool
}

type arm64 struct {
    HasAES   bool
    HasPMULL bool
}

type s390x struct {
    HasAES     bool
    HasAESCBC  bool
    HasAESCTR  bool
    HasGHASH   bool
    HasAESGCM  bool
}

var X86 = x86{
    HasAES:       false, // 你可以根据自己的需求检测是否支持
    HasPCLMULQDQ: false,
}

var ARM64 = arm64{
    HasAES:   false,
    HasPMULL: false,
}

var S390X = s390x{
    HasAES:     false,
    HasAESCBC:  false,
    HasAESCTR:  false,
    HasGHASH:   false,
    HasAESGCM:  false,
}

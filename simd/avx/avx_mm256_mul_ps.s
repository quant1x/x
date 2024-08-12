//go:build !noasm && amd64
// AUTO-GENERATED BY GOAT -- DO NOT EDIT

TEXT ·mm256_mul_ps(SB), $0-32
        MOVQ a+0(FP), DI
        MOVQ b+8(FP), SI
        MOVQ c+16(FP), DX
        BYTE $0x55               // push        rbp
        WORD $0x8948; BYTE $0xe5 // mov rbp, rsp
        LONG $0xf8e48348         // and rsp, -8
        LONG $0x0710fcc5         // vmovups     ymm0, ymmword ptr [rdi]
        LONG $0x0659fcc5         // vmulps      ymm0, ymm0, ymmword ptr [rsi]
        LONG $0x0211fcc5         // vmovups     ymmword ptr [rdx], ymm0
        WORD $0x8948; BYTE $0xec // mov rsp, rbp
        BYTE $0x5d               // pop rbp
        WORD $0xf8c5; BYTE $0x77 // vzeroupper
        BYTE $0xc3               // ret

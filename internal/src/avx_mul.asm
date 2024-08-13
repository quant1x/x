	.text
	.intel_syntax noprefix
	.file	"avx_mul.c"
	.globl	mm256_mul_ps1                   # -- Begin function mm256_mul_ps1
	.p2align	4, 0x90
	.type	mm256_mul_ps1,@function
mm256_mul_ps1:                          # @mm256_mul_ps1
# %bb.0:
	push	rbp
	mov	rbp, rsp
	and	rsp, -8
	vmulps	ymm0, ymm1, ymm0
	mov	rsp, rbp
	pop	rbp
	ret
.Lfunc_end0:
	.size	mm256_mul_ps1, .Lfunc_end0-mm256_mul_ps1
                                        # -- End function
	.ident	"Ubuntu clang version 14.0.0-1ubuntu1.1"
	.section	".note.GNU-stack","",@progbits
	.addrsig

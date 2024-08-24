mod=github.com/quant1x/x
arch=amd64
pkg=simd
asm_path=labs/asm
mkfile_path=$(CURDIR)
output_path=$(mkfile_path)/simd

all: hello f32x8 b1x8

hello:
	@echo $(CURDIR)

b1x8:
	go run $(mod)/$(asm_path)/avx2/boolean -out $(output_path)/b8x32_$(arch).s -stubs $(output_path)/b8x32_$(arch).go -pkg $(pkg)

f32x8:
	go run $(mod)/$(asm_path)/avx2/f32x8 -out $(output_path)/f32x8_$(arch).s -stubs $(output_path)/f32x8_$(arch).go -pkg $(pkg)
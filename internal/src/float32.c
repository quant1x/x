#include <immintrin.h>
#include <stdint.h>
#include <stdbool.h>


void* address_seek(void *addr, int_t offset)
{
    return addr + offset;
}


void f32x4_add(float const *a, float const *b, float *c)
{
    __m128 vec1 = _mm_load_ps(a);
    __m128 vec2 = _mm_load_ps(b);
    __m128 res  = _mm_add_ps(vec1, vec2);
    _mm_storeu_ps(c, res);
}

void f32x8_add(float const *a, float const *b, float *c)
{
    __m256 vec1 = _mm256_load_ps(a);
    __m256 vec2 = _mm256_load_ps(b);
    __m256 res  = _mm256_add_ps(vec1, vec2);
    _mm256_storeu_ps(c, res);
}

void MultiplyAndAdd(float* arg1, float* arg2, float* arg3, float* result) {
    __m256 vec1 = _mm256_load_ps(arg1);
    __m256 vec2 = _mm256_load_ps(arg2);
    __m256 vec3 = _mm256_load_ps(arg3);
    __m256 res  = _mm256_fmadd_ps(vec1, vec2, vec3);
    _mm256_storeu_ps(result, res);
}

void float32_add(float *a, float *b, float *c, int64_t n)
{
    for (int i = 0; i < n; i++)
    {
        c[i] = a[i] + b[i];
    }
}
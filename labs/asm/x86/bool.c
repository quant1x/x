#include <immintrin.h>
#include <stdint.h>

__m256i mm256_loadu_epi8(const int8_t *a)
{
    __m256i r = _mm256_loadu_epi8(a);
    return r;
}

package algorithms

import "fmt"

// FindPeaksValleys 查找波峰波谷
func FindPeaksValleys(highList, lowList []float64) (peaks, valleys []int, err error) {
	n := len(highList)
	if n != len(lowList) || n < 3 {
		return nil, nil, fmt.Errorf("输入序列长度不匹配或过短（需要至少3个点）")
	}

	// 使用 int8 节省内存（diff 值仅为 -1,0,1）
	diffHigh := make([]int8, n) // n 个元素，最后一个为 0
	diffLow := make([]int8, n)

	// 1. 计算一阶符号差分（n-1 项），最后一项补 0
	for i := 0; i < n-1; i++ {
		// 内联 compare
		if highList[i+1] > highList[i] {
			diffHigh[i] = 1
		} else if highList[i+1] < highList[i] {
			diffHigh[i] = -1
		} // else 0（默认）

		if lowList[i+1] > lowList[i] {
			diffLow[i] = 1
		} else if lowList[i+1] < lowList[i] {
			diffLow[i] = -1
		}
	}
	// 最后一个差分设为 0（无后续）
	diffHigh[n-1] = 0
	diffLow[n-1] = 0

	// 2. 处理平台（diff == 0 的点），合并 high 和 low 处理
	for i := 0; i < n-1; i++ {
		// 处理 high
		if diffHigh[i] == 0 {
			if i == 0 {
				for j := 1; j < n-1; j++ {
					if diffHigh[j] != 0 {
						diffHigh[i] = diffHigh[j]
						break
					}
				}
			} else if i == n-2 {
				diffHigh[i] = diffHigh[i-1]
			} else {
				diffHigh[i] = diffHigh[i+1]
			}
		}

		// 处理 low
		if diffLow[i] == 0 {
			if i == 0 {
				for j := 1; j < n-1; j++ {
					if diffLow[j] != 0 {
						diffLow[i] = diffLow[j]
						break
					}
				}
			} else if i == n-2 {
				diffLow[i] = diffLow[i-1]
			} else {
				diffLow[i] = diffLow[i+1]
			}
		}
	}

	// 3. 检测波峰波谷：二阶差分
	// 预分配空间（通常极值点不会超过 n/3）
	peaks = make([]int, 0, n/3)
	valleys = make([]int, 0, n/3)

	for i := 0; i < n-1; i++ {
		dHigh := int(diffHigh[i+1]) - int(diffHigh[i])
		dLow := int(diffLow[i+1]) - int(diffLow[i])

		if dHigh == -2 {
			peaks = append(peaks, i+1)
		}
		if dLow == 2 {
			valleys = append(valleys, i+1)
		}
	}

	return peaks, valleys, nil
}

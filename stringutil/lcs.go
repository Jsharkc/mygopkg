package stringutil

func LongestCommonSubsequence(str1, str2 string) int {
	rune1 := []rune(str1)
	rune2 := []rune(str2)

	m, n := len(rune1), len(rune2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i, c1 := range rune1 {
		for j, c2 := range rune2 {
			if c1 == c2 {
				dp[i+1][j+1] = dp[i][j] + 1
			} else {
				dp[i+1][j+1] = max(dp[i][j+1], dp[i+1][j])
			}
		}
	}
	return dp[m][n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

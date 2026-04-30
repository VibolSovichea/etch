package theme

import (
	"math/rand"
	"strings"
)

func ApplyScanlines(output string) string {
	lines := strings.Split(output, "\n")
	for i := range lines {
		if i%2 == 1 {
			lines[i] = "\033[2m" + lines[i] + "\033[22m"
		}
	}
	return strings.Join(lines, "\n")
}

func ApplyNoise(output string, density float64) string {
	if density <= 0 {
		return output
	}
	noiseChars := []rune{'░', '·', '▪', '▫'}
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		runes := []rune(line)
		for j := len(runes) - 1; j >= 0; j-- {
			if runes[j] != ' ' {
				break
			}
			if rand.Float64() < density {
				ch := noiseChars[rand.Intn(len(noiseChars))]
				runes[j] = ch
				lines[i] = "\033[2m" + string(runes) + "\033[22m"
				break
			}
		}
	}
	return strings.Join(lines, "\n")
}

func ApplyEffects(t *Theme, output string) string {
	if t.Effects.Scanline {
		output = ApplyScanlines(output)
	}
	if t.Effects.Noise {
		output = ApplyNoise(output, 0.03)
	}
	return output
}

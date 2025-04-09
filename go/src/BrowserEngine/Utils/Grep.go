package Utils

import (
	"errors"
)

func Grep(pattern, array []byte) (int, error) {
	var (
		/* The variables are:
		sfttab: shift table
		patlen: the length of the pattern
		arrlen: the length of the array of bytes
		*/
		patlen int = len(pattern)
		arrlen int = len(array)
		idx    int = patlen - 1
		jdx    int
	)
	sfttab := make(map[byte]int)

	if patlen == 0 {
		return -(1 << 31), errors.New("[ERROR] Pattern is length of 0.")
	}

	if patlen > arrlen {
		return -(1 << 31), errors.New("[ERROR] Length of the array < Length of the pattern")
	}

	/* Create shift table */
	for i := range patlen {
		sfttab[pattern[i]] = max(1, patlen-i-1)
	}

	for idx < arrlen {
		jdx = 0
		for jdx < patlen && pattern[patlen-jdx-1] == array[idx-jdx] {
			jdx++
		}

		if jdx == patlen {
			return idx - patlen + 1, nil
		} else {
			shift, ok := sfttab[array[idx+jdx]]
			if !ok {
				shift = patlen
			}

			if shift == 0 {
				shift = patlen - 1
			}
			idx += (shift - jdx)
		}
	}
	/* That return means that we didn't found anything */
	return 0, nil
}

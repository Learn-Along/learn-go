package types

type Filter []bool

// // Coalesces a map of bool arrays into an ANDED array of bool such that 
// // say, age > 7 and name = "John" returns only true when both are true
// func (f *Filter) Coalesce() []bool {
// 	d := []bool{}

// 	defaultValue := true
// 	for _, list := range *f {
// 		for i, value := range list {
// 			if i < len(d) {
// 				d[i] = d[i] && value
// 			} else {
// 				d = append(d, defaultValue && value)
// 			}
// 		}

// 		// after first loop, default is false as it means
// 		// that index is missing in initial array, hence AND = false
// 		defaultValue = false
// 	}

// 	return d
// }
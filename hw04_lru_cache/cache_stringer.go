package hw04_lru_cache

import (
	"fmt"
)

// String - наглядное представление значения KeyValue-структуры.
func (kv KeyValue) String() string {
	return fmt.Sprintf("%q->%q", kv.key, kv.value)
}

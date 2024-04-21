package requests

import (
	"fmt"
	"sort"
	"sync"
)

func FindWithDB(app *App, stemRequest []string) map[string][]int {
	fmt.Println("Поиск по базе данных")
	result := make(map[string][]int)

	// Составляем map keyword -> [id1, id2, id3]
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(stemRequest))
	for _, word := range stemRequest {
		go func() {
			localIndex := make([]int, 0)
			for id, _ := range app.Db.Entries {
				keywords := app.Db.Entries[id].Keywords
				for _, keyword := range keywords {
					if keyword == word {
						localIndex = append(localIndex, id)
						break
					}
				}
			}
			mu.Lock()
			result[word] = localIndex
			mu.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()

	return result

}

func FindWithIndex(app *App, stemRequest []string) map[string][]int {
	fmt.Println("Поиск по индекс файлу")
	result := make(map[string][]int)
	for _, word := range stemRequest {
		result[word] = app.Db.Index[word]
	}
	return result
}

func processResult(app *App, searchResult map[string][]int, limit int) []int {
	// Находим пересечение id в полученых списках. Чтобы потом выбрать те айди, которые чаще всего встречались
	intersection := make(map[int]int)
	for keyword := range searchResult {
		for _, id := range searchResult[keyword] {
			intersection[id]++
		}
	}

	// Создаем обратный intersection, чтобы отсортировать по количестве встречаемых раз
	pairs := make([][2]int, 0, len(intersection))
	for k, v := range intersection {
		pairs = append(pairs, [2]int{k, v})
	}

	// Сортируем по количеству раз, если количество раз равны, то сортируем по кол-во слов в комиксе
	// Делаем это из той логики, что чем меньше слов в комиксе, тем каждое слово важнее
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i][1] != pairs[j][1] {
			return pairs[i][1] > pairs[j][1]
		}
		return len(app.Db.Entries[pairs[i][0]].Keywords) < len(app.Db.Entries[pairs[j][0]].Keywords)
	})

	// Получаем список из 10 релевантных комиксов
	resultLen := min(len(pairs), limit)
	keys := make([]int, resultLen)
	for i := 0; i < resultLen; i++ {
		keys[i] = pairs[i][0]
	}
	return keys
}

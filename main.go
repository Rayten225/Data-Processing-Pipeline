package main

import (
	"fmt"
	"time"
)

// Константы для размера буфера и интервала времени
const bufferSize = 5
const flushInterval = 2 * time.Second

func notNegativeFunc(dataSourceInt []int) <-chan int {
	output := make(chan int)
	go func() {
		for _, v := range dataSourceInt {
			if v > 0 {
				output <- v
			}
		}
		close(output)
	}()
	return output
}

func notMultipleThree(nums <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for v := range nums {
			if v%3 != 0 && v != 0 {
				output <- v
			}
		}
	}()
	return output
}

// Стадия буферизации
func bufferStage(input <-chan int) <-chan int {
	output := make(chan int)
	buffer := make([]int, 0, bufferSize)

	go func() {
		ticker := time.NewTicker(flushInterval)
		defer ticker.Stop()
		defer close(output)

		for {
			select {
			case v, ok := <-input:
				if ok {
					buffer = append(buffer, v)
					if len(buffer) == bufferSize {
						for _, item := range buffer {
							output <- item
						}
						buffer = buffer[:0]
					}
				} else {
					// Закрытие входного канала, опустошаем буфер
					for _, item := range buffer {
						output <- item
					}
					return
				}
			case <-ticker.C:
				// Таймер истёк, опустошаем буфер
				if len(buffer) > 0 {
					for _, item := range buffer {
						output <- item
					}
					buffer = buffer[:0]
				}
			}
		}
	}()
	return output
}

func main() {
	var col int
	var startDataSource []int

	fmt.Print("Введите кол-во чисел: ")
	_, err := fmt.Scan(&col)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < col; i++ {
		var num int
		fmt.Printf("Введите %d число: ", i+1)
		_, err := fmt.Scan(&num)
		if err != nil {
			fmt.Println(err)
		} else {
			startDataSource = append(startDataSource, num)
		}
	}

	nums := notNegativeFunc(startDataSource)
	filtered := notMultipleThree(nums)
	buffered := bufferStage(filtered)

	for data := range buffered {
		fmt.Println("Итоговое значение:", data)
	}
}

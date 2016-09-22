//
// Copyright (c) 2016 ЗАО Геликон Про http://www.gelicon.biz
//
// Пакет реализует управление версиями вида 1.2.3-SNAPSHOT
//
package version

import (
	"strings"
	"regexp"
	"strconv"
)

const (
	SnapsotSuffix = "-SNAPSHOT"
)

// Версия представляет собой набор элементов, разделенных точкой и, опционально, имеющих суффикс -SNAPSHOT.
//
// Каждый значимый элемент версии должен начинается с цифры. Если элемент имеет нечисловой суффикс,
// то числовое значение берется до начала этого суффикса.
//
// Версия с суффиком -SNAPSHOT далее будет называться рабочей версией,
// а версия без этого суффикса будет называться релизом.
type Ver struct {
	number   []int
	suffix   []string
	level    int
	snapshot bool
}

// Инициализирует версию заданного уровня. Уровеню 1 соответствует версия вида 1,
// уровню 2 -- 1.0, уровню 3 соответствует версия вида 1.0.0 и т.д.
//
// После инициализации версия всегда становиться рабочей (SNAPSHOT).
func (ver *Ver) Start(level int) *Ver {

	if level < 1 {
		level = 1
	}

	number := make([]int, level)
	suffix := make([]string, level)

	if (level > 0) {
		number[0] = 1
	}

	ver.number = number
	ver.suffix = suffix
	ver.level = level
	ver.snapshot = true

	return ver
}

// Переключение указанного компонента версии (level) на следующее значение и устанавливка
// флага рабочей версии (SNAPSHOT). Все компоненты версии с уровнем выше указанного обнуляются.
//
// При вызове Next(1) производится переключение '1.0.0' в '2.0.0',
// при вызове Next(2) производится переключение '1.2.3' в '1.3.0' и т.д.
func (ver *Ver) Next(level int) *Ver {

	if level <= 0 {
		return ver
	}

	for level > ver.level {
		ver.number = append(ver.number, 0)
		ver.suffix = append(ver.suffix, "")

		copy(ver.number[ver.level + 1:], ver.number[ver.level:])
		copy(ver.suffix[ver.level + 1:], ver.suffix[ver.level:])

		ver.number[ver.level] = 0
		ver.suffix[ver.level] = ""

		ver.level++
	}

	ver.number[level - 1]++

	for i := level; i < ver.level; i++ {
		ver.number[i] = 0
	}

	ver.snapshot = true

	return ver
}

// Отсечение компонентов версии, далее указанного уровня (level).
func (ver *Ver) Trunk(level int) *Ver {

	if level <= 0 {
		ver.number = nil
		ver.suffix = nil
		ver.level = 0
		return ver
	}

	if level > ver.level {
		level = ver.level
	}

	ver.number = ver.number[:level]
	ver.suffix = ver.suffix[:level]

	ver.level = level

	return ver
}

// Переключение самого старшего компонента версии и устанавка флага рабочей версии (SNAPSHOT).
func (ver *Ver) Snapshot() *Ver {

	if ver.level == 0 {
		return ver
	}

	ver.snapshot = true
	ver.number[ver.level - 1]++

	return ver
}

// Снятие флага флага рабочей версии (SNAPSHOT). Если флаг не был установлен, то производится
// переключение самого старшего компонента версии на следующее значение.
func (ver *Ver) Release() *Ver {

	if ver.snapshot {
		ver.snapshot = false
		return ver
	}

	if ver.level == 0 {
		return ver
	}

	ver.number[ver.level - 1]++

	return ver
}

// Формироваие строки версии
func (ver Ver) Render() string {

	parts := make([]string, len(ver.suffix))

	for i := range parts {
		if i < ver.level {
			parts[i] = strconv.Itoa(ver.number[i]) + ver.suffix[i]
		} else {
			parts[i] = ver.suffix[i]
		}
	}

	s := strings.Join(parts, ".")

	if ver.snapshot {
		s = s + SnapsotSuffix
	}

	return s
}

// Получение количества компонент версии
func (ver Ver) Level() int {
	return ver.level
}

// Инициализация версии по строковому значению
func (ver *Ver) Parse(value string) *Ver {

	trimmedValue := strings.TrimSuffix(value, SnapsotSuffix)

	snapshot := len(value) != len(trimmedValue)

	parts := strings.Split(trimmedValue, ".")

	if value == "" {
		parts = nil
	}

	number := make([]int, len(parts))
	suffix := make([]string, len(parts))

	re := regexp.MustCompile(`^(\d+)(.*)`)

	num := true

	level := 0

	for k, v := range parts {

		if !num {
			number[k] = -1
			suffix[k] = v
			continue
		}

		match := re.FindStringSubmatch(v)

		if match == nil || len(match[1]) > 7 {
			num = false
			number[k] = -1
			suffix[k] = v
			continue
		}

		number[k], _ = strconv.Atoi(match[1])
		suffix[k] = match[2]

		level++
	}

	ver.number = number
	ver.suffix = suffix
	ver.level = level
	ver.snapshot = snapshot

	return ver
}

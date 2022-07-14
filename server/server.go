// Учебная программа для отработки взаимодействия Сервера и Клиента из модуля 35 SkillFactory
//
// Разработайте сетевую службу по аналогии с сервером времени, которая бы каждому
// подключившемуся клиенту показывала раз в 3 секунды случайную Go-поговорку. Поговорки возьмите с сайта.
//
// * Служба должна поддерживать множественные одновременные подключения.
// * Служба не должна завершать соединение с клиентом.
// * Вы должны проверить работу приложения с помощью telnet.
//
// Егор Логинов, GO-11, SF, Модуль 35.8.1
package main

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
)

// База поговорок для ответа сервера
var proverbs = []string{
	"Don't communicate by sharing memory, share memory by communicating.",
	"Concurrency is not parallelism.",
	"Channels orchestrate; mutexes serialize.",
	"The bigger the interface, the weaker the abstraction.",
	"Make the zero value useful.",
	"interface{} says nothing.",
	"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
	"A little copying is better than a little dependency.",
	"Syscall must always be guarded with build tags.",
	"Cgo must always be guarded with build tags.",
	"Cgo is not Go.",
	"With the unsafe package there are no guarantees.",
	"Clear is better than clever.",
	"Reflection is never clear.",
	"Errors are values.",
	"Don't just check errors, handle them gracefully.",
	"Design the architecture, name the components, document the details.",
	"Documentation is for users.",
	"Don't panic.",
}

// Таймаут отдачи поговорки в секундах
const resTime int = 3

// Служба будет слушать запросы на всех IP-адресах компьютера на порту 12345.
// Например, 127.0.0.1:12345
const addr = "0.0.0.0:12345"

// Протокол сетевой службы.
const proto = "tcp4"

func main() {

	// Детерминация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Запуск сетевой службы по протоколу TCP на порту 12345.
	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Подключения обрабатываются в бесконечном цикле.
	// Иначе после обслуживания первого подключения сервер завершит работу.
	for {
		// Принимаем подключение.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Вызов обработчика подключения.
		// Отдельная go-рутина, иначе подключения будут
		// обрабатываться последоваельно (для единственного клиента)
		go handleConn(conn)
	}
}

// Обработчик. Вызывается для каждого соединения.
func handleConn(conn net.Conn) {
	// Закрытие соединения (по условию, сервер сам не закрывает соединение,
	// но предположим, что закрыть соединение может клиент, а то как-то некрасиво).
	defer conn.Close()

	// Канал синхронизации для закрытия соединения
	f := make(chan int)

	go func() {
		// Чтение сообщения от клиента.
		reader := bufio.NewReader(conn)
		b, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
		}
		// Удаление символов конца строки.
		msg := strings.TrimSuffix(string(b), "\n")
		// Если получили "finished" - закрываем соединение.
		if msg == "finished" {
			close(f)
		}
	}()

	// Основная петля отправки пословиц
	for {
		// Создаем таймер
		t := time.NewTimer(time.Second * time.Duration(resTime))
		select {
		// При срабатывании таймера отдаем клиенту пословицу
		case <-t.C:
			{
				conn.Write([]byte(getProverb() + "\n"))
			}
		case <-f:
			return
		}
	}
}

// getProverb возвращает случайную поговорку
func getProverb() string {

	return proverbs[rand.Intn(len(proverbs))]

}

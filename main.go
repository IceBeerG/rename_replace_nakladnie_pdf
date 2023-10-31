package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/ledongthuc/pdf"
)

var (
	pause, p, f, newname, NewPath, OldPath, newTitle, foo string
	searchString1, searchString2, searchString3           string
	nextChars1, nextChars2, nextChars3, nextChars4        int
	kernel32                                              = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleTitleW                                  = kernel32.NewProc("SetConsoleTitleW")
)

func main() {
	logFilePath := "nakladnie.log" // Имя файла для логирования ошибок
	logFilePath = filepath.Join(filepath.Dir(os.Args[0]), logFilePath)

	// Открываем файл для записи логов
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Ошибка открытия файла", err, getLine())
	}
	defer logFile.Close()

	// Устанавливаем файл в качестве вывода для логгера
	log.SetOutput(logFile)

	log.Println("\nНачало работы приложения")
	newTitle = "Переименование и распределение накладных(lukyanov_va)" // Вводим новое имя окна
	setConsoleTitle(newTitle)                                          // Устанавливаем новое имя окна

	fmt.Println("Программа перемещает накладные в формате pdf")
	fmt.Println("из папки запуска программа и распределяет по ответственным.")
	fmt.Println("Чтобы она это сделала, нужно подготовить txt файлы co списком OO.")
	fmt.Println("Файл all_OO.txt заполнить всеми OO c запятой в конце как в накладной.")
	fmt.Println("Файлы системотехников нужно называть по шаблону systech_XXXXX.txt,")
	fmt.Println("где XXXXX это например фамилия системотехника")
	fmt.Println("B конце файлов обязательно добавляем пустую последнюю строку, \nиначе не считается последнее значение")
	fmt.Println("В случае закрытия окна без оповещения об окончании работы, смотрите файл nakladnie.log")
	fmt.Println("\nВведите любой символ на нажмите Enter для начала работы")
	fmt.Println()
	fmt.Println()
	fmt.Scan(&pause)
	start := time.Now() // старт отсчета времени выполнения

	m, err := filepath.Glob("systech_*.txt") //выясняем количество файлов по шаблону
	if err != nil {
		log.Fatal("[ERR]", err, getLine())
	}

	renameFromPdf() // функция переименования файлов

	for _, val := range m { // перебираем все файлы systech_*.txt
		val = strings.Replace(val, "systech_", "", -1) // убираем из названия файла systech_
		val = strings.Replace(val, ".txt", "", -1)     // убираем из названия файла .txt
		_, err := os.Stat("./" + val)                  // проверяем существует ли папка
		if os.IsNotExist(err) {                        // если папка не существует
			err = os.Mkdir(val, 0777) // Создаем папку с правами доступа 0777
			if err != nil {
				log.Fatal("[ERR] Ошибка создания папки: ", val, err, getLine())
				return
			}
			fmt.Println("Папка " + val + " успешно создана")
		} else if err != nil {
			log.Fatal("[ERR]", err, getLine())
			return
		}

		// fmt.Println("Папка " + val + " уже существует")

		per(val) // функция распределения файлов по системотехникам
	}

	fmt.Printf("\n\n")
	duration := time.Since(start) //получаем время, прошедшее с момента старта
	fmt.Println(" Переименование и перемещение завершено!\n", "Время выполнения: ", duration, "\n\n", p)
	fmt.Println()
	fmt.Scan(&pause)
	log.Println("Программа закончила свою работу")
}

func findWord(content, searchString string, nextChars int) string {
	r := regexp.MustCompile(searchString + `(.{0,` + fmt.Sprint(nextChars) + `})`)
	match := r.FindStringSubmatch(content)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

func setConsoleTitle(title string) { // для смены заголовока программы
	ptrTitle, _ := syscall.UTF16PtrFromString(title)
	_, _, _ = procSetConsoleTitleW.Call(uintptr(unsafe.Pointer(ptrTitle)))
}

// получение текста из PDF файла
func ReadPlainTextFromPDF(pdfpath string) (text string, err error) {
	f, r, err := pdf.Open(pdfpath) // открываем файл PDF

	if err != nil {
		log.Fatal("[ERR] ", err, getLine())
		return
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText() // получаем текст в байтовом виде
	if err != nil {
		log.Fatal("[ERR] ", err, getLine())
		return
	}

	buf.ReadFrom(b)
	text = buf.String() // текст в привычном виде
	return
}

// функция переименования файла сверяя с текстом в файле pdf
func per(systech string) {
	var ss []string
	f = "systech_" + systech + ".txt" // файлы ответственных со списком ОО
	p = "Введите любой символ и Enter для выхода из программы"

	// Открываем текущую директорию
	dir, err := os.Open(".")
	if err != nil {
		log.Fatal("[ERR] He могу открыть директорию\n", err, getLine())
		return
	}
	defer dir.Close()

	// Получаем список файлов и папок
	files, err := dir.ReadDir(-1)
	if err != nil {
		log.Fatal("[ERR] He могу получить список файлов\n", err, getLine())
		return
	}

	// открываем файл "systech_systech.txt"
	fileoo, err := os.Open(f)
	if err != nil {
		log.Fatal("[ERR] He могу открыть файл: "+f+". \n", err, getLine())
		return
	}
	defer fileoo.Close()

	// читаем строки из файла "systech_systech.txt"
	reader := bufio.NewReader(fileoo)
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal("[ERR] Ошибка чтения строки в файле ", fileoo, "\n", err, getLine())
				return
			}
		}
		sline := len(line)
		line = line[:sline-2] // убираем лишние символы из названия ОО
		ss = append(ss, line) // добавляем значение в слайс со списком ОО
	}

	for _, file := range files { // перебираем все файлы в папке
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") { // если файл pdf
			for i := 0; i < len(ss); i++ { // перебираем все значения в слайсе со списком ОО
				if strings.Contains(file.Name(), ss[i]) { //проверяем есть ли совпадения в тексте pdf с названием ОО
					OldPath := "./" + file.Name()
					NewPath := "./" + systech + "/" + file.Name()
					fmt.Println(NewPath)
					log.Printf("Перемещаем файл %s в папку %s\n", file.Name(), systech)
					err := os.Rename(OldPath, NewPath) // перемещаем файл в папку системотехника
					if err != nil {
						log.Fatal("[ERR] Ошибка перемещения файла\n", err, getLine())
					}
					i = len(ss) // выходим из цикла так как нашли и переименовали файл
				}

			}
		}
	}
}

func renameFromPdf() {
	var ss []string                                               // слайс со списком всех объектов обслуживания
	p = "Введите любой символ и Enter для продолжения или выхода" // Сообщение для паузы по завершении перемещения или при ошибке
	foo = "all_OO.txt"                                            // файл со списком всех объектов обслуживания
	searchString1 = "НАКЛАДНАЯ № "                                //искомая строка
	nextChars1 = 15                                               // количество символов после найденой строки
	searchString2 = "нителя"                                      //искомая строка
	nextChars2 = 10                                               // количество символов после найденой строки
	searchString3 = "на внутреннее перемещение"                   //искомая строка
	nextChars3 = 15                                               // количество символов после найденой строки
	nextChars4 = 10                                               // количество символов после найденой строки

	//открываем папку в которой запускается программа
	dir, err := os.Open(".")
	if err != nil {
		log.Fatal("[ERR] Не могу открыть директорию\n", err, getLine())
		return
	}
	defer dir.Close()

	// Получаем список файлов и папок
	files, err := dir.ReadDir(-1)
	if err != nil {
		log.Fatal("[ERR] Не могу получить список файлов\n", err, getLine())
		return
	}

	// открываем файл со списком всех объектов обслуживания
	fileoo, err := os.Open(foo)
	if err != nil {
		log.Fatal("[ERR] He могу открыть файл: "+foo+". \n", err, getLine())
		return
	}
	defer fileoo.Close()

	// читаем строки из файла со списком всех объектов обслуживания
	reader := bufio.NewReader(fileoo)
	for {
		line, err := reader.ReadString('\n') // считываем строку

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal("[ERR] Ошибка считывания строки\n", err, getLine())
				return
			}
		}
		sline := len(line)
		line = line[:sline-2] // убираем лишние символы из названия ОО
		ss = append(ss, line) // добавляем значение в слайс со списком ОО
	}

	for _, file := range files { // перебираем все файлы
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") { // если файл PDF
			OldPath = "./" + file.Name()                      // имя файла источника
			content, err := ReadPlainTextFromPDF(file.Name()) // получаем текст из PDF файла
			if err != nil {
				log.Fatal("[ERR] Не могу прочитать файл\n", err, getLine())
			}

			if strings.Contains(content, "отпуск материалов") {
				result1 := findWord(content, searchString1, nextChars1)
				result2 := findWord(content, searchString2, nextChars2)
				newname = "(M-15) " + strings.Replace(result1, "/", "_", -1) + " от " + result2
			} else if strings.Contains(content, "основных средств") {
				result3 := findWord(content, searchString3, nextChars3)
				result4 := findWord(content, result3, nextChars4)
				newname = "(ОС-2) " + strings.Replace(result3, "/", "_", -1) + " от " + result4
			}

			for i := 0; i < len(ss); i++ { // перебираем значения в слайсе
				if strings.Contains(content, ss[i]) { // если в тексте PDF встречается значение из слайса
					s := strings.Replace(ss[i], ",", "", -1)
					newname = s + "_" + newname // новое имя файла с названием ОО
				}
			}
			NewPath = "./" + newname + ".pdf"
			if NewPath != OldPath {
				fmt.Println(newname) // выводим новое имя файла
				fmt.Println(NewPath)
				logString := fmt.Sprintf("Файл %s переименовываем в %s", file.Name(), newname)
				log.Println(logString)
				err = os.Rename(OldPath, NewPath) // переименовываем файл
				if err != nil {
					log.Fatal("[ERR] Ошибка переименования файла\n", err, getLine())
				} else {
					log.Println("Имена совпадают, переименовывание не требуется.")
				}
			}
		}
	}
}

// получение строки кода где возникла ошибка
func getLine() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

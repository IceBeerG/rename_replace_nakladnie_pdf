# rename_replace_nakladnie_pdf

Переименовывает накладные и раскладывает их по ответственным.
Поддерживаемые накладные:
- Унифицированная форма № ОС-2 Утверждена постановлением Госкомстата России от 21.01.2003 № 7
- Типовая межотраслевая форма № М-15 Утверждена постановлением Госкомстата России от 30.10.97 № 71а

Запуск

1. Скопируйте все файлы на свой локальный компьютер и распакуйте.
2. Установить Golang https://go.dev/
3. Открываем коммандную строку и переходим в распакованную папку.
4. Выполняем команду go build -o rename_replace_nakladnie.exe main.go.
5. В файл all_OO.txt добавляем все Объекты Обслуживания. В конце файла обязательно переходим на пустую строку.
6. Редактируем systech__XXXXX.txt. Вместо XXXXX вписываем имя или фамимлию ответственного. Внутри прописываем названия ОО, которые закрепленый за человеком. В конце файла обязательно переходим на пустую строку.
7. Создаем дополнительные файлы systech__XXXXX.txt если требуется.
8. Запускаем исполняемый файл rename_replace_nakladnie.exe.



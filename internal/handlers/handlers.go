package handlers

import (
    "html/template"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/Droff-hub/golf1-sprint6-final/internal/service"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    err = tmpl.Execute(w, nil)
    if err != nil {
        http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
        return
    }
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
        return
    }

    // ВАЖНО: тесты ожидают поле с именем "myFile"
    file, header, err := r.FormFile("myFile")
    if err != nil {
        http.Error(w, "Ошибка получения файла", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
        return
    }

    converted, err := service.Convert(string(data))
    if err != nil {
        http.Error(w, "Ошибка конвертации", http.StatusInternalServerError)
        return
    }

    // Генерируем имя для выходного файла
    ext := filepath.Ext(header.Filename)
    timestamp := time.Now().UTC().String()
    timestamp = strings.Map(func(r rune) rune {
        if r == ':' || r == ' ' {
            return '_'
        }
        return r
    }, timestamp)
    outputFilename := "converted_" + timestamp + ext

    // Сохраняем результат в файл
    outFile, err := os.Create(outputFilename)
    if err != nil {
        http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
        return
    }
    defer outFile.Close()

    _, err = outFile.WriteString(converted)
    if err != nil {
        http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
        return
    }

    // Возвращаем результат конвертации (обрезаем пробелы/переносы)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(strings.TrimSpace(converted)))
}
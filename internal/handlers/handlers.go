package handlers

import (
    "fmt"
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
    fmt.Println("=== UploadHandler hit ===")

    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        fmt.Println("ParseMultipartForm error:", err)
        http.Error(w, "Ошибка парсинга формы: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Выводим все ключи формы
    if r.MultipartForm != nil {
        fmt.Println("Form fields:")
        for key := range r.MultipartForm.File {
            fmt.Printf("  - %s\n", key)
        }
    } else {
        fmt.Println("r.MultipartForm is nil")
    }

    // Пробуем получить файл по имени "file"
    file, header, err := r.FormFile("file")
    if err != nil {
        fmt.Printf("FormFile(\"file\") error: %v\n", err)
        http.Error(w, "Ошибка получения файла: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer file.Close()
    fmt.Printf("File received: %s, size: ?\n", header.Filename)

    data, err := io.ReadAll(file)
    if err != nil {
        fmt.Println("ReadAll error:", err)
        http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
        return
    }
    fmt.Printf("Read %d bytes\n", len(data))

    converted, err := service.Convert(string(data))
    if err != nil {
        fmt.Println("Convert error:", err)
        http.Error(w, "Ошибка конвертации", http.StatusInternalServerError)
        return
    }
    fmt.Printf("Converted: %s\n", converted)

    // Сохраняем файл
    ext := filepath.Ext(header.Filename)
    timestamp := time.Now().UTC().String()
    timestamp = strings.Map(func(r rune) rune {
        if r == ':' || r == ' ' {
            return '_'
        }
        return r
    }, timestamp)
    outputFilename := "converted_" + timestamp + ext
    outFile, err := os.Create(outputFilename)
    if err != nil {
        fmt.Println("Create file error:", err)
        http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
        return
    }
    defer outFile.Close()
    outFile.WriteString(converted)

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(converted))
}
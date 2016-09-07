package main

import (
    "fmt"
    "flag"
    "runtime"
    "net/http"
    "os"
    "path/filepath"
    "io"
    "archive/tar"
)

func getArchiveExtension() string {
    switch runtime.GOOS {
    case "windows":
        return "msi"
    case "mac":
        return "pkg"
    default:
        return "tar.gz"
    }
}

func downloadGoVersion(dest string, version string) {
    var extension string
    extension = getArchiveExtension()
    var filename string
    filename = fmt.Sprintf("go%s.%s-%s.%s", version, runtime.GOOS, runtime.GOARCH, extension)
    out, err := os.Create(dest + string(filepath.Separator) + "tmp" + string(filepath.Separator) + filename)
    if err != nil {
        fmt.Println("Error on archive file:", err)
        os.Exit(1)
    }
    defer out.Close()

    var url string
    url = fmt.Sprintf("https://storage.googleapis.com/golang/go%s.%s-%s.%s", version, runtime.GOOS, runtime.GOARCH, extension)
    fmt.Println("Download Golang Version: ", url)
    resp, err2 := http.Get(url)
    defer resp.Body.Close()
    if err2 != nil {
        fmt.Println("Version unavailable:", err2)
        os.Exit(1)
    }
    _, err3 := io.Copy(out, resp.Body)
    if err3 != nil {
        fmt.Println("Bad Archive:", err3)
        os.Exit(1)
    }
}


func createEnvDir(dest string) {
    err := os.Mkdir(dest, 0777)
    if err != nil {
        fmt.Println("Error:", err)
    }

    err2 := os.Mkdir(dest + string(filepath.Separator) + "bin", 0777)
    if err2 != nil {
        fmt.Println("Error:", err2)
    }

    err3 := os.Mkdir(dest + string(filepath.Separator) + "lib", 0777)
    if err3 != nil {
        fmt.Println("Error:", err3)
    }

    err4 := os.Mkdir(dest + string(filepath.Separator) + "tmp", 0777)
    if err4 != nil {
        fmt.Println("Error:", err4)
    }

    err5 := os.Mkdir(dest + string(filepath.Separator) + "lib" + "go", 0777)
    if err5 != nil {
        fmt.Println("Error:", err5)
    }
}

func setupGoCompiler(dest string, version string) {
    fmt.Println("Setup a go install")
    var extension string
    extension = getArchiveExtension()
    var filename string
    filename = fmt.Sprintf("go%s.%s-%s.%s", version, runtime.GOOS, runtime.GOARCH, extension)
    var tarball string
    tarball = filepath.Join(dest, "tmp", filename)
    reader, err := os.Open(tarball)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    defer reader.Close()
    tarReader := tar.NewReader(reader)
    for {
        header, err2 := tarReader.Next()
        if err == io.EOF {
            break
        } else if err2 != nil {
            fmt.Println("Error2:", err2)
            os.Exit(1)
        }
        fmt.Println("Extracting:", header.Name)
        var target string
        target = filepath.Join(dest, "lib", "go", header.Name)
        path := filepath.Join(target, header.Name)
        info := header.FileInfo()
        if info.IsDir() {
            if err = os.MkdirAll(path, info.Mode()); err != nil {
                fmt.Println("Error3:", err)
                os.Exit(1)
            }
            continue
        }
        file, err:= os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
        if err != nil {
            fmt.Println("Error4:", err)
            os.Exit(1)
        }
        defer file.Close()
        _, err =io.Copy(file, tarReader)
        if err != nil {
            fmt.Println("Error5:", err)
            os.Exit(1)
        }
    }
}

func main() {
    var goversion string
    var dest string
    flag.StringVar(&goversion, "go-version", "1.7", "Golang version to use in the virtual environment")
    flag.StringVar(&dest, "dest", "", "virtual environment base directory")
    flag.Parse()
    cwd, _ := os.Getwd()
    fmt.Println("Setup a new virtualenv named", dest, "in", cwd)
    createEnvDir(dest)
    downloadGoVersion(dest, goversion)
    setupGoCompiler(dest, goversion)
}

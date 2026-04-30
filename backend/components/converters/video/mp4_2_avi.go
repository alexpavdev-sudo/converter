package video

import (
	"converter/components/converters"
	"converter/entities"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

type MP4ToAVIHandler struct {
	converters.BaseConverter
	outputPath string
	ffmpegPath string
}

func NewMP4ToAVIHandler(outputPath string) *MP4ToAVIHandler {
	return &MP4ToAVIHandler{
		outputPath: outputPath,
		ffmpegPath: findFFmpeg(),
	}
}

func (h *MP4ToAVIHandler) GetOutputPath() string {
	return h.outputPath
}

func (h *MP4ToAVIHandler) Rollback() error {
	return converters.BaseConverter{}.Rollback(h.GetOutputPath())
}

func findFFmpeg() string {
	paths := []string{
		"/usr/bin/ffmpeg",
		"/usr/local/bin/ffmpeg",
		"/usr/bin/ffmpeg.exe",
		"ffmpeg",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "ffmpeg"
}

func (h *MP4ToAVIHandler) Convert(file entities.File, perm fs.FileMode) (int64, error) {
	var size int64 = 0
	inputPath := file.PathFull()
	// Проверяем существование входного файла
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return size, fmt.Errorf("input file not found: %s", inputPath)
	}

	// Проверяем, что ffmpeg доступен
	if _, err := exec.LookPath(h.ffmpegPath); err != nil {
		return size, fmt.Errorf("ffmpeg not found: %v", err)
	}

	// Создаем пустой файл с нужными правами
	f, err := os.OpenFile(h.GetOutputPath(), os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return size, fmt.Errorf("failed to create output file: %v", err)
	}
	f.Close()

	// Параметры конвертации MP4 в AVI
	// Используем Xvid для видео (хорошее качество) и MP3 для аудио
	args := []string{
		"-i", inputPath, // входной файл
		"-c:v", "libxvid", // видео кодек Xvid
		"-q:v", "5", // качество видео (2-31, 2 - наилучшее, 5 - хорошее)
		"-c:a", "mp3", // аудио кодек MP3
		"-b:a", "192k", // битрейт аудио 192 kbps
		"-ar", "44100", // частота дискретизации 44.1 kHz
		"-ac", "2", // стерео
		"-y", // перезаписывать выходной файл
		h.GetOutputPath(),
	}

	cmd := exec.Command(h.ffmpegPath, args...)

	// Получаем вывод команды для отладки
	stdout, err := cmd.StderrPipe()
	if err != nil {
		return size, fmt.Errorf("failed to get stderr pipe: %v", err)
	}

	// Запускаем команду
	if err := cmd.Start(); err != nil {
		return size, fmt.Errorf("failed to start ffmpeg: %v", err)
	}

	// Читаем вывод (для логирования)
	output := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if n > 0 {
			output = append(output, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	// Ждем завершения команды
	if err := cmd.Wait(); err != nil {
		return size, fmt.Errorf("ffmpeg conversion failed: %v, output: %s", err, string(output))
	}

	// Проверяем, что выходной файл создан
	fileInfo, err := os.Stat(h.GetOutputPath())
	if err != nil {
		return size, fmt.Errorf("output file was not created: %v", err)
	}

	// Устанавливаем права доступа
	if err := os.Chmod(h.GetOutputPath(), perm); err != nil {
		fmt.Printf("Warning: could not set permissions: %v\n", err)
	}

	size = fileInfo.Size()

	// Проверка, что файл не пустой
	if size == 0 {
		return size, fmt.Errorf("converted file is empty")
	}

	return size, nil
}

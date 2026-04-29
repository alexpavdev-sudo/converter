package video

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

type MP4ToAVIHandler struct {
	ffmpegPath string
}

func NewMP4ToAVIHandler() *MP4ToAVIHandler {
	return &MP4ToAVIHandler{
		ffmpegPath: findFFmpeg(),
	}
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

func (h *MP4ToAVIHandler) Convert(inputPath, outputPath string, perm fs.FileMode) (int64, error) {
	var size int64 = 0

	// Проверяем существование входного файла
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return size, fmt.Errorf("input file not found: %s", inputPath)
	}

	// Проверяем, что ffmpeg доступен
	if _, err := exec.LookPath(h.ffmpegPath); err != nil {
		return size, fmt.Errorf("ffmpeg not found: %v", err)
	}

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
		outputPath,
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
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return size, fmt.Errorf("output file was not created: %v", err)
	}

	// Устанавливаем права доступа
	if err := os.Chmod(outputPath, perm); err != nil {
		fmt.Printf("Warning: could not set permissions: %v\n", err)
	}

	size = fileInfo.Size()

	// Проверка, что файл не пустой
	if size == 0 {
		return size, fmt.Errorf("converted file is empty")
	}

	return size, nil
}

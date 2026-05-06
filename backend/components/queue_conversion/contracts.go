package queue_conversion

type ConverterQueue interface {
	Push(fileId uint) error
}

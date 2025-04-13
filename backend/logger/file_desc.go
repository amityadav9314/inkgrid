package logger

type FileDescriptor uint8

const (
	UNDEFINED_FILE_DESCRIPTOR FileDescriptor = iota
	STANDARD_INPUT
	STANDARD_OUTPUT
	STANDARD_ERROR
	LOG_FILE
)

func (fileDescriptor FileDescriptor) String() string {
	switch fileDescriptor {
	case UNDEFINED_FILE_DESCRIPTOR:
		return "UNDEFINED_FILE_DESCRIPTOR"
	case STANDARD_INPUT:
		return "stdin"
	case STANDARD_OUTPUT:
		return "stdout"
	case STANDARD_ERROR:
		return "stderr"
	case LOG_FILE:
		return "logfile"
	default:
		return ""
	}
}

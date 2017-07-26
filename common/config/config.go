package config


type Config interface {
	Import() error
}


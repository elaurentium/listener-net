package utils

import (
	"fmt"
	"io"
	"os"
)


type ProgressReader struct {
	io.Reader
	Total int64
	Length int64
	Progress int64
}


func (p *ProgressReader) Read(buf []byte) (int, error) {
	n,err := p.Reader.Read(buf)

	if n > 0 {
		p.Total += int64(n)
		percentage := float64(p.Total) / float64(p.Length) * float64(100)
		if percentage - float64(p.Progress) > 2 {
			p.Progress = int64(percentage)
			fmt.Fprintf(os.Stderr, "\r%.2f%%", float64(p.Progress))
			if p.Progress > 98.0 {
				p.Progress = 100
				fmt.Fprintf(os.Stderr, "\r%.2f%%\n", float64(p.Progress))
			}
		}
	}

	return n, err
}
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"syscall"
)

// TODO: Make it work if the file ends with \n
// TODO: Make it streaming

func main() {
	nLines := flag.Int("n", 10, "Number of lines from end that we would like to print")
	flag.Parse()

	var filePath string
	if len(flag.Args()) > 0 {
		filePath = flag.Args()[0]
	} else {
		fmt.Println("Please pass valid file in the arg list")
		return
	}

	var fStat syscall.Stat_t
	err := syscall.Stat(filePath, &fStat)
	if err != nil {
		fmt.Printf("Error gets stat of file:[%s]. Error: [%s] \n", filePath, err)
		return
	}
	if fStat.Mode&syscall.S_IFMT != syscall.S_IFREG {
		fmt.Printf("Please input regular file for tail command. \n")
		return
	}

	fSize := fStat.Size

	fd, err := syscall.Open(filePath, syscall.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Error in opening file: [%s]. Error: [%s] \n", filePath, err)
		return
	}

	offSetNum := 1
	var ans int64
	for currLines := 0; currLines < *nLines; offSetNum++ {
		currBlock := int64((100 * offSetNum))
		var currOffSet int64
		if fSize >= currBlock {
			currOffSet = fSize - currBlock
		} else {
			currOffSet = 0
		}
		_, err := syscall.Seek(fd, currOffSet, io.SeekStart)
		if err != nil {
			log.Printf("Error in seek. Error: [%s] \n", err)
			return
		}

		rBuffer := make([]byte, 100)
		nB, err := syscall.Read(fd, rBuffer)
		if err != nil {
			fmt.Printf("Error in reading file: [%s] \n", err)
		}
		if nB == 0 {
			break
		}
		for i := nB - 1; i >= 0; i-- {
			v := rBuffer[i]
			if v == '\n' {
				currLines++
			}
			if currLines == *nLines {
				ans = currOffSet + int64(i+1)
				break
			}
		}
		syscall.Seek(fd, 0, 0) // Go Back to the start of the file
	}

	ansBytes := make([]byte, fSize-ans)
	syscall.Seek(fd, ans, io.SeekStart)
	_, _ = syscall.Read(fd, ansBytes)
	fmt.Print(string(ansBytes))

	defer syscall.Close(fd)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strconv"
	"syscall"
	"time"
)

type FileInfo struct {
	inode            uint64
	perm             string
	refCount         int
	owner            string
	userGroup        string
	size             int64
	lastModifiedTime string
	fileName         string
}

func getOwnerName(ownerId uint32) string {
	u, err := user.LookupId(strconv.FormatUint(uint64(ownerId), 10))
	if err != nil {
		log.Println("Couldnt get username from ownerId")
		return ""
	}
	return u.Name
}

func getFileLastModifiedTime(mTime syscall.Timespec) string {
	hMTime := time.Unix(mTime.Sec, mTime.Nsec)
	return hMTime.Format("Jan 2 15:04")
}

func getFileTypeByte(perm uint16) string {
	flType := perm & syscall.S_IFMT
	switch flType {
	case syscall.S_IFREG:
		return "-"
	case syscall.S_IFDIR:
		return "d"
	case syscall.S_IFLNK:
		return "l"
	case syscall.S_IFIFO:
		return "p"
	case syscall.S_IFSOCK:
		return "s"
	case syscall.S_IFBLK:
		return "b"
	case syscall.S_IFCHR:
		return "c"
	}
	return "-"
}

func getPermissionString(mode uint16) string {
	var perms string
	// Check if its a dir, file or symlink
	perms += getFileTypeByte(mode)
	// Check Owner Perm
	// Owner
	if mode&0400 != 0 {
		perms += "r"
	} else {
		perms += "-"
	}
	if mode&0200 != 0 {
		perms += "w"
	} else {
		perms += "-"
	}
	if mode&0100 != 0 {
		perms += "x"
	} else {
		perms += "-"
	}

	// Group
	if mode&0040 != 0 {
		perms += "r"
	} else {
		perms += "-"
	}
	if mode&0020 != 0 {
		perms += "w"
	} else {
		perms += "-"
	}
	if mode&0010 != 0 {
		perms += "x"
	} else {
		perms += "-"
	}

	// Others
	if mode&0004 != 0 {
		perms += "r"
	} else {
		perms += "-"
	}
	if mode&0002 != 0 {
		perms += "w"
	} else {
		perms += "-"
	}
	if mode&0001 != 0 {
		perms += "x"
	} else {
		perms += "-"
	}

	return perms

}

func getUserGroupName(groupId uint32) string {
	g, err := user.LookupGroupId(strconv.FormatUint(uint64(groupId), 10))
	if err != nil {
		log.Println("Couldnt get groupName from groupId")
		return ""
	}
	return g.Name
}

func main() {
	lArg := flag.Bool("l", false, "List detailed information about files and directory")
	flag.Parse()

	dir := flag.Arg(0)

	if len(dir) == 0 {
		dir = "."
	}

	// TODO: For SymLink dont go inside the symlink

	de, err := os.ReadDir(dir)

	if err != nil {
		log.Printf("Error in reading directory. Error: [%s] \n", err)
		return
	}

	ans := make([]string, 0)

	for _, dirEnt := range de {
		if !*lArg {
			ans = append(ans, dirEnt.Name())
		} else {
			// Call stat
			// inode permission refCount owner group size lastModifiedTime name
			var fStat syscall.Stat_t
			var fileInfo FileInfo
			err := syscall.Stat(path.Join(dir, dirEnt.Name()), &fStat)
			if err != nil {
				log.Printf("Error in getting stat of file. Error: [%s] \n", err)
				return
			}

			fileInfo.fileName = dirEnt.Name()
			fileInfo.inode = fStat.Ino
			fileInfo.size = fStat.Size
			fileInfo.refCount = int(fStat.Nlink)
			fileInfo.owner = getOwnerName(fStat.Uid)
			fileInfo.userGroup = getUserGroupName(fStat.Gid)
			fileInfo.lastModifiedTime = getFileLastModifiedTime(fStat.Mtimespec)
			fileInfo.perm = getPermissionString(fStat.Mode)

			str := fmt.Sprintf("%d %-10s %d %-15s %-5s %-10d %s %s", fileInfo.inode, fileInfo.perm, fileInfo.refCount, fileInfo.owner, fileInfo.userGroup, fileInfo.size, fileInfo.lastModifiedTime, fileInfo.fileName)
			ans = append(ans, str)

			fmt.Println(str)

		}
	}
	if !*lArg {
		for _, e := range ans {
			fmt.Println(e)
		}
	}

}

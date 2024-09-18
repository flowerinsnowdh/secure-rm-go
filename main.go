/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package main

import (
	"crypto/rand"
	"fmt"
	"github.com/flowerinsnowdh/securerm-go/common"
	"os"
	"path"
	"strings"
)

var args *common.Args

func init() {
	args = &common.Args{
		Files: make([]string, 0, len(os.Args)),
	}

	for _, arg := range os.Args[1:] {
		switch strings.ToLower(arg) {
		case "-h", "--help":
			printUsage()
			return
		case "-r", "--recursion":
			args.Recursion = true
		case "-v", "--verbose":
			args.Verbose = true
		case "-vr", "-rv":
			args.Recursion = true
			args.Verbose = true
		default:
			args.Files = append(args.Files, arg)
		}
	}
}

func main() {
	for _, file := range args.Files {
		secureRemove(file) // 递归删除所有参数中的目录
	}
}

func secureRemove(baseFile string) {
	var err error
	var info os.FileInfo
	if info, err = os.Stat(baseFile); err != nil { // 获取文件详情失败
		_, _ = fmt.Fprintln(os.Stderr, "[81] failed to remove '", baseFile, "':", err)
		return
	}
	if info.IsDir() { // 是目录
		if !args.Recursion { // 没有使用递归 flag
			_, _ = fmt.Fprintln(os.Stderr, baseFile, "is a directory and can only remove by --recursion")
			os.Exit(1)
		}
		// 列出目录下的所有文件
		var entries []os.DirEntry
		if entries, err = os.ReadDir(baseFile); err != nil { // 失败，跳出当前方法
			_, _ = fmt.Fprintln(os.Stderr, "[82] failed to remove '", baseFile, "':", err)
			return
		}
		for _, entry := range entries { // 遍历目录下的文件，逐一删除
			var subF os.FileInfo
			if subF, err = entry.Info(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "[83] failed to remove '", baseFile, "':", err)
				return
			}
			secureRemove(path.Join(baseFile, subF.Name()))
		}
		if err = os.Remove(baseFile); err != nil { // 最终删除目录
			_, _ = fmt.Fprintln(os.Stderr, "[84] failed to remove directory '", baseFile, "':", err)
			return
		}
	} else {
		err = first(baseFile, info.Size())
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "\n[85] failed to fill '", baseFile, "' with random:", err)
			return
		}
		err = second(baseFile, info.Size())
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "\n[86] failed to fill '", baseFile, "' with 1:", err)
			return
		}
		err = third(baseFile, info.Size())
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "\n[87] failed to remove '", baseFile, "' with 0:", err)
			return
		}
		if err = os.Remove(baseFile); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "\n[88] failed to delete '", baseFile, "':", err)
			return
		}
		if args.Verbose {
			fmt.Println("OK")
		}
	}
}

func fill(file string, n int64, msg string, fillAction func([]byte)) error {
	if args.Verbose {
		fmt.Print(msg)
	}
	var f *os.File
	var err error
	if f, err = os.OpenFile(file, os.O_WRONLY, 0); err != nil {
		return err
	}
	var buffer []byte = make([]byte, 4096)
	for i := int64(0); i < n; i += 4096 {
		fillAction(buffer)
		if (n - i) < 4096 {
			_, err = f.Write(buffer[:(n - i)])
		} else {
			_, err = f.Write(buffer)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func first(file string, n int64) error {
	return fill(file, n, "removing '"+file+"'.", func(data []byte) {
		_, _ = rand.Read(data)
	})
}

func second(file string, n int64) error {
	return fill(file, n, ".", func(data []byte) {
		for i := 0; i < len(data); i++ {
			data[i] = 0xFF
		}
	})
}

func third(file string, n int64) error {
	return fill(file, n, ".", func(data []byte) {
		for i := 0; i < len(data); i++ {
			data[i] = 0
		}
	})
}

func printUsage() {
	fmt.Println(os.Args[0], "[flags] [file] [file2] [file3]...")
	fmt.Println("flags:")
	fmt.Println("    -h, --help: print this help menu")
	fmt.Println("    -r, --recursion: recursion remove")
	fmt.Println("    -v, --verbose: print verbose")
}

package lib

import (
	"github.com/golang/glog"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type FileSystem struct {
	Config
	pathfs.FileSystem
}

func (self *FileSystem) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	glog.Info("FileSystem#GetAttr(): name = ", name)
	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	} else if _, ok := self.Config[name]; ok {
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0755,
			Size: uint64(123),
		}, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (self *FileSystem) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	if name == "" {
		c = []fuse.DirEntry{}
		for k := range self.Config {
			c = append(c, fuse.DirEntry{Name: k, Mode: fuse.S_IFREG})
		}
		return c, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (self *FileSystem) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	glog.Info("FileSystem#Open()")
	if v, ok := self.Config[name]; ok {
		locals := Locals{Image: v.Image}
		return nodefs.NewDataFile([]byte(Render(&locals))), fuse.OK
	}
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}
	return nil, fuse.ENOENT
}

func MountFileSystem(conf Config, mnt string) (*fuse.Server, error) {
	glog.Info("FileSystem#OpenDir()")
	nfs := pathfs.NewPathNodeFs(&FileSystem{Config: conf, FileSystem: pathfs.NewDefaultFileSystem()}, nil)
	server, _, err := nodefs.MountRoot(mnt, nfs.Root(), nil)
	return server, err
}

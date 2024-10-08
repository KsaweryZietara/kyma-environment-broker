package main

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kyma-project/kyma-environment-broker/internal/schemamigrator/mocks"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func Test_migrationScript_copyFile(t *testing.T) {
	t.Run("Should return error while opening source file fails", func(t *testing.T) {
		// given
		mfs := &mocks.FileSystem{}
		mfs.On("Open", "src").Return(nil, fmt.Errorf("failed to open file"))
		ms := &migrationScript{
			fs: mfs,
		}

		// when
		err := ms.copyFile("src", "dst")

		// then
		assert.Error(t, err)
	})
	t.Run("Should return error while creating destination file fails", func(t *testing.T) {
		// given
		mfs := &mocks.FileSystem{}
		mfs.On("Open", "src").Return(&os.File{}, nil)
		mfs.On("Create", "dst").Return(nil, fmt.Errorf("failed to create file"))
		ms := &migrationScript{
			fs: mfs,
		}

		// when
		err := ms.copyFile("src", "dst")

		// then
		assert.Error(t, err)
	})
	t.Run("Should return error while copying file fails", func(t *testing.T) {
		// given
		mfs := &mocks.FileSystem{}
		mfs.On("Open", "src").Return(&os.File{}, nil)
		mfs.On("Create", "dst").Return(&os.File{}, nil)
		mfs.On("Copy", &os.File{}, &os.File{}).Return(int64(0), fmt.Errorf("failed to copy file"))
		ms := &migrationScript{
			fs: mfs,
		}

		// when
		err := ms.copyFile("src", "dst")

		// then
		assert.Error(t, err)
	})
	t.Run("Should return error while returning FileInfo fails", func(t *testing.T) {
		// given
		mfs := &mocks.FileSystem{}
		mfi := &mocks.MyFileInfo{}
		mfs.On("Open", "src").Return(&os.File{}, nil)
		mfs.On("Create", "dst").Return(&os.File{}, nil)
		mfs.On("Copy", &os.File{}, &os.File{}).Return(int64(65), nil)
		mfs.On("Stat", "src").Return(mfi, fmt.Errorf("failed to get FileInfo"))
		ms := &migrationScript{
			fs: mfs,
		}

		// when
		err := ms.copyFile("src", "dst")

		// then
		assert.Error(t, err)
	})
	t.Run("Should return error while changing the mode of the file fails", func(t *testing.T) {
		// given
		mfs := &mocks.FileSystem{}
		mfi := &mocks.MyFileInfo{}
		mfs.On("Open", "src").Return(&os.File{}, nil)
		mfs.On("Create", "dst").Return(&os.File{}, nil)
		mfs.On("Copy", &os.File{}, &os.File{}).Return(int64(65), nil)
		mfs.On("Stat", "src").Return(mfi, nil)
		mfi.On("Mode").Return(fs.FileMode(0666))
		mfs.On("Chmod", "dst", fs.FileMode(0666)).Return(fmt.Errorf("failed to change file mode"))
		ms := &migrationScript{
			fs: mfs,
		}

		// when
		err := ms.copyFile("src", "dst")

		// then
		assert.Error(t, err)
	})
	t.Run("Should not return error", func(t *testing.T) {
		// given
		mfs := &mocks.FileSystem{}
		mfi := &mocks.MyFileInfo{}
		mfs.On("Open", "src").Return(&os.File{}, nil)
		mfs.On("Create", "dst").Return(&os.File{}, nil)
		mfs.On("Copy", &os.File{}, &os.File{}).Return(int64(65), nil)
		mfs.On("Stat", "src").Return(mfi, nil)
		mfi.On("Mode").Return(fs.FileMode(0666))
		mfs.On("Chmod", "dst", fs.FileMode(0666)).Return(nil)
		ms := &migrationScript{
			fs: mfs,
		}

		// when
		err := ms.copyFile("src", "dst")

		// then
		assert.Nil(t, err)
	})
}

func Test_migrationScript_copyDir(t *testing.T) {
	t.Run("Should not return error", func(t *testing.T) {
		// given
		mfs := &mocks.FileSystem{}
		mfs.On("ReadDir", "src").Return([]fs.DirEntry{}, nil)
		ms := &migrationScript{
			fs: mfs,
		}

		// when
		err := ms.copyDir("src", "dst")

		// then
		assert.Nil(t, err)

	})
	t.Run("Should return error while reading directory fails", func(t *testing.T) {
		// given
		mfs := &mocks.FileSystem{}
		mfs.On("ReadDir", "src").Return(nil, fmt.Errorf("failed to read directory"))
		ms := &migrationScript{
			fs: mfs,
		}

		// when
		err := ms.copyDir("src", "dst")

		// then
		assert.Error(t, err)

	})
}

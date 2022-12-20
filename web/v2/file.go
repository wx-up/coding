package v2

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// FileUploader 结构体内部定义各种参数
type FileUploader struct {
	// FileField 对应表单提交的字段名
	FileField string

	// DstPathFunc 计算文件存储的目标路径（ 文件名出现冲突会导致覆盖，交给用户去生成文件名 ）
	DstPathFunc func(part *multipart.FileHeader) string
}

func (f *FileUploader) Handle() HandleFunc {
	// 这里可以额外做一些检测
	// if f.FileField == "" {
	// 	// 这种方案默认值我其实不是很喜欢
	// 	// 因为我们需要教会用户说，这个 file 是指什么意思
	// 	f.FileField = "file"
	// }
	return func(ctx *Context) {
		src, srcHeader, err := ctx.Req.FormFile(f.FileField)
		if err != nil {
			ctx.RespStatusCode = 400
			ctx.RespData = []byte("上传失败，未找到数据")
			return
		}
		defer src.Close()
		dst, err := os.OpenFile(f.DstPathFunc(srcHeader),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			return
		}
		defer dst.Close()

		// 从 src 中以 buf 为单元拷贝数据到 dst 直到拷贝完
		// buf 为 nil 时，有默认值为 32kb
		// 这个值会影响性能（ 太大的话占用内存高，上传效率高，太小的话占用内存底，上传效率低 ）
		_, err = io.CopyBuffer(dst, src, nil)
		if err != nil {
			ctx.RespStatusCode = 500
			ctx.RespData = []byte("上传失败")
			return
		}
		ctx.RespData = []byte("上传成功")
	}
}

// HandleFunc 这种设计方案也是可以的，但是不如上一种灵活。
// 它可以直接用来注册路由，而上一种可以在返回 HandleFunc 之前检测一下传入的字段
func (f *FileUploader) HandleFunc(ctx *Context) {
	src, srcHeader, err := ctx.Req.FormFile(f.FileField)
	if err != nil {
		ctx.RespStatusCode = 400
		ctx.RespData = []byte("上传失败，未找到数据")
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(f.DstPathFunc(srcHeader),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		ctx.RespStatusCode = 500
		ctx.RespData = []byte("上传失败")
		return
	}
	defer dst.Close()

	_, err = io.CopyBuffer(dst, src, nil)
	if err != nil {
		ctx.RespStatusCode = 500
		ctx.RespData = []byte("上传失败")
		return
	}
	ctx.RespData = []byte("上传成功")
}

type FileDownloader struct {
	Dir string
}

func (fd *FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		// http://localhost:8081/download?file=xxx.jpg
		filename, err := ctx.ParamValue("file")
		if err != nil {
			ctx.RespData = []byte("文件不存在")
			ctx.RespStatusCode = http.StatusNotFound
			return
		}

		// filepath.Clean 处理传入的 path
		path := filepath.Join(fd.Dir, filepath.Clean(filename))

		// os.IsNotExist
		if _, err = os.Stat(path); err != nil {
			ctx.RespData = []byte("文件不存在")
			ctx.RespStatusCode = http.StatusNotFound
			return
		}

		// 获取路径的最后一段，在当前场景下也就是文件名
		fn := filepath.Base(path)

		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		// 文件的打开、读取、返回给前端等操作委托给 http 包，不自己实现
		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}

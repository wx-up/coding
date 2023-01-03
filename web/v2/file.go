package v2

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/golang-lru/v2"
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

// StaticResource 静态资源处理
type StaticResource struct {
	// 静态资源目录
	dir string
	// 文件名后缀对应 content-type
	extensionContentTypeMap map[string]string

	// 缓存
	cache *lru.Cache[string, *fileCacheItem]

	maxFileSize int // 单个文件的最大大小
}

type fileCacheItem struct {
	fileName    string
	fileSize    int
	contentType string
	data        []byte
}

type StaticResourceOption func(*StaticResource)

// WithMoreExtension 更多 extension type
func WithMoreExtension(extMap map[string]string) StaticResourceOption {
	return func(resource *StaticResource) {
		for k, v := range extMap {
			resource.extensionContentTypeMap[k] = v
		}
	}
}

// WithFileCache 静态文件将会被缓存
// maxFileSizeThreshold 超过这个大小的文件，就被认为是大文件，我们将不会缓存
// maxCacheFileCnt 最多缓存多少个文件
// 所以我们最多缓存 maxFileSizeThreshold * maxCacheFileCnt（ 控制内存消耗，否则会增加内存的压力 ）
func WithFileCache(maxFileSizeThreshold int, maxCacheFileCnt int) StaticResourceOption {
	return func(h *StaticResource) {
		c, err := lru.New[string, *fileCacheItem](maxCacheFileCnt)
		if err != nil {
			log.Printf("创建缓存失败，将不会缓存静态资源")
		}
		h.maxFileSize = maxFileSizeThreshold
		h.cache = c
	}
}

func NewStaticResource(dir string, opts ...StaticResourceOption) *StaticResource {
	handler := &StaticResource{
		dir: dir,
		extensionContentTypeMap: map[string]string{
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "image/pdf",
		},
	}
	for _, opt := range opts {
		opt(handler)
	}
	return handler
}

func (sr *StaticResource) Handler() HandleFunc {
	return func(ctx *Context) {
		// 获取请求路径上的文件名
		// 路由：/img/:file
		path, err := ctx.PathValue("file")
		if err != nil {
			ctx.RespData = []byte("文件不存在")
			ctx.RespStatusCode = http.StatusNotFound
			return
		}

		// 获取文件
		cacheItem, err := sr.readFile(path)
		if err != nil {
			ctx.RespData = []byte("文件不存在")
			ctx.RespStatusCode = http.StatusNotFound
			return
		}

		// 缓存并返回结果
		sr.cacheFile(cacheItem)
		ctx.Resp.Header().Set("Content-Type", cacheItem.contentType)
		ctx.Resp.Header().Set("Content-Length", fmt.Sprintf("%d", cacheItem.fileSize))
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = cacheItem.data
	}
}

func (sr *StaticResource) readFile(path string) (*fileCacheItem, error) {
	// 缓存中存在就直接返回
	if item, ok := sr.readFileFromCache(path); ok {
		return item, nil
	}

	// 判断是否是支持的文件类型
	ext := getFileExt(path)
	t, ok := sr.extensionContentTypeMap[ext]
	if !ok {
		return nil, errors.New("不支持的文件类型")
	}

	// 获取服务器上对应文件的地址
	path = filepath.Join(sr.dir, path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// 读取文件数据
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	item := &fileCacheItem{
		fileName:    path,
		fileSize:    len(data),
		contentType: t,
		data:        data,
	}
	return item, nil
}

// cacheFile 缓存文件
func (sr *StaticResource) cacheFile(item *fileCacheItem) {
	if sr.cache == nil {
		return
	}

	// 文件大小超过最大限制，则不缓存
	if item.fileSize > sr.maxFileSize {
		return
	}
	sr.cache.Add(item.fileName, item)
}

// readFileFromCache 从缓存中获取
func (sr *StaticResource) readFileFromCache(path string) (*fileCacheItem, bool) {
	if sr.cache == nil {
		return nil, false
	}
	return sr.cache.Get(path)
}

// getFileExt 获取文件的后缀
func getFileExt(name string) string {
	index := strings.LastIndex(name, ".")
	if index == len(name)-1 {
		return ""
	}
	return name[index+1:]
}
